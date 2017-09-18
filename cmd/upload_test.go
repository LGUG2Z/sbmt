package cmd_test

import (
	. "github.com/lgug2z/sbmt/cmd"

	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/spf13/afero"
)

var _ = Describe("Upload", func() {
	fs := afero.NewMemMapFs()
	f := Flags{
		LocalFolder:   "/local",
		DecryptFolder: "/decrypt",
		EncryptRemote: "remote:",
	}
	r := MockRclone{Fs: fs}

	BeforeEach(func() {
		fs = afero.NewMemMapFs()
		r.Fs = fs

		Expect(fs.MkdirAll("/local/sub", os.ModePerm)).To(Succeed())
		Expect(fs.MkdirAll("/decrypt", os.ModePerm)).To(Succeed())
	})

	Describe("When there are files in the local folder", func() {
		It("Should correctly map the upload location of a file in the root of the local folder", func() {
			file := "/local/a.txt"

			expected := []LocalRemoteMapping{LocalRemoteMapping{Source: file, Destination: "remote:/a.txt"}}
			actual := GetUploadMappings([]string{file}, f.LocalFolder, f.EncryptRemote)
			Expect(expected).To(Equal(actual))
		})

		It("Should correctly map the upload location of a file in a subfolder of the local folder", func() {
			file := "/local/x/y/z/a.txt"

			expected := []LocalRemoteMapping{LocalRemoteMapping{Source: file, Destination: "remote:/x/y/z/a.txt"}}
			actual := GetUploadMappings([]string{file}, f.LocalFolder, f.EncryptRemote)
			Expect(expected).To(Equal(actual))
		})

		It("Should move the local file to the remote and remove it from the local", func() {
			_, err := fs.Create("/local/sub/a.txt")
			Expect(err).ToNot(HaveOccurred())

			_ = captureStdout(func() {
				Expect(Upload(fs, r, f)).To(Succeed())
			})

			exists, err := afero.Exists(fs, "/local/sub/a.txt")
			Expect(err).ToNot(HaveOccurred())
			Expect(exists).To(BeFalse())
		})

		It("Should leave empty folders untouched after completing a move operation", func() {
			_, err := fs.Create("/local/sub/b.txt")
			Expect(err).ToNot(HaveOccurred())

			_ = captureStdout(func() {
				Expect(Upload(fs, r, f)).To(Succeed())
			})

			exists, err := afero.Exists(fs, "/local/sub")
			Expect(err).ToNot(HaveOccurred())
			Expect(exists).To(BeTrue())
		})
	})
})
