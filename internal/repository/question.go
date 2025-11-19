package repository

import (
	"qna-api/internal/model"
)

// Методы для вопросов
func (r *Repository) GetAllQuestions() ([]model.Question, error) {
	var questions []model.Question
	result := r.db.Find(&questions)
	return questions, result.Error
}

func (r *Repository) GetQuestionByID(id int) (*model.Question, error) {
	var question model.Question
	result := r.db.Preload("Answers").First(&question, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &question, nil
}

func (r *Repository) CreateQuestion(question *model.Question) error {
	result := r.db.Create(question)
	return result.Error
}

func (r *Repository) DeleteQuestion(id int) error {
	result := r.db.Delete(&model.Question{}, id)
	return result.Error
}
