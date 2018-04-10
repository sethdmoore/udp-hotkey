package application

import (
	"errors"
	"fmt"
)

func init() {
	fmt.Printf("something\n")
}

func ServerStart() error {
	return errors.New("Server unavailable on this platform")
}
