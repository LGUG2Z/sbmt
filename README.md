# sbmt
Plexdrive, Rclone and UnionFS mount management made easier.

## Overview
The Seedbox Mount Tool (sbmt) is a single binary intended to help ensure the integrity
of Plexdrive, Rclone and UnionFS mounts on a seedbox and synchronise data between local
and cloud storage.

This project takes heavy inspiration from work already done by [Mads Lundt](https://github.com/madslundt/docker-cloud-media-scripts)
and all the contributors over at [hoarding.me](https://hoarding.me/).

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
