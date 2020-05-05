package hotkeys

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"

	"github.com/sethdmoore/serial-hotkey/constants"
	"github.com/sethdmoore/serial-hotkey/util"
	"github.com/sethdmoore/serial-hotkey/windows"
)

// Hotkey is exported so we can reference fields.
// It contains all the fields necessary to send hotkeys over the serial wire
type HotKey struct {
	//Id               int    // Unique id
	Modifiers        uint16 // Mask of modifiers, uin16 == 2 bytes
	KeyCode          uint8  // Key code, e.g. 'A', uint8 == 1 bytes
	KeyLinux         uint8  // Linux mapping of same hotkey
	KeySerial        []byte // Single Key combination to send over serial
	KeyHeldSerial    []byte // Packet for holding key, used with ModNoRepeat
	KeyReleaseSerial []byte // Packet for releasing key, used with ModNoRepeat
	KeyWindowsString string // Friendly windows name of key
	KeyHeld          bool   // Whether the key is held down already
}

// NewHotKey registers a new HotKey struct
func NewHotKey(modifiers uint16, keycode uint8) *HotKey {
	// HotKey constructor
	var h HotKey
	var err error
	mod := &bytes.Buffer{}

	// buf ends up being the packet sent over serial
	// 0th byte is the mode, see constants/constants.go
	// 1-2nd byte is the Modifier keys being held, see variables/variables_windows.go
	// 3rd byte is the linux keycap to act on
	buf := make([]byte, 4)

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

	// modifiers is uint16, so only 2 bytes
	binary.BigEndian.PutUint16(buf[1:], h.Modifiers)

	// KeyLinux is uint8, so only 1 byte
	buf[3] = h.KeyLinux

	// we set the mode last so we can set the struct's Held and Release fields
	if h.Modifiers&windows.ModNoRepeat != 0 {
		buf[0] = constants.Held

		// empty byte slice doesn't want to be copied, has something to do with len(src)
		// instead, just append the slice ot our
		h.KeyHeldSerial = append(h.KeyHeldSerial, buf...)
		//fmt.Printf("ATTN: %x\nCOMP: %x\n", buf, h.KeyHeldSerial)

		buf[0] = constants.Release
		// same deal here, using assignment '=' seems to reference a pointer.
		// Exploit append to structure our packet
		h.KeyReleaseSerial = append(h.KeyReleaseSerial, buf...)
	}
	// don't worry about pointers here, since this is the final time we
	// manipulate index 0 of buf
	buf[0] = constants.Press

	h.KeySerial = buf
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
