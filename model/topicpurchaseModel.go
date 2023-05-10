package model

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ TopicpurchaseModel = (*customTopicpurchaseModel)(nil)

type (
	// TopicpurchaseModel is an interface to be customized, add more methods here,
	// and implement the added methods in customTopicpurchaseModel.
	TopicpurchaseModel interface {
		topicpurchaseModel
	}

	customTopicpurchaseModel struct {
		*defaultTopicpurchaseModel
	}
)

// NewTopicpurchaseModel returns a model for the database table.
func NewTopicpurchaseModel(conn sqlx.SqlConn, c cache.CacheConf) TopicpurchaseModel {
	return &customTopicpurchaseModel{
		defaultTopicpurchaseModel: newTopicpurchaseModel(conn, c),
	}
}
