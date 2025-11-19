package handler

import (
	"encoding/json"
	"net/http"
	"qna-api/internal/model"
)

func (h *Handler) CreateAnswer(w http.ResponseWriter, r *http.Request) {
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

func (h *Handler) GetAnswer(w http.ResponseWriter, r *http.Request) {
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

func (h *Handler) DeleteAnswer(w http.ResponseWriter, r *http.Request) {
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
