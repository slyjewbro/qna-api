package handler

import (
	"encoding/json"
	"net/http"
	"qna-api/internal/model"
)

func (h *Handler) GetQuestions(w http.ResponseWriter, r *http.Request) {
	questions, err := h.service.GetAllQuestions()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to get questions")
		return
	}

	writeJSON(w, http.StatusOK, questions)
}

func (h *Handler) GetQuestion(w http.ResponseWriter, r *http.Request) {
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

func (h *Handler) CreateQuestion(w http.ResponseWriter, r *http.Request) {
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

func (h *Handler) DeleteQuestion(w http.ResponseWriter, r *http.Request) {
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
