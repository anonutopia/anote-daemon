package main

// import (
// 	"fmt"
// 	"log"

// 	tgbotapi "gopkg.in/telegram-bot-api.v4"
// )

// func initBot() *tgbotapi.BotAPI {
// 	bot, err := tgbotapi.NewBotAPI(conf.Telegram)
// 	if err != nil {
// 		log.Panic(err)
// 	}

// 	bot.Debug = true

// 	log.Printf("Authorized on account %s", bot.Self.UserName)

// 	msg := tgbotapi.NewMessage(-1001325718529, "Robot successfully started - anote-daemon.")
// 	bot.Send(msg)

// 	return bot
// }

// func sendGroupsMessageInvestment(investment float64) {
// 	msg := tgbotapi.NewMessage(-1001397587839, fmt.Sprintf("We just had a new Anote purchase - %.2f EUR.", investment))
// 	bot.Send(msg)

// 	msg = tgbotapi.NewMessage(-1001161265502, fmt.Sprintf("We just had a new Anote purchase - %.2f EUR.", investment))
// 	bot.Send(msg)

// 	msg = tgbotapi.NewMessage(-1001361489843, fmt.Sprintf("We just had a new Anote purchase - %.2f EUR.", investment))
// 	bot.Send(msg)
// }

// func sendGroupsMessagePrice(newPrice float64) {
// 	msgHr := tgbotapi.NewMessage(-1001161265502, fmt.Sprintf("Cijena Anote upravo je narasla na %.8f EUR.", newPrice))
// 	bot.Send(msgHr)

// 	msgEn := tgbotapi.NewMessage(-1001361489843, fmt.Sprintf("The price of Anote has just increased to %.8f EUR.", newPrice))
// 	bot.Send(msgEn)
// }

// func logTelegram(message string) {
// 	msg := tgbotapi.NewMessage(-1001325718529, message)
// 	bot.Send(msg)
// }

// type TelegramUpdate struct {
// 	UpdateID int `json:"update_id"`
// 	Message  struct {
// 		MessageID int `json:"message_id"`
// 		From      struct {
// 			ID           int    `json:"id"`
// 			IsBot        bool   `json:"is_bot"`
// 			FirstName    string `json:"first_name"`
// 			Username     string `json:"username"`
// 			LanguageCode string `json:"language_code"`
// 		} `json:"from"`
// 		Chat struct {
// 			ID                          int    `json:"id"`
// 			Title                       string `json:"title"`
// 			Type                        string `json:"type"`
// 			AllMembersAreAdministrators bool   `json:"all_members_are_administrators"`
// 		} `json:"chat"`
// 		Date     int    `json:"date"`
// 		Text     string `json:"text"`
// 		Entities []struct {
// 			Offset int    `json:"offset"`
// 			Length int    `json:"length"`
// 			Type   string `json:"type"`
// 		} `json:"entities"`
// 	} `json:"message"`
// }
