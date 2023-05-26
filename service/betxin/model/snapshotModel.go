package model

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ SnapshotModel = (*customSnapshotModel)(nil)

type (
	// SnapshotModel is an interface to be customized, add more methods here,
	// and implement the added methods in customSnapshotModel.
	SnapshotModel interface {
		snapshotModel
	}

	customSnapshotModel struct {
		*defaultSnapshotModel
	}
)

// NewSnapshotModel returns a model for the database table.
func NewSnapshotModel(conn sqlx.SqlConn, c cache.CacheConf) SnapshotModel {
	return &customSnapshotModel{
		defaultSnapshotModel: newSnapshotModel(conn, c),
	}
}
