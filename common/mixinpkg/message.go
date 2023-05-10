package mixinpkg

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"log"

	"github.com/fox-one/mixin-sdk-go"
)

var text = flag.String("text", "hello world", "text message")

// 连接成功
// bet Yes No
// 关注社群
func SendMessage(receiptID ...string) {
	ctx := context.Background()

	me, err := MixinClient.UserMe(ctx)
	if err != nil {
		log.Fatalln(err)
	}

	if me.App == nil {
		log.Fatalln("use a bot keystore instead")
	}

	sessions, err := MixinClient.FetchSessions(ctx, receiptID)
	if err != nil {
		log.Fatalln(err)
	}

	_ = sessions

	var reqs []*mixin.MessageRequest
	for _, rcid := range receiptID {
		req := &mixin.MessageRequest{
			ConversationID: mixin.UniqueConversationID(MixinClient.ClientID, rcid),
			RecipientID:    rcid,
			MessageID:      mixin.RandomTraceID(),
			Category:       mixin.MessageCategoryPlainText,
			Data:           base64.StdEncoding.EncodeToString([]byte(*text)),
		}
		if err := MixinClient.EncryptMessageRequest(req, sessions); err != nil {
			log.Fatalln(err)
		}
		reqs = append(reqs, req)
	}

	_, err = MixinClient.SendEncryptedMessages(ctx, reqs)
	if err != nil {
		log.Fatalln(err)
	}
}

func SendCard(ctx context.Context, tid, title, intro string, receiptIds ...string) {
	me, err := MixinClient.UserMe(ctx)
	if err != nil {
		log.Fatalln(err)
	}

	if len(intro) > 45 {
		intro = intro[:44] + "..."
	}
	sessions, err := MixinClient.FetchSessions(ctx, receiptIds)
	if err != nil {
		log.Fatalln(err)
	}

	_ = sessions

	card := &mixin.AppCardMessage{
		AppID:       MixinClient.ClientID,
		IconURL:     fmt.Sprintf("https://betxin.one/main/topic/%s", tid),
		Title:       "title",
		Description: "intro",
		Action:      "https://mixin.one/messenger",
		Shareable:   true,
	}

	cardJSON, err := json.Marshal(card)
	if err != nil {
		return
	}

	cardEncoded := base64.RawStdEncoding.EncodeToString(cardJSON)

	req := &mixin.MessageRequest{
		ConversationID: mixin.UniqueConversationID(MixinClient.ClientID, me.App.CreatorID),
		RecipientID:    me.App.CreatorID,
		MessageID:      mixin.RandomTraceID(),
		Category:       mixin.MessageCategoryAppCard,
		Data:           cardEncoded,
	}

	if err := MixinClient.EncryptMessageRequest(req, sessions); err != nil {
		log.Fatalln(err)
	}

	err = MixinClient.EncryptMessageRequest(req, sessions)
	if err != nil {
		log.Fatal(err)
		return
	}
}

func SendButtonGroup(ctx context.Context) {
	me, err := MixinClient.UserMe(ctx)
	if err != nil {
		log.Fatalln(err)
	}

	// bet Yes 100 No
	var buttons mixin.AppButtonGroupMessage
	buttons = append(buttons, mixin.AppButtonMessage{
		Label: "",
		Color: "",
		Action: "",
	})

	buttonsJson, err := json.Marshal(buttons)
	if err != nil {
		return
	}

	cardEncoded := base64.RawStdEncoding.EncodeToString(buttonsJson)

	req := &mixin.MessageRequest{
		ConversationID: mixin.UniqueConversationID(MixinClient.ClientID, me.App.CreatorID),
		RecipientID:    me.App.CreatorID,
		MessageID:      mixin.RandomTraceID(),
		Category:       mixin.MessageCategoryAppCard,
		Data:           cardEncoded,
	}

	err = MixinClient.SendMessage(ctx, req)
	if err != nil {
		log.Fatal(err)
	}
}