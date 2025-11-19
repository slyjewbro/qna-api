package repository

import (
	"qna-api/internal/model"
)

func (r *Repository) GetAll() ([]model.Question, error) {
	var questions []model.Question
	result := r.db.Find(&questions)
	return questions, result.Error
}

func (r *Repository) GetByID(id int) (*model.Question, error) {
	var question model.Question
	result := r.db.Preload("Answers").First(&question, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &question, nil
}

func (r *Repository) Create(question *model.Question) error {
	result := r.db.Create(question)
	return result.Error
}

func (r *Repository) Delete(id int) error {
	result := r.db.Delete(&model.Question{}, id)
	return result.Error
}
