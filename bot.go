package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"
	"net/http"

	pk "github.com/Virtoriyyya/pokemon-quiz-bot/pokemonwhogame"
	tgbotapi "gopkg.in/telegram-bot-api.v4"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	file, _ := os.Open("pokemon.csv")
	pk.AllPokemon, _ = csv.NewReader(file).ReadAll()
	pk.StoredAnswers = make(map[int64]pk.Pokemon)

	bot := setupBot()
	updates := bot.ListenForWebhook("/" + bot.Token)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	log.Printf("Listening on port %s", port)
	go http.ListenAndServe(fmt.Sprintf(":%s", port), nil)

	for update := range updates {

		if update.Message == nil {
			continue
		}


		if update.Message.Chat.IsGroup() || update.Message.Chat.IsSuperGroup() {
			log.Printf("Request from chat: %s", update.Message.Chat.Title)
		} else {
			log.Printf("Request from user: %s", update.Message.Chat.UserName)
		}
		if bot.Debug {
			log.Printf("Update: %v", update.Message.Text)
		}

		command := update.Message.Command()
		switch {

		case strings.EqualFold(command, "who"):
			pk.WhosThatPokemon(bot, update)

		case strings.EqualFold(command, "its"):
			pk.Its(bot, update)

		case strings.EqualFold(command, "debug") && update.Message.Chat.ID == 36992723:
			rsp := fmt.Sprintf("Switching debug mode to %t", !bot.Debug)
			log.Printf(rsp)
			bot.Debug = !bot.Debug

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, rsp)
			msg.ReplyToMessageID = update.Message.MessageID
			bot.Send(msg)
		}
	}
}


func setupBot() *tgbotapi.BotAPI {
	bot, err := tgbotapi.NewBotAPI("7184008496:AAHpq8nKb7uVIsbDBTpMqCV9NJL1kAtmovg")
	if err != nil {
		log.Fatal(err)
	}

	bot.Debug = false
	log.Printf("Authorized on account %s", bot.Self.UserName)

	_, err = bot.SetWebhook(tgbotapi.NewWebhook("https://0016-46-63-26-132.ngrok-free.app/" + bot.Token))
	
	if err != nil {
		log.Fatal(err)
	}
	info, err := bot.GetWebhookInfo()
	if err != nil {
		log.Fatal(err)
	}
	if info.LastErrorDate != 0 {
		log.Printf("Telegram callback failed: %s", info.LastErrorMessage)
	}
	return bot
}
