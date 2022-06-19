package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/xelaj/go-dry"
	"github.com/xelaj/mtproto/telegram"

	utils "github.com/xelaj/mtproto/examples/example_utils"
)

type App struct {
	ID   int    `json:"id"`
	Hash string `json:"hash"`
}

func main() {
	appStorage := utils.PrepareAppStorageForExamples()
	sessionFile := filepath.Join(appStorage, "session.json")
	appCredentials := "credentials.json"
	publicKeys := "tg_public_keys.pem"

	jsonFile, err := os.Open(appCredentials)
	dry.PanicIfErr(err)

	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	var app App

	json.Unmarshal([]byte(byteValue), &app)

	client, err := telegram.NewClient(telegram.ClientConfig{
		SessionFile:     sessionFile,
		ServerHost:      "149.154.167.50:443",
		PublicKeysFile:  publicKeys,
		AppID:           app.ID,
		AppHash:         app.Hash,
		InitWarnChannel: true,
	})
	utils.ReadWarningsToStdErr(client.Warnings)
	dry.PanicIfErr(err)

	channel, err := client.ContactsResolveUsername("ithueti")
	dry.PanicIfErr(err)

	chat := channel.Chats[0].(*telegram.Channel)

	messages, err := client.MessagesSearch(&telegram.MessagesSearchParams{
		Peer: &telegram.InputPeerChannel{
			ChannelID:  chat.ID,
			AccessHash: chat.AccessHash,
		},
		Q:         "#деградироватьпришёл",
		Filter:    &telegram.InputMessagesFilterEmpty{},
		MinDate:   -1,
		MaxDate:   -1,
		OffsetID:  0,
		AddOffset: 0,
		Limit:     1000,
		MaxID:     100000,
		MinID:     0,
	})
	dry.PanicIfErr(err)

	var ids []int32
	msg := messages.(*telegram.MessagesChannelMessages).Messages
	for _, m := range msg {
		ids = append(ids, m.(*telegram.MessageObj).ID)
	}

	jsArray, err := json.Marshal(ids)
	dry.PanicIfErr(err)

	println(string(jsArray))
}
