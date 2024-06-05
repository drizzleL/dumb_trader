package bot

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"
)

var httpClient *http.Client

func init() {
	proxy := "socks5://127.0.0.1:1081"
	_proxy, _ := url.Parse(proxy)

	tr := &http.Transport{
		Proxy:           http.ProxyURL(_proxy),
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	httpClient = &http.Client{
		Transport: tr,
		Timeout:   time.Second * 5,
	}
}

// chat信息
// https://api.telegram.org/bot%s/getUpdates

var baseUrl = "https://api.telegram.org/bot%s/sendMessage"
var baseVideoUrl = "https://api.telegram.org/bot%s/sendVideo"
var token string

func realUrl() string {
	return fmt.Sprintf(baseUrl, token)
}

func realVideoUrl() string {
	return fmt.Sprintf(baseVideoUrl, token)
}

func Init(botToken string) {
	token = botToken
}

func Send() {
	c := Chat{
		// ChatId: 6950368320,
		ChatId: -1002137525475,
		Text:   "行吗 我说",
	}
	b, _ := json.Marshal(c)
	buf := bytes.NewBuffer(b)
	resp, err := httpClient.Post(realUrl(), "application/json", buf)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	b, err = io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	log.Println(string(b))

}

func SendVideo() {
	c := Chat{
		// ChatId: 6950368320,
		ChatId: -1002137525475,
		Video:  "http://i.giphy.com/13IC4LVeP5NGNi.gif",
	}
	b, _ := json.Marshal(c)
	buf := bytes.NewBuffer(b)
	resp, err := httpClient.Post(realVideoUrl(), "application/json", buf)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	b, err = io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	log.Println(string(b))

}

type Chat struct {
	Video  string `json:"video,omitempty"`
	ChatId int    `json:"chat_id,omitempty"`
	Text   string `json:"text,omitempty"`
}
