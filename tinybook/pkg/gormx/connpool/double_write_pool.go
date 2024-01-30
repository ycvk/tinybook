package connpool

import (
	"context"
	"database/sql"
	"errors"
	"go.uber.org/atomic"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

const (
	PatternSrcOnly  = "src_only"
	PatternSrcFirst = "src_first"
	PatternDstFirst = "dst_first"
	PatternDstOnly  = "dst_only"
)

type DoubleWritePool struct {
	src     gorm.ConnPool
	dst     gorm.ConnPool
	pattern *atomic.String // 模式

	log *zap.Logger
}

func NewDoubleWritePool(src *gorm.DB, dst *gorm.DB, pattern string, log *zap.Logger) *DoubleWritePool {
	return &DoubleWritePool{
		src:     src.ConnPool,
		dst:     dst.ConnPool,
		pattern: atomic.NewString(pattern),
		log:     log,
	}

}

func (d *DoubleWritePool) BeginTx(ctx context.Context, opts *sql.TxOptions) (gorm.ConnPool, error) {
	pattern := d.pattern.Load()
	switch pattern {
	case PatternSrcOnly:
		tx, err := d.src.(gorm.TxBeginner).BeginTx(ctx, opts)
		return &DoubleWriteTx{
			src:     tx,
			pattern: pattern,
			log:     d.log,
		}, err
	case PatternSrcFirst:
		tx, err := d.src.(gorm.TxBeginner).BeginTx(ctx, opts)
		if err != nil {
			return nil, err
		}
		tx2, err := d.dst.(gorm.TxBeginner).BeginTx(ctx, opts) // src 成功后再写 dst, dst 成功与否无所谓
		if err != nil {
			d.log.Error("double write src first dst begin tx failed", zap.Error(err))
		}
		return &DoubleWriteTx{
			src:     tx,
			dst:     tx2,
			pattern: pattern,
			log:     d.log,
		}, nil
	case PatternDstFirst:
		tx, err := d.dst.(gorm.TxBeginner).BeginTx(ctx, opts)
		if err != nil {
			return nil, err
		}
		tx2, err := d.src.(gorm.TxBeginner).BeginTx(ctx, opts) // dst 成功后再写 src, src 成功与否无所谓
		if err != nil {
			d.log.Error("double write dst first src begin tx failed", zap.Error(err))
		}
		return &DoubleWriteTx{
			src:     tx2,
			dst:     tx,
			pattern: pattern,
			log:     d.log,
		}, nil
	case PatternDstOnly:
		tx, err := d.dst.(gorm.TxBeginner).BeginTx(ctx, opts)
		return &DoubleWriteTx{
			dst:     tx,
			pattern: pattern,
			log:     d.log,
		}, err
	default:
		return nil, errors.New("unknown pattern") // 未知模式
	}
}

func (d *DoubleWritePool) PrepareContext(ctx context.Context, query string) (*sql.Stmt, error) {
	// 这个方法没办法实现 因为无法返回两个 stmt
	panic("double write not support prepare")
}

func (d *DoubleWritePool) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	switch d.pattern.Load() {
	case PatternSrcOnly:
		return d.src.ExecContext(ctx, query, args...)
	case PatternSrcFirst:
		res, err := d.src.ExecContext(ctx, query, args...) // 先写 src
		if err == nil {
			_, err := d.dst.ExecContext(ctx, query, args...) // 再写 dst
			if err != nil {
				d.log.Error("double write src first dst exec failed", zap.Error(err), zap.String("query", query))
			}
		}
		return res, err
	case PatternDstFirst:
		res, err := d.dst.ExecContext(ctx, query, args...) // 先写 dst
		if err == nil {
			_, err := d.src.ExecContext(ctx, query, args...) // 再写 src
			if err != nil {
				d.log.Error("double write dst first src exec failed", zap.Error(err), zap.String("query", query))
			}
		}
		return res, err
	case PatternDstOnly:
		return d.dst.ExecContext(ctx, query, args...)
	default:
		return nil, errors.New("unknown pattern") // 未知模式
	}
}

func (d *DoubleWritePool) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	switch d.pattern.Load() {
	case PatternSrcOnly, PatternSrcFirst:
		return d.src.QueryContext(ctx, query, args...)
	case PatternDstFirst, PatternDstOnly:
		return d.dst.QueryContext(ctx, query, args...)
	default:
		return nil, errors.New("unknown pattern") // 未知模式
	}
}

func (d *DoubleWritePool) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	switch d.pattern.Load() {
	case PatternSrcOnly, PatternSrcFirst:
		return d.src.QueryRowContext(ctx, query, args...)
	case PatternDstFirst, PatternDstOnly:
		return d.dst.QueryRowContext(ctx, query, args...)
	default:
		panic("unknown pattern") // 未知模式
	}
}

