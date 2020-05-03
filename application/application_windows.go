package application

import (
	"errors"
	"fmt"
	"log"
	"time"
	"unsafe"

	"github.com/davecgh/go-spew/spew"
	"github.com/sethdmoore/serial-hotkey/hotkeys"
	"github.com/sethdmoore/serial-hotkey/serial"
	"github.com/sethdmoore/serial-hotkey/windows"
)

func handleHeldKey(keystate *windows.KeyState, key *hotkeys.HotKey, serial chan<- []byte) {
	wincalls := windows.Get()
	log.Printf("New Thread: Key %d being held...\n", keystate.KeyCode)

	// XXX: not sure if this is thread safe to do
	key.KeyHeld = true
	serial <- []byte(fmt.Sprintf("down:%s\n", key.KeySerial))
	for {
		time.Sleep(10 * time.Millisecond)
		r1, _, _ := wincalls.KeyState.Call(uintptr(keystate.KeyCode))
		if r1 == 0 {
			log.Printf("Key %d released!\n", keystate.KeyCode)

			serial <- []byte(fmt.Sprintf("up:%s\n", key.KeySerial))
			break
		}
	}
	// XXX: not sure about doing this
	key.KeyHeld = false
}

// ClientStart should never be called from Windows as it is not supported
func ClientStart() error {
	return errors.New("Client unavailable on this platform")
}

func serialWriter(serialPort string, input <-chan []byte) {
	var msg []byte

	port, err := serial.Connect(serialPort)
	if err != nil {
		log.Fatalf("Could not connect to %s, %v", serialPort, err)
	}

	defer port.Close()

	for {
		msg = <-input
		_, err := port.Write(msg)

		if err != nil {
			log.Fatalf("FATAL: port.Write: %v", err)
		}
	}
}

// ServerStart starts the hotkey server on Windows
func ServerStart(serialPort string) error {
	var msg windows.MSG
	wincalls := windows.Get()
	keys := hotkeys.Keys
	var keystate windows.KeyState

	serialChan := make(chan []byte)

	// thread for serial connection and handling
	go serialWriter(serialPort, serialChan)

	fmt.Println("running")

	for {
		r1, _, err := wincalls.GetMSG.Call(uintptr(unsafe.Pointer(&msg)), 0, 0, 0, 1)
		if r1 == 0 {
			log.Printf("Error parsing windows message: %v, %v", r1, err)
			continue
		}

		// Registered id is in the WPARAM field:
		if id := msg.WPARAM; id != 0 {
			keystate.KeyCode = keys[id].KeyCode

			if keys[id].Modifiers&windows.ModNoRepeat != 0 && !keys[id].KeyHeld {
				//fmt.Printf("Creating a new thread...\n")
				go handleHeldKey(&keystate, keys[id], serialChan)
			}

			if keys[id].Modifiers&windows.ModNoRepeat == 0 {
				fmt.Println("Hotkey pressed:", keys[id])
				serialChan <- []byte(fmt.Sprintf("press:%s\n", keys[id].KeySerial))
				/*
					if err != nil {
						log.Fatalf("port.Write: %v", err)
					}
				*/
			}

			if id == 3 { // CTRL+ALT+X = Exit
				fmt.Println("CTRL+ALT+X pressed, goodbye...")
				return nil
			}
		} else {
			spew.Dump(msg)
		}

		// Not sure if this section is required
		// MSDN documentation is shy on this
		wincalls.TranslateMSG.Call(uintptr(unsafe.Pointer(&msg)))
		wincalls.DispatchMSG.Call(uintptr(unsafe.Pointer(&msg)))

		time.Sleep(time.Millisecond * 50)
	}
	return errors.New("The main loop exited")
}
