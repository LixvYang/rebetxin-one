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

type SendBtnGroupLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

type Memo struct {
	Tid    string `json:"tid"`
	Select int    `json:"select" comment:"0yes 1no"`
}

func NewSendBtnGroupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendBtnGroupLogic {
	return &SendBtnGroupLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SendBtnGroupLogic) SendBtnGroup(in *pb.SendBtnGroupReq) (*pb.SendBtnGroupResp, error) {
	l.SendButtonGroup(l.ctx, in.GetTid(), in.GetAssetId(), in.GetReceiptId())

	return &pb.SendBtnGroupResp{}, nil
}

func (l *SendBtnGroupLogic) SendButtonGroup(ctx context.Context, tid, AssetId string, receiptId string) {
	sessions, err := l.svcCtx.MixinClient.FetchSessions(ctx, []string{receiptId})
	if err != nil {
		log.Fatalln(err)
	}

	_ = sessions

	memo := new(Memo)
	memo.Tid = tid
	memo.Select = 0
	yesMemoByte, _ := json.Marshal(memo)
	memo.Select = 1
	noMemoByte, _ := json.Marshal(memo)

	yesAction := fmt.Sprintf("mixin://pay?recipient=%s&asset=%s&amount=%d&memo=%s&trace=%s", l.svcCtx.MixinClient.ClientID, AssetId, 100, base64.StdEncoding.EncodeToString(yesMemoByte), "")
	noAction := fmt.Sprintf("mixin://pay?recipient=%s&asset=%s&amount=%d&memo=%s&trace=%s", l.svcCtx.MixinClient.ClientID, AssetId, 100, base64.StdEncoding.EncodeToString(noMemoByte), "")
	btnGrpMsg := &mixin.AppButtonGroupMessage{
		{
			Label:  "BET YES",
			Action: fmt.Sprintf(yesAction, mixin.RandomTraceID()),
			Color:  "#6BC0CE",
		},
		{
			Label:  "BET NO",
			Action: fmt.Sprintf(noAction, mixin.RandomTraceID()),
			Color:  "#C8697D",
		},
	}

	btnGrpMsgJson, _ := json.Marshal(btnGrpMsg)

	r := &mixin.MessageRequest{
		ConversationID: mixin.UniqueConversationID(l.svcCtx.MixinClient.ClientID, receiptId),
		RecipientID:    receiptId,
		MessageID:      mixin.RandomTraceID(),
		Category:       mixin.MessageCategoryAppButtonGroup,
		Data:           base64.RawStdEncoding.EncodeToString([]byte(btnGrpMsgJson)),
	}

	err = l.svcCtx.MixinClient.SendMessage(ctx, r)
	if err != nil {
		fmt.Println(err)
	}
}
