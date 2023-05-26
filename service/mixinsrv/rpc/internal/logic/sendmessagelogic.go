package logic

import (
	"context"
	"encoding/base64"
	"log"

	"github.com/fox-one/mixin-sdk-go"
	"github.com/lixvyang/rebetxin-one/service/mixinsrv/rpc/internal/svc"
	"github.com/lixvyang/rebetxin-one/service/mixinsrv/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type SendMessageLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSendMessageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendMessageLogic {
	return &SendMessageLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SendMessageLogic) SendMessage(in *pb.SendMessageReq) (*pb.SendMessageResp, error) {
	l.Message(l.ctx, in.GetReceiptId(), in.GetContent())

	return &pb.SendMessageResp{}, nil
}

func (l *SendMessageLogic) Message(ctx context.Context, receiptID string, text string) {
	me, err := l.svcCtx.MixinClient.UserMe(ctx)
	if err != nil {
		log.Fatalln(err)
	}

	if me.App == nil {
		log.Fatalln("use a bot keystore instead")
	}

	sessions, err := l.svcCtx.MixinClient.FetchSessions(ctx, []string{receiptID})
	if err != nil {
		log.Fatalln(err)
	}

	_ = sessions

	req := &mixin.MessageRequest{
		ConversationID: mixin.UniqueConversationID(l.svcCtx.MixinClient.ClientID, receiptID),
		RecipientID:    receiptID,
		MessageID:      mixin.RandomTraceID(),
		Category:       mixin.MessageCategoryPlainText,
		Data:           base64.StdEncoding.EncodeToString([]byte(text)),
	}
	if err := l.svcCtx.MixinClient.EncryptMessageRequest(req, sessions); err != nil {
		log.Fatalln(err)
	}

	_, err = l.svcCtx.MixinClient.SendEncryptedMessages(ctx, []*mixin.MessageRequest{req})
	if err != nil {
		log.Fatalln(err)
	}
}
