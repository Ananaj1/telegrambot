package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Updates struct {
	Ok     bool `json:"ok"`
	Result []struct {
		UpdateID int `json:"update_id"`
		Message  struct {
			MessageID int `json:"message_id"`
			From      struct {
				ID        int    `json:"id"`
				FirstName string `json:"first_name"`
				LastName  string `json:"last_name"`
				Username  string `json:"username"`
			} `json:"from"`
			Chat struct {
				ID        int64  `json:"id"`
				FirstName string `json:"first_name"`
				LastName  string `json:"last_name"`
				Username  string `json:"username"`
				Type      string `json:"type"`
			} `json:"chat"`
			Date     int    `json:"date"`
			Text     string `json:"text"`
			Entities []struct {
				Type   string `json:"type"`
				Offset int    `json:"offset"`
				Length int    `json:"length"`
			} `json:"entities"`
		} `json:"message"`
	} `json:"result"`
}

type TeleAPI struct {
	apiUrl    string
	token     string
	getMsg    string
	sendMsg   string
	sendPhoto string
	offset    int
	timeout   int
	limit     int
}

func (t *TeleAPI) GetUpdates() {

	url := t.apiUrl + t.token + t.getMsg +
		"?offset=" + strconv.Itoa(t.offset) +
		"&timeout=" + strconv.Itoa(t.timeout) +
		"&limit=" + strconv.Itoa(t.limit)

	resp, err := http.Get(url)
	if err != nil {
		log.Println(err)
		time.Sleep(3 * time.Second)
		t.GetUpdates()
	} else {
		defer resp.Body.Close()

		updates := new(Updates)
		json.NewDecoder(resp.Body).Decode(&updates)

		if updates.Ok {
			for _, val := range updates.Result {
				t.SendMessage(
					val.Message.Chat.ID,
					val.Message.Chat.FirstName,
					val.Message.Text,
				)
				t.offset = val.UpdateID + 1
			}

			if t.offset > 0 {
				t.GetUpdates()
			} else {
				time.Sleep(3 * time.Second)
				t.GetUpdates()
			}
		}
	}
}

func (t *TeleAPI) SendMessage(chatID int64, name string, text string) {
	send := func(botMethod string, botMsg string) {
		jsonStr := []byte(`{"chat_id": ` + strconv.FormatInt(chatID, 10) +
			`, ` + botMsg + `}`)

		req, _ := http.NewRequest(
			"POST",
			t.apiUrl+t.token+botMethod,
			bytes.NewBuffer(jsonStr),
		)
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			log.Println(err)
		}
		defer resp.Body.Close()
	}

	// Handle message from user
	switch strings.ToLower(text) {
	case "/start":
		send(t.sendMsg, `"text": "Hello, *`+name+
			`*. I'm *GopherBot*.\nLet's play ping-pong",`+
			`"parse_mode": "markdown"`)
	case "kiss", "Kiss", "поцеловать", "Поцеловать":
		send(t.sendMsg, `"text": "о боже, *`+name+
			`* поцеловал  *WeatherBot*",`+
			`"parse_mode": "markdown"`)

	case "ping":
		send(t.sendMsg, `"text": "pong"`)
	case "pong":
		send(t.sendMsg, `"text": "ping"`)
	case "hi", "hello":
		send(t.sendMsg, `"text": "Hello"`)
	case "Хе", "хе":
		send(t.sendMsg, `"text": "Именно так"`)
	case "привет", "Привет":
		send(t.sendMsg, `"text": "Приветик"`)
	case "время", "Сколько время":
		send(t.sendMsg, `"text": "Я не могу"`)
	case "Фото", "фото":
		send(t.sendPhoto, `"photo": `+
			`"https://cdn-icons-png.flaticon.com/512/4712/4712109.png"`)
		send(t.sendMsg, `"text": "на месте"`)
	case "Да бот?", "да бот?":
		send(t.sendMsg, `"text": "Именно так как ты сказал"`)
	case "Молодец бот", "молодец бот":
		send(t.sendMsg, `"text": "Спасибо ☺ "`)
	default:
		send(t.sendPhoto, `"photo": `+
			`""`)

	}
}

func main() {
	fmt.Println("Running...")

	var teleApi = &TeleAPI{
		apiUrl:    "https://api.telegram.org/bot",
		token:     "6297302605:AAFXOApPySY8dJNy2BTQjpcdCu7zqwxp7vA",
		getMsg:    "/getUpdates",
		sendMsg:   "/sendMessage",
		sendPhoto: "/sendPhoto",
		offset:    0,
		timeout:   25,
		limit:     1,
	}

	teleApi.GetUpdates()
}
