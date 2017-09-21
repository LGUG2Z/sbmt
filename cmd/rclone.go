package cmd

import (
	"os"
	"os/exec"
	"strings"
	"syscall"
	"time"
)

// Rclone represents an interface to mount, unmount and verify an Rclone mount.
type Rclone struct {
	Paths Paths
}

// MoveTo acts as a wrapper around rclone's 'moveto' command.
func (r Rclone) MoveTo(source, destination string) ([]byte, error) {
	command := exec.Command("rclone", "moveto", source, destination, "--verbose")
	output, err := command.CombinedOutput()
	if err != nil {
		return []byte{}, err
	}

	return output, nil
}

// Mount acts as a wrapper around the rclone 'mount' command.
func (r Rclone) Mount() error {
	command := exec.Command(
		"rclone",
		"mount",
		"--allow-other",
		"--max-read-ahead",
		"2G",
		"--dir-cache-time",
		"1m0s",
		r.Paths.DecryptRemote,
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

// Unmount checks if rclone is currently mounted at a given path and unmounts it if it is.
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

// Mounted queries the mount status of an rclone mount at a given path.
func (r Rclone) Mounted() (bool, error) {
	activeMounts, err := exec.Command("mount").CombinedOutput()
	if err != nil {
		return false, err
	}

	return strings.Contains(string(activeMounts), r.Paths.Decrypt), nil
}
