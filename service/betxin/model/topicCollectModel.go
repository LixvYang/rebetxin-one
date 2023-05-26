package model

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ TopicCollectModel = (*customTopicCollectModel)(nil)

type (
	// TopicCollectModel is an interface to be customized, add more methods here,
	// and implement the added methods in customTopicCollectModel.
	TopicCollectModel interface {
		topicCollectModel
	}

	customTopicCollectModel struct {
		*defaultTopicCollectModel
	}
)

// NewTopicCollectModel returns a model for the database table.
func NewTopicCollectModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) TopicCollectModel {
	return &customTopicCollectModel{
		defaultTopicCollectModel: newTopicCollectModel(conn, c, opts...),
	}
}
