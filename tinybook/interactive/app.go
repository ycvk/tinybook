package main

import (
	"tinybook/tinybook/interactive/events"
	"tinybook/tinybook/pkg/grpcx"
)

type App struct {
	consumers []events.Consumer
	server    *grpcx.Server
}
