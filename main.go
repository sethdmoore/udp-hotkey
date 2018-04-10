package main

import (
	"fmt"
	"runtime"

	"github.com/sethdmoore/serial-hotkey/server"
)

func main() {

	// Lock the main thread so windows doesn't miss messages
	runtime.LockOSThread()
	if runtime.GOOS == "windows" {
		server.Start()
	} else {
		fmt.Println("client mode not implemented yet")
	}
}
