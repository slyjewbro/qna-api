package service

import (
	"qna-api/internal/model"
	"qna-api/internal/repository"
)

type Service struct {
	repo *repository.Repository
}

func NewService(repo *repository.Repository) *Service {
	return &Service{repo: repo}
}

type IQuestionService interface {
	GetAllQuestions() ([]model.Question, error)
	GetQuestion(id int) (*model.Question, error)
	CreateQuestion(req model.CreateQuestionRequest) (*model.Question, error)
	DeleteQuestion(id int) error
}

type IAnswerService interface {
	CreateAnswer(questionID int, req model.CreateAnswerRequest) (*model.Answer, error)
	GetAnswer(id int) (*model.Answer, error)
	DeleteAnswer(id int) error
}
