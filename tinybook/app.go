package main

import (
	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
	"tinybook/tinybook/internal/events"
	"tinybook/tinybook/internal/job"
)

type App struct {
	server    *gin.Engine
	consumers []events.Consumer
	cron      *cron.Cron
	scheduler *job.Scheduler
}
