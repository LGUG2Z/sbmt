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

// UnionFS represents an interface to mount, unmount and verify a UnionFS mount.
type UnionFS struct {
	Paths Paths
}

// Mount acts as a wrapper around the UnionFS 'mount' command.
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

// Unmount checks if UnionFS is currently mounted at a given path and unmounts it if it is.
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

// Mounted queries the mount status of a UnionFS mount at a given path.
func (u UnionFS) Mounted() (bool, error) {
	activeMounts, err := exec.Command("mount").CombinedOutput()
	if err != nil {
		return false, err
	}

	return strings.Contains(string(activeMounts), u.Paths.Union), nil
}
