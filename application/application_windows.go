package application

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"log"
	"time"
	"unsafe"

	"github.com/davecgh/go-spew/spew"
	"github.com/sethdmoore/serial-hotkey/constants"
	"github.com/sethdmoore/serial-hotkey/hotkeys"
	//"github.com/sethdmoore/serial-hotkey/serial"
	"github.com/sethdmoore/serial-hotkey/types"
	"github.com/sethdmoore/serial-hotkey/windows"
	"net"
)

func handleHeldKey(keystate *windows.KeyState, key *hotkeys.HotKey, serial chan<- types.Packet) {
	wincalls := windows.Get()
	log.Printf("New Thread: Key %d being held...\n", key.KeyCode)

	// XXX: not sure if this is thread safe to do
	key.KeyHeld = true
	serial <- key.KeyHeldSerial
	for {
		time.Sleep(10 * time.Millisecond)
		r1, _, _ := wincalls.KeyState.Call(uintptr(keystate.KeyCode))
		if r1 == 0 {
			log.Printf("Key %d released!\n", keystate.KeyCode)

			serial <- key.KeyReleaseSerial
			break
		}
	}
	// XXX: not sure about doing this
	key.KeyHeld = false
}

// ClientStart should never be called from Windows as it is not supported
func ClientStart(_ string) error {
	return errors.New("Client unavailable on this platform")
}

func serialWriter(address string, input <-chan types.Packet) {
	var msg types.Packet

	conn, err := net.Dial("udp", address)

	if err != nil {
		log.Fatalf("Could not connect to %s, %v", address, err)
	}

	defer conn.Close()

	//var buf bytes.Buffer

	for {
		msg = <-input
		// Must reinitialize the byte buffer and the encoder
		// otherwise the buffer continues to fill
		var buf bytes.Buffer
		enc := gob.NewEncoder(&buf)

		//err = gob.NewEncoder(&buf).Encode(msg)
		err = enc.Encode(msg)
		if err != nil {
			log.Fatalf("Could not encode packet into message: %v", err)
		}
		//spew.Dump(msg)
		//fmt.Println(buf)
		//spew.Dump(msg)
		_, err := conn.Write(buf.Bytes())
		if err != nil {
			log.Fatalf("FATAL: port.Write: %v", err)
		}
	}
}

// ServerStart starts the hotkey server on Windows
func ServerStart(address string) error {
	var msg windows.MSG
	wincalls := windows.Get()
	keys := hotkeys.Keys
	var keystate windows.KeyState

	serialChan := make(chan types.Packet)

	// thread for udp connection and handling
	go serialWriter(address, serialChan)

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

			if keys[id].Modifiers&constants.ModNoRepeat != 0 && !keys[id].KeyHeld {
				go handleHeldKey(&keystate, keys[id], serialChan)
			}

			if keys[id].Modifiers&constants.ModNoRepeat == 0 {
				fmt.Printf("Hotkey pressed: %s\n", keys[id].KeyWindowsString)
				serialChan <- keys[id].KeySerial
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
