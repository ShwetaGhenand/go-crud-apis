package users

type Config struct {
	Port int
	DBConfig
	Secret string
}
