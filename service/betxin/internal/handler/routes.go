// Code generated by goctl. DO NOT EDIT.
package handler

import (
	"net/http"

	category "github.com/lixvyang/rebetxin-one/service/betxin/internal/handler/category"
	collect "github.com/lixvyang/rebetxin-one/service/betxin/internal/handler/collect"
	oauth "github.com/lixvyang/rebetxin-one/service/betxin/internal/handler/oauth"
	purchase "github.com/lixvyang/rebetxin-one/service/betxin/internal/handler/purchase"
	refund "github.com/lixvyang/rebetxin-one/service/betxin/internal/handler/refund"
	snapshot "github.com/lixvyang/rebetxin-one/service/betxin/internal/handler/snapshot"
	topic "github.com/lixvyang/rebetxin-one/service/betxin/internal/handler/topic"
	topicAdmin "github.com/lixvyang/rebetxin-one/service/betxin/internal/handler/topicAdmin"
	user "github.com/lixvyang/rebetxin-one/service/betxin/internal/handler/user"
	"github.com/lixvyang/rebetxin-one/service/betxin/internal/svc"

	"github.com/zeromicro/go-zero/rest"
)

func RegisterHandlers(server *rest.Server, serverCtx *svc.ServiceContext) {
	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodGet,
				Path:    "/category",
				Handler: category.ListCategoryHandler(serverCtx),
			},
			{
				Method:  http.MethodGet,
				Path:    "/category/:id",
				Handler: category.GetCategoryHandler(serverCtx),
			},
		},
		rest.WithPrefix("/api/v1"),
	)

	server.AddRoutes(
		rest.WithMiddlewares(
			[]rest.Middleware{serverCtx.Admin},
			[]rest.Route{
				{
					Method:  http.MethodPost,
					Path:    "/category",
					Handler: category.CreateCategoryHandler(serverCtx),
				},
			}...,
		),
		rest.WithJwt(serverCtx.Config.Auth.AccessSecret),
		rest.WithPrefix("/api/v1"),
	)

	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodPost,
				Path:    "/refund",
				Handler: refund.CreateRefundHandler(serverCtx),
			},
			{
				Method:  http.MethodGet,
				Path:    "/refund/:tid",
				Handler: refund.GetRefundByTidHandler(serverCtx),
			},
			{
				Method:  http.MethodGet,
				Path:    "/refund",
				Handler: refund.ListRefundHandler(serverCtx),
			},
		},
		rest.WithJwt(serverCtx.Config.Auth.AccessSecret),
		rest.WithPrefix("/api/v1"),
	)

	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodGet,
				Path:    "/topics/:cid",
				Handler: topic.ListTopicByCidHandler(serverCtx),
			},
			{
				Method:  http.MethodGet,
				Path:    "/topic/:tid",
				Handler: topic.GetTopicByTidHandler(serverCtx),
			},
			{
				Method:  http.MethodGet,
				Path:    "/topic/search",
				Handler: topic.SearchTopicHandler(serverCtx),
			},
		},
		rest.WithPrefix("/api/v1"),
	)

	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodPatch,
				Path:    "/topic/:tid",
				Handler: topic.UpdateTopicPriceHandler(serverCtx),
			},
		},
		rest.WithJwt(serverCtx.Config.Auth.AccessSecret),
		rest.WithPrefix("/api/v1"),
	)

	server.AddRoutes(
		rest.WithMiddlewares(
			[]rest.Middleware{serverCtx.Admin},
			[]rest.Route{
				{
					Method:  http.MethodPost,
					Path:    "/topic",
					Handler: topicAdmin.CreateTopicHandler(serverCtx),
				},
				{
					Method:  http.MethodDelete,
					Path:    "/topic",
					Handler: topicAdmin.DeleteTopicHandler(serverCtx),
				},
				{
					Method:  http.MethodGet,
					Path:    "/topic",
					Handler: topicAdmin.ListTopicHandler(serverCtx),
				},
				{
					Method:  http.MethodPut,
					Path:    "/topic",
					Handler: topicAdmin.UpdateTopicHandler(serverCtx),
				},
				{
					Method:  http.MethodPost,
					Path:    "/topic/:tid/stop",
					Handler: topicAdmin.StopTopicHandler(serverCtx),
				},
			}...,
		),
		rest.WithPrefix("/api/v1"),
	)

	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodPost,
				Path:    "/purchase",
				Handler: purchase.CreateTopicPurchaseHandler(serverCtx),
			},
			{
				Method:  http.MethodGet,
				Path:    "/purchase/uid/:uid/tid/:tid",
				Handler: purchase.GetTopicPurchaseHandler(serverCtx),
			},
			{
				Method:  http.MethodGet,
				Path:    "/purchase",
				Handler: purchase.ListTopicPurchaseHandler(serverCtx),
			},
		},
		rest.WithJwt(serverCtx.Config.Auth.AccessSecret),
		rest.WithPrefix("/api/v1"),
	)

	server.AddRoutes(
		rest.WithMiddlewares(
			[]rest.Middleware{serverCtx.Admin},
			[]rest.Route{
				{
					Method:  http.MethodGet,
					Path:    "/purchases",
					Handler: purchase.GetTopicPurchasesHandler(serverCtx),
				},
			}...,
		),
		rest.WithPrefix("/api/v1"),
	)

	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodGet,
				Path:    "/collect",
				Handler: collect.ListTopicCollectHandler(serverCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/collect",
				Handler: collect.CreateTopicCollectHandler(serverCtx),
			},
			{
				Method:  http.MethodDelete,
				Path:    "/collect",
				Handler: collect.DeleteTopicCollectHandler(serverCtx),
			},
		},
		rest.WithJwt(serverCtx.Config.Auth.AccessSecret),
		rest.WithPrefix("/api/v1"),
	)

	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodPost,
				Path:    "/user/signin",
				Handler: oauth.SigninHandler(serverCtx),
			},
		},
		rest.WithPrefix("/api/v1"),
	)

	server.AddRoutes(
		rest.WithMiddlewares(
			[]rest.Middleware{serverCtx.Admin},
			[]rest.Route{
				{
					Method:  http.MethodPost,
					Path:    "/user",
					Handler: user.CreateUserHandler(serverCtx),
				},
			}...,
		),
		rest.WithPrefix("/api/v1"),
	)

	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodGet,
				Path:    "/user",
				Handler: user.GetUserHandler(serverCtx),
			},
		},
		rest.WithJwt(serverCtx.Config.Auth.AccessSecret),
		rest.WithPrefix("/api/v1"),
	)

	server.AddRoutes(
		rest.WithMiddlewares(
			[]rest.Middleware{serverCtx.Admin},
			[]rest.Route{
				{
					Method:  http.MethodGet,
					Path:    "/user",
					Handler: user.GetUserListHandler(serverCtx),
				},
			}...,
		),
		rest.WithJwt(serverCtx.Config.Auth.AccessSecret),
		rest.WithPrefix("/api/v1/admin"),
	)

	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodPost,
				Path:    "/snapshot",
				Handler: snapshot.CreateSnapshotHandler(serverCtx),
			},
			{
				Method:  http.MethodGet,
				Path:    "/snapshot",
				Handler: snapshot.GetSnapshotHandler(serverCtx),
			},
		},
		rest.WithJwt(serverCtx.Config.Auth.AccessSecret),
		rest.WithPrefix("/api/v1"),
	)
}
