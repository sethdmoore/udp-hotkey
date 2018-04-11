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

func ServerStart() error {

	serialPort := "COM1"

	port, err := serial.Connect(serialPort)
	if err != nil {
		log.Fatalf("Could not connect to %s, %v", serialPort, err)
	}
	defer port.Close()

	wincalls := windows.Get()

	var msg windows.MSG

	fmt.Println("running")

	keys := hotkeys.Keys

	var keystate windows.KeyState

	for {
		r1, _, err := wincalls.GetMSG.Call(uintptr(unsafe.Pointer(&msg)), 0, 0, 0, 1)
		if r1 == 0 {
			log.Printf("Error parsing windows message: %v, %v", r1, err)
		}

		// Registered id is in the WPARAM field:
		if id := msg.WPARAM; id != 0 {
			keystate.KeyCode = keys[id].KeyCode

			if keys[id].Modifiers&windows.ModNoRepeat != 0 {
				log.Printf("Key %d being held...\n", keystate.KeyCode)
				_, err := port.Write([]byte(fmt.Sprintf("down:%s\n", keys[id].KeySerial)))
				if err != nil {
					log.Fatalf("port.Write: %v", err)
				}
			inner:
				for {
					time.Sleep(10 * time.Millisecond)
					r1, _, _ := wincalls.KeyState.Call(uintptr(keystate.KeyCode))
					if r1 == 0 {
						log.Printf("Key %d released!\n", keystate.KeyCode)

						_, err := port.Write([]byte(fmt.Sprintf("up:%s\n", keys[id].KeySerial)))
						if err != nil {
							log.Fatalf("port.Write: %v", err)
						}
						break inner
					}

				}

			} else {
				fmt.Println("Hotkey pressed:", keys[id])
				_, err := port.Write([]byte(fmt.Sprintf("press:%s\n", keys[id].KeySerial)))
				if err != nil {
					log.Fatalf("port.Write: %v", err)
				}
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
