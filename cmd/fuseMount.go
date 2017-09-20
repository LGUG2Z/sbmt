package cmd

type FuseMount interface {
	Mount() error
	Unmount() error
	Mounted() (bool, error)
}
