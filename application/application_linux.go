package application

import (
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/sethdmoore/keybd_event"
	"github.com/sethdmoore/serial-hotkey/constants"
	"github.com/sethdmoore/serial-hotkey/serial"
	"github.com/sethdmoore/serial-hotkey/types"
	"time"
)

func ServerStart(serialPort string) error {
	return errors.New("Server unavailable on this platform")
}

func ClientStart() error {
	kb, err := keybd_event.NewKeyBinding()
	if err != nil {
		return err
	}

	var packet types.Packet

	port, err := serial.Connect("/dev/pts/1")
	if err != nil {
		return err
	}

	for {
		err := binary.Read(port, binary.BigEndian, &packet)
		if err != nil {
			fmt.Printf("ERR: problem reading from serial: %v\n", err)
			time.Sleep(30 * time.Second)
			continue
		}

		fmt.Printf("DEBUG: %x\n", packet)

		if packet.Modifiers&constants.ModAlt != 0 {
			kb.HasALT(true)
		}

		if packet.Modifiers&constants.ModCtrl != 0 {
			kb.HasCTRL(true)
		}

		if packet.Modifiers&constants.ModShift != 0 {
			kb.HasSHIFT(true)
		}

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
