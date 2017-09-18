package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"time"
)

type Rclone struct {
	Paths Paths
}

func (r Rclone) MoveTo(source, destination string) ([]byte, error) {
	command := exec.Command(RcloneCmd, MoveTo, source, destination, "--verbose")
	output, err := command.CombinedOutput()
	if err != nil {
		return []byte{}, err
	}

	return output, nil
}

func (r Rclone) Mount() error {
	command := exec.Command(
		"rclone",
		"mount",
		"--allow-other",
		"--max-read-ahead",
		"2G",
		"--dir-cache-time",
		"1m0s",
		r.Paths.Mount,
		r.Paths.Decrypt,
	)

	command.Env = os.Environ()

	if err := command.Start(); err != nil {
		return err
	}

	for i := 0; i < 10; i++ {
		isMounted, err := r.Mounted()
		if err != nil {
			return err
		}

		if !isMounted {
			time.Sleep(1000 * time.Millisecond)
		}

		if isMounted {
			return nil
		}
	}

	return ErrCouldNotVerifyMount(r.Paths.Decrypt)
}

func (r Rclone) Unmount() error {
	isMounted, err := r.Mounted()
	if err != nil {
		return err
	}

	if isMounted {
		if err := syscall.Unmount(r.Paths.Decrypt, 0); err != nil {
			return err
		}
	}

	return nil
}

func (r Rclone) Mounted() (bool, error) {
	activeMounts, err := exec.Command("mount").CombinedOutput()
	if err != nil {
		return false, err
	}

	return strings.Contains(string(activeMounts), fmt.Sprintf("%s", r.Paths.Decrypt)), nil
}
