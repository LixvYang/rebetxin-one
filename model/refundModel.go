package model

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ RefundModel = (*customRefundModel)(nil)

type (
	// RefundModel is an interface to be customized, add more methods here,
	// and implement the added methods in customRefundModel.
	RefundModel interface {
		refundModel
	}

	customRefundModel struct {
		*defaultRefundModel
	}
)

// NewRefundModel returns a model for the database table.
func NewRefundModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) RefundModel {
	return &customRefundModel{
		defaultRefundModel: newRefundModel(conn, c, opts...),
	}
}
