package events

type InconsistentEvent struct {
	ID        int64
	Direction string
	Type      string
}

var (
	InconsistentEventTypeNotEqual   = "not_equal"   // 不相等
	InconsistentEventTypeTargetMiss = "target_miss" // 目标缺失
	InconsistentEventTypeBaseMiss   = "base_miss"   // 源缺失
)
