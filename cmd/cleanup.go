package cmd

import (
	"fmt"

	"os"

	"strings"

	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

var cleanupCmd = &cobra.Command{
	Use:   "cleanup",
	Short: "Clean up encrypted files deleted from a UnionFS mount on Google Drive",
	Long:  cleanupLong,
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		if !hasRequiredFlags(cmd, cleanupFlags) {
			fmt.Println(ErrMissingRequiredFlags)
			os.Exit(1)
		}

		fs := afero.NewOsFs()

		if err := Cleanup(fs, cleanupFlags); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

// Cleanup looks for files hidden by 'deletion' from a read-only UnionFS mount and removes them from a remote.
func Cleanup(fs afero.Fs, f Flags) error {
	unionDeletionsFolder := fmt.Sprintf("%s/%s/", f.UnionFolder, ".unionfs")
	var unionDeletions []string

	w := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		isDir, err := afero.IsDir(fs, path)
		if err != nil {
			return err
		}

		if !isDir {
			unionDeletions = append(unionDeletions, path)
		}

		return nil
	}

	if err := afero.Walk(fs, unionDeletionsFolder, w); err != nil {
		return err
	}

	unionDeletionPaths := getUnionFSDeletions(unionDeletions)
	mountDeletionPaths := getMountPaths(unionDeletionPaths, f.DecryptFolder, unionDeletionsFolder)

	for _, d := range mountDeletionPaths {
		if err := delete(fs, d); err != nil {
			return err
		}
		fmt.Printf("%s deleted.\n", d)
	}

	for _, d := range unionDeletionPaths {
		if err := cleanup(fs, d); err != nil {
			return err
		}
		fmt.Printf("%s cleaned up.\n", d)
	}

	for {
		emptyFolderPaths, err := getEmptyHiddenFolders(fs, unionDeletionsFolder)
		if err != nil {
			return err
		}

		if len(emptyFolderPaths) == 0 {
			break
		}

		for _, d := range emptyFolderPaths {
			if e := cleanup(fs, d); e != nil {
				return e
			}
			fmt.Printf("%s cleaned up.\n", d)
		}
	}

	//emptyFolderPaths, err := getEmptyHiddenFolders(fs, unionDeletionsFolder)
	//if err != nil {
	//	return err
	//}
	//
	//for _, d := range emptyFolderPaths {
	//	if e := cleanup(fs, d); e != nil {
	//		return e
	//	}
	//	fmt.Printf("%s cleaned up.\n", d)
	//}

	isEmpty, err := afero.IsEmpty(fs, unionDeletionsFolder)
	if err != nil {
		return err
	}

	if isEmpty {
		if err := fs.Remove(unionDeletionsFolder); err != nil {
			return err
		}
	}

	return nil
}

func getEmptyHiddenFolders(fs afero.Fs, hiddenRoot string) ([]string, error) {
	var emptyFolders []string

	w := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		isDir, err := afero.IsDir(fs, path)
		if err != nil {
			return err
		}

		if isDir {
			isEmpty, err := afero.IsEmpty(fs, path)
			if err != nil {
				return err
			}

			if isEmpty {
				if path != hiddenRoot {
					emptyFolders = append(emptyFolders, path)
				}
			}
		}

		return nil
	}

	if err := afero.Walk(fs, hiddenRoot, w); err != nil {
		return []string{}, err
	}

	return emptyFolders, nil
}

func cleanup(fs afero.Fs, toCleanup string) error {
	return fs.Remove(toCleanup)
}

func delete(fs afero.Fs, toDelete string) error {
	return fs.Remove(toDelete)
}

func getUnionFSDeletions(absPaths []string) []string {
	var forCleanup []string

	for _, p := range absPaths {
		if strings.HasSuffix(p, SuffixUnionFSHidden) {
			forCleanup = append(forCleanup, p)
		}
	}

	return forCleanup
}

func getMountPaths(absPaths []string, mountPath, hiddenRoot string) []string {
	var mntPaths []string

	for _, p := range absPaths {
		relativePath := strings.TrimPrefix(p, hiddenRoot)
		unhiddenPath := strings.TrimSuffix(relativePath, SuffixUnionFSHidden)
		mPath := fmt.Sprintf("%s/%s", mountPath, unhiddenPath)
		mntPaths = append(mntPaths, mPath)
	}

	return mntPaths
}

var cleanupFlags Flags

func init() {
	RootCmd.AddCommand(cleanupCmd)
	cleanupCmd.Flags().StringVar(&cleanupFlags.UnionFolder, "union", "", "location of the unionfs folder")
	cleanupCmd.Flags().StringVar(&cleanupFlags.DecryptFolder, "decrypt", "", "location of the decrypted data folder")
}
