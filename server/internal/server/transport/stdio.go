package transport

import (
	"io"
	"os"
)

type stdrwc struct{}

var _ io.ReadWriteCloser = stdrwc{}

func (stdrwc) Read(p []byte) (int, error) {
	return os.Stdin.Read(p)
}

func (stdrwc) Write(p []byte) (int, error) {
	return os.Stdout.Write(p)
}

func (stdrwc) Close() error {
	if err := os.Stdin.Close(); err == nil {
		return os.Stdout.Close()
	} else {
		return err
	}
}
