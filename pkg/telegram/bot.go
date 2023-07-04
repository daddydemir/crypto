package telegram

import (
	"github.com/daddydemir/crypto/config"
	"github.com/daddydemir/crypto/config/log"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"strconv"
)

func Pre(m1, m2 string) {
	bot, _ := tgbotapi.NewBotAPI(config.Get("BOT_TOKEN"))
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	bot.GetUpdatesChan(u)
	sendMessage(bot, m1, m2)

}

func sendMessage(bot *tgbotapi.BotAPI, m1, m2 string) {
	chatId := config.Get("CHAT_ID")
	id_int, _ := strconv.Atoi(chatId)
	id_64 := int64(id_int)
	msg1 := tgbotapi.NewMessage(id_64, "\t BIGGER \n"+m1)
	msg2 := tgbotapi.NewMessage(id_64, "\t SMALLER\n"+m2)
	_, err := bot.Send(msg1)
	bot.Send(msg2)
	if err != nil {
		log.Errorln("Mesaj gönderilirken bir hata oluştu:", err)
	}
}
