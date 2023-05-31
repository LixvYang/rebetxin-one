package oauth

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/golang-jwt/jwt"
	"github.com/lixvyang/rebetxin-one/common/errorx"
	"github.com/lixvyang/rebetxin-one/service/betxin/internal/svc"
	"github.com/lixvyang/rebetxin-one/service/betxin/internal/types"
	"github.com/lixvyang/rebetxin-one/service/betxin/model"
	"github.com/pandodao/passport-go/auth"
	"github.com/pandodao/passport-go/eip4361"
	"github.com/pandodao/passport-go/mvm"

	"github.com/zeromicro/go-zero/core/logx"
)

type SigninLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSigninLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SigninLogic {
	return &SigninLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SigninLogic) Signin(req *types.SigninReq) (resp *types.SigninResp, err error) {
	if req.LoginMethod != "mixin_token" && req.LoginMethod != "mvm" {
		logx.Errorw("req.LoginMethod: ", logx.LogField{Key: "Error: ", Value: err})
		return nil, errorx.NewDefaultError("Login Type Error!")
	}

	authorizer := auth.New([]string{
		"30aad5a5-e5f3-4824-9409-c2ff4152724e",
	}, []string{
		"localhost:4000",
		"localhost:4000/*",
	})

	switch req.LoginMethod {
	case "mixin_token":
		// 2. 访问用户信息
		userInfo, err := authorizer.Authorize(l.ctx, &auth.AuthorizationParams{
			Method:     auth.AuthMethodMixinToken,
			MixinToken: req.MixinToken,
		})
		if err != nil {
			fmt.Println("333")
			return nil, errorx.NewDefaultError("authorizer.Authorize error!")
		}
		log.Println("444")

		user := model.User{
			AvatarUrl:      sql.NullString{String: userInfo.AvatarURL, Valid: true},
			IdentityNumber: userInfo.IdentityNumber,
			SessionId:      sql.NullString{String: userInfo.SessionID, Valid: true},
			Uid:            userInfo.UserID,
			FullName:       sql.NullString{String: userInfo.FullName, Valid: true},
			Biography:      sql.NullString{String: userInfo.Biography, Valid: true},
		}
		fmt.Println(222)
		// 3. 查询用户是否已经存在
		// 3.1 若不存在则创建
		// 3.2 若存在则更新数据
		// 创建
		fmt.Println(333)
		_, err = l.svcCtx.UserModel.FindOneByUid(l.ctx, user.Uid)
		if err != nil {
			if err == model.ErrNotFound {
				// 创建
				_, err = l.svcCtx.UserModel.Insert(l.ctx, &user)
				if err != nil {
					logx.Errorw("UserModel.Insert failed", logx.LogField{Key: "err", Value: err.Error()})
					return nil, errorx.NewDefaultError("UserModel.Insert(l.ctx, &user) error")
				}
			}
		} else {
			// 已经存在
			oUser, err := l.svcCtx.UserModel.FindOneByUid(l.ctx, user.Uid)
			if err != nil {
				logx.Errorw("UserModel.FindOneByUid", logx.LogField{Key: "err", Value: err.Error()})
			}

			user.Id = oUser.Id
			err = l.svcCtx.UserModel.Update(l.ctx, &user)
			if err != nil {
				logx.Errorw("UserModel.Update err: ", logx.LogField{Key: "err", Value: err.Error()})
				return nil, err
			}
		}

		jwtToken, err := l.getJwtToken(l.svcCtx.Config.Auth.AccessSecret, time.Now().Unix(), l.svcCtx.Config.Auth.AccessExpire, user.Uid)
		if err != nil {
			logx.Errorw("l.getJwtToken failed", logx.LogField{Key: "err", Value: err.Error()})
			return nil, err
		}

		return &types.SigninResp{
			Token: jwtToken,
		}, nil

	case "mvm":
		message, err := eip4361.Parse(req.SignedMsg)
		if err != nil {
			return nil, err
		}

		if err := message.Validate(time.Now()); err != nil {
			return nil, err
		}

		if err := eip4361.Verify(message, req.Sign); err != nil {
			return nil, err
		}

		// get the public key from the message, and use it to login
		jwtToken, err := l.LoginMvm(l.ctx, message.Address)
		if err != nil {
			return nil, err
		}

		return &types.SigninResp{
			Token: jwtToken,
		}, nil

	default:
		return nil, errorx.NewDefaultError("Error")
	}
	return
}

func (l *SigninLogic) getJwtToken(secretKey string, iat int64, seconds int64, uid string) (string, error) {
	claims := make(jwt.MapClaims)
	claims["exp"] = iat + seconds
	claims["iat"] = iat
	claims["uid"] = uid
	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims = claims
	return token.SignedString([]byte(secretKey))
}

func (l *SigninLogic) LoginMvm(ctx context.Context, pubkey string) (string, error) {
	addr := common.HexToAddress(pubkey)
	mvmUser, err := mvm.GetBridgeUser(ctx, addr)
	if err != nil {
		return "", err
	}

	contractAddr, err := mvm.GetUserContract(ctx, mvmUser.UserID)
	if err != nil {
		fmt.Printf("err mvm.GetUserContract: %v\n", err)
		return "", err
	}

	// if contractAddr is not 0x000..00, it means the user has already registered a mvm account
	emptyAddr := common.Address{}
	if contractAddr != emptyAddr {
		return "", err
	}

	user := model.User{
		IsMvmUser:  1,
		Uid:        mvmUser.UserID,
		FullName:   sql.NullString{String: mvmUser.FullName, Valid: true},
		Contract:   sql.NullString{String: mvmUser.Contract, Valid: true},
		PrivateKey: sql.NullString{String: mvmUser.Key.PrivateKey, Valid: true},
		ClientId:   sql.NullString{String: mvmUser.Key.ClientID, Valid: true},
		CreatedAt:  mvmUser.CreatedAt,
		SessionId:  sql.NullString{String: mvmUser.Key.SessionID, Valid: true},
	}
	fmt.Printf("%+v", user)

	_, err = l.svcCtx.UserModel.Insert(ctx, &user)
	if err != nil {
		logx.Errorw("UserModel.Insert: ", logx.LogField{Key: "err", Value: err.Error()})
		return "", err
	}

	jwtToken, err := l.getJwtToken(l.svcCtx.Config.Auth.AccessSecret, time.Now().Unix(), l.svcCtx.Config.Auth.AccessExpire, user.Uid)
	if err != nil {
		logx.Errorw("l.getJwtToken failed", logx.LogField{Key: "err", Value: err.Error()})
		return "", err
	}

	return jwtToken, nil
}
