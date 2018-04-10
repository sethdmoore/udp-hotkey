package serial

import (
	"io"

	"github.com/jacobsa/go-serial/serial"
)

func Connect(PortName string) (io.ReadWriteCloser, error) {
	var conn io.ReadWriteCloser
	var err error
	options := serial.OpenOptions{
		PortName:        PortName,
		BaudRate:        9600,
		DataBits:        8,
		StopBits:        1,
		MinimumReadSize: 4,
	}

	conn, err = serial.Open(options)
	if err != nil {
		return nil, err
	}

	return conn, err
}
