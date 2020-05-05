package util

import (
	//"syscall"
	"errors"
	"fmt"

	//"github.com/sethdmoore/serial-hotkey/linux"
	"github.com/sethdmoore/serial-hotkey/keymaps/linuxkeys"
	"github.com/sethdmoore/serial-hotkey/keymaps/windowskeys"
	//linuxkey "github.com/micmonay/keybd_event/keybd_linux"
)

//var KeyMap KeyMapBase

func WinKeyToLinux(key uint8) (uint8, error) {
	for name, val := range windowskeys.Keys {
		if val == key {
			if val, ok := linuxkeys.Keys[name]; ok {
				linuxkey := linuxkeys.Keys[name]
				return linuxkey, nil
			} else {
				return 0, errors.New(fmt.Sprintf("%s with value %d does not map to any known linux virtual key :(", name, val))
			}
		}
	}
	return 0, errors.New(fmt.Sprintf("Key ID %d does not map to any named Windows virtual key :(", key))
}

func KeyToByteArray() []byte {
	return []byte{}
}
