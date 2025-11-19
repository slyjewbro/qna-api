package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"qna-api/internal/service"

	"github.com/gorilla/mux"
)

type Handler struct {
	service *service.Service
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) InitRoutes() *mux.Router {
	router := mux.NewRouter()

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
