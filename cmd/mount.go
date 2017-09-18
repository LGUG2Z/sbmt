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
	Long:  MountLong,
	Run: func(cmd *cobra.Command, args []string) {
		paths := Paths{
			Decrypt:   mountFlags.DecryptFolder,
			Mount:     mountFlags.DecryptRemote,
			Local:     mountFlags.LocalFolder,
			PlexDrive: mountFlags.PlexDriveFolder,
			Union:     mountFlags.UnionFolder,
		}

		if err := Mount(Rclone{paths}, UnionFS{paths}, PlexDrive{paths}); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func Mount(rclone, unionFS, plexDrive MounterUnmounter) error {
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

func hasBrokenMounts(unionFS, rclone, plexDrive MounterUnmounter) (bool, error) {
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

func remount(unionFS, rclone, plexDrive MounterUnmounter) error {
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

type Printer struct{}

var mountFlags Flags

func init() {
	RootCmd.AddCommand(mountCmd)
	mountCmd.Flags().StringVarP(&mountFlags.UnionFolder, "union-folder", "u", "", "location of the unionfs folder")
	mountCmd.Flags().StringVarP(&mountFlags.PlexDriveFolder, "plexdrive-folder", "p", "", "location of the plexdrive folder")
	mountCmd.Flags().StringVarP(&mountFlags.LocalFolder, "local-folder", "l", "", "location of the local folder (union read-write)")
	mountCmd.Flags().StringVarP(&mountFlags.DecryptFolder, "decrypt-folder", "d", "", "location of the decrypted plexdrive folder (union read-only)")
	mountCmd.Flags().StringVarP(&mountFlags.DecryptRemote, "decrypt-remote", "m", "", "name of the remote to use to decrypt data from plexdrive (with trailing :)")
}
