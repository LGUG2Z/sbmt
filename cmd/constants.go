package cmd

import (
	"errors"
	"fmt"
)

const SuffixUnionFSHidden = "_HIDDEN~"

var (
	ErrMissingRequiredFlags      = errors.New("required flags for this command are missing")
	ErrSbmtProcessAlreadyRunning = func(pid int) error {
		return fmt.Errorf("an sbmt process is already running with pid %v", pid)
	}
	ErrCouldNotVerifyMount = func(mount string) error { return fmt.Errorf("could not verify successful mount of %s", mount) }
)
