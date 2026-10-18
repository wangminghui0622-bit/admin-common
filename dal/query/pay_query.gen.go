package query

import (
	"gorm.io/gen"
	"gorm.io/gorm"
)

// PayQuery pay 模块专用查询集合。
// 该文件为 admin-common 新增的 pay 查询文件，避免改动共享的 gen.go。
type PayQuery struct {
	PayDemoOrder      payDemoOrder
	PayOrder          payOrder
	PayOrderExtension payOrderExtension
	PayRefund         payRefund
}

// UsePay 构造 pay 模块查询对象。
func UsePay(db *gorm.DB, opts ...gen.DOOption) *PayQuery {
	return &PayQuery{
		PayDemoOrder:      newPayDemoOrder(db, opts...),
		PayOrder:          newPayOrder(db, opts...),
		PayOrderExtension: newPayOrderExtension(db, opts...),
		PayRefund:         newPayRefund(db, opts...),
	}
}
