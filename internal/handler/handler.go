package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"qna-api/internal/model"
	"qna-api/internal/service"

	"github.com/gorilla/mux"
)

type Handler struct {
	service service.ServiceInterface
}

func NewHandler(service service.ServiceInterface) *Handler {
	return &Handler{service: service}
}

func (h *Handler) InitRoutes() *mux.Router {
	router := mux.NewRouter()

	// Root and health routes
	router.HandleFunc("/", h.rootHandler).Methods("GET")
	router.HandleFunc("/health", h.healthCheck).Methods("GET")

	// Questions routes
	router.HandleFunc("/questions", h.GetQuestions).Methods("GET")
	router.HandleFunc("/questions", h.CreateQuestion).Methods("POST")
	router.HandleFunc("/questions/{id}", h.GetQuestion).Methods("GET")
	router.HandleFunc("/questions/{id}", h.DeleteQuestion).Methods("DELETE")

	// Answers routes
	router.HandleFunc("/questions/{id}/answers", h.CreateAnswer).Methods("POST")
	router.HandleFunc("/answers/{id}", h.GetAnswer).Methods("GET")
	router.HandleFunc("/answers/{id}", h.DeleteAnswer).Methods("DELETE")

	return router
}

// Health check handler - работает без service
func (h *Handler) healthCheck(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{
		"status":  "healthy",
		"service": "Q&A API",
	})
}

// Root handler - работает без service
func (h *Handler) rootHandler(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"message": "Q&A API Service",
		"version": "1.0.0",
	})
}

// GetQuestions - получить все вопросы
func (h *Handler) GetQuestions(w http.ResponseWriter, r *http.Request) {
	if h.service == nil {
		writeJSON(w, http.StatusOK, []interface{}{})
		return
	}

	questions, err := h.service.GetAllQuestions()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to get questions")
		return
	}

	writeJSON(w, http.StatusOK, questions)
}

// CreateQuestion - создать вопрос
func (h *Handler) CreateQuestion(w http.ResponseWriter, r *http.Request) {
	if h.service == nil {
		writeError(w, http.StatusServiceUnavailable, "Service not available")
		return
	}

	var req model.CreateQuestionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Text == "" {
		writeError(w, http.StatusBadRequest, "Question text is required")
		return
	}

	question, err := h.service.CreateQuestion(req)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to create question")
		return
	}

	writeJSON(w, http.StatusCreated, question)
}

// GetQuestion - получить вопрос по ID
func (h *Handler) GetQuestion(w http.ResponseWriter, r *http.Request) {
	if h.service == nil {
		writeError(w, http.StatusServiceUnavailable, "Service not available")
		return
	}

	id, err := getIDFromRequest(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "Invalid question ID")
		return
	}

	question, err := h.service.GetQuestion(id)
	if err != nil {
		writeError(w, http.StatusNotFound, "Question not found")
		return
	}

	writeJSON(w, http.StatusOK, question)
}

// DeleteQuestion - удалить вопрос
func (h *Handler) DeleteQuestion(w http.ResponseWriter, r *http.Request) {
	if h.service == nil {
		writeError(w, http.StatusServiceUnavailable, "Service not available")
		return
	}

	id, err := getIDFromRequest(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "Invalid question ID")
		return
	}

	if err := h.service.DeleteQuestion(id); err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to delete question")
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"message": "Question deleted successfully"})
}

// CreateAnswer - создать ответ
func (h *Handler) CreateAnswer(w http.ResponseWriter, r *http.Request) {
	if h.service == nil {
		writeError(w, http.StatusServiceUnavailable, "Service not available")
		return
	}

	questionID, err := getIDFromRequest(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "Invalid question ID")
		return
	}

	var req model.CreateAnswerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Text == "" || req.UserID == "" {
		writeError(w, http.StatusBadRequest, "Answer text and user ID are required")
		return
	}

	answer, err := h.service.CreateAnswer(questionID, req)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to create answer")
		return
	}

	writeJSON(w, http.StatusCreated, answer)
}

// GetAnswer - получить ответ по ID
func (h *Handler) GetAnswer(w http.ResponseWriter, r *http.Request) {
	if h.service == nil {
		writeError(w, http.StatusServiceUnavailable, "Service not available")
		return
	}

	id, err := getIDFromRequest(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "Invalid answer ID")
		return
	}

	answer, err := h.service.GetAnswer(id)
	if err != nil {
		writeError(w, http.StatusNotFound, "Answer not found")
		return
	}

	writeJSON(w, http.StatusOK, answer)
}

// DeleteAnswer - удалить ответ
func (h *Handler) DeleteAnswer(w http.ResponseWriter, r *http.Request) {
	if h.service == nil {
		writeError(w, http.StatusServiceUnavailable, "Service not available")
		return
	}

	id, err := getIDFromRequest(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "Invalid answer ID")
		return
	}

	if err := h.service.DeleteAnswer(id); err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to delete answer")
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"message": "Answer deleted successfully"})
}

// Utility functions
func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}

func getIDFromRequest(r *http.Request) (int, error) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		return 0, err
	}
	return id, nil
}
