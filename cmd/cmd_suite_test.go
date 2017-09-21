package cmd_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"bytes"
	"io"
	"os"
	"testing"
)

func TestCmd(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Cmd Suite")
}

func captureStdout(f func()) (string, error) {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	f()

	if err := w.Close(); err != nil {
		return "", err
	}

	os.Stdout = old

	var buf bytes.Buffer

	_, err := io.Copy(&buf, r)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}
