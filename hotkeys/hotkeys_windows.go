package hotkeys

import (
	"bytes"
	"fmt"
	"log"

	"github.com/sethdmoore/serial-hotkey/constants"
	"github.com/sethdmoore/serial-hotkey/types"
	"github.com/sethdmoore/serial-hotkey/util"
	"github.com/sethdmoore/serial-hotkey/windows"
)

// Hotkey is exported so we can reference fields.
// It contains all the fields necessary to send hotkeys over the serial wire
type HotKey struct {
	Modifiers        uint16       // Mask of modifiers, uin16 == 2 bytes
	KeyCode          uint8        // Key code, e.g. 'A', uint8 == 1 bytes
	KeyLinux         uint8        // Linux mapping of same hotkey
	KeySerial        types.Packet // Single Key combination to send over serial
	KeyHeldSerial    types.Packet // Packet for holding key, used with ModNoRepeat
	KeyReleaseSerial types.Packet // Packet for releasing key, used with ModNoRepeat
	KeyWindowsString string       // Friendly windows name of key
	KeyHeld          bool         // Whether the key is held down already
}

// NewHotKey registers a new HotKey struct
func NewHotKey(modifiers uint16, keycode uint8) *HotKey {
	// HotKey constructor
	var h HotKey
	var err error
	mod := &bytes.Buffer{}

	//h.Id = id
	h.Modifiers = modifiers
	h.KeyCode = keycode

	if h.Modifiers&constants.ModAlt != 0 {
		mod.WriteString("Alt+")
	}
	if h.Modifiers&constants.ModCtrl != 0 {
		mod.WriteString("Ctrl+")
	}
	if h.Modifiers&constants.ModShift != 0 {
		mod.WriteString("Shift+")
	}
	if h.Modifiers&constants.ModWin != 0 {
		mod.WriteString("Win+")
	}

	h.KeyLinux, err = util.WinKeyToLinux(h.KeyCode)
	if err != nil {
		log.Fatalf("Error: problem mapping Windows key %d to Linux: .. %v", h.KeyCode, err)
	}

	// we set the mode last so we can set the struct's Held and Release fields
	if h.Modifiers&constants.ModNoRepeat != 0 {
		h.KeyHeldSerial = types.Packet{
			Action:    constants.KeyHeld,
			Modifiers: h.Modifiers,
			KeyCode:   h.KeyLinux,
		}

		h.KeyReleaseSerial = types.Packet{
			Action:    constants.KeyRelease,
			Modifiers: h.Modifiers,
			KeyCode:   h.KeyLinux,
		}
	}

	h.KeySerial = types.Packet{
		Action:    constants.KeyPress,
		Modifiers: h.Modifiers,
		KeyCode:   h.KeyLinux,
	}
	// set the human readable version

	h.KeyWindowsString = fmt.Sprintf("%s%c", mod, h.KeyCode)
	// for debugging
	// fmt.Printf("Press: %x\nHeld: %x\nRelease: %x\n", h.KeySerial, h.KeyHeldSerial, h.KeyReleaseSerial)

	return &h
}

// Keys contains a map of windows keyIds to HotKey structs
var Keys map[int16]*HotKey

func init() {

	Keys = map[int16]*HotKey{
		1: NewHotKey(constants.ModCtrl, 'O'),                   // ALT+CTRL+O
		2: NewHotKey(constants.ModAlt+constants.ModShift, 'M'), // ALT+SHIFT+M
		3: NewHotKey(constants.ModAlt+constants.ModCtrl, 'X'),  // ALT+CTRL+X
		4: NewHotKey(constants.ModNoRepeat, 127),
		5: NewHotKey(constants.ModNoRepeat+constants.ModCtrl, 127),
		6: NewHotKey(constants.ModNoRepeat+constants.ModShift, 127),
		7: NewHotKey(constants.ModNoRepeat+constants.ModAlt, 127),
		//8: NewHotKey(constants.ModAlt, 112),
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
