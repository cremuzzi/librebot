package main

import (
    "log"
    "os"
    "strings"
    "github.com/go-telegram-bot-api/telegram-bot-api"
    MQTT "github.com/eclipse/paho.mqtt.golang"
)

func getEnv(key, fallback string) string {
    if value, ok := os.LookupEnv(key); ok {
        return value
    }
    return fallback
}

func handleLuzAction(theMessage string) string {
    msgArgs := strings.Fields(theMessage)
    var action string
    if len(msgArgs) < 3 {
        return "falta especificar el numero de la sala y la accion\nejemplo: /luz 41 off"
    }
    switch msgArgs[2] {
        case "ON","on","On","oN","1":
            action = "1"
        case "OFF","off","Off","0":
            action = "0"
        default:
            return "la accion tiene que ser on o off\nejemplo: /luz 41 off"
    } 
    lightSwitcher(msgArgs[1], action)
    return "cambiando la luz"
}

func lightSwitcher(roomid string, msg string) {
    var topic string = "l/"+roomid+"/s"
    var broker string = getEnv("MQTT_BROKER","tcp://iot.eclipse.org:1883")
 
    log.Printf("broker %s", broker)
    log.Printf("topic %s", topic)
    log.Printf("msg %s", msg)

    opts := MQTT.NewClientOptions()
    opts.AddBroker(broker)
    opts.SetClientID("librebot")

    client := MQTT.NewClient(opts)
    if token := client.Connect(); token.Wait() && token.Error() != nil {
        panic(token.Error())
    }
    log.Printf("publishing mqtt")
    token := client.Publish(topic, byte(0), false, msg)
    token.Wait()

    client.Disconnect(250)
    log.Printf("finished publishing mqtt")
}

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
                msg.Text = "/m tus marcadas\n/luz prender o apagar la luz\n/help este mensaje"
            case "m":
                msg.Text = "retrieving your time attendance from today"
            case "luz":
                msg.Text = handleLuzAction(update.Message.Text)
            default:
                msg.Text = "I don't know that command"
            }
            bot.Send(msg)
        }

    }
}
