package application

import (
	//"encoding/binary"
	"errors"
	"fmt"
	"github.com/sethdmoore/keybd_event"
	"github.com/sethdmoore/serial-hotkey/constants"
	//"github.com/sethdmoore/serial-hotkey/serial"
	"bytes"
	"encoding/gob"
	"github.com/sethdmoore/serial-hotkey/types"
	"net"
	//"time"
)

func ServerStart(serialPort string) error {
	return errors.New("Server unavailable on this platform")
}

func ClientStart(serialPath string) error {
	kb, err := keybd_event.NewKeyBinding()
	if err != nil {
		return err
	}

	var packet types.Packet

	//port, err := serial.Connect(serialPath)
	conn, err := net.ListenPacket("udp", ":1111")
	if err != nil {
		return err
	}
	defer conn.Close()

	//buf := &bytes.Buffer{}
	//var buf bytes.Buffer
	buf := make([]byte, 512)
	//reader := bytes.NewReader(&buf)

	for {
		n, _, err := conn.ReadFrom(buf)
		//gob.Decode(buf)
		if err != nil {
			fmt.Printf("ERR: problem reading packet: %v\n", err)
			continue
		}

		// All this to decode a buffer into an io.Reader so we can gob.Decode into our packet type
		// That's kind of complicated. https://stackoverflow.com/a/26150948
		err = gob.NewDecoder(bytes.NewReader(buf[:n])).Decode(&packet)
		if err != nil {
			fmt.Printf("ERR: could not decode packet: %v\n", err)
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
