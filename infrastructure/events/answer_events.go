package events

import (
	"doc_bot/infrastructure/repositories"
	"doc_bot/libs/liblog"
)

// AnswerQuerier интерфейс для ответов на комманды и сообщения
type AnswerQuerier interface {
	FindAnswer(question string) (string, error)
}

// AnswerQueries представляет выборки ответов
type AnswerQueries struct {
	logger     liblog.Logger
	answerRepo *repositories.AnswerRepository
}

// NewAnswerQueries конструктор AnswerQueries
func NewAnswerQueries(
	answerRepo *repositories.AnswerRepository,
	logger liblog.Logger,
) *AnswerQueries {
	return &AnswerQueries{answerRepo: answerRepo, logger: logger}
}
