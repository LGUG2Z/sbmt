package cmd_test

import (
	. "github.com/lgug2z/sbmt/cmd"

	"os"

	"fmt"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/spf13/afero"
)

func UnionFsDelete(fs afero.Fs, unionPath, deletionPath string) error {

	relativePath := strings.TrimPrefix(deletionPath, unionPath)
	hiddenFilePath := fmt.Sprintf("%s/%s%s%s", unionPath, ".unionfs", relativePath, SuffixUnionFSHidden)
	normalFolderPath := fmt.Sprintf("%s/%s%s", unionPath, ".unionfs", relativePath)

	isDir, err := afero.IsDir(fs, deletionPath)
	if err != nil {
		return err
	}

	// If we are dealing with a folder
	if isDir {
		isEmpty, err := afero.IsEmpty(fs, deletionPath)
		if err != nil {
			return err
		}

		// If the folder is empty, create a _HIDDEN~ folder
		if isEmpty {
			if err := fs.MkdirAll(hiddenFilePath, os.ModePerm); err != nil {
				return err
			}
			// If the folder is not empty
		} else {
			// Make a folder without the _HIDDEN~ suffix
			if err := fs.MkdirAll(normalFolderPath, os.ModePerm); err != nil {
				return err
			}

			// Make _HIDDEN~ files inside that folder of its contents
			files, err := afero.ReadDir(fs, deletionPath)
			if err != nil {
				return err
			}

			for _, f := range files {
				err := UnionFsDelete(fs, unionPath, fmt.Sprintf("%s/%s", deletionPath, f.Name()))
				if err != nil {
					return err
				}
			}

			// Make a _HIDDEN~ version of the now empty directory
			if err := fs.MkdirAll(hiddenFilePath, os.ModePerm); err != nil {
				return err
			}
		}
		// If we are dealing with a file, create the _HIDDEN~ file with the normal directory names
	} else {
		_, err := fs.Create(hiddenFilePath)
		if err != nil {
			return err
		}
	}

	// Remove the path from the main union folder
	if err := fs.Remove(deletionPath); err != nil {
		return err
	}

	return nil
}

