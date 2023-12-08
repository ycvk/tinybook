## wrk压测命令

### signup接口

10个线程，持续10s，20个连接，每个连接每次发送1个请求

```shell
wrk -t10 -d10s -c20 -s signup.lua http://localhost/users/signup
```

### login接口

10个线程，持续3s，20个连接，每个连接每次发送1个请求

```shell
wrk -t10 -d3s -c20 -s login.lua http://localhost/users/login
```

### profile接口

10个线程，持续3s，20个连接，每个连接每次发送1个请求

```shell
wrk -t10 -d3s -c20 -s profile.lua http://localhost/users/profile
```

### rank接口

10个线程，持续3s，20个连接，每个连接每次发送1个请求

```shell
wrk -t10 -d3s -c20 -s rank.lua http://localhost/articles/rank/10
```