type DoubleWriteTx struct {
	src gorm.Tx
	dst gorm.Tx

	pattern string // 模式
	log     *zap.Logger
}

func (d *DoubleWriteTx) Commit() error { //实现commit方法 事务提交
	switch d.pattern {
	case PatternSrcOnly:
		return d.src.Commit()
	case PatternSrcFirst:
		err := d.src.Commit() // 先提交 src
		if err == nil {
			if d.dst != nil { // dst 为 nil 时, 表明dst开启事务失败, 不需要再提交 dst
				err := d.dst.Commit() // 再提交 dst
				if err != nil {
					d.log.Error("double write tx src first dst commit failed", zap.Error(err))
				}
			}
		}
		return err
	case PatternDstFirst:
		err := d.dst.Commit() // 先提交 dst
		if err == nil {
			if d.src != nil { // src 为 nil 时, 表明src开启事务失败, 不需要再提交 src
				err := d.src.Commit() // 再提交 src
				if err != nil {
					d.log.Error("double write tx dst first src commit failed", zap.Error(err))
				}
			}
		}
		return err
	case PatternDstOnly:
		return d.dst.Commit()
	default:
		return errors.New("unknown pattern") // 未知模式
	}
}

func (d *DoubleWriteTx) Rollback() error {
	switch d.pattern {
	case PatternSrcOnly:
		return d.src.Rollback()
	case PatternSrcFirst:
		err := d.src.Rollback() // 先提交 src
		if err == nil {
			if d.dst != nil { // dst 为 nil 时, 表明dst开启事务失败, 不需要再回滚 dst
				err := d.dst.Rollback() // 再提交 dst
				if err != nil {
					d.log.Error("double write tx src first dst commit failed", zap.Error(err))
				}
			}
		}
		return err
	case PatternDstFirst:
		err := d.dst.Rollback() // 先提交 dst
		if err == nil {
			if d.src != nil { // src 为 nil 时, 表明src开启事务失败, 不需要再回滚 src
				err := d.src.Rollback() // 再提交 src
				if err != nil {
					d.log.Error("double write tx dst first src commit failed", zap.Error(err))
				}
			}
		}
		return err
	case PatternDstOnly:
		return d.dst.Rollback()
	default:
		return errors.New("unknown pattern") // 未知模式
	}
}

func (d *DoubleWriteTx) PrepareContext(ctx context.Context, query string) (*sql.Stmt, error) {
	// 这个方法没办法实现 因为无法返回两个 stmt
	panic("double write not support prepare")
}

func (d *DoubleWriteTx) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	switch d.pattern {
	case PatternSrcOnly:
		return d.src.ExecContext(ctx, query, args...)
	case PatternSrcFirst:
		res, err := d.src.ExecContext(ctx, query, args...) // 先写 src
		if err == nil && d.dst != nil {                    // dst 为 nil 时, 表明dst开启事务失败, 不需要再写 dst
			_, err := d.dst.ExecContext(ctx, query, args...) // 再写 dst
			if err != nil {
				d.log.Error("double write src first dst exec failed", zap.Error(err), zap.String("query", query))
			}
		}
		return res, err
	case PatternDstFirst:
		res, err := d.dst.ExecContext(ctx, query, args...) // 先写 dst
		if err == nil && d.src != nil {                    // src 为 nil 时, 表明src开启事务失败, 不需要再写 src
			_, err := d.src.ExecContext(ctx, query, args...) // 再写 src
			if err != nil {
				d.log.Error("double write dst first src exec failed", zap.Error(err), zap.String("query", query))
			}
		}
		return res, err
	case PatternDstOnly:
		return d.dst.ExecContext(ctx, query, args...)
	default:
		return nil, errors.New("unknown pattern") // 未知模式
	}
}

func (d *DoubleWriteTx) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	switch d.pattern {
	case PatternSrcOnly, PatternSrcFirst:
		return d.src.QueryContext(ctx, query, args...)
	case PatternDstFirst, PatternDstOnly:
		return d.dst.QueryContext(ctx, query, args...)
	default:
		return nil, errors.New("unknown pattern") // 未知模式
	}
}

func (d *DoubleWriteTx) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	switch d.pattern {
	case PatternSrcOnly, PatternSrcFirst:
		return d.src.QueryRowContext(ctx, query, args...)
	case PatternDstFirst, PatternDstOnly:
		return d.dst.QueryRowContext(ctx, query, args...)
	default:
		panic("unknown pattern") // 未知模式
	}
}
