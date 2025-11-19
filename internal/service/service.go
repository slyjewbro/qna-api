package service

import (
	"qna-api/internal/model"
	"qna-api/internal/repository"
)

// ServiceInterface определяет контракт для сервиса
type ServiceInterface interface {
	// Question methods
	GetAllQuestions() ([]model.Question, error)
	GetQuestion(id int) (*model.Question, error)
	CreateQuestion(req model.CreateQuestionRequest) (*model.Question, error)
	DeleteQuestion(id int) error

	// Answer methods
	CreateAnswer(questionID int, req model.CreateAnswerRequest) (*model.Answer, error)
	GetAnswer(id int) (*model.Answer, error)
	DeleteAnswer(id int) error
}

// ServiceImpl - реализация сервиса
type ServiceImpl struct {
	repo repository.RepositoryInterface // Используем интерфейс
}

// NewService создает новый экземпляр сервиса
func NewService(repo repository.RepositoryInterface) ServiceInterface {
	return &ServiceImpl{repo: repo}
}
