package gormx

import (
	"github.com/prometheus/client_golang/prometheus"
	"gorm.io/gorm"
	"time"
)

type Callbacks struct {
	vector *prometheus.SummaryVec
}

func (c *Callbacks) Name() string {
	return "prometheus"
}

func (c *Callbacks) Initialize(db *gorm.DB) error {
	_ = db.Callback().Create().Before("*").
		Register("prometheus_create_before", c.Before())
	_ = db.Callback().Create().After("*").
		Register("prometheus_create_after", c.After("INSERT"))

	_ = db.Callback().Update().Before("*").
		Register("prometheus_update_before", c.Before())
	_ = db.Callback().Update().After("*").
		Register("prometheus_update_after", c.After("UPDATE"))

	_ = db.Callback().Delete().Before("*").
		Register("prometheus_delete_before", c.Before())
	_ = db.Callback().Delete().After("*").
		Register("prometheus_delete_after", c.After("DELETE"))

	_ = db.Callback().Query().Before("*").
		Register("prometheus_query_before", c.Before())
	_ = db.Callback().Query().After("*").
		Register("prometheus_query_after", c.After("SELECT"))

	_ = db.Callback().Raw().Before("*").
		Register("prometheus_raw_before", c.Before())
	_ = db.Callback().Raw().After("*").
		Register("prometheus_raw_after", c.After("RAW"))

	_ = db.Callback().Row().Before("*").
		Register("prometheus_row_before", c.Before())
	_ = db.Callback().Row().After("*").
		Register("prometheus_row_after", c.After("ROW"))
	return nil
}

func NewCallbacks(opts prometheus.SummaryOpts) *Callbacks {
	vec := prometheus.NewSummaryVec(opts,
		[]string{"type", "table"})
	prometheus.MustRegister(vec)
	return &Callbacks{
		vector: vec,
	}
}

func (c *Callbacks) Before() func(db *gorm.DB) {
	return func(db *gorm.DB) {
		now := time.Now()
		db.Set("gorm:started_at", now)
	}
}

func (c *Callbacks) After(typ string) func(db *gorm.DB) {
	return func(db *gorm.DB) {
		startedAt, ok := db.Get("gorm:started_at")
		if !ok {
			return
		}
		duration := time.Since(startedAt.(time.Time)).Milliseconds()
		c.vector.WithLabelValues(typ, db.Statement.Table).Observe(float64(duration))
	}
}
