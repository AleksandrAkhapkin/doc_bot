package restapi

import (
	"doc_bot/libs/liberror"
	"doc_bot/libs/liblog"
	"doc_bot/pkg/application/endpoints"

	b "gopkg.in/tucnak/telebot.v2"
)

// CommandLineHandlers представляет обработчики Shop.
type CommandLineHandlers struct {
	Endpoints *endpoints.CommandLineEndpoints
	Bot       *b.Bot
	Logger    liblog.Logger
}

// NewCommandLineHandlers возвращает новый инстанс CommandLineHandlers.
func NewCommandLineHandlers(
	logger liblog.Logger,
	endpoints *endpoints.CommandLineEndpoints,
	bot *b.Bot,
) *CommandLineHandlers {
	return &CommandLineHandlers{
		Endpoints: endpoints,
		Bot:       bot,
		Logger:    logger,
	}
}

func (handlers *CommandLineHandlers) response(m *b.Message, res interface{}) {
	mes, err := handlers.Bot.Send(m.Sender, res)
	if err != nil {
		handlers.Logger.WithField("user", m.Sender).WithField("mes", mes)
		return
	}
}
func (handlers *CommandLineHandlers) errorResponse(m *b.Message, err error) {
	if err, ok := err.(*liberror.Error); ok {
		mes, err := handlers.Bot.Send(m.Sender, err.Err)
		if err != nil {
			handlers.Logger.WithField("user", m.Sender).WithField("mes", mes)
			return
		}
		return
	}
	mes, err := handlers.Bot.Send(
		m.Sender, "Со мной что то случилось, я не могу ответить!")
	if err != nil {
		handlers.Logger.WithField("user", m.Sender).WithField("mes", mes)
		return
	}
}
