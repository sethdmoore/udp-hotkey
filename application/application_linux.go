package application

import (
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/sethdmoore/keybd_event"
	"github.com/sethdmoore/serial-hotkey/constants"
	"github.com/sethdmoore/serial-hotkey/serial"
	"github.com/sethdmoore/serial-hotkey/types"
	"strconv"
	"strings"
)

func ServerStart(serialPort string) error {
	return errors.New("Server unavailable on this platform")
}
func parseText(input string, k *keybd_event.KeyBinding) (string, error) {
	input_list := strings.Split(input, ":")
	if len(input_list) != 2 {
		return "", errors.New("Input was not a valid command")
	}

	action, keys := input_list[0], input_list[1]

	modifier_list := strings.Split(keys, "+")

	key := modifier_list[len(modifier_list)-1]
	key_int, err := strconv.Atoi(key)
	if err != nil {
		return "", errors.New(fmt.Sprintf("Could not convert %s to key integer: %v", key, err))
	}

	//key_list := []int{key_int}

	k.SetKeys(key_int)
	modifier_list = modifier_list[:len(modifier_list)-1]

	for _, mod := range modifier_list {
		switch mod {
		case "Alt":
			k.HasALT(true)
		case "Shift":
			k.HasSHIFT(true)
		case "Ctrl":
			k.HasCTRL(true)
		}
	}
	return action, nil
}

func ClientStart() error {
	kb, err := keybd_event.NewKeyBinding()
	if err != nil {
		return err
	}

	var packet types.Packet

	port, err := serial.Connect("/dev/pts/2")
	// Compose bufio ReadN methods from our serial lib's ReadWriteCloser

	if err != nil {
		return err
	}

	for {
		err := binary.Read(port, binary.BigEndian, &packet)
		if err != nil {
			fmt.Printf("ERR: problem reading from serial: %v\n", err)
			continue
		}

		fmt.Printf("DEBUG: %x\n", packet)

		kb.SetKeys(int(packet.KeyCode))

		switch packet.Action {
		case constants.KeyHeld:
			err = kb.PressKeys()
			if err != nil {
				return err
			}
		case constants.KeyRelease:
			err = kb.ReleaseKeys()
			if err != nil {
				return err
			}
		case constants.KeyPress:
			err = kb.Launching()
			if err != nil {
				return err
			}
		}
	}
	return nil

}
