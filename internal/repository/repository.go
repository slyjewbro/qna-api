package repository

import (
	"qna-api/internal/model"

	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// Интерфейсы вопросов
type IQuestionRepository interface {
	GetAllQuestions() ([]model.Question, error)
	GetQuestionByID(id int) (*model.Question, error)
	CreateQuestion(question *model.Question) error
	DeleteQuestion(id int) error
}

// Интерфейсы ответов
type IAnswerRepository interface {
	CreateAnswer(answer *model.Answer) error
	GetAnswerByID(id int) (*model.Answer, error)
	GetAnswersByQuestionID(questionID int) ([]model.Answer, error)
	DeleteAnswer(id int) error
}
