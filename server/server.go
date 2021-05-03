package server

type Server interface {
	Serve()
	Clean() error
}
