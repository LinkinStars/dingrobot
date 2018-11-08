package dingrobot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// Roboter is the interface implemented by Robot that can send multiple types of messages.
type Roboter interface {
	SendText(content string, atMobiles []string, isAtAll bool) error
	SendLink(title, text, messageURL, picURL string) error
	SendMarkdown(title, text string, atMobiles []string, isAtAll bool) error
}

// Robot represents a dingtalk custom robot that can send messages to groups.
type Robot struct {
	Webhook string
}

// New returns a dingtalk robot.
func New(webhook string) Roboter {
	return Robot{Webhook: webhook}
}

// SendText send a text type message.
func (r Robot) SendText(content string, atMobiles []string, isAtAll bool) error {
	return r.send(&textMessage{
		MsgType: msgTypeText,
		Text: textParam{
			Content: content,
		},
		At: atParam{
			AtMobiles: atMobiles,
			IsAtAll:   isAtAll,
		},
	})
}

// SendLink send a link type message.
func (r Robot) SendLink(title, text, messageURL, picURL string) error {
	return r.send(&linkMessage{
		MsgType: msgTypeLink,
		Link: linkParam{
			Title:      title,
			Text:       text,
			MessageURL: messageURL,
			PicURL:     picURL,
		},
	})
}

// SendMarkdown send a markdown type message.
func (r Robot) SendMarkdown(title, text string, atMobiles []string, isAtAll bool) error {
	return r.send(&markdownMessage{
		MsgType: msgTypeMarkdown,
		Markdown: markdownParam{
			Title: title,
			Text:  text,
		},
		At: atParam{
			AtMobiles: atMobiles,
			IsAtAll:   isAtAll,
		},
	})
}

type messageResponse struct {
	Errcode int
	Errmsg  string
}

func (r Robot) send(msg interface{}) error {
	m, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	resp, err := http.Post(r.Webhook, "application/json", bytes.NewReader(m))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var mr messageResponse
	err = json.Unmarshal(data, &mr)
	if err != nil {
		return err
	}
	if mr.Errcode != 0 {
		return fmt.Errorf("dingrobot send failed: %v", mr.Errmsg)
	}

	return nil
}