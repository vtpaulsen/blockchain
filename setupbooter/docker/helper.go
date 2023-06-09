package setupbooter

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
	"time"
)

// Retrieves ip and port from a stdout
func getIPandPort(stdout io.ReadCloser) (string, string) {
	reader := bufio.NewReader(stdout)
	for {
		line, err := reader.ReadString('\n')
		logOnError(err, "Error while reading output from peer (trying to fetch ip, port)")

		line = strings.TrimSpace(line)
		if strings.Contains(line, LISTENER_IP_PRINT_MESSAGE) {
			line = strings.Replace(line, LISTENER_IP_PRINT_MESSAGE, "", 1)
			ip, port, err := net.SplitHostPort(line)
			logOnError(err, "Error while splitting host and port on line '"+line+"'")

			return ip, port
		}
	}
}

func PrintIfError(err error, str string) {
	if err != nil {
		fmt.Println(str + ": " + err.Error())
	}
}

// Dumps docker-compose log to stdout on error
func logOnError(err error, msg string) {
	PrintIfError(err, msg)
	if err == nil {
		return
	}
	DumpDockerComposeLogs()
	panic("Fatal error - Dumped docker-compose logs - " + msg)
}

func DumpDockerComposeLogs() {
	fmt.Println(composeBuffer.String())
}

// Throttles for given milliseconds
func throttle(duration int) {
	time.Sleep(time.Millisecond * time.Duration(duration))
}

// Tries fun for timeout milliseconds. If not done, errFun is called
func tryUntilTimeout(fun func(), errFun func(), timeout int) {
	maxRuntime, err := time.ParseDuration(strconv.Itoa(timeout) + "ms")
	logOnError(err, "Error while creating timeout duration")

	succeeded := make(chan bool)
	// After timeout send false to the channel
	time.AfterFunc(maxRuntime, func() {
		succeeded <- false
	})

	go func() {
		fun()
		succeeded <- true
	}()

	if <-succeeded {
		return
	} else {
		errFun()
	}
}
