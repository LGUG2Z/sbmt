package cmd

import (
	"fmt"

	"os"

	"github.com/spf13/cobra"
)

var unmountCmd = &cobra.Command{
	Use:   "unmount",
	Short: "Unmount any active UnionFS, Rclone and Plexdrive mounts",
	Long:  unmountLong,
	Run: func(cmd *cobra.Command, args []string) {
		if !hasRequiredFlags(cmd, unmountFlags) {
			fmt.Println(ErrMissingRequiredFlags)
			os.Exit(1)
		}

		paths := Paths{
			Decrypt:   unmountFlags.DecryptFolder,
			Plexdrive: unmountFlags.PlexdriveFolder,
			Union:     unmountFlags.UnionFolder,
		}

		if err := Unmount(Rclone{paths}, UnionFS{paths}, Plexdrive{paths}); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

// Unmount unmounts a series of UnionFS, Rclone and Plexdrive mounts in that order to ensure that no mount has a busy
// status which stops the successful completion of an unmounting operation of another mount.
func Unmount(rclone, unionFS, plexdrive FuseMount) error {
	for {
		if err := unmountAll(unionFS, rclone, plexdrive); err != nil {
			return err
		}

		hasActiveMounts, err := hasActiveMounts(unionFS, rclone, plexdrive)
		if err != nil {
			return err
		}

		if !hasActiveMounts {
			break
		}
	}

	fmt.Println("All active mounts successfully unmounted.")

	return nil
}

func hasActiveMounts(unionFS, rclone, plexdrive FuseMount) (bool, error) {
	isMountedUnionFS, err := unionFS.Mounted()
	if err != nil {
		return false, err
	}

	isMountedRclone, err := rclone.Mounted()
	if err != nil {
		return false, err
	}

	isMountedPlexdrive, err := plexdrive.Mounted()
	if err != nil {
		return false, err
	}

	return isMountedUnionFS || isMountedRclone || isMountedPlexdrive, nil
}

var unmountFlags Flags

func init() {
	RootCmd.AddCommand(unmountCmd)
	unmountCmd.Flags().StringVar(&unmountFlags.UnionFolder, "union", "", "location of the unionfs mount folder")
	unmountCmd.Flags().StringVar(&unmountFlags.PlexdriveFolder, "plexdrive", "", "location of the plexdrive mount folder")
	unmountCmd.Flags().StringVar(&unmountFlags.DecryptFolder, "decrypt", "", "location of the rclone decryption mount folder")
}
