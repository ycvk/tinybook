//go:build k8s

package config

var Config = config{
	DB: DBConfig{
		Host: "root:root@tcp(tinybook-mysql:3308)/tinybook?charset=utf8mb4&parseTime=True&loc=Local",
	},
	Redis: RedisConfig{
		Host: "tinybook-redis:6380",
	},
}
