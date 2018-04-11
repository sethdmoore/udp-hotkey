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
		application.ServerStart()
	} else {
		fmt.Printf("Starting...\n")
		err := application.ClientStart()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
		}
	}
}
