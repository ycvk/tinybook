package ioc

import (
	"github.com/godruoyi/go-snowflake"
	"time"
)

func InitSnowflake() {
	snowflake.SetMachineID(1)
	snowflake.SetStartTime(time.Date(2014, 9, 1, 0, 0, 0, 0, time.UTC))
}