var _ = Describe("Cleanup", func() {
	var fs afero.Fs
	var f Flags

	BeforeEach(func() {
		fs = afero.NewMemMapFs()
		Expect(fs.MkdirAll("/union/sub", os.ModePerm)).To(Succeed())
		Expect(fs.MkdirAll("/mount/sub", os.ModePerm)).To(Succeed())

		f.DecryptFolder = "/mount"
		f.UnionFolder = "/union"
	})

	Describe("When a UnionFS deletion is simulated", func() {
		It("Should create hidden files for the deleted files in the .unionfs folder", func() {
			_, err := fs.Create("/union/a")
			Expect(err).ToNot(HaveOccurred())

			Expect(UnionFsDelete(fs, f.UnionFolder, "/union/a")).To(Succeed())
			exists1, err := afero.Exists(fs, "/union/a")
			Expect(err).ToNot(HaveOccurred())
			Expect(exists1).To(BeFalse())

			exists2, err := afero.Exists(fs, "/union/.unionfs/a_HIDDEN~")
			Expect(err).ToNot(HaveOccurred())
			Expect(exists2).To(BeTrue())
		})

		It("Should create folder for the deleted folder in the .unionfs folder with its contents as hidden", func() {
			_, err := fs.Create("/union/sub/a")
			Expect(err).ToNot(HaveOccurred())

			Expect(UnionFsDelete(fs, f.UnionFolder, "/union/sub")).To(Succeed())

			exists, err := afero.Exists(fs, "/union/sub")
			Expect(err).ToNot(HaveOccurred())
			Expect(exists).To(BeFalse())

			exists, err = afero.Exists(fs, "/union/sub/a")
			Expect(err).ToNot(HaveOccurred())
			Expect(exists).To(BeFalse())

			isDir, err := afero.IsDir(fs, "/union/.unionfs/sub_HIDDEN~")
			Expect(err).ToNot(HaveOccurred())
			Expect(isDir).To(BeTrue())

			exists, err = afero.Exists(fs, "/union/.unionfs/sub/a_HIDDEN~")
			Expect(err).ToNot(HaveOccurred())
			Expect(exists).To(BeTrue())
		})

		It("Should handle deeply nested files correctly", func() {
			_, err := fs.Create("/union/sub/a/b/c.txt")
			Expect(err).ToNot(HaveOccurred())

			Expect(UnionFsDelete(fs, f.UnionFolder, "/union/sub")).To(Succeed())

			exists, err := afero.Exists(fs, "/union/.unionfs/sub/a/b/c.txt_HIDDEN~")
			Expect(err).ToNot(HaveOccurred())
			Expect(exists).To(BeTrue())

			isDir, err := afero.IsDir(fs, "/union/.unionfs/sub_HIDDEN~")
			Expect(err).ToNot(HaveOccurred())
			Expect(isDir).To(BeTrue())

			isDir, err = afero.IsDir(fs, "/union/.unionfs/sub/a_HIDDEN~")
			Expect(err).ToNot(HaveOccurred())
			Expect(isDir).To(BeTrue())

			isDir, err = afero.IsDir(fs, "/union/.unionfs/sub/a/b_HIDDEN~")
			Expect(err).ToNot(HaveOccurred())
			Expect(isDir).To(BeTrue())
		})
	})

	Describe("When a file has been deleted from the UnionFS mount", func() {
		It("Should remove those files from the rclone mount folder", func() {
			_, err := fs.Create("/union/sub/a")
			Expect(err).ToNot(HaveOccurred())
			_, err = fs.Create("/mount/sub/a")
			Expect(err).ToNot(HaveOccurred())

			Expect(UnionFsDelete(fs, f.UnionFolder, "/union/sub/a")).To(Succeed())

			_ = captureStdout(func() {
				Expect(Cleanup(fs, f)).To(Succeed())
			})

			exists, err := afero.Exists(fs, "/mount/sub/a")
			Expect(err).ToNot(HaveOccurred())
			Expect(exists).To(BeFalse())
		})

		It("Cleanup should not affect files not having had an attempted deletion by UnionFS", func() {
			_, err := fs.Create("/union/a")
			Expect(err).ToNot(HaveOccurred())
			_, err = fs.Create("/mount/a")
			Expect(err).ToNot(HaveOccurred())

			_, err = fs.Create("/union/sub/b")
			Expect(err).ToNot(HaveOccurred())
			_, err = fs.Create("/mount/sub/b")
			Expect(err).ToNot(HaveOccurred())

			Expect(UnionFsDelete(fs, f.UnionFolder, "/union/a")).To(Succeed())

			_ = captureStdout(func() {
				Expect(Cleanup(fs, f)).To(Succeed())
			})

			exists, err := afero.Exists(fs, "/mount/sub/b")
			Expect(err).ToNot(HaveOccurred())
			Expect(exists).To(BeTrue())
		})

		It("Cleanup should remove the hidden file created by UnionFS after removal from the mount", func() {
			_, err := fs.Create("/union/a")
			Expect(err).ToNot(HaveOccurred())
			_, err = fs.Create("/mount/a")
			Expect(err).ToNot(HaveOccurred())

			Expect(UnionFsDelete(fs, f.UnionFolder, "/union/a")).To(Succeed())

			_ = captureStdout(func() {
				Expect(Cleanup(fs, f)).To(Succeed())
			})

			exists, err := afero.Exists(fs, "/union/.unionfs/a_HIDDEN~")
			Expect(err).ToNot(HaveOccurred())
			Expect(exists).To(BeFalse())
		})

		Describe("When a hidden folder becomes empty after a cleanup", func() {
			It("Should deleted at the end of the cleanup operation", func() {
				_, err := fs.Create("/union/sub/a")
				Expect(err).ToNot(HaveOccurred())
				_, err = fs.Create("/mount/sub/a")
				Expect(err).ToNot(HaveOccurred())

				Expect(UnionFsDelete(fs, f.UnionFolder, "/union/sub/a")).To(Succeed())

				_ = captureStdout(func() {
					Expect(Cleanup(fs, f)).To(Succeed())
				})

				exists, err := afero.Exists(fs, "/union/.unionfs/sub")
				Expect(err).ToNot(HaveOccurred())
				Expect(exists).To(BeFalse())
			})
		})

		Describe("When the .unionfs folder is empty at the end of a cleanup", func() {
			It("Should be deleted", func() {
				_, err := fs.Create("/union/sub/a")
				Expect(err).ToNot(HaveOccurred())
				_, err = fs.Create("/mount/sub/a")
				Expect(err).ToNot(HaveOccurred())

				Expect(UnionFsDelete(fs, f.UnionFolder, "/union/sub/a")).To(Succeed())

				_ = captureStdout(func() {
					Expect(Cleanup(fs, f)).To(Succeed())
				})

				exists, err := afero.Exists(fs, "/union/.unionfs")
				Expect(err).ToNot(HaveOccurred())
				Expect(exists).To(BeFalse())
			})
		})
	})
})
