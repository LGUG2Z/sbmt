package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"
	"syscall"
	"time"
)

type MounterUnmounter interface {
	Mount() error
	Unmount() error
	Mounted() (bool, error)
}

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

type UnionFS struct {
	Paths Paths
}

func (u UnionFS) Mount() error {
	v := path.Base(u.Paths.Union)

	command := exec.Command(
		"unionfs",
		"-o",
		"cow,allow_other",
		"-o",
		fmt.Sprintf("volname=%s", v),
		fmt.Sprintf("%s=RW:%s=RO", u.Paths.Local, u.Paths.Decrypt),
		u.Paths.Union,
	)

	command.Env = os.Environ()

	if err := command.Start(); err != nil {
		return err
	}

	for i := 0; i < 10; i++ {
		isMounted, err := u.Mounted()
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

	return ErrCouldNotVerifyMount(u.Paths.Union)
}

func (u UnionFS) Unmount() error {
	isMounted, err := u.Mounted()
	if err != nil {
		return err
	}

	if isMounted {
		if err := syscall.Unmount(u.Paths.Union, 0); err != nil {
			return err
		}
	}

	return nil
}

func (u UnionFS) Mounted() (bool, error) {
	activeMounts, err := exec.Command("mount").CombinedOutput()
	if err != nil {
		return false, err
	}

	return strings.Contains(string(activeMounts), fmt.Sprintf("%s", u.Paths.Union)), nil
}
