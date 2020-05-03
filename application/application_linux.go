package application

import (
	"errors"
	"fmt"
	"github.com/hpcloud/tail"
	"github.com/sethdmoore/keybd_event"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	HOLD    = "down"
	RELEASE = "up"
	PRESS   = "press"
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
	//t, err := tail.TailFile("")
	kb, err := keybd_event.NewKeyBinding()
	if err != nil {
		return err
	}

	time.Sleep(2 * time.Second)

	t, err := tail.TailFile("/tmp/windowsfile", tail.Config{Follow: true, Location: &tail.SeekInfo{0, os.SEEK_END}})
	if err != nil {
		return err
	}

	time.Sleep(1 * time.Second)
	fmt.Printf("Debug... no lines..\n")
	for line := range t.Lines {
		mode, err := parseText(line.Text, &kb)
		if err != nil {
			fmt.Printf("Warning: %v\n", err)
			continue
		}
		fmt.Printf("line.Text: %s\n", line.Text)

		switch mode {
		case HOLD:
			err = kb.PressKeys()
			if err != nil {
				return err
			}
		case RELEASE:
			err = kb.ReleaseKeys()
			if err != nil {
				return err
			}
		case PRESS:
			err = kb.Launching()
			if err != nil {
				return err
			}
		}
	}
	return nil

}
