package logic

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"

	"github.com/fox-one/mixin-sdk-go"
	"github.com/lixvyang/rebetxin-one/service/mixinsrv/rpc/internal/svc"
	"github.com/lixvyang/rebetxin-one/service/mixinsrv/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type SendCardLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSendCardLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendCardLogic {
	return &SendCardLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SendCardLogic) SendCard(in *pb.SendCardReq) (*pb.SendCardResp, error) {
	l.SendCards(l.ctx, in.GetTitle(), in.GetIntro(), in.GetIconUrl(), in.GetTid(), in.GetReceiptId())

	return &pb.SendCardResp{}, nil
}

func (l *SendCardLogic) SendCards(ctx context.Context, title, intro, iconUrl, tid string, receiptId string) {
	sessions, err := l.svcCtx.MixinClient.FetchSessions(ctx, []string{receiptId})
	if err != nil {
		log.Fatalln(err)
	}

	_ = sessions

	card := &mixin.AppCardMessage{
		AppID:       l.svcCtx.MixinClient.ClientID,
		IconURL:     iconUrl,
		Title:       title,
		Description: intro,
		Action:      fmt.Sprintf("https://betxin.one/main/topic/%s", tid),
		Shareable:   true,
	}

	cardJson, _ := json.Marshal(card)

	r := &mixin.MessageRequest{
		ConversationID: mixin.UniqueConversationID(l.svcCtx.MixinClient.ClientID, receiptId),
		RecipientID:    receiptId,
		MessageID:      mixin.RandomTraceID(),
		Category:       mixin.MessageCategoryAppCard,
		Data:           base64.RawStdEncoding.EncodeToString([]byte(cardJson)),
	}

	err = l.svcCtx.MixinClient.SendMessage(ctx, r)
	if err != nil {
		fmt.Println(err)
	}
}
