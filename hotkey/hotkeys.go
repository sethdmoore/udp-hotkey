package hotkey

import (
	"bytes"
	"fmt"

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
