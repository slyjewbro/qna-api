package repository

import "qna-api/internal/model"

// RepositoryInterface определяет контракт для репозитория
type RepositoryInterface interface {
	// Question methods
	GetAllQuestions() ([]model.Question, error)
	GetQuestionByID(id int) (*model.Question, error)
	CreateQuestion(question *model.Question) error
	DeleteQuestion(id int) error

	// Answer methods
	CreateAnswer(answer *model.Answer) error
	GetAnswerByID(id int) (*model.Answer, error)
	GetAnswersByQuestionID(questionID int) ([]model.Answer, error)
	DeleteAnswer(id int) error
}
