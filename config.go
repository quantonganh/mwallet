package mwallet

// Config represents an app config
type Config struct {
	DB struct {
		Host string
		Port int
		User string
		Password string
		Name string
	}
}