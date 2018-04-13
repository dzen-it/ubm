package ubm

// Server Listen and Serve requests
type Server interface {
	Serve() error
}
