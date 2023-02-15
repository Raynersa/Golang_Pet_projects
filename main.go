package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"
)

type Update struct {
	UpdateId int `json:"update_id"`
	Message  struct {
		MessageID int `json:"message_id"`
		From      struct {
			ID        int    `json:"id"`
			IsBot     bool   `json:"is_bot"`
			FirstName string `json:"first_name"`
			UserName  string `json:"username"`
			LangCode  string `json:"language_code"`
		} `json:"from"`
		Chat struct {
			ChatID    int    `json:"id"`
			FirstName string `json:"first_name"`
			UserName  string `json:"username"`
			Type      string `json:"type"`
		} `json:"chat"`
		Date    int    `json:"date"`
		Text    string `json:"text"`
		Sticker struct {
			Width       int    `json:"width"`
			Height      int    `json:"height"`
			Emoji       string `json:"emoji"`
			SetName     string `json:"set_name"`
			IsAnimated  bool   `json:"is_animated"`
			IsVideo     bool   `json:"is_video"`
			TypeSticker string `json:"type"`
			Thumb       struct {
				FileID       string `json:"text"`
				FileUniqueId string `json:"file_unique_id"`
				FileSize     int    `json:"file_size"`
				Width        int    `json:"width"`
				Height       int    `json:"height"`
			} `json:"thumb"`
			FileID       string `json:"text"`
			FileUniqueId string `json:"file_unique_id"`
			FileSize     int    `json:"file_size"`
		} `json:"sticker"`
		Entities []struct {
			OffSet int    `json:"offset"`
			Length int    `json:"length"`
			Type   string `json:"type"`
		} `json:"entities"`
	} `json:"message"`
}

type GetUpdatesResponse struct {
	Ok     bool     `json:"ok"`
	Result []Update `json:"result"`
}

type Message struct {
	ChatID int    `json:"chat_id"`
	Text   string `json:"text"`
}

type BotMessage struct {
	ChatId int    `json:"chat_id"`
	Text   string `json:"text"`
	// ReplyMarkup ReplyKeyboardMarkup `json:"reply_markup"`
	Keyboard InlineKeyboardMarkup `json:"reply_markup"`
}

type ReplyKeyboardMarkup struct {
	Keyboard [][]KeyboardButton `json:"keyboard"`
	Resize   bool               `json:"resize_keyboard"`
}

type InlineKeyboardMarkup struct {
	Inline_keyboard [][]InlineKeyboardButton `json:"inline_keyboard"`
}

type InlineKeyboardButton struct {
	Text          string `json:"text"`
	Callback_data string `json:"callback_data"`
}

type KeyboardButton struct {
	Text string `json:"text"`
}

func getUpdates(update Update, botUrl string) {
	msg := Message{ChatID: update.Message.Chat.ChatID, Text: update.Message.Text}

	jsonString, enn := json.Marshal(msg)
	if enn != nil {
		panic(enn)
	}
	fmt.Println(string(jsonString))

	client := http.Client{}
	resp, err := client.Post(botUrl+"/sendMessage", "Application/json", bytes.NewReader(jsonString))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	fmt.Println(string(body))
}

func respond_bot(update Update, botUrl string) error {
	var botMessage BotMessage
	botMessage.ChatId = update.Message.Chat.ChatID
	message := "Nen"
	// botMessage.ReplyMarkup.Resize = true
	if update.Message.Text == "Start" {
		botMessage.Text = "Start"
		// botMessage.ReplyMarkup.Keyboard = [][]KeyboardButton{{{Text: message}}}
		botMessage.Keyboard.Inline_keyboard = [][]InlineKeyboardButton{{{Text: message, Callback_data: ""}}}
	}

	buf, err := json.Marshal(botMessage)
	if err != nil {
		return err
	}
	fmt.Println(string(buf))
	_, err = http.Post(botUrl+"/sendMessage", "application/json", bytes.NewBuffer(buf))
	if err != nil {
		return err
	}
	return nil
}

func main() {
	offset := 0
	BotAPI := "https://api.telegram.org/bot"
	botToken := "6114897006:AAF8GZM4pG7ruyKZ9WfJbema_EcoPteqcdQ"
	botUrl := BotAPI + botToken
	for {
		res, err := http.Get(botUrl + "/getUpdates" + "?offset=" + strconv.Itoa(offset))
		if err != nil {
			fmt.Printf("error making http request: %s\n", err)
			os.Exit(1)
		}

		body, err := ioutil.ReadAll(res.Body) // response body is []byte
		fmt.Printf(string(body), "\n")
		var resp GetUpdatesResponse
		if err := json.Unmarshal(body, &resp); err != nil { // Parse []byte to the go struct pointer
			fmt.Println("Can not unmarshal JSON")
			fmt.Printf("error %s", err)
		}

		if resp.Ok == true {
			for _, update := range resp.Result {

				getUpdates(update, botUrl)
				respond_bot(update, botUrl)
				offset = update.UpdateId + 1
			}
		}
		time.Sleep(1 * time.Second)
	}
}
