package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"time"
)

type PlexDrive struct {
	Paths Paths
}

func (plx PlexDrive) Mount() error {
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

func (plx PlexDrive) Unmount() error {
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

func (plx PlexDrive) Mounted() (bool, error) {
	activeMounts, err := exec.Command("mount").CombinedOutput()
	if err != nil {
		return false, err
	}

	return strings.Contains(string(activeMounts), fmt.Sprintf("%s", plx.Paths.PlexDrive)), nil
}
