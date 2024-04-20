package fixer

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"tinybook/tinybook/pkg/migrator"
)

type OverrideFixer[T migrator.Entity] struct {
	base   *gorm.DB
	target *gorm.DB

	columns []string
}

// NewOverrideFixerWithColumns 通过指定的 columns 创建 OverrideFixer
func NewOverrideFixerWithColumns[T migrator.Entity](base *gorm.DB, target *gorm.DB, columns []string) *OverrideFixer[T] {
	return &OverrideFixer[T]{
		base:    base,
		target:  target,
		columns: columns,
	}
}

// NewOverrideFixer 创建 OverrideFixer
func NewOverrideFixer[T migrator.Entity](base *gorm.DB, target *gorm.DB) *OverrideFixer[T] {
	rows, err := base.Model(new(T)).Order("id").Rows()
	if err != nil {
		return nil
	}
	strings, err := rows.Columns()
	if err != nil {
		return nil
	}
	return &OverrideFixer[T]{
		base:    base,
		target:  target,
		columns: strings,
	}
}

func (f *OverrideFixer[T]) Fix(ctx context.Context, id int64) error {
	var src T
	// 从 base 中获取记录
	err := f.base.WithContext(ctx).Where("id=?", id).First(&src).Error
	switch {
	// 从 base 中获取记录成功, Upsert 到 target 中
	case err == nil:
		return f.target.WithContext(ctx).Model(&src).Clauses(clause.OnConflict{
			DoUpdates: clause.AssignmentColumns(f.columns),
		}).Create(&src).Error
		// 从 base 中获取记录失败, 未找到记录, 从 target 中删除
	case errors.Is(err, gorm.ErrRecordNotFound):
		return f.target.WithContext(ctx).Model(&src).Delete("id=?", id).Error
	default:
		return err
	}
}
