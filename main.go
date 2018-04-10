package main

import (
	"runtime"

	"github.com/sethdmoore/serial-hotkey/server"
)

func main() {

	// Lock the main thread so windows doesn't miss messages
	runtime.LockOSThread()
	server.Start()
}
