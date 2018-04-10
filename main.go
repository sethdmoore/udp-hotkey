package main

import (
	"fmt"
	"runtime"

	"github.com/sethdmoore/serial-hotkey/application"
)

func main() {

	// Lock the main thread so windows doesn't miss messages
	runtime.LockOSThread()
	if runtime.GOOS == "windows" {
		application.ServerStart()
	} else {
		fmt.Println("client mode not implemented yet")
	}
}
