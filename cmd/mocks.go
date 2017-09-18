package cmd

import (
	"fmt"

	"github.com/spf13/afero"
)

type MockRclone struct {
	Fs             afero.Fs
	MountError     error
	UnmountError   error
	IsMountedError error
	IsMountedBool  bool
	Paths          Paths
}

type MockUnionFS struct {
	Paths          Paths
	MountError     error
	UnmountError   error
	IsMountedError error
	IsMountedBool  bool
}

type MockPlexDrive struct {
	Paths          Paths
	MountError     error
	UnmountError   error
	IsMountedError error
	IsMountedBool  bool
}

func (plx MockPlexDrive) Mount() error {
	return plx.MountError
}

func (u MockUnionFS) Mount() error {
	return u.MountError
}

func (f MockRclone) Mount() error {
	return f.MountError
}

func (plx MockPlexDrive) Unmount() error {
	return plx.UnmountError
}

func (u MockUnionFS) Unmount() error {
	return u.UnmountError
}

func (f MockRclone) Unmount() error {
	return f.UnmountError
}

func (plx MockPlexDrive) Mounted() (bool, error) {
	return plx.IsMountedBool, plx.IsMountedError
}

func (u MockUnionFS) Mounted() (bool, error) {
	return u.IsMountedBool, u.IsMountedError
}

func (f MockRclone) Mounted() (bool, error) {
	return f.IsMountedBool, f.IsMountedError
}

func (f MockRclone) MoveTo(source, destination string) ([]byte, error) {
	if err := f.Fs.Remove(source); err != nil {
		return []byte{}, nil
	}

	s := fmt.Sprintf("Moved %s to %s\n", source, destination)
	return []byte(s), nil
}