package cmd

import (
	"fmt"

	"os"

	"strings"

	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

var uploadCmd = &cobra.Command{
	Use:   "upload",
	Short: "Upload newly created files to an encrypted Rclone Google Drive remote",
	Args:  cobra.NoArgs,
	Long:  uploadLong,
	Run: func(cmd *cobra.Command, args []string) {
		if !hasRequiredFlags(cmd, mountFlags) {
			fmt.Println(ErrMissingRequiredFlags)
			os.Exit(1)
		}

		var fs afero.Fs = afero.NewOsFs()
		var r Rclone

		if err := Upload(fs, r, uploadFlags); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

// LocalRemoteMapping represents the location of a file on a local filesystem and the location it will be uploaded to
// on an encrypted remote.
type LocalRemoteMapping struct {
	Source, Destination string
}

// Upload iterates through all files in the local read-write folder of a UnionFS mount and uploads each one to a
// location determined by a LocalRemoteMapping on an encrypted remote. Once an upload is successfully completed the
// original file is removed from the local filesystem.
func Upload(fs afero.Fs, r RcloneDispatcher, f Flags) error {
	if err := isRunning(); err != nil {
		return err
	}

	var localFiles []string

	w := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		isDir, err := afero.IsDir(fs, path)
		if err != nil {
			return err
		}

		if !isDir {
			localFiles = append(localFiles, path)
		}

		return nil
	}

	if err := afero.Walk(fs, f.LocalFolder, w); err != nil {
		return err
	}

	uploadMappings := GetUploadMappings(localFiles, f.LocalFolder, f.EncryptRemote)

	for _, m := range uploadMappings {
		output, err := r.MoveTo(m.Source, m.Destination)
		if err != nil {
			return err
		}

		fmt.Print(string(output))
	}

	return nil
}

// GetUploadMappings creates a list of LocalRemoteMapping objects for a list of file paths on a local filesystem.
func GetUploadMappings(localFiles []string, localRoot, remoteMount string) []LocalRemoteMapping {
	var uploadLocations []LocalRemoteMapping

	for _, f := range localFiles {
		relativePath := strings.TrimPrefix(f, localRoot)
		uploadPath := fmt.Sprintf("%s%s", remoteMount, relativePath)
		uploadLocations = append(uploadLocations, LocalRemoteMapping{Source: f, Destination: uploadPath})
	}

	return uploadLocations
}

var uploadFlags Flags

func init() {
	RootCmd.AddCommand(uploadCmd)

	uploadCmd.Flags().StringVar(&uploadFlags.LocalFolder, "local", "", "location of the local folder to upload from")
	uploadCmd.Flags().StringVar(&uploadFlags.EncryptRemote, "encrypt-remote", "", "name of the remote mount to upload to (with ':' suffix)")
}
