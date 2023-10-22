//go:build !k8s

package config

var Config = config{
	DB: DBConfig{
		Host: "root:root@tcp(127.0.0.1:3306)/ycvk?charset=utf8mb4&parseTime=True&loc=Local",
	},
	Redis: RedisConfig{
		Host: "localhost:6379",
	},
}
