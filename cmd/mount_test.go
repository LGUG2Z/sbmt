package cmd_test

import (
	. "github.com/lgug2z/sbmt/cmd"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Mount", func() {
	p := Paths{
		Decrypt:       "/decrypt",
		Local:         "/local",
		DecryptRemote: "remote:",
		Plexdrive:     "/plexdrive",
		Union:         "/union",
	}

	rclone := MockRclone{Paths: p}
	unionFS := MockUnionFS{Paths: p}
	plexdrive := MockPlexdrive{Paths: p}

	Describe("If any one mount is not active", func() {
		It("Should try to unmount all mounts before mounting them all again", func() {
			plexdrive.IsMountedBool = false
			rclone.IsMountedBool = true
			unionFS.IsMountedBool = true

			output, err := captureStdout(func() {
				Expect(Mount(rclone, unionFS, plexdrive)).To(Succeed())
			})

			Expect(err).ToNot(HaveOccurred())

			Expect(output).To(ContainSubstring("Broken mount detected. Remounting all."))
		})
	})

	Describe("If all mounts are not active", func() {
		It("Should not take any action", func() {
			plexdrive.IsMountedBool = true
			rclone.IsMountedBool = true
			unionFS.IsMountedBool = true

			output, err := captureStdout(func() {
				Expect(Mount(rclone, unionFS, plexdrive)).To(Succeed())
			})

			Expect(err).ToNot(HaveOccurred())

			Expect(output).ToNot(ContainSubstring("Broken mount detected. Remounting all."))
		})
	})

	Describe("If a remount cannot be verified", func() {
		It("Should throw an error", func() {
			plexdrive.IsMountedBool = true
			rclone.IsMountedBool = true
			unionFS.IsMountedBool = false
			rclone.MountError = ErrCouldNotVerifyMount(rclone.Paths.Decrypt)

			_, err := captureStdout(func() {
				err := Mount(rclone, unionFS, plexdrive)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal(rclone.MountError.Error()))
			})

			Expect(err).ToNot(HaveOccurred())
		})
	})
})
