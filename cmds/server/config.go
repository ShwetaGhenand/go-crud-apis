package server

type Config struct {
	Port int
	DBConfig
	Secret string
}
