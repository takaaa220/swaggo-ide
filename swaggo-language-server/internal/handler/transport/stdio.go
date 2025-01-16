package transport

import (
	"fmt"
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
	stdinErr := os.Stdin.Close()
	stdoutErr := os.Stdout.Close()

	if stdinErr != nil && stdoutErr != nil {
		return fmt.Errorf("stdin error: %v, stdout error: %v", stdinErr, stdoutErr)
	}
	if stdinErr != nil {
		return stdinErr
	}
	return stdoutErr
}
