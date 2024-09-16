package app

type App interface {
	Run() error
	Shutdown() error
}
