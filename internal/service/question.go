package service

import (
	"qna-api/internal/model"
	"time"
)

func (s *Service) GetAllQuestions() ([]model.Question, error) {
	return s.repo.GetAll()
}

func (s *Service) GetQuestion(id int) (*model.Question, error) {
	return s.repo.GetByID(id)
}

func (s *Service) CreateQuestion(req model.CreateQuestionRequest) (*model.Question, error) {
	question := &model.Question{
		Text:      req.Text,
		CreatedAt: time.Now(),
	}

	if err := s.repo.Create(question); err != nil {
		return nil, err
	}

	return question, nil
}

func (s *Service) DeleteQuestion(id int) error {
	return s.repo.Delete(id)
}
