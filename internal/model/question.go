package model

import (
	"time"
)

type Question struct {
	ID        int       `json:"id" gorm:"primaryKey"`
	Text      string    `json:"text" gorm:"type:text;not null"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	Answers   []Answer  `json:"answers,omitempty" gorm:"foreignKey:QuestionID;constraint:OnDelete:CASCADE"`
}

type CreateQuestionRequest struct {
	Text string `json:"text" validate:"required,min=1"`
}
