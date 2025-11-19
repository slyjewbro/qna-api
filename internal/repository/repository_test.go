package repository

import (
	"testing"

	"qna-api/internal/model"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupTestDB() *gorm.DB {
	// Используем тестовую PostgreSQL базу
	dsn := "host=localhost user=postgres password=password dbname=qna_test port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		// Если PostgreSQL не доступен, пропускаем тест
		return nil
	}

	// Auto migrate models
	db.AutoMigrate(&model.Question{}, &model.Answer{})

	// Очищаем таблицы перед тестом
	db.Exec("TRUNCATE TABLE answers CASCADE")
	db.Exec("TRUNCATE TABLE questions CASCADE")

	return db
}

func TestCreateAndGetQuestion(t *testing.T) {
	db := setupTestDB()
	if db == nil {
		t.Skip("PostgreSQL not available, skipping test")
		return
	}

	repo := NewRepository(db)

	// Test create question
	question := &model.Question{Text: "Test question"}
	err := repo.CreateQuestion(question)
	assert.NoError(t, err)
	assert.NotZero(t, question.ID)

	// Test get question
	found, err := repo.GetQuestionByID(question.ID)
	assert.NoError(t, err)
	assert.Equal(t, "Test question", found.Text)
}

func TestCreateAnswer(t *testing.T) {
	db := setupTestDB()
	if db == nil {
		t.Skip("PostgreSQL not available, skipping test")
		return
	}

	repo := NewRepository(db)

	// Create question first
	question := &model.Question{Text: "Test question"}
	repo.CreateQuestion(question)

	// Create answer
	answer := &model.Answer{
		QuestionID: question.ID,
		UserID:     "user-123",
		Text:       "Test answer",
	}

	err := repo.CreateAnswer(answer)
	assert.NoError(t, err)
	assert.NotZero(t, answer.ID)
}

func TestCascadeDelete(t *testing.T) {
	db := setupTestDB()
	if db == nil {
		t.Skip("PostgreSQL not available, skipping test")
		return
	}

	repo := NewRepository(db)

	// Create question with answers
	question := &model.Question{Text: "Test question"}
	repo.CreateQuestion(question)

	answer := &model.Answer{
		QuestionID: question.ID,
		UserID:     "user-123",
		Text:       "Test answer",
	}
	repo.CreateAnswer(answer)

	// Delete question
	err := repo.DeleteQuestion(question.ID)
	assert.NoError(t, err)

	// Verify answer is also deleted
	_, err = repo.GetAnswerByID(answer.ID)
	assert.Error(t, err)
}
