package main

import (
    "github.com/go-telegram-bot-api/telegram-bot-api"
    "log"
    "os"
)

func main() {

    token, ok := os.LookupEnv("API_TOKEN")
    if !ok {
        log.Panic("API_TOKEN missing")
    }

    bot, err := tgbotapi.NewBotAPI(token)
    if err != nil {
        log.Panic(err)
    }

    bot.Debug = true

    log.Printf("Authorized on account %s", bot.Self.UserName)

    u := tgbotapi.NewUpdate(0)
    u.Timeout = 60

    updates, err := bot.GetUpdatesChan(u)

    for update := range updates {
        if update.Message == nil {
            continue
        }

        log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

        if update.Message.IsCommand() {
            msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
            switch update.Message.Command() {
            case "help":
                msg.Text = "type /m or /help for this message"
            case "m":
                msg.Text = "retrieving your time attendance from today"
            default:
                msg.Text = "I don't know that command"
            }
            bot.Send(msg)
        }

    }
}
