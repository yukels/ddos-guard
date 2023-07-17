package server

// Config rest server configuration
type Config struct {
	Hostname          string
	Port              int
	ReadHeaderTimeout int64
	ReadTimeout       int64
	WriteTimeout      int64
	MaxHeaderBytes    int
	MaxBytes          int64
}
