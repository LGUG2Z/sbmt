package cmd

import (
	"fmt"

	"os"

	"strings"

	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

// uploadCmd represents the upload command
var uploadCmd = &cobra.Command{
	Use:   "upload",
	Short: "A brief description of your command",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		var fs afero.Fs = afero.NewOsFs()
		var r Rclone

		if err := Upload(fs, r, uploadFlags); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

type LocalRemoteMapping struct {
	Source, Destination string
}

func Upload(fs afero.Fs, r RcloneDispatcher, f Flags) error {
	if err := isRunning(); err != nil {
		return err
	}

	var localFiles []string

	w := func(path string, info os.FileInfo, err error) error {
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

	uploadMappings := GetUploadMappings(localFiles, f.LocalFolder, f.RemoteMount)

	for _, m := range uploadMappings {
		output, err := r.MoveTo(m.Source, m.Destination)
		if err != nil {
			return err
		}

		fmt.Printf(string(output))
	}

	return nil
}

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

	uploadCmd.Flags().StringVarP(&uploadFlags.LocalFolder, "local-folder", "l", "", "location of the local (rw) folder")
	uploadCmd.Flags().StringVarP(&uploadFlags.RemoteMount, "remote-mount", "r", "", "name of the remote mount to upload to (with ':')")
}
