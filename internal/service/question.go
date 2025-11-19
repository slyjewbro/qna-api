package service

import (
	"qna-api/internal/model"
	"time"
)

func (s *Service) GetAllQuestions() ([]model.Question, error) {
	return s.repo.GetAllQuestions()
}

func (s *Service) GetQuestion(id int) (*model.Question, error) {
	return s.repo.GetQuestionByID(id)
}

func (s *Service) CreateQuestion(req model.CreateQuestionRequest) (*model.Question, error) {
	question := &model.Question{
		Text:      req.Text,
		CreatedAt: time.Now(),
	}

	if err := s.repo.CreateQuestion(question); err != nil {
		return nil, err
	}

	return question, nil
}

func (s *Service) DeleteQuestion(id int) error {
	return s.repo.DeleteQuestion(id)
}
