package cmd

import (
	ps "github.com/mitchellh/go-ps"
	"github.com/spf13/cobra"
)

// Paths represents all the paths and remotes required for mounting Plexdrive, decrypting with Rclone and presenting
// a unified view with UnionFS. Paths are used directly by types implementing the FuseMount interface.
type Paths struct {
	Decrypt,
	DecryptRemote,
	Local,
	PlexDrive,
	Union string
}

// RcloneDispatcher represents any object that is capable of executing a set of rclone commands.
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
