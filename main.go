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
		application.ServerStart("COM1")
	} else {
		fmt.Printf("Starting...\n")
		err := application.ClientStart()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
		}
	}
}
