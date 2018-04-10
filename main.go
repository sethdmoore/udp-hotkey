package main

import (
	"encoding/binary"
	"fmt"
	"log"
	"runtime"
	"time"
	"unsafe"

	"github.com/davecgh/go-spew/spew"
	"github.com/sethdmoore/serial-hotkey/hotkey"
	"github.com/sethdmoore/serial-hotkey/serial"
	"github.com/sethdmoore/serial-hotkey/util"
	"github.com/sethdmoore/serial-hotkey/windows"
)

func main() {
	serialPort := "COM1"

	// Lock the main thread so windows doesn't miss messages
	runtime.LockOSThread()

	port, err := serial.Connect(serialPort)
	if err != nil {
		log.Fatalf("Could not connect to %s, %v", serialPort, err)
	}
	defer port.Close()

	wincalls := windows.Get()

	keybuff := make([]byte, binary.MaxVarintLen64)

	var msg windows.MSG
	var ks windows.KeyState

	keys := map[int16]*hotkey.Hotkey{
		1: &hotkey.Hotkey{1, windows.ModAlt + windows.ModCtrl, 'O'},  // ALT+CTRL+O
		2: &hotkey.Hotkey{2, windows.ModAlt + windows.ModShift, 'M'}, // ALT+SHIFT+M
		3: &hotkey.Hotkey{3, windows.ModAlt + windows.ModCtrl, 'X'},  // ALT+CTRL+X
		//4: &Hotkey{4, ModNoRepeat + ModAlt, 127},
		4: &hotkey.Hotkey{4, windows.ModNoRepeat, 127},
		//5: &Hotkey{5, ModNoRepeat, 'Y'},
	}

	fmt.Println("running")

	// Register hotkeys:
	for _, v := range keys {
		r1, r2, err := wincalls.RegHotKey.Call(
			0, uintptr(v.Id), uintptr(v.Modifiers), uintptr(v.KeyCode))
		if r1 == 1 {
			fmt.Println("Registered", v)

		} else {

			if err != nil {
				log.Fatalf("Could not register! %v %v %v", r1, r2, err)
			}
			fmt.Println("Failed to register", v, ", error:", err)
		}
	}

	for {
		r1, _, err := wincalls.GetMSG.Call(uintptr(unsafe.Pointer(&msg)), 0, 0, 0, 1)
		if r1 == 0 {
			log.Printf("Error parsing windows message: %v, %v", r1, err)
		}

		// Registered id is in the WPARAM field:
		if id := msg.WPARAM; id != 0 {
			ks.KeyCode = keys[id].KeyCode
			_ = binary.PutVarint(keybuff, int64(ks.KeyCode))

			linuxkey, err := util.WinKeyToLinux(ks.KeyCode)
			if err != nil {
				log.Printf("Warning: problem mapping Windows key %d to Linux: .. %v", ks.KeyCode, err)
				continue
			}

			if keys[id].Modifiers&windows.ModNoRepeat != 0 {
				log.Printf("Key %d being held...\n", ks.KeyCode)
				_, err := port.Write([]byte(fmt.Sprintf("down:%d\n", linuxkey)))
				if err != nil {
					log.Fatalf("port.Write: %v", err)
				}
			inner:
				for {
					time.Sleep(10 * time.Millisecond)
					r1, _, _ := wincalls.KeyState.Call(uintptr(ks.KeyCode))
					if r1 == 0 {
						log.Printf("Key %d released!\n", ks.KeyCode)

						_, err := port.Write([]byte(fmt.Sprintf("up:%d\n", linuxkey)))
						if err != nil {
							log.Fatalf("port.Write: %v", err)
						}
						break inner
					}

				}

			} else {
				fmt.Println("Hotkey pressed:", keys[id])
				_, err := port.Write([]byte(fmt.Sprintf("press:%d\n", linuxkey)))
				if err != nil {
					log.Fatalf("port.Write: %v", err)
				}
			}

			if id == 3 { // CTRL+ALT+X = Exit
				fmt.Println("CTRL+ALT+X pressed, goodbye...")
				return
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
}
