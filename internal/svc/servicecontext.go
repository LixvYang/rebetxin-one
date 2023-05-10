package svc

import (
	"context"
	"time"

	"github.com/lixvyang/rebetxin-one/internal/config"
	"github.com/lixvyang/rebetxin-one/internal/middleware"
	"github.com/lixvyang/rebetxin-one/model"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/rest"
)

const (
	Time_LAYOUT = "2006-01-02 15:04:05"
)

type TopicCollectMap struct {
	QueryTime time.Time
	// uid
	TopicCollectMap map[string]map[string]bool
}

type ServiceContext struct {
	Config             config.Config
	Admin              rest.Middleware
	CategoryModel      model.CategoryModel
	RefundModel        model.RefundModel
	TopicModel         model.TopicModel
	TopicCollectModel  model.TopicCollectModel
	TopicPurchaseModel model.TopicpurchaseModel
	UserModel          model.UserModel
	CategoryMap        map[int64]*model.Category

	TopicCollectMap *TopicCollectMap
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMysql(c.Mysql.DNS)
	svc := new(ServiceContext)
	svc.Config = c
	svc.Admin = middleware.NewAdminMiddleware().Handle
	svc.CategoryModel = model.NewCategoryModel(conn, c.CacheRedis)
	svc.RefundModel = model.NewRefundModel(conn, c.CacheRedis)
	svc.TopicModel = model.NewTopicModel(conn, c.CacheRedis)
	svc.TopicCollectModel = model.NewTopicCollectModel(conn, c.CacheRedis)
	svc.TopicPurchaseModel = model.NewTopicpurchaseModel(conn, c.CacheRedis)
	svc.UserModel = model.NewUserModel(conn, c.CacheRedis)

	categoryMap := make(map[int64]*model.Category)
	categoryList, err := svc.CategoryModel.List(context.Background())
	if err != nil {
		logx.Errorw("categoryRPC.ListCategory", logx.LogField{Key: "Error: ", Value: err})
		panic(err)
	}
	for _, cate := range categoryList {
		categoryMap[cate.Id] = &model.Category{
			Id:           cate.Id,
			CategoryName: cate.CategoryName,
		}
	}

	topicCollectMap := &TopicCollectMap{
		QueryTime:       time.Now().AddDate(0, 0, -1),
		TopicCollectMap: make(map[string]map[string]bool),
	}
	svc.TopicCollectMap = topicCollectMap

	return svc
}

func (l *ServiceContext) StringToTime(s string) time.Time {
	timeObj, err := time.Parse(Time_LAYOUT, s)
	if err != nil {
		logx.Errorw("time.Parse(Time_LAYOUT, s)", logx.LogField{Key: "Error: ", Value: err.Error()})
		return time.Time{}
	}
	return timeObj
}

func (l *ServiceContext) TimeToString(t time.Time) string {
	return t.Format(Time_LAYOUT)
}
