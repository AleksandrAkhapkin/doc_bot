package telegram

import (
	"time"

	b "gopkg.in/tucnak/telebot.v2"

	"doc_bot/config"
)

// NewBot создает новый инстанс бота
func NewBot(cnf config.Telegram) (*b.Bot, error) {
	s := b.Settings{
		Token:  cnf.Token,
		Poller: &b.LongPoller{Timeout: 10 * time.Second},
	}
	return b.NewBot(s)
}
