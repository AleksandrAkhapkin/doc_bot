package events

// FindAnswer находит ответ совпадающий с вопросом
func (a *AnswerQueries) FindAnswer(question string) (string, error) {
	return a.answerRepo.FindAnswer(question)
}
