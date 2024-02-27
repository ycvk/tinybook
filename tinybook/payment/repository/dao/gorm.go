package dao

import (
	"context"
	"gorm.io/gorm"
	"time"
	"tinybook/tinybook/payment/domain"
)

type PaymentGORMDAO struct {
	db *gorm.DB
}

func (p *PaymentGORMDAO) GetPayment(ctx context.Context, bizTradeNO string) (Payment, error) {
	var res Payment
	err := p.db.WithContext(ctx).Where("biz_trade_no = ?", bizTradeNO).First(&res).Error
	return res, err
}

func (p *PaymentGORMDAO) FindExpiredPayment(
	ctx context.Context,
	offset int, limit int, t time.Time) ([]Payment, error) {
	var res []Payment
	err := p.db.WithContext(ctx).Where("status = ? AND utime < ?",
		// 我的 IDE 有问题，AsUint8 会报错
		uint8(domain.PaymentStatusInit), t.UnixMilli()).
		Offset(offset).Limit(limit).Find(&res).Error
	return res, err
}

func (p *PaymentGORMDAO) UpdateTxnIDAndStatus(ctx context.Context,
	bizTradeNo string,
	txnID string, status domain.PaymentStatus) error {
	return p.db.WithContext(ctx).Model(&Payment{}).
		Where("biz_trade_no = ?", bizTradeNo).
		Updates(map[string]any{
			"txn_id": txnID,
			"status": status.AsUint8(),
			"utime":  time.Now().UnixMilli(),
		}).Error
}

func NewPaymentGORMDAO(db *gorm.DB) PaymentDAO {
	return &PaymentGORMDAO{db: db}
}

func (p *PaymentGORMDAO) Insert(ctx context.Context, pmt Payment) error {
	now := time.Now().UnixMilli()
	pmt.Utime = now
	pmt.Ctime = now
	return p.db.WithContext(ctx).Create(&pmt).Error
}
