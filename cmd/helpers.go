package cmd

import (
	ps "github.com/mitchellh/go-ps"
	"github.com/spf13/cobra"
)

type Paths struct {
	Decrypt,
	Local,
	Mount,
	PlexDrive,
	Union string
}

const (
	RcloneCmd = "rclone"
	MoveTo    = "moveto"
)

type RcloneDispatcher interface {
	MoveTo(source, destination string) ([]byte, error)
}

func isRunning() error {
	var sbmtProcesses []ps.Process

	processes, _ := ps.Processes()
	for _, p := range processes {
		if p.Executable() == "sbmt" {
			sbmtProcesses = append(sbmtProcesses, p)
		}
	}

	if len(sbmtProcesses) > 1 {
		return ErrSbmtProcessAlreadyRunning(sbmtProcesses[0].Pid())
	}

	return nil
}

func hasRequiredFlags(cmd *cobra.Command, f Flags) bool {
	if cmd.Use == "mount" {
		return len(f.PlexDriveFolder) > 0 &&
			len(f.LocalFolder) > 0 &&
			len(f.UnionFolder) > 0 &&
			len(f.DecryptFolder) > 0 &&
			len(f.DecryptRemote) > 0
	}

	if cmd.Use == "upload" {
		return len(f.LocalFolder) > 0 && len(f.EncryptRemote) > 0

	}

	if cmd.Use == "cleanup" {
		return len(f.DecryptFolder) > 0 && len(f.UnionFolder) > 0
	}

	return false
}
