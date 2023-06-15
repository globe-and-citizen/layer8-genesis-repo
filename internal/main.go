package internal

type ServerImpl interface {
	Serve() error
	Shutdown()
}
