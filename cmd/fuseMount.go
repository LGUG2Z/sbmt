package cmd

// FuseMount represents any type of FUSE mount that can be mounted, unmounted and have its mount status queried.
type FuseMount interface {
	Mount() error
	Unmount() error
	Mounted() (bool, error)
}
