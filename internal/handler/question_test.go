package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"qna-api/internal/model"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockService struct {
	mock.Mock
}

func (m *MockService) GetAllQuestions() ([]model.Question, error) {
	args := m.Called()
	return args.Get(0).([]model.Question), args.Error(1)
}

func (m *MockService) GetQuestion(id int) (*model.Question, error) {
	args := m.Called(id)
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

func TestCreateQuestion(t *testing.T) {
	mockService := new(MockService)
	handler := NewHandler(mockService)

	question := &model.Question{
		ID:   1,
		Text: "Test question",
	}

	mockService.On("CreateQuestion", mock.Anything).Return(question, nil)

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
