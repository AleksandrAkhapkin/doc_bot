package domain

import "doc_bot/libs/liberror"

var (
	// ErrNotFindAnswer Ответ не найден
	ErrNotFindAnswer = &liberror.Error{
		Err: "Даже не знаю что на это ответить",
	}
)
