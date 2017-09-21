package cmd

const (
	rootLong = `The Seedbox Mount Tool (sbmt) is a single binary intended to help ensure the integrity
of Plexdrive, Rclone and UnionFS mounts on a seedbox and synchronise data between
local and cloud storage.`

	mountLong = `Mounts encrypted files stored on Google Drive using Plexdrive, Rclone and UnionFS.

The encrypted remote files on Google Drive are mounted by Plexdrive in the folder
specified using the --plexdrive flag. A local Rclone mount then reads the
encrypted files in the plexdrive folder and shows their unencrypted representation
in the folder specified using the --decrypt flag. A UnionFS mount is then
created at the location specified by the --union flag, which shows a unified
representation of all of the files stored in Google Drive (in their unencrypted
form, in the --decrypt location), and files that have been downloaded to
the local machine, in a folder specified using the --local flag.

The flow of data is as follows:

GDrive (enc) -> Plexdrive (enc) -> Rclone (dec) -> UnionFS (dec).

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
	--plexdrive /plexdrive \
	--decrypt-remote plexdrive-decrypted: \
	--decrypt /decrypt \
	--local /local \
	--union /union
`
	cleanupLong = `Cleans up files that have been deleted on a UnionFS mount on an encrypted GDrive remote.

Files deleted from read-only section of a UnionFS mount are not actually deleted but
rather hidden from view. A record of these hidden files is kept in a hidden subfolder of
the unionfs mount location (.unionfs). The cleanup command will iterate through all the
files of this hidden subfolder, find the corresponding files in the location identified
using the --decrypt flag and remove them from Google Drive.

Example:

sbmt cleanup --decrypt /decrypt --union /union
`
	uploadLong = `Uploads any newly created files to an encrypted Rclone remote.

Files created in the read-write section of a UnionFS mount will be iterated over,
encrypted, and uploaded to the Rclone remote specified using the --encrypt-remote
flag.

The upload command uses the rclone 'moveto' command, which ensures that the file
to be uploaded will automatically be removed from the local read-write folder as
soon as the upload is confirmed as having been successful.

Example:

sbmt upload --local /local --encrypt-remote encrypted-remote:
`
	unmountLong = `Unmounts any active UnionFS, Rclone and Plexdrive mounts at the specified locations.

The mounts are unmounted in the following order to ensure that there are no problems
with busy resources that prevent a successful unmounting operation:

UnionFS -> Rclone -> Plexdrive

Example:

sbmt unmount --union /union --decrypt /decrypt --plexdrive /plexdrive
`
)
