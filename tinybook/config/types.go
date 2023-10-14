package config

type config struct {
	DB    DBConfig
	Redis RedisConfig
}

type DBConfig struct {
	Host string
}

type RedisConfig struct {
	Host string
}
