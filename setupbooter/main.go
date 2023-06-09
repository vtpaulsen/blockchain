package main

import (
	"os"
	"runtime"
	setupbooter "setupbooter/docker"
	"testing"
)

func main() {
	t := &testing.T{}

	var setup []setupbooter.PeerInstance

	switch os.Args[1] {
	case "SetupA":
		setup = setupbooter.SetupA
	case "SetupB":
		setup = setupbooter.SetupB
	case "SetupC":
		setup = setupbooter.SetupC
	case "SetupD":
		setup = setupbooter.SetupD
	}

	test := setupbooter.DockerTester(t, setup)
	test.OutputToStdout = true
	test.StartMultiContainerEnv()

	// prevent exiting
	runtime.Goexit()
}
