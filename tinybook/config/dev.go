//go:build !k8s

package config

var Config = config{
	DB: DBConfig{
		Host: "root:root@tcp(127.0.0.1:3308)/tinybook?charset=utf8mb4&parseTime=True&loc=Local&time_zone=Asia/Shanghai",
	},
	Redis: RedisConfig{
		Host: "localhost:6379",
	},
}
