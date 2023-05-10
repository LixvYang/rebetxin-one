package model

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ TopicModel = (*customTopicModel)(nil)

type (
	// TopicModel is an interface to be customized, add more methods here,
	// and implement the added methods in customTopicModel.
	TopicModel interface {
		topicModel
	}

	customTopicModel struct {
		*defaultTopicModel
	}
)

// NewTopicModel returns a model for the database table.
func NewTopicModel(conn sqlx.SqlConn, c cache.CacheConf) TopicModel {
	return &customTopicModel{
		defaultTopicModel: newTopicModel(conn, c),
	}
}
