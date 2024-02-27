package domain

type Amount struct {
	// 如果要支持国际化，那么这个是不能少的
	Currency string
	// 这里我们遵循微信的做法，就用 int64 来记录分数。
	// 那么对于不同的货币来说，这个字段的含义就不同。
	// 比如说一些货币没有分，只有整数。
	Total int64
}

type Payment struct {
	Amt Amount
	// 代表业务，业务方决定怎么生成，
	// 我们这边不管。
	BizTradeNO string
	// 订单本身的描述
	Description string

	Status PaymentStatus
	// 第三方那边返回的 ID
	TxnID string
}

type PaymentStatus uint8

func (s PaymentStatus) AsUint8() uint8 {
	return uint8(s)
}

const (
	PaymentStatusUnknown = iota
	PaymentStatusInit
	PaymentStatusSuccess
	PaymentStatusFailed
	PaymentStatusRefund
)
