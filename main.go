package main

import (
	"bytes"
	"fmt"
	"log"
	"runtime"
	"syscall"
	"time"
	"unsafe"

	"github.com/davecgh/go-spew/spew"
	"github.com/jacobsa/go-serial/serial"
)

// hopefully fix bugs?

const (
	ModAlt = 1 << iota
	ModCtrl
	ModShift
	ModWin
	ModNoRepeat = 16384
)

type Hotkey struct {
	Id        int // Unique id
	Modifiers int // Mask of modifiers
	KeyCode   int // Key code, e.g. 'A'
}

type MSG struct {
	HWND   uintptr
	UINT   uintptr
	WPARAM int16
	LPARAM int64
	DWORD  int32
	POINT  struct{ X, Y int64 }
}

type Keystate struct {
	KeyCode int
}

func (h *Hotkey) String() string {
	mod := &bytes.Buffer{}
	if h.Modifiers&ModAlt != 0 {
		mod.WriteString("Alt+")
	}
	if h.Modifiers&ModCtrl != 0 {
		mod.WriteString("Ctrl+")
	}
	if h.Modifiers&ModShift != 0 {
		mod.WriteString("Shift+")
	}
	if h.Modifiers&ModWin != 0 {
		mod.WriteString("Win+")
	}
	return fmt.Sprintf("Hotkey[Id: %d, %s%c]", h.Id, mod, h.KeyCode)
}

func grabHotkeys() {
	// https://stackoverflow.com/questions/38646794/implement-a-global-hotkey-in-golang

}

func main() {

	//hopefully fix all bugs :(
	runtime.LockOSThread()
	var counter int
	options := serial.OpenOptions{
		PortName:        "COM1",
		BaudRate:        9600,
		DataBits:        8,
		StopBits:        1,
		MinimumReadSize: 4,
	}

	var m MSG
	var ks Keystate
	user32 := syscall.MustLoadDLL("user32")
	defer user32.Release()

	reghotkey := user32.MustFindProc("RegisterHotKey")
	reghotkey.Call()

	keys := map[int16]*Hotkey{
		1: &Hotkey{1, ModAlt + ModCtrl, 'O'},  // ALT+CTRL+O
		2: &Hotkey{2, ModAlt + ModShift, 'M'}, // ALT+SHIFT+M
		3: &Hotkey{3, ModAlt + ModCtrl, 'X'},  // ALT+CTRL+X
		//4: &Hotkey{4, ModNoRepeat + ModAlt, 127},
		4: &Hotkey{4, ModNoRepeat, 127},
		//5: &Hotkey{5, ModNoRepeat, 'Y'},
	}

	//peekmsg := user32.MustFindProc("PeekMessageW")
	peekmsg := user32.MustFindProc("GetMessageW")
	keystate := user32.MustFindProc("GetAsyncKeyState")

	translate := user32.MustFindProc("TranslateMessage")
	dispatch := user32.MustFindProc("DispatchMessageW")

	// Register hotkeys:
	for _, v := range keys {
		r1, _, err := reghotkey.Call(
			0, uintptr(v.Id), uintptr(v.Modifiers), uintptr(v.KeyCode))
		if r1 == 1 {
			fmt.Println("Registered", v)
		} else {
			fmt.Println("Failed to register", v, ", error:", err)
		}
	}

	for {
		var msg = &m
		var ks = &ks
		fmt.Printf("Size: %v\n", unsafe.Sizeof(msg))
		_ = "breakpoint"
		r1, r2, err := peekmsg.Call(uintptr(unsafe.Pointer(msg)), 0, 0, 0, 1)
		fmt.Printf("%v\n", r1)
		if r1 == 0 {
			fmt.Printf("Something is wrong: %v, %v", r1, err)
		} else {
			fmt.Printf("%v\n", r2)
		}

		fmt.Printf("Size peekmsg: %v\n", unsafe.Sizeof(peekmsg))

		// Registered id is in the WPARAM field:

		if id := msg.WPARAM; id != 0 {
			ks.KeyCode = keys[id].KeyCode
			counter++
			//spew.Dump(msg)

			fmt.Println("Hotkey pressed:", keys[id])
			//spew.Dump(keys[id].Modifiers)
			//fmt.Printf("%v\n", keys[id].Modifiers&ModNoRepeat != 0)

			if keys[id].Modifiers&ModNoRepeat != 0 {
				fmt.Println("Button being held...")
			inner:
				for {
					//spew.Dump(ks)
					//r1, _, err := keystate.Call(uintptr(ks.KeyCode))
					time.Sleep(10 * time.Millisecond)
					r1, _, _ := keystate.Call(uintptr(ks.KeyCode))
					//fmt.Printf("%v\n", r1)
					//spew.Dump(r1)

					//uint64(r1)
					if r1 == 0 {
						fmt.Printf("You released!\n")
						break inner
					}

				}

			}

			if id == 3 { // CTRL+ALT+X = Exit
				fmt.Println("CTRL+ALT+X pressed, goodbye...")
				return
			}
		} else {
			spew.Dump(msg)
		}

		translate.Call(uintptr(unsafe.Pointer(msg)))
		dispatch.Call(uintptr(unsafe.Pointer(msg)))
		fmt.Println("about to sleep")
		fmt.Printf("Count: %v\n", counter)
		time.Sleep(time.Millisecond * 50)
	}

	fmt.Println("running")
	port, err := serial.Open(options)
	if err != nil {
		log.Fatalf("serial.Open: %v", err)
	}

	defer port.Close()

	//b := []byte{0x00, 0x01}
	b := []byte("foo bar")
	n, err := port.Write(b)
	if err != nil {
		log.Fatalf("port.Write: %v", err)
	}

	fmt.Println("Wrote", n, "bytes.")
}
