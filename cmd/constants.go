package cmd

import "fmt"

const SuffixUnionFSHidden = "_HIDDEN~"

var (
	ErrSbmProcessAlreadyRunning = func(pid int) error {
		return fmt.Errorf("An sbmt process is already running with pid %v. Not continuing.", pid)
	}
	ErrFailedToMount       = func(folder string) error { return fmt.Errorf("Failed to mount %s.\n", folder) }
	ErrCouldNotVerifyMount = func(mount string) error { return fmt.Errorf("Could not verify successful mount of %s.", mount) }
)
