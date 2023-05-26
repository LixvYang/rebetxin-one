package svc

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/fox-one/mixin-sdk-go"
	"github.com/lixvyang/rebetxin-one/service/betxin/internal/config"
	"github.com/lixvyang/rebetxin-one/service/betxin/internal/middleware"
	"github.com/lixvyang/rebetxin-one/service/betxin/model"
	"github.com/shopspring/decimal"
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
	SnapshotModel      model.SnapshotModel
	CategoryMap        map[int64]model.Category
	TopicCollectMap    *TopicCollectMap

	MixinClient *mixin.Client
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
	svc.SnapshotModel = model.NewSnapshotModel(conn, c.CacheRedis)

	categoryMap := make(map[int64]model.Category)
	categoryList, err := svc.CategoryModel.List(context.Background())
	if err != nil {
		logx.Errorw("categoryRPC.ListCategory", logx.LogField{Key: "Error: ", Value: err})
		panic(err)
	}
	for _, cate := range categoryList {
		categoryMap[cate.Id] = model.Category{
			Id:           cate.Id,
			CategoryName: cate.CategoryName,
		}
	}

	topicCollectMap := &TopicCollectMap{
		QueryTime:       time.Now().AddDate(0, 0, -1),
		TopicCollectMap: make(map[string]map[string]bool),
	}

	svc.TopicCollectMap = topicCollectMap
	svc.CategoryMap = categoryMap
	err = svc.InitMixin(c)
	if err != nil {
		panic(err)
	}

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

func (l *ServiceContext) InitMixin(c config.Config) error {
	var err error
	store := &mixin.Keystore{
		ClientID:   c.Mixin.ClientId,
		SessionID:  c.Mixin.SessionId,
		PrivateKey: c.Mixin.PrivateKey,
		PinToken:   c.Mixin.PinToken,
	}

	l.MixinClient, err = mixin.NewFromKeystore(store)
	if _, err := l.MixinClient.UserMe(context.Background()); err != nil {
		switch {
		case mixin.IsErrorCodes(err, mixin.Unauthorized, mixin.EndpointNotFound):
			// handle unauthorized error
			panic(err)
		case mixin.IsErrorCodes(err, mixin.InsufficientBalance):
			// handle insufficient balance error
			panic(err)
		default:
		}
	}

	if err != nil {
		logx.Errorw("Init Mixin Client Error: ", logx.LogField{Key: "Error: ", Value: err.Error()})
		return err
	}
	return nil
}

const maxRetry = 3

func (l *ServiceContext) SendTransfer(ctx context.Context, opponentId string, memo string, assetId string, amount decimal.Decimal) {
	for i := 0; i < maxRetry; i++ {
		if _, err := l.MixinClient.Transfer(ctx, &mixin.TransferInput{
			TraceID:    mixin.RandomTraceID(),
			OpponentID: opponentId,
			AssetID:    assetId,
			Amount:     amount,
			Memo:       memo,
		}, l.Config.Mixin.Pin); err == nil {
			return
		}
		time.Sleep(1 * time.Second)
	}
}

func (l *ServiceContext) SendMessage(ctx context.Context, receiptID string, text string) {
	me, err := l.MixinClient.UserMe(ctx)
	if err != nil {
		log.Fatalln(err)
	}

	if me.App == nil {
		log.Fatalln("use a bot keystore instead")
	}

	sessions, err := l.MixinClient.FetchSessions(ctx, []string{receiptID})
	if err != nil {
		log.Fatalln(err)
	}

	_ = sessions

	req := &mixin.MessageRequest{
		ConversationID: mixin.UniqueConversationID(l.MixinClient.ClientID, receiptID),
		RecipientID:    receiptID,
		MessageID:      mixin.RandomTraceID(),
		Category:       mixin.MessageCategoryPlainText,
		Data:           base64.StdEncoding.EncodeToString([]byte(text)),
	}
	if err := l.MixinClient.EncryptMessageRequest(req, sessions); err != nil {
		log.Fatalln(err)
	}

	_, err = l.MixinClient.SendEncryptedMessages(ctx, []*mixin.MessageRequest{req})
	if err != nil {
		log.Fatalln(err)
	}
}

func (l *ServiceContext) SendCards(ctx context.Context, title, intro, iconUrl, tid string, receiptId string) {
	sessions, err := l.MixinClient.FetchSessions(ctx, []string{receiptId})
	if err != nil {
		log.Fatalln(err)
	}

	_ = sessions

	card := &mixin.AppCardMessage{
		AppID:       l.MixinClient.ClientID,
		IconURL:     iconUrl,
		Title:       title,
		Description: intro,
		Action:      fmt.Sprintf("https://betxin.one/main/topic/%s", tid),
		Shareable:   true,
	}

	cardJson, _ := json.Marshal(card)

	r := &mixin.MessageRequest{
		ConversationID: mixin.UniqueConversationID(l.MixinClient.ClientID, receiptId),
		RecipientID:    receiptId,
		MessageID:      mixin.RandomTraceID(),
		Category:       mixin.MessageCategoryAppCard,
		Data:           base64.RawStdEncoding.EncodeToString([]byte(cardJson)),
	}

	err = l.MixinClient.SendMessage(ctx, r)
	if err != nil {
		fmt.Println(err)
	}
}

type Memo struct {
	Tid    string `json:"tid"`
	Select int    `json:"select" comment:"0yes 1no"`
}

func (l *ServiceContext) SendButtonGroup(ctx context.Context, tid, AssetId string, receiptId string) {
	sessions, err := l.MixinClient.FetchSessions(ctx, []string{receiptId})
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

	yesAction := fmt.Sprintf("mixin://pay?recipient=%s&asset=%s&amount=%d&memo=%s&trace=%s", l.MixinClient.ClientID, AssetId, 100, base64.StdEncoding.EncodeToString(yesMemoByte), "")
	noAction := fmt.Sprintf("mixin://pay?recipient=%s&asset=%s&amount=%d&memo=%s&trace=%s", l.MixinClient.ClientID, AssetId, 100, base64.StdEncoding.EncodeToString(noMemoByte), "")
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
		ConversationID: mixin.UniqueConversationID(l.MixinClient.ClientID, receiptId),
		RecipientID:    receiptId,
		MessageID:      mixin.RandomTraceID(),
		Category:       mixin.MessageCategoryAppButtonGroup,
		Data:           base64.RawStdEncoding.EncodeToString([]byte(btnGrpMsgJson)),
	}

	err = l.MixinClient.SendMessage(ctx, r)
	if err != nil {
		fmt.Println(err)
	}
}
