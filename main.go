package main

import (
	"fmt"
	"github.com/sethdmoore/serial-hotkey/application"
	"runtime"
)

func main() {

	// Lock the main thread so windows doesn't miss messages
	runtime.LockOSThread()
	if runtime.GOOS == "windows" {
		application.ServerStart("seth.home:1111")
	} else {
		fmt.Printf("Starting...\n")
		err := application.ClientStart("/dev/pts/2")
		if err != nil {
			fmt.Printf("Error: %v\n", err)
		}
	}
}
