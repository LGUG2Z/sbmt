package cmd

const (
	MountLong = `Mounts encrypted files stored on Google Drive using Plexdrive, Rclone and UnionFS.

The encrypted remote files on Google Drive are mounted by Plexdrive in the folder
specified using the --plexdrive-folder flag. A local Rclone mount then reads the
encrypted files in the plexdrive folder and shows their unencrypted representation
in the folder specified using the --decrypt-folder flag. A UnionFS mount is then
created at the location specified by the --union-folder flag, which shows a unified
representation of all of the files stored in Google Drive (in their unencrypted
form, in the --decrypt-folder location), and files that have been downloaded to
the local machine, in a folder specified using the --local-folder flag.

The flow of data is as follows:

GDrive (enc) -> Plexdrive Mount (enc) -> Rclone Mount (dec) -> UnionFS Mount (dec).

In order to show an unencrypted representation of the data mounted by Plexdrive,
a separate Rclone remote will have to be created, in which the remote location
is set to the local folder where Plexdrive will be mounted. The encryption type
and the encryption passwords should be the same as those set for the encrypted
Google Drive remote. The name of this remote must be passed to the mount command
using the --decrypt-remote flag.

If the mount command is run and all mounts are intact, no further action will be
taken. Conversely, if any single mount is broken, every mount will be unmounted
and remounted to re-establish mount integrity.

Example:

sbmt mount \
	--plexdrive-folder /plexdrive \
	--decrypt-remote plexdrive-decrypted: \
	--decrypt-folder /decrypt \
	--local-folder /local \
	--union-folder /union
`
	CleanupLong = ``
	UploadLong  = ``
)
