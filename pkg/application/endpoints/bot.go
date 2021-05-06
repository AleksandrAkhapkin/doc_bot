package endpoints

import (
	"doc_bot/infrastructure/events"

	b "gopkg.in/tucnak/telebot.v2"
)

// CommandLineEndpoints представляет конечные точки Shop вне транспортного уровня.
type CommandLineEndpoints struct {
	queries events.AnswerQuerier
}

// NewCommandLineEndpoints возвращает новый инстанс CommandLineEndpoints.
func NewCommandLineEndpoints(
	queries events.AnswerQuerier,
) *CommandLineEndpoints {
	return &CommandLineEndpoints{queries: queries}
}

// TextQuestions конечная точка для HelloCommandLine
func (e *CommandLineEndpoints) TextQuestions(
	message *b.Message,
) (string, error) {
	return e.queries.FindAnswer(message.Text)
}
