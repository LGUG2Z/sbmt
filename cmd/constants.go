package cmd

import (
	"errors"
	"fmt"
)

// SuffixUnionFSHidden is the Union FS suffix appended to read-only files that have had an attempted deletion.
const SuffixUnionFSHidden = "_HIDDEN~"

var (
	// ErrMissingRequiredFlags is returned when required flags for a command are missing.
	ErrMissingRequiredFlags = errors.New("required flags for this command are missing")

	// ErrSbmtProcessAlreadyRunning is returned if the upload command is called while another sbmt process is running.
	ErrSbmtProcessAlreadyRunning = func(pid int) error {
		return fmt.Errorf("an sbmt process is already running with pid %v", pid)
	}

	// ErrCouldNotVerifyMount is returned if it is not possible to verify a successful mounting operation.
	ErrCouldNotVerifyMount = func(mount string) error { return fmt.Errorf("could not verify successful mount of %s", mount) }
)
