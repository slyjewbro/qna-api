package model

import (
	"time"
)

type Answer struct {
	ID         int       `json:"id" gorm:"primaryKey"`
	QuestionID int       `json:"question_id" gorm:"not null;index"`
	UserID     string    `json:"user_id" gorm:"type:varchar(36);not null;index"`
	Text       string    `json:"text" gorm:"type:text;not null"`
	CreatedAt  time.Time `json:"created_at" gorm:"autoCreateTime"`
}

type CreateAnswerRequest struct {
	UserID string `json:"user_id" validate:"required,min=1"`
	Text   string `json:"text" validate:"required,min=1"`
}
