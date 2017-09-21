package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// mountCmd represents the mount command
var mountCmd = &cobra.Command{
	Use:   "mount",
	Short: "Set up and ensure integrity of Plexdrive, Rclone and UnionFS mounts",
	Long:  mountLong,
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		if !hasRequiredFlags(cmd, mountFlags) {
			fmt.Println(ErrMissingRequiredFlags)
			os.Exit(1)
		}

		paths := Paths{
			Decrypt:       mountFlags.DecryptFolder,
			DecryptRemote: mountFlags.DecryptRemote,
			Local:         mountFlags.LocalFolder,
			PlexDrive:     mountFlags.PlexDriveFolder,
			Union:         mountFlags.UnionFolder,
		}

		if err := Mount(Rclone{paths}, UnionFS{paths}, Plexdrive{paths}); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

// Mount verifies the integrity of a connected series of Plexdrive, Rclone and UnionFS mounts. If a mount is broken,
// everything will be forcibly unmounted and remounted.
func Mount(rclone, unionFS, plexDrive FuseMount) error {
	hasBrokenMounts, err := hasBrokenMounts(unionFS, rclone, plexDrive)
	if err != nil {
		return err
	}

	if hasBrokenMounts {
		fmt.Println("Broken mount detected. Remounting all.")
		if err := remount(unionFS, rclone, plexDrive); err != nil {
			return err
		}
	}

	fmt.Print("All mounts are active.")

	return nil
}

func hasBrokenMounts(unionFS, rclone, plexDrive FuseMount) (bool, error) {
	isMountedUnionFS, err := unionFS.Mounted()
	if err != nil {
		return false, err
	}

	isMountedRclone, err := rclone.Mounted()
	if err != nil {
		return false, err
	}

	isMountedPlexDrive, err := plexDrive.Mounted()
	if err != nil {
		return false, err
	}

	return !isMountedUnionFS || !isMountedRclone || !isMountedPlexDrive, nil
}

func remount(unionFS, rclone, plexDrive FuseMount) error {
	if err := unionFS.Unmount(); err != nil {
		return err
	}

	if err := rclone.Unmount(); err != nil {
		return err
	}

	if err := plexDrive.Unmount(); err != nil {
		return err
	}

	if err := plexDrive.Mount(); err != nil {
		return err
	}

	if err := rclone.Mount(); err != nil {
		return err
	}

	if err := unionFS.Mount(); err != nil {
		return err
	}

	return nil
}

var mountFlags Flags

func init() {
	RootCmd.AddCommand(mountCmd)
	mountCmd.Flags().StringVar(&mountFlags.UnionFolder, "union", "", "location of the unionfs mount folder")
	mountCmd.Flags().StringVar(&mountFlags.PlexDriveFolder, "plexdrive", "", "location of the plexdrive mount folder")
	mountCmd.Flags().StringVar(&mountFlags.LocalFolder, "local", "", "location of the local folder (union read-write)")
	mountCmd.Flags().StringVar(&mountFlags.DecryptFolder, "decrypt", "", "location of the decrypted plexdrive folder (union read-only)")
	mountCmd.Flags().StringVar(&mountFlags.DecryptRemote, "decrypt-remote", "", "name of the remote to use to decrypt data from plexdrive (with ':' suffix)")
}
