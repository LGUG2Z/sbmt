package cmd

type MounterUnmounter interface {
	Mount() error
	Unmount() error
	Mounted() (bool, error)
}
