package hotkeys

import (
	"bytes"
	"fmt"
	"log"

	"github.com/sethdmoore/serial-hotkey/util"
	"github.com/sethdmoore/serial-hotkey/windows"
)

type hotKey struct {
	Id               int    // Unique id
	Modifiers        int    // Mask of modifiers
	KeyCode          int    // Key code, e.g. 'A'
	KeyLinux         int    // Linux mapping of same hotkey
	KeySerial        string // Key combination to send over serial
	KeyWindowsString string // Friendly windows name of key
}

func NewHotKey(id int, modifiers int, keycode int) *hotKey {
	// HotKey constructor
	var h hotKey
	var err error
	mod := &bytes.Buffer{}

	h.Id = id
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
	h.KeySerial = fmt.Sprintf("%s%c", mod, h.KeyLinux)
	h.KeyWindowsString = fmt.Sprintf("%s%c", mod, h.KeyCode)

	return &h
}

var Keys map[int16]*hotKey

func init() {

	Keys = map[int16]*hotKey{
		//1: &Hotkey{1, windows.ModAlt + windows.ModCtrl, 'O'},  // ALT+CTRL+O
		1: NewHotKey(1, windows.ModAlt+windows.ModCtrl, 'O'),  // ALT+CTRL+O
		2: NewHotKey(2, windows.ModAlt+windows.ModShift, 'M'), // ALT+SHIFT+M
		3: NewHotKey(3, windows.ModAlt+windows.ModCtrl, 'X'),  // ALT+CTRL+X
		//4: &Hotkey{4, ModNoRepeat + ModAlt, 127},
		4: NewHotKey(4, windows.ModNoRepeat, 127),
		//5: &Hotkey{5, ModNoRepeat, 'Y'},
	}
	// We basically export this map, so avoid warnings
	//_ = Keys

	wincalls := windows.Get()

	// Register hotkeys:
	for _, v := range Keys {
		r1, r2, err := wincalls.RegHotKey.Call(
			0, uintptr(v.Id), uintptr(v.Modifiers), uintptr(v.KeyCode))
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
