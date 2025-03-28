package handler

import (
	"api_quiz/dto"
	"api_quiz/internal/usecase"
	"api_quiz/utils/helper"
	"api_quiz/utils/middleware"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type SubmissionHandler struct {
	submissionUC usecase.SubmissionUseCase
}

func NewSubmissionHandler(submissionUC usecase.SubmissionUseCase) *SubmissionHandler {
	return &SubmissionHandler{submissionUC}
}

func (h *SubmissionHandler) GetAllSubmission(w http.ResponseWriter, r *http.Request) {
	_, ok := r.Context().Value(middleware.UserContextKey).(*helper.JWTClaims)
	if !ok {
		helper.WriteError(w, http.StatusUnauthorized, "not token provide")
		return
	}

	response, err := h.submissionUC.GetAllSubmission()
	if err != nil {
		helper.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	helper.WriteJSON(w, http.StatusOK, response)
}
func (h *SubmissionHandler) GetSubmissionById(w http.ResponseWriter, r *http.Request) {
	_, ok := r.Context().Value(middleware.UserContextKey).(*helper.JWTClaims)
	if !ok {
		helper.WriteError(w, http.StatusUnauthorized, "not token provide")
		return
	}

	params := mux.Vars(r)
	submissionId, _ := strconv.Atoi(params["submissionid"])

	response, err := h.submissionUC.GetSubmissionById(uint(submissionId))
	if err != nil {
		helper.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	helper.WriteJSON(w, http.StatusOK, response)
}

func (h *SubmissionHandler) CreateSubmission(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(middleware.UserContextKey).(*helper.JWTClaims)
	if !ok {
		helper.WriteError(w, http.StatusUnauthorized, "not token provide")
		return
	}

	params := mux.Vars(r)
	quizId, _ := strconv.Atoi(params["quizid"])

	var input dto.Submission

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		helper.WriteError(w, http.StatusBadRequest, "invalid body")
		return
	}
	input.QuizID = uint(quizId)
	input.UserID = claims.UserID

	response, err := h.submissionUC.CreateSubmission(&input)
	if err != nil {
		helper.WriteError(w, http.StatusInternalServerError, err.Error())
	}

	helper.WriteJSON(w, http.StatusCreated, response)
}
func (h *SubmissionHandler) UpdateSubmission(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(middleware.UserContextKey).(*helper.JWTClaims)
	if !ok {
		helper.WriteError(w, http.StatusUnauthorized, "not token provide")
		return
	}

	params := mux.Vars(r)
	submisionId, _ := strconv.Atoi(params["submissionid"])

	var input dto.SubmissionUpdate

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		helper.WriteError(w, http.StatusBadRequest, "invalid body")
		return
	}

	input.SubmissionID = uint(submisionId)
	response, err := h.submissionUC.UpdateSubmision(&input, claims.UserID)
	if err != nil {
		switch err {
		case helper.ErrQuizNotFound:
			helper.WriteError(w, http.StatusNotFound, err.Error())
			return
		case helper.ErrUnauhorized:
			helper.WriteError(w, http.StatusUnauthorized, err.Error())
			return
		default:
			helper.WriteError(w, http.StatusInternalServerError, err.Error())
			return

		}
	}

	helper.WriteJSON(w, http.StatusCreated, response)
}

func (h *SubmissionHandler) DeleteSubmission(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(middleware.UserContextKey).(*helper.JWTClaims)
	if !ok {
		helper.WriteError(w, http.StatusUnauthorized, "not token provide")
		return
	}

	params := mux.Vars(r)
	submisionId, _ := strconv.Atoi(params["submissionid"])

	if err := h.submissionUC.DeleteSubmision(uint(submisionId), claims.UserID); err != nil {
		switch err {
		case helper.ErrQuizNotFound:
			helper.WriteError(w, http.StatusNotFound, err.Error())
			return
		case helper.ErrUnauhorized:
			helper.WriteError(w, http.StatusUnauthorized, err.Error())
			return
		default:
			helper.WriteError(w, http.StatusInternalServerError, err.Error())
			return

		}
	}
}
