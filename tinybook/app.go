package main

import (
	"geek_homework/tinybook/internal/events"
	"geek_homework/tinybook/internal/job"
	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
)

type App struct {
	server    *gin.Engine
	consumers []events.Consumer
	cron      *cron.Cron
	scheduler *job.Scheduler
}
