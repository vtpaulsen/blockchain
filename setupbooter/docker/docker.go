package setupbooter

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"syscall"
	"testing"
)

type PeerContainer struct {
	id                int
	ownPublicKey      string
	targetIP          string
	targetPort        string
	shellInput        io.WriteCloser
	shellOutputBuffer *bytes.Buffer
}

var composeBuffer = new(strings.Builder)
var bufferUpdates = make(chan string)
var stopBufferUpdates = false

// Service name of docker compose service
var dockerComposeService = "peer"

// Image name of the built image running in the docker container
var dockerImageName = "dissy"

// Prefix of every docker container running the service
var dockerAttachPrefix = dockerImageName + "_" + dockerComposeService + "_"

// Maps from container id to a PeerContainer
var idToPeerContainer map[int]PeerContainer = make(map[int]PeerContainer)

// Maps from a container id to a public key
var idToPublicKey map[int]string = make(map[int]string)

type dockerTester struct {
	Test                    *testing.T
	Peers                   []PeerInstance
	OutputToStdout          bool
	MessageThrottleDuration int
}

func DockerTester(t *testing.T, peers []PeerInstance) dockerTester {
	peers = appendContainerIds(peers)
	peers = topologicalSort(peers)
	return dockerTester{
		Test:                    t,
		Peers:                   peers,
		OutputToStdout:          false,
		MessageThrottleDuration: 0,
	}
}

func bootEnvironment(dt dockerTester) {
	wfc := func() {
		waitForConnectivity(dt.Test, dt.Peers)
	}
	errFunc := func() {
		logOnError(errors.New("docker testing environment timed out on boot"), "Environment boot timeout")
	}

	// Wait for connectivity for max 30 seconds.
	tryUntilTimeout(wfc, errFunc, 30*1000)
}

// Executes a docker attach command
func dockerAttach(number int) *exec.Cmd {
	return exec.Command("docker", "attach", dockerAttachPrefix+strconv.Itoa(number))
}

/*
* Executes a docker compose command
* If genesisUsers is set to true, then the first 10 containers
* will use the 10 hardcoded pk/sk of the genesis block
 */
func dockerCompose(count int) *exec.Cmd {
	cmd := exec.Command("docker-compose", "--compatibility", "--project-name", dockerImageName, "up", "--scale", dockerComposeService+"="+strconv.Itoa(count))
	return cmd
}

func copyToComposeBuffer(stdout io.ReadCloser, outputToStdout bool) {
	reader := bufio.NewReader(stdout)
	for {
		str, err := reader.ReadString('\n')
		if err == io.EOF {
			return
		}
		logOnError(err, "Error copying to buffer")
		composeBuffer.WriteString(str)
		if outputToStdout {
			fmt.Print(str)
		}
		if !stopBufferUpdates {
			bufferUpdates <- str
		}
	}
}

// Append container ids to all peer instances
// by their natural given order in the array
func appendContainerIds(peers []PeerInstance) []PeerInstance {
	var result []PeerInstance
	for i, pi := range peers {
		result = append(result, PeerInstance{
			containerId:     i + 1,
			ShouldConnectTo: pi.ShouldConnectTo,
		})
	}
	return result
}

// Wait until everything is connected as instructed in peers []PeerInstance
func waitForConnectivity(t *testing.T, peers []PeerInstance) {
	// Wait for all containers to be ready (booted with all images up and running)
	waitForAllContainers(len(peers))
	// Start connecting all peers, once we know everything is up and running
	initializeConnections(peers)
}

// Call with main process where cmd = docker-compose process
func killAllContainers(cmd exec.Cmd) {
	// syscall to get pgid of main docker-compose process
	pgid, err := syscall.Getpgid(cmd.Process.Pid)
	if err == nil {
		// If pid is negative, but not -1, sig shall be sent to all processes
		// (excluding an unspecified set of system processes) whose process group ID
		// is equal to the absolute value of pid, and for which the process has permission to send a signal.
		// man page: https://linux.die.net/man/3/kill
		syscall.Kill(-pgid, syscall.SIGTERM) // (thus we give negative pgid)
		cmd.Wait()
	}
}

// Src: https://github.com/acarl005/stripansi/blob/master/stripansi.go
func removeAnsiEscChars(str string) string {
	const ansi = "[\u001B\u009B][[\\]()#;?]*(?:(?:(?:[a-zA-Z\\d]*(?:;[a-zA-Z\\d]*)*)?\u0007)|(?:(?:\\d{1,4}(?:;\\d{0,4})*)?[\\dA-PRZcf-ntqry=><~]))"
	var re = regexp.MustCompile(ansi)
	return re.ReplaceAllString(str, "")
}

// Waits until all containers are booted and ready to receive input
func waitForAllContainers(instances int) {
	countOfReadyContainers := 0
	for {
		bufferUpdate := <-bufferUpdates
		bufferUpdate = removeAnsiEscChars(bufferUpdate)
		// READY_MESSAGE is the message that the container prints when booted
		if strings.Contains(bufferUpdate, READY_MESSAGE) {
			//re1 := regexp.MustCompile("peer_([0-9]*)\\s*\\|\\s*" + READY_MESSAGE + "(.*)")
			//peerPK := re1.FindStringSubmatch(bufferUpdate)

			//containerId, err := strconv.Atoi(peerPK[1])
			//logOnError(err, "Failed to convert container id to int")

			//idToPublicKey[containerId] = peerPK[2]
			countOfReadyContainers++
		}
		// If we have all public keys as amount of instances, we are done
		//if len(idToPublicKey) == instances && countOfReadyContainers == instances {
		if countOfReadyContainers == instances {
			stopBufferUpdates = true
			return
		}
	}
}

// Connect all the peers that should be connected and get their ip:port.
func initializeConnections(peers []PeerInstance) {
	for i := 0; i < len(peers); i++ {
		cmd := dockerAttach(peers[i].containerId)
		stdin, err := cmd.StdinPipe()
		logOnError(err, "Error while creating stdin pipe")

		stdout, err := cmd.StdoutPipe()
		logOnError(err, "Error while creating stdout pipe")

		cmd.Start()

		initPeerConnection(peers[i], stdin, stdout)
	}

}

func initPeerConnection(peer PeerInstance, stdin io.WriteCloser, stdout io.ReadCloser) {

	// if shouldConnectTo is 0, then we shouldn't connect to a machine
	// thus we send \n\n to get our listener ip and port
	if peer.ShouldConnectTo == 0 {
		io.WriteString(stdin, "\n\n")
	} else {
		// Because peers are topologically sorted, we know that it is in the map, and that we can connect to it now
		peerToConnectTo := idToPeerContainer[peer.ShouldConnectTo]
		io.WriteString(stdin, peerToConnectTo.targetIP+"\n")
		io.WriteString(stdin, peerToConnectTo.targetPort+"\n")
	}

	ip, port := getIPandPort(stdout)
	peerContainer := PeerContainer{
		id:                peer.ShouldConnectTo,
		ownPublicKey:      idToPublicKey[peer.containerId],
		targetIP:          ip,
		targetPort:        port,
		shellInput:        stdin,
		shellOutputBuffer: new(bytes.Buffer),
	}

	// Keep buffering everything in the background until EOF
	go io.Copy(peerContainer.shellOutputBuffer, stdout)

	idToPeerContainer[peer.containerId] = peerContainer
}
