package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"qna-api/internal/model"
	"qna-api/internal/service"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockService реализует service.ServiceInterface
type MockService struct {
	mock.Mock
}

func (m *MockService) GetAllQuestions() ([]model.Question, error) {
	args := m.Called()
	return args.Get(0).([]model.Question), args.Error(1)
}

func (m *MockService) GetQuestion(id int) (*model.Question, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Question), args.Error(1)
}

func (m *MockService) CreateQuestion(req model.CreateQuestionRequest) (*model.Question, error) {
	args := m.Called(req)
	return args.Get(0).(*model.Question), args.Error(1)
}

func (m *MockService) DeleteQuestion(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockService) CreateAnswer(questionID int, req model.CreateAnswerRequest) (*model.Answer, error) {
	args := m.Called(questionID, req)
	return args.Get(0).(*model.Answer), args.Error(1)
}

func (m *MockService) GetAnswer(id int) (*model.Answer, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Answer), args.Error(1)
}

func (m *MockService) DeleteAnswer(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

func TestHealthCheck(t *testing.T) {
	// Health check не требует service, можно передать nil
	// Но так как мы используем интерфейс, нужно передать nil явно
	var svc service.ServiceInterface = nil
	handler := NewHandler(svc)

	req := httptest.NewRequest("GET", "/health", nil)
	rr := httptest.NewRecorder()

	router := handler.InitRoutes()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var response map[string]string
	json.Unmarshal(rr.Body.Bytes(), &response)

	assert.Equal(t, "healthy", response["status"])
}

func TestCreateQuestion_Success(t *testing.T) {
	mockService := new(MockService)
	handler := NewHandler(mockService)

	expectedQuestion := &model.Question{
		ID:   1,
		Text: "Test question",
	}

	mockService.On("CreateQuestion", mock.AnythingOfType("model.CreateQuestionRequest")).
		Return(expectedQuestion, nil)

	reqBody := map[string]string{"text": "Test question"}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("POST", "/questions", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	router := handler.InitRoutes()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)

	var response model.Question
	json.Unmarshal(rr.Body.Bytes(), &response)

	assert.Equal(t, 1, response.ID)
	assert.Equal(t, "Test question", response.Text)

	mockService.AssertExpectations(t)
}

func TestGetQuestions_Success(t *testing.T) {
	mockService := new(MockService)
	handler := NewHandler(mockService)

	expectedQuestions := []model.Question{
		{ID: 1, Text: "Question 1"},
		{ID: 2, Text: "Question 2"},
	}

	// Настраиваем mock
	mockService.On("GetAllQuestions").Return(expectedQuestions, nil)

	// Выполняем запрос
	req := httptest.NewRequest("GET", "/questions", nil)
	rr := httptest.NewRecorder()

	router := handler.InitRoutes()
	router.ServeHTTP(rr, req)

	// Проверяем результат
	assert.Equal(t, http.StatusOK, rr.Code)

	var response []model.Question
	json.Unmarshal(rr.Body.Bytes(), &response)

	assert.Len(t, response, 2)
	assert.Equal(t, "Question 1", response[0].Text)
	assert.Equal(t, "Question 2", response[1].Text)

	mockService.AssertExpectations(t)
}

func TestGetQuestion_NotFound(t *testing.T) {
	mockService := new(MockService)
	handler := NewHandler(mockService)

	// Настраиваем mock для возврата ошибки
	mockService.On("GetQuestion", 999).Return(nil, assert.AnError)

	// Выполняем запрос
	req := httptest.NewRequest("GET", "/questions/999", nil)
	rr := httptest.NewRecorder()

	router := handler.InitRoutes()
	router.ServeHTTP(rr, req)

	// Проверяем что получили 404
	assert.Equal(t, http.StatusNotFound, rr.Code)

	mockService.AssertExpectations(t)
}

func TestCreateAnswer_Success(t *testing.T) {
	mockService := new(MockService)
	handler := NewHandler(mockService)

	expectedAnswer := &model.Answer{
		ID:         1,
		QuestionID: 1,
		UserID:     "user-123",
		Text:       "Test answer",
	}

	// Настраиваем mock
	mockService.On("CreateAnswer", 1, mock.AnythingOfType("model.CreateAnswerRequest")).
		Return(expectedAnswer, nil)

	// Подготавливаем запрос
	reqBody := map[string]string{
		"user_id": "user-123",
		"text":    "Test answer",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("POST", "/questions/1/answers", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	// Устанавливаем параметры маршрута для mux
	req = mux.SetURLVars(req, map[string]string{"id": "1"})

	// Выполняем запрос
	router := handler.InitRoutes()
	router.ServeHTTP(rr, req)

	// Проверяем результат
	assert.Equal(t, http.StatusCreated, rr.Code)

	var response model.Answer
	json.Unmarshal(rr.Body.Bytes(), &response)

	assert.Equal(t, 1, response.ID)
	assert.Equal(t, "Test answer", response.Text)
	assert.Equal(t, "user-123", response.UserID)

	mockService.AssertExpectations(t)
}
