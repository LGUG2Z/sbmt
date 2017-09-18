package cmd_test

import (
	. "github.com/lgug2z/sbmt/cmd"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Mount", func() {
	p := Paths{
		Decrypt:   "/decrypt",
		Local:     "/local",
		Mount:     "remote:",
		PlexDrive: "/plexdrive",
		Union:     "/union",
	}

	rclone := MockRclone{Paths: p}
	unionFS := MockUnionFS{Paths: p}
	plexDrive := MockPlexDrive{Paths: p}

	Describe("If any one mount is not active", func() {
		It("Should try to unmount all mounts before mounting them all again", func() {
			plexDrive.IsMountedBool = false
			rclone.IsMountedBool = true
			unionFS.IsMountedBool = true

			output := captureStdout(func() {
				Expect(Mount(rclone, unionFS, plexDrive)).To(Succeed())
			})

			Expect(output).To(ContainSubstring("Broken mount detected. Remounting all."))
		})
	})

	Describe("If all mounts are not active", func() {
		It("Should not take any action", func() {
			plexDrive.IsMountedBool = true
			rclone.IsMountedBool = true
			unionFS.IsMountedBool = true

			output := captureStdout(func() {
				Expect(Mount(rclone, unionFS, plexDrive)).To(Succeed())
			})

			Expect(output).ToNot(ContainSubstring("Broken mount detected. Remounting all."))
		})
	})

	Describe("If a remount cannot be verified", func() {
		It("Should throw an error", func() {
			plexDrive.IsMountedBool = true
			rclone.IsMountedBool = true
			unionFS.IsMountedBool = false
			rclone.MountError = ErrCouldNotVerifyMount(rclone.Paths.Decrypt)

			_ = captureStdout(func() {
				err := Mount(rclone, unionFS, plexDrive)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal(rclone.MountError.Error()))
			})
		})
	})
})
