package setupbooter

import (
	"io"
	"strconv"
	"syscall"
)

type PeerInstance struct {
	containerId     int
	ShouldConnectTo int
}

// Start a multi container environment with given peers.
// Automatically connects peers after how []PeerInstance is constructed
func (dt dockerTester) StartMultiContainerEnv() {
	cmd := dockerCompose(len(dt.Peers))
	stdout, err := cmd.StdoutPipe()
	PrintIfError(err, "Error while creating stdout pipe for docker docker compose")

	// Copy the compose-cmd output to a buffer
	go copyToComposeBuffer(stdout, dt.OutputToStdout)
	// Sets process group id of child to Pgid
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	cmd.Start()

	// We kill all containers on test cleanup
	dt.Test.Cleanup(func() {
		defer killAllContainers(*cmd)
	})

	// Wait for all peers to be connected before we return
	bootEnvironment(dt)
}

// Send a raw message from container id
func (dt dockerTester) SendMessageFromContainer(containerId int, msg string) {
	peer := idToPeerContainer[containerId]
	_, err := io.WriteString(peer.shellInput, msg+"\n")
	logOnError(err, "Error while writing '"+msg+"' to stdin on container "+strconv.Itoa(containerId))
	// Lets throttle for the amount of ms given (default 0)
	throttle(dt.MessageThrottleDuration)
}

func (dt dockerTester) StopContainer(containerId int) {
	idToPeerContainer[containerId].shellInput.Write([]byte("quit\n"))
	delete(idToPeerContainer, containerId)
}
