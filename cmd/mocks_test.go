package cmd

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/spf13/afero"
)

var _ = Describe("Mocks", func() {
	Describe("When given a local source file and a remote destination path", func() {
		var r MockRclone
		It("Should execute the MoveTo command", func() {
			r.Fs = afero.NewMemMapFs()
			source := "/local/sub/folder/a.txt"
			_, err := r.Fs.Create(source)
			Expect(err).ToNot(HaveOccurred())

			destination := "remote:sub/folder/a.txt"

			output, err := r.MoveTo(source, destination)
			Expect(err).ToNot(HaveOccurred())

			exists, err := afero.Exists(r.Fs, source)
			Expect(err).ToNot(HaveOccurred())
			Expect(exists).To(BeFalse())

			Expect(string(output)).To(Equal("Moved /local/sub/folder/a.txt to remote:sub/folder/a.txt\n"))
		})
	})
})
