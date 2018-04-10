package hotkeys

import (
	"bytes"
	"fmt"
	"log"

	"github.com/sethdmoore/serial-hotkey/windows"
)

type Hotkey struct {
	Id        int // Unique id
	Modifiers int // Mask of modifiers
	KeyCode   int // Key code, e.g. 'A'
}

func (h *Hotkey) String() string {
	mod := &bytes.Buffer{}
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
	return fmt.Sprintf("Hotkey[Id: %d, %s%c]", h.Id, mod, h.KeyCode)
}

var Keys map[int16]*Hotkey

func init() {

	Keys = map[int16]*Hotkey{
		1: &Hotkey{1, windows.ModAlt + windows.ModCtrl, 'O'},  // ALT+CTRL+O
		2: &Hotkey{2, windows.ModAlt + windows.ModShift, 'M'}, // ALT+SHIFT+M
		3: &Hotkey{3, windows.ModAlt + windows.ModCtrl, 'X'},  // ALT+CTRL+X
		//4: &Hotkey{4, ModNoRepeat + ModAlt, 127},
		4: &Hotkey{4, windows.ModNoRepeat, 127},
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
			fmt.Println("Registered", v)

		} else {

			if err != nil {
				log.Fatalf("Could not register! %v %v %v", r1, r2, err)
			}
			fmt.Println("Failed to register", v, ", error:", err)
		}
	}

}
