package cmd

import ps "github.com/mitchellh/go-ps"

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
		return ErrSbmProcessAlreadyRunning(sbmtProcesses[0].Pid())
	}

	return nil
}
