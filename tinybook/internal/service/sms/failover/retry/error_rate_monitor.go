package retry

import (
	"log/slog"
	"sync"
	"time"
)

type ErrorMonitor interface {
	RecordResult(success bool) // 记录结果
	CheckErrorRate() bool      // 检查错误率是否超过阈值
}

// TimeStampedResult 时间戳结果，记录每次调用的成功或失败与时间
type TimeStampedResult struct {
	success bool
	time    time.Time
}

// ErrorRateMonitor 错误率监控器
type ErrorRateMonitor struct {
	results                   []TimeStampedResult // 存储时间戳结果的环形缓冲区
	windowStart               time.Time           // 窗口开始时间
	timeWindowDuration        time.Duration       // 时间窗口大小
	errorRate                 float64             // 当前错误率
	threshold                 float64             // 当前阈值
	initialErrorRateThreshold float64             // 初始错误率阈值 比如 0.3 代表 30% 的错误率阈值
	maxErrorRateThreshold     float64             // 最大错误率阈值 比如 0.6 代表 60% 的最大错误率阈值
	rwMutex                   sync.RWMutex
}

func NewErrorRateMonitor(initialErrRateThreshold float64, maxErrRateThreshold float64, duration time.Duration) *ErrorRateMonitor {
	erm := &ErrorRateMonitor{
		windowStart:               time.Now(),
		timeWindowDuration:        duration, // 时间窗口大小，决定了监控的粒度
		threshold:                 initialErrRateThreshold,
		initialErrorRateThreshold: initialErrRateThreshold,
		maxErrorRateThreshold:     maxErrRateThreshold,
	}
	// 定时调整错误率和阈值
	go erm.adjustErrorRateAndThreshold()
	return erm
}

// adjustErrorRateAndThreshold 定期调整错误率和阈值
func (erm *ErrorRateMonitor) adjustErrorRateAndThreshold() {
	ticker := time.NewTicker(erm.timeWindowDuration) // 监控的粒度
	for {
		<-ticker.C
		slog.Info("开始调整错误率和阈值")
		erm.rwMutex.Lock()
		erm.cleanUpOldResults()  // 清理过时的结果
		erm.calculateErrorRate() // 计算当前错误率
		erm.adjustThreshold()    // 自适应调整阈值
		erm.rwMutex.Unlock()
	}
}

// cleanUpOldResults 清理过时的结果
func (erm *ErrorRateMonitor) cleanUpOldResults() {
	cutoff := time.Now().Add(-erm.timeWindowDuration) // 计算过时的时间
	newResults := make([]TimeStampedResult, 0)
	// 遍历所有结果，只保留在窗口内的结果
	for _, result := range erm.results {
		if result.time.After(cutoff) {
			newResults = append(newResults, result)
		}
	}
	erm.results = newResults
}

// calculateErrorRate 计算当前错误率
func (erm *ErrorRateMonitor) calculateErrorRate() {
	var failures int
	// 遍历所有结果，计算错误率
	for _, result := range erm.results {
		if !result.success {
			failures++
		}
	}
	total := len(erm.results)
	if total > 0 {
		erm.errorRate = float64(failures) / float64(total)
		slog.Info("错误率", "errorRate", erm.errorRate)
	}
}

// adjustThreshold 自适应调整阈值
func (erm *ErrorRateMonitor) adjustThreshold() {
	// 这只是一个示例，根据结果数量动态调整阈值，以后可以根据更复杂的逻辑来调整
	erm.threshold = erm.initialErrorRateThreshold + float64(len(erm.results))/1000.0
	if erm.threshold > erm.maxErrorRateThreshold {
		erm.threshold = erm.maxErrorRateThreshold
	}
}

// RecordResult 异步记录结果
func (erm *ErrorRateMonitor) RecordResult(success bool) {
	erm.rwMutex.Lock() // 写锁
	defer erm.rwMutex.Unlock()
	erm.results = append(erm.results, TimeStampedResult{success, time.Now()})
}

// CheckErrorRate 检查当前错误率是否超过阈值
func (erm *ErrorRateMonitor) CheckErrorRate() bool {
	erm.rwMutex.RLock() // 读锁
	defer erm.rwMutex.RUnlock()
	return erm.errorRate > erm.threshold
}
