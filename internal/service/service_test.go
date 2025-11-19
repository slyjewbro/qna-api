package service

import (
	"testing"

	"qna-api/internal/model"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockRepository реализует repository.RepositoryInterface
type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) GetAllQuestions() ([]model.Question, error) {
	args := m.Called()
	return args.Get(0).([]model.Question), args.Error(1)
}

func (m *MockRepository) GetQuestionByID(id int) (*model.Question, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Question), args.Error(1)
}

func (m *MockRepository) CreateQuestion(question *model.Question) error {
	args := m.Called(question)
	return args.Error(0)
}

func (m *MockRepository) DeleteQuestion(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockRepository) CreateAnswer(answer *model.Answer) error {
	args := m.Called(answer)
	return args.Error(0)
}

func (m *MockRepository) GetAnswerByID(id int) (*model.Answer, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Answer), args.Error(1)
}

func (m *MockRepository) GetAnswersByQuestionID(questionID int) ([]model.Answer, error) {
	args := m.Called(questionID)
	return args.Get(0).([]model.Answer), args.Error(1)
}

func (m *MockRepository) DeleteAnswer(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

func TestService_CreateQuestion(t *testing.T) {
	mockRepo := new(MockRepository)
	service := NewService(mockRepo)

	// Настраиваем mock
	mockRepo.On("CreateQuestion", mock.AnythingOfType("*model.Question")).
		Return(nil).
		Run(func(args mock.Arguments) {
			question := args.Get(0).(*model.Question)
			question.ID = 1 // Симулируем присвоение ID
		})

	// Вызываем метод service
	req := model.CreateQuestionRequest{Text: "Test question"}
	result, err := service.CreateQuestion(req)

	// Проверяем результат
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "Test question", result.Text)
	assert.NotZero(t, result.ID)

	mockRepo.AssertExpectations(t)
}

func TestService_GetAllQuestions(t *testing.T) {
	mockRepo := new(MockRepository)
	service := NewService(mockRepo)

	expectedQuestions := []model.Question{
		{ID: 1, Text: "Question 1"},
		{ID: 2, Text: "Question 2"},
	}

	// Настраиваем mock
	mockRepo.On("GetAllQuestions").Return(expectedQuestions, nil)

	// Вызываем метод service
	result, err := service.GetAllQuestions()

	// Проверяем результат
	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, "Question 1", result[0].Text)

	mockRepo.AssertExpectations(t)
}

func TestService_CreateAnswer(t *testing.T) {
	mockRepo := new(MockRepository)
	service := NewService(mockRepo)

	// Настраиваем mock для проверки существования вопроса
	mockRepo.On("GetQuestionByID", 1).Return(&model.Question{ID: 1, Text: "Test question"}, nil)
	mockRepo.On("CreateAnswer", mock.AnythingOfType("*model.Answer")).
		Return(nil).
		Run(func(args mock.Arguments) {
			answer := args.Get(0).(*model.Answer)
			answer.ID = 1 // Симулируем присвоение ID
		})

	// Вызываем метод service
	req := model.CreateAnswerRequest{
		UserID: "user-123",
		Text:   "Test answer",
	}
	result, err := service.CreateAnswer(1, req)

	// Проверяем результат
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "Test answer", result.Text)
	assert.Equal(t, "user-123", result.UserID)
	assert.Equal(t, 1, result.QuestionID)

	mockRepo.AssertExpectations(t)
}

func TestService_CreateAnswer_QuestionNotFound(t *testing.T) {
	mockRepo := new(MockRepository)
	service := NewService(mockRepo)

	// Настраиваем mock для возврата ошибки (вопрос не найден)
	mockRepo.On("GetQuestionByID", 999).Return(nil, assert.AnError)
	// НЕ настраиваем CreateAnswer, так как он не должен вызываться

	// Вызываем метод service
	req := model.CreateAnswerRequest{
		UserID: "user-123",
		Text:   "Test answer",
	}
	result, err := service.CreateAnswer(999, req)

	// Проверяем что получили ошибку
	assert.Error(t, err)
	assert.Nil(t, result)

	mockRepo.AssertExpectations(t)
}

func TestService_GetAnswer(t *testing.T) {
	mockRepo := new(MockRepository)
	service := NewService(mockRepo)

	expectedAnswer := &model.Answer{
		ID:         1,
		QuestionID: 1,
		UserID:     "user-123",
		Text:       "Test answer",
	}

	// Настраиваем mock
	mockRepo.On("GetAnswerByID", 1).Return(expectedAnswer, nil)

	// Вызываем метод service
	result, err := service.GetAnswer(1)

	// Проверяем результат
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "Test answer", result.Text)
	assert.Equal(t, "user-123", result.UserID)

	mockRepo.AssertExpectations(t)
}

func TestService_DeleteQuestion(t *testing.T) {
	mockRepo := new(MockRepository)
	service := NewService(mockRepo)

	// Настраиваем mock
	mockRepo.On("DeleteQuestion", 1).Return(nil)

	// Вызываем метод service
	err := service.DeleteQuestion(1)

	// Проверяем результат
	assert.NoError(t, err)

	mockRepo.AssertExpectations(t)
}

func TestService_DeleteAnswer(t *testing.T) {
	mockRepo := new(MockRepository)
	service := NewService(mockRepo)

	// Настраиваем mock
	mockRepo.On("DeleteAnswer", 1).Return(nil)

	// Вызываем метод service
	err := service.DeleteAnswer(1)

	// Проверяем результат
	assert.NoError(t, err)

	mockRepo.AssertExpectations(t)
}
