package config

type Config struct {
	Port   int
	Appid  string
	Secret string
}

var Conf = new(Config)
