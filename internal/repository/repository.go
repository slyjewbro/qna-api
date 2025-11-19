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

type IQuestionRepository interface {
	GetAll() ([]model.Question, error)
	GetByID(id int) (*model.Question, error)
	Create(question *model.Question) error
	Delete(id int) error
}

type IAnswerRepository interface {
	Create(answer *model.Answer) error
	GetByID(id int) (*model.Answer, error)
	GetByQuestionID(questionID int) ([]model.Answer, error)
	Delete(id int) error
}
