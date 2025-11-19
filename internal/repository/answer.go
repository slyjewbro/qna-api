package repository

import (
	"qna-api/internal/model"
)

func (r *Repository) Create(answer *model.Answer) error {
	// Проверяем существование вопроса
	var question model.Question
	if err := r.db.First(&question, answer.QuestionID).Error; err != nil {
		return err
	}

	result := r.db.Create(answer)
	return result.Error
}

func (r *Repository) GetByID(id int) (*model.Answer, error) {
	var answer model.Answer
	result := r.db.First(&answer, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &answer, nil
}

func (r *Repository) GetByQuestionID(questionID int) ([]model.Answer, error) {
	var answers []model.Answer
	result := r.db.Where("question_id = ?", questionID).Find(&answers)
	return answers, result.Error
}

func (r *Repository) Delete(id int) error {
	result := r.db.Delete(&model.Answer{}, id)
	return result.Error
}
