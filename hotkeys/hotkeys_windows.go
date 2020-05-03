package hotkeys

import (
	"bytes"
	"fmt"
	"log"

	"github.com/sethdmoore/serial-hotkey/util"
	"github.com/sethdmoore/serial-hotkey/windows"
)

type HotKey struct {
	//Id               int    // Unique id
	Modifiers        int    // Mask of modifiers
	KeyCode          int    // Key code, e.g. 'A'
	KeyLinux         int    // Linux mapping of same hotkey
	KeySerial        string // Key combination to send over serial
	KeyWindowsString string // Friendly windows name of key
	KeyHeld          bool   // Whether the key is held down already
}

// NewHotKey registers a new HotKey struct
func NewHotKey(modifiers int, keycode int) *HotKey {
	// HotKey constructor
	var h HotKey
	var err error
	mod := &bytes.Buffer{}

	//h.Id = id
	h.Modifiers = modifiers
	h.KeyCode = keycode

	if h.Modifiers&windows.ModAlt != 0 {
		mod.WriteString("Alt+")
	}
	if h.Modifiers&windows.ModCtrl != 0 {
		mod.WriteString("Ctrl+")
	}
	if h.Modifiers&windows.ModShift != 0 {
		mod.WriteString("Shift+")
	}
	if h.Modifiers&windows.ModWin != 0 {
		mod.WriteString("Win+")
	}

	h.KeyLinux, err = util.WinKeyToLinux(h.KeyCode)
	if err != nil {
		log.Fatalf("Error: problem mapping Windows key %d to Linux: .. %v", h.KeyCode, err)
	}
	h.KeySerial = fmt.Sprintf("%s%d", mod, h.KeyLinux)
	h.KeyWindowsString = fmt.Sprintf("%s%c", mod, h.KeyCode)

	return &h
}

var Keys map[int16]*HotKey

func init() {

	Keys = map[int16]*HotKey{
		1: NewHotKey(windows.ModCtrl, 'O'),                 // ALT+CTRL+O
		2: NewHotKey(windows.ModAlt+windows.ModShift, 'M'), // ALT+SHIFT+M
		3: NewHotKey(windows.ModAlt+windows.ModCtrl, 'X'),  // ALT+CTRL+X
		4: NewHotKey(windows.ModNoRepeat, 127),
		5: NewHotKey(windows.ModNoRepeat+windows.ModCtrl, 127),
		6: NewHotKey(windows.ModNoRepeat+windows.ModShift, 127),
		7: NewHotKey(windows.ModNoRepeat+windows.ModAlt, 127),
	}

	wincalls := windows.Get()

	for id, v := range Keys {
		r1, r2, err := wincalls.RegHotKey.Call(
			//0, uintptr(v.Id), uintptr(v.Modifiers), uintptr(v.KeyCode))
			0, uintptr(id), uintptr(v.Modifiers), uintptr(v.KeyCode))
		if r1 == 1 {
			fmt.Printf("Registered %s\n", v.KeyWindowsString)

		} else {

			if err != nil {
				log.Fatalf("Could not register! %v %v %v", r1, r2, err)
			}
			fmt.Println("Failed to register", v, ", error:", err)
		}
	}
}
