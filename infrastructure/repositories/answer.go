package repositories

import (
	"doc_bot/pkg/domain"

	"github.com/jackc/pgconn"
	"gorm.io/gorm"
)

// AnswerRepository управляет сущностью Shop.
type AnswerRepository struct {
	db *gorm.DB
}

// NewAnswerRepository возвращает новый инстанс AnswerRepository.
func NewAnswerRepository(
	db *gorm.DB,
) (*AnswerRepository, error) {
	return &AnswerRepository{db: db}, nil
}

// FindAnswer находит ответ соответсвующий вопросу
func (f *AnswerRepository) FindAnswer(question string) (string, error) {
	questionAnswer := QuestionsAnswers{}
	err := f.db.Where(QuestionsAnswers{Question: question}).Take(&questionAnswer).Error
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok &&
			pgErr.ConstraintName != "users__login_and_event_id_uniq" {
			return "", err
		}
		return "", domain.ErrNotFindAnswer
	}
	return questionAnswer.Answer, nil
}
