package cmd

import (
	"os"
	"os/exec"
	"strings"
	"syscall"
	"time"
)

// Plexdrive represents an interface to mount, unmount and verify a Plexdrive mount.
type Plexdrive struct {
	Paths Paths
}

// Mount acts as a wrapper around the Plexdrive 'mount' command.
func (plx Plexdrive) Mount() error {
	command := exec.Command(
		"plexdrive",
		"mount",
		"-o",
		"allow_other",
		plx.Paths.PlexDrive,
	)

	command.Env = os.Environ()

	if err := command.Start(); err != nil {
		return err
	}

	for i := 0; i < 10; i++ {
		isMounted, err := plx.Mounted()
		if err != nil {
			return err
		}

		if !isMounted {
			time.Sleep(100 * time.Millisecond)
		}

		if isMounted {
			return nil
		}
	}

	return ErrCouldNotVerifyMount(plx.Paths.PlexDrive)
}

// Unmount checks if Plexdrive is currently mounted at a given path and unmounts it if it is.
func (plx Plexdrive) Unmount() error {
	isMounted, err := plx.Mounted()
	if err != nil {
		return err
	}

	if isMounted {
		if err := syscall.Unmount(plx.Paths.PlexDrive, 0); err != nil {
			return err
		}
	}

	return nil
}

// Mounted queries the mount status of Plexdrive mount at a given path.
func (plx Plexdrive) Mounted() (bool, error) {
	activeMounts, err := exec.Command("mount").CombinedOutput()
	if err != nil {
		return false, err
	}

	return strings.Contains(string(activeMounts), plx.Paths.PlexDrive), nil
}
