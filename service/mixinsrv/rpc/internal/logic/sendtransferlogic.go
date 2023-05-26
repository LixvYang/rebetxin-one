package logic

import (
	"context"
	"time"

	"github.com/fox-one/mixin-sdk-go"
	"github.com/lixvyang/rebetxin-one/service/mixinsrv/rpc/internal/svc"
	"github.com/lixvyang/rebetxin-one/service/mixinsrv/rpc/pb"
	"github.com/shopspring/decimal"

	"github.com/zeromicro/go-zero/core/logx"
)

type SendTransferLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSendTransferLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendTransferLogic {
	return &SendTransferLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// -----------------------mixinsrv-----------------------
func (l *SendTransferLogic) SendTransfer(in *pb.SendTransferReq) (*pb.SendTransferResp, error) {
	amount, _ := decimal.NewFromString(in.GetAmount())
	l.Transfer(l.ctx, in.GetOpponentId(), in.GetMemo(), in.GetAssetId(), amount)

	return &pb.SendTransferResp{}, nil
}

const maxRetry = 3

func (l *SendTransferLogic) Transfer(ctx context.Context, opponentId string, memo string, assetId string, amount decimal.Decimal) {
	for i := 0; i < maxRetry; i++ {
		if _, err := l.svcCtx.MixinClient.Transfer(ctx, &mixin.TransferInput{
			TraceID:    mixin.RandomTraceID(),
			OpponentID: opponentId,
			AssetID:    assetId,
			Amount:     amount,
			Memo:       memo,
		}, l.svcCtx.Pin); err == nil {
			return
		}
		time.Sleep(1 * time.Second)
	}
}
