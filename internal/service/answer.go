package service

import (
	"qna-api/internal/model"
	"time"
)

func (s *Service) CreateAnswer(questionID int, req model.CreateAnswerRequest) (*model.Answer, error) {
	answer := &model.Answer{
		QuestionID: questionID,
		UserID:     req.UserID,
		Text:       req.Text,
		CreatedAt:  time.Now(),
	}

	if err := s.repo.CreateAnswer(answer); err != nil {
		return nil, err
	}

	return answer, nil
}

func (s *Service) GetAnswer(id int) (*model.Answer, error) {
	return s.repo.GetAnswerByID(id)
}

func (s *Service) DeleteAnswer(id int) error {
	return s.repo.DeleteAnswer(id)
}
