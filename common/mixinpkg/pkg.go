package mixinpkg

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"

	"github.com/fox-one/mixin-sdk-go"
	"github.com/shopspring/decimal"
)

const maxRetry = 3

func SendTransfer(ctx context.Context, opponentId string, memo string, assetId string, amount decimal.Decimal) {
	for i := 0; i < maxRetry; i++ {
		if _, err := MixinClient.Transfer(ctx, &mixin.TransferInput{
			TraceID:    mixin.RandomTraceID(),
			OpponentID: opponentId,
			AssetID:    assetId,
			Amount:     amount,
			Memo:       memo,
		}, Config.Mixin.Pin); err == nil {
			return
		}
	}
}

func SendMessage(ctx context.Context, receiptID string, text string) {
	me, err := MixinClient.UserMe(ctx)
	if err != nil {
		log.Fatalln(err)
	}

	if me.App == nil {
		log.Fatalln("use a bot keystore instead")
	}

	sessions, err := MixinClient.FetchSessions(ctx, []string{receiptID})
	if err != nil {
		log.Fatalln(err)
	}

	_ = sessions

	req := &mixin.MessageRequest{
		ConversationID: mixin.UniqueConversationID(MixinClient.ClientID, receiptID),
		RecipientID:    receiptID,
		MessageID:      mixin.RandomTraceID(),
		Category:       mixin.MessageCategoryPlainText,
		Data:           base64.StdEncoding.EncodeToString([]byte(text)),
	}
	if err := MixinClient.EncryptMessageRequest(req, sessions); err != nil {
		log.Fatalln(err)
	}

	_, err = MixinClient.SendEncryptedMessages(ctx, []*mixin.MessageRequest{req})
	if err != nil {
		log.Fatalln(err)
	}
}

func SendCards(ctx context.Context, title, intro, iconUrl, tid string, receiptId string) {
	sessions, err := MixinClient.FetchSessions(ctx, []string{receiptId})
	if err != nil {
		log.Fatalln(err)
	}

	_ = sessions

	card := &mixin.AppCardMessage{
		AppID:       MixinClient.ClientID,
		IconURL:     iconUrl,
		Title:       title,
		Description: intro,
		Action:      fmt.Sprintf("https://betxin.one/main/topic/%s", tid),
		Shareable:   true,
	}

	cardJson, _ := json.Marshal(card)

	r := &mixin.MessageRequest{
		ConversationID: mixin.UniqueConversationID(MixinClient.ClientID, receiptId),
		RecipientID:    receiptId,
		MessageID:      mixin.RandomTraceID(),
		Category:       mixin.MessageCategoryAppCard,
		Data:           base64.RawStdEncoding.EncodeToString([]byte(cardJson)),
	}

	err = MixinClient.SendMessage(ctx, r)
	if err != nil {
		fmt.Println(err)
	}
}

type Memo struct {
	Tid    string `json:"tid"`
	Select int    `json:"select" comment:"0yes 1no"`
}

func SendButtonGroup(ctx context.Context, tid, AssetId string, receiptId string) {
	sessions, err := MixinClient.FetchSessions(ctx, []string{receiptId})
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

	yesAction := fmt.Sprintf("mixin://pay?recipient=%s&asset=%s&amount=%d&memo=%s&trace=%s", MixinClient.ClientID, AssetId, 100, base64.StdEncoding.EncodeToString(yesMemoByte))
	noAction := fmt.Sprintf("mixin://pay?recipient=%s&asset=%s&amount=%d&memo=%s&trace=%s", MixinClient.ClientID, AssetId, 100, base64.StdEncoding.EncodeToString(noMemoByte))
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
		ConversationID: mixin.UniqueConversationID(MixinClient.ClientID, receiptId),
		RecipientID:    receiptId,
		MessageID:      mixin.RandomTraceID(),
		Category:       mixin.MessageCategoryAppButtonGroup,
		Data:           base64.RawStdEncoding.EncodeToString([]byte(btnGrpMsgJson)),
	}

	err = MixinClient.SendMessage(ctx, r)
	if err != nil {
		fmt.Println(err)
	}
}
