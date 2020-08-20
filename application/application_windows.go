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
	"github.com/sethdmoore/serial-hotkey/types"
	"github.com/sethdmoore/serial-hotkey/windows"
	"net"
)

// handleHeldKey watches for a key to be released by polling the w32API keymap on an interval
// when released, it uses the chan to send a release action packet
func handleHeldKey(keystate *windows.KeyState, key *hotkeys.HotKey, serial chan<- types.Packet) {
	wincalls := windows.Get()
	log.Printf("New Thread: Key %d being held...\n", key.KeyCode)

	// XXX: not sure if this is thread safe to do
	// XXX: but also not sure if it's physically possible to double hold a hot key
	key.KeyHeld = true
	// write the keycode to the channel
	serial <- key.KeyHeldSerial
	// watch and wait for release
	for {
		time.Sleep(10 * time.Millisecond)
		// poll the keymap API, we only care about r1
		r1, _, _ := wincalls.KeyState.Call(uintptr(keystate.KeyCode))
		// some W32API C lib weirdness with multiple return. r1 is what we want
		if r1 == 0 {
			log.Printf("Key %d released!\n", keystate.KeyCode)
			serial <- key.KeyReleaseSerial
			// essentially kill the thread
			break
		}
	}
	// XXX: not sure about doing this
	// XXX: but you'd probably need multiple keyboards in order to break it...
	key.KeyHeld = false
}

// ClientStart should never be called from Windows as it is not supported
func ClientStart(_ string) error {
	return errors.New("Client unavailable on this platform")
}

func packetWriter(address string, input <-chan types.Packet) {

	conn, err := net.Dial("udp", address)

	if err != nil {
		log.Fatalf("Could not connect to %s, %v", address, err)
	}

	defer conn.Close()

	//var buf bytes.Buffer

	var buf bytes.Buffer
	var msg types.Packet

	for {
		// reset variables so we don't leave garbage lying around
		buf = bytes.Buffer{}
		msg = types.Packet{}

		msg = <-input
		spew.Dump(msg)

		// Must reinitialize the byte buffer and the encoder
		// otherwise the buffer continues to fill
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
	go packetWriter(address, serialChan)

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

			// handle holding hotkeys with a thread
			// this way we don't block execution while waiting for the key to be released
			if keys[id].Modifiers&constants.ModNoRepeat != 0 && !keys[id].KeyHeld {
				go handleHeldKey(&keystate, keys[id], serialChan)
			}

			// handle single fire hotkeys
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
