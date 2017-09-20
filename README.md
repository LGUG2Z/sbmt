# sbmt
Plexdrive, Rclone and UnionFS mount management made easier.

## Overview
The Seedbox Mount Tool (sbmt) is a single binary intended to help ensure the integrity
of Plexdrive, Rclone and UnionFS mounts on a seedbox and synchronise data between local
and cloud storage.

This project takes heavy inspiration from work already done by [Mads Lundt](https://github.com/madslundt/docker-cloud-media-scripts)
Jamie and all the contributors over at [hoarding.me](https://hoarding.me/) and the various
scripts and tips I found posted on the [PlexACD subreddit](https://reddit.com/r/plexacd).

With `sbmt` I am trying to hit the sweet spot between the approach of Dockerising everything
and maintaining an array of bash scripts to enable the smooth running of a seedbox running
[Plex](https://plex.tv) and the seamless synchronisation of local and remotely stored media.
An additional consideration that has had a big influence on the direction of `sbmt` is
the desire to easily be able to switch from one seedbox to another and get up and running
again as quickly as possible.

With a single binary managing mounting and ensuring mount integrity, dealing with synchronising
of local and remote data as well as ensuring any required cleanups of local and remote data,
cron jobs can be copied easily from one seedbox to another, either maintaining an old folder
structure or implementing a new folder structure just by changing the flags passed to different
commands.

## Installation
The latest version of bfm can be installed using `go get`.

```
go get -u github.com/LGUG2Z/sbmt
```

Make sure `$GOPATH` is set correctly that and that `$GOPATH/bin` is in your `$PATH`.

The `sbmt` executable will be installed under the `$GOPATH/bin` directory.

## Prior Setup
Before scheduling `sbmt` commands to be run on a seedbox, there are some pieces of prior
setup required:

1. Create an [unencrypted remote for Google Drive](https://rclone.org/drive/) using Rclone. 
Files uploaded to this remote will be stored unencrypted.
2. Create an [encrypted remote for Google Drive](https://rclone.org/drive/) using Rclone.
Files uploaded to this remote will be stored encrypted.
3. Configure [Plexdrive](https://github.com/dweidenfeld/plexdrive) to be able to access your Google Drive account.
Encrypted files viewed on this mount remain encrypted.
4. Create a remote to [decrypt the encrypted files mounted by Plexdrive](https://github.com/dweidenfeld/plexdrive/issues/206).
Encrypted files viewed on this mount will be shown in their unencrypted form.

## Usage

The three commands of `sbmt` are `mount`, `upload` and `cleanup`.

### mount

The encrypted remote files on Google Drive are mounted by Plexdrive in the folder
specified using the `--plexdrive-folder` flag. A local Rclone mount then reads the
encrypted files in the plexdrive folder and shows their unencrypted representation
in the folder specified using the `--decrypt-folder` flag. A UnionFS mount is then
created at the location specified by the `--union-folder` flag, which shows a unified
representation of all of the files stored in Google Drive (in their unencrypted
form, in the `--decrypt-folder` location), and files that have been downloaded to
the local machine, in a folder specified using the `--local-folder` flag.

The flow of data is as follows:

GDrive (enc) -> Plexdrive Mount (enc) -> Rclone Mount (dec) -> UnionFS Mount (dec).

In order to show an unencrypted representation of the data mounted by Plexdrive,
a separate Rclone remote will have to be created, in which the remote location
is set to the local folder where Plexdrive will be mounted. The encryption type
and the encryption passwords should be the same as those set for the encrypted
Google Drive remote. The name of this remote must be passed to the mount command
using the `--decrypt-remote` flag.

If the mount command is run and all mounts are intact, no further action will be
taken. Conversely, if any single mount is broken, every mount will be unmounted
and remounted to re-establish mount integrity.

Example:

```bash
sbmt mount \
    --plexdrive-folder /plexdrive \
    --decrypt-remote plexdrive-decrypted: \
    --decrypt-folder /decrypt \
    --local-folder /local \
    --union-folder /union
```

### cleanup

Files deleted from read-only section of a UnionFS mount are not actually deleted but
rather hidden from view. A record of these hidden files is kept in a hidden subfolder of
the unionfs mount location (.unionfs). The cleanup command will iterate through all the
files of this hidden subfolder, find the corresponding files in the location identified
using the `--decrypt-folder` flag and remove them from Google Drive.

Example:
```bash
sbmt cleanup \
    --decrypt-folder /decrypt \
    --union-folder /union
```

### upload

Files created in the read-write section of a UnionFS mount will be iterated over,
encrypted, and uploaded to the Rclone remote specified using the `--encrypt-remote`
flag.

The upload command uses the rclone `moveto` command, which ensures that the file
to be uploaded will automatically be removed from the local read-write folder as
soon as the upload is confirmed as having been successful.

Example:

```bash
sbmt upload \
    --local-folder /local \
    --encrypt-remote encrypted-remote:
```
