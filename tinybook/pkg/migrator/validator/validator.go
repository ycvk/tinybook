package validator

import (
	"context"
	"errors"
	"github.com/samber/lo"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"gorm.io/gorm"
	"strconv"
	"time"
	"tinybook/tinybook/pkg/migrator"
	events2 "tinybook/tinybook/pkg/migrator/events"
)

type Validator[T migrator.Entity] struct {
	base   *gorm.DB
	target *gorm.DB
	log    *zap.Logger

	producer  events2.Producer
	direction string
	batchSize int
}

func (v *Validator[T]) Validate(ctx context.Context) error {
	var eg errgroup.Group
	eg.Go(func() error {
		return v.validateBaseToTarget(ctx)
	})
	eg.Go(func() error {
		return v.validateTargetToBase(ctx)
	})
	return eg.Wait()
}

// ValidateBaseToTarget 验证 base 和 target 是否一致 (base -> target)
func (v *Validator[T]) validateBaseToTarget(ctx context.Context) error {
	offset := -1
	for {
		offset++
		var src T
		err := v.base.WithContext(ctx).Order("id").Offset(offset).Limit(1).Find(&src).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) { // 未找到记录
				return nil
			}
			v.log.Error("base -> target validate failed", zap.Error(err))
			continue // 未知错误
		}
		var dst T
		err = v.target.WithContext(ctx).Where("id=?", src.ID()).Find(&dst).Error
		switch {
		case err == nil: // 找到记录
			if !src.CompareTo(dst) {
				v.log.Error("base -> target compare failed, src: "+strconv.FormatInt(src.ID(), 10), zap.Error(err))
				v.notify(src.ID(), events2.InconsistentEventTypeNotEqual)
				continue
			}
		case errors.Is(err, gorm.ErrRecordNotFound): // 未找到记录
			v.log.Error("base -> target not found, src: "+strconv.FormatInt(src.ID(), 10), zap.Error(err))
			// 通知不一致事件 (base 中存在但 target 中不存在的记录)
			v.notify(src.ID(), events2.InconsistentEventTypeTargetMiss)
			continue
		default:
			v.log.Error("base -> target find failed", zap.Error(err))
			continue // 未知错误
		}
	}
}

func (v *Validator[T]) validateTargetToBase(ctx context.Context) error {
	offset := -v.batchSize
	for {
		offset += v.batchSize
		var dst []T
		err := v.target.WithContext(ctx).Select("id").Order("id").Offset(offset).Limit(v.batchSize).Find(&dst).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) { // 未找到记录
				return nil
			}
			v.log.Error("target -> base validate failed", zap.Error(err))
			continue // 未知错误
		}
		var src []T
		// 提取 dst 中的 id
		ids := lo.Map[T, int64](dst, func(item T, index int) int64 {
			return item.ID()
		})
		err = v.base.WithContext(ctx).Select("id").Where("id in ?", ids).Find(&src).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) { // 未找到记录
				// 通知不一致事件 (target 中存在但 base 中不存在的记录)
				v.notifyBaseMissing(dst)
				continue
			}
			v.log.Error("target -> base find failed", zap.Error(err))
		}
		// 寻找差集 dst - src 即 target 中存在但 base 中不存在的记录
		miss := lo.Filter[T](dst, func(item T, index int) bool {
			return !lo.Contains(ids, item.ID())
		})
		// 通知不一致事件 (target 中存在但 base 中不存在的记录)
		v.notifyBaseMissing(miss)
		// 如果 dst 的数量小于 batchSize, 说明已经到了最后一批
		if len(src) < v.batchSize {
			return nil
		}
	}
}

func (v *Validator[T]) notifyBaseMissing(ts []T) {
	for _, item := range ts {
		v.notify(item.ID(), events2.InconsistentEventTypeBaseMiss)
	}
}

// notify 通知不一致事件
func (v *Validator[T]) notify(id int64, ty string) {
	timeout, cancelFunc := context.WithTimeout(context.Background(), time.Second)
	defer cancelFunc()
	err := v.producer.ProduceInconsistentEvent(timeout, events2.InconsistentEvent{
		ID:        id,
		Type:      ty,
		Direction: v.direction,
	})
	if err != nil {
		v.log.Error("produce inconsistent event failed "+strconv.FormatInt(id, 10), zap.Error(err))
	}
}
