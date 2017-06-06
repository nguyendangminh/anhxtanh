package main

import (
	"io/ioutil"
	"encoding/json"
	"log"

	"github.com/michlabs/fbbot"
	fptai "github.com/fpt-corp/fptai-sdk-go"
)

const (
	FPTAI_TOKEN = "your_fptai_application_token"

	FB_PAGE_ACCESS_TOKEN = "your_fb_page_access_token"
	FB_VERIFY_TOKEN = "your_fb_verify_token"
)

var client *fptai.Client
var PORT int = 1203
var ErrMsg string = "Xin lỗi, tôi không biết. Biển học vô bờ."

type Einstein struct {
	Answers []Answer 
}
type Answer struct {
	Intent string `json:"intent"`
	Text string `json:"answer"`
}

func (t *Einstein) Load(path string) error {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(data, &t.Answers); err != nil {
		return err
	}
	return nil
}

func (t *Einstein) Answer(intent string) string {
	for _, answer := range t.Answers {
		if answer.Intent == intent {
			return answer.Text
		}
	}
	return ErrMsg
}

func (t *Einstein) HandleMessage(bot *fbbot.Bot, msg *fbbot.Message) {
	bot.TypingOn(msg.Sender)
	resp, err := client.RecognizeIntents(msg.Text)
	if err != nil || len(resp.Intents) == 0 {
		bot.SendText(msg.Sender, ErrMsg)
		return
	}
	intent := resp.Intents[0].Name
	bot.SendText(msg.Sender, t.Answer(intent)) 
}

func init() {
	client = fptai.NewClient(FPTAI_TOKEN)
}

func main() {
	var ei Einstein
	if err := ei.Load("data.json"); err != nil {
		log.Fatal(err)
	}

	bot := fbbot.New(PORT, FB_VERIFY_TOKEN, FB_PAGE_ACCESS_TOKEN)
	bot.AddMessageHandler(&ei)
	bot.Run()
}