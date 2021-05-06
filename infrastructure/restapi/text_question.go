package restapi

import b "gopkg.in/tucnak/telebot.v2"

// TextQuestion создает обработчик для получения ответа.
func (handlers *CommandLineHandlers) TextQuestion(m *b.Message) {
	res, err := handlers.Endpoints.TextQuestions(m)
	if err != nil {
		handlers.errorResponse(m, err)
	}
	handlers.response(m, res)
}
