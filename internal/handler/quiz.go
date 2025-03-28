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

type QuizHandler struct {
	quizUC usecase.QuizUseCase
}

func NewQuizHandler(quizUC usecase.QuizUseCase) *QuizHandler {
	return &QuizHandler{quizUC}
}

func (h *QuizHandler) GetAllQuiz(w http.ResponseWriter, r *http.Request) {
	_, ok := r.Context().Value(middleware.UserContextKey).(*helper.JWTClaims)
	if !ok {
		helper.WriteError(w, http.StatusUnauthorized, "not token provide")
		return
	}

	response, err := h.quizUC.GetAllQuiz()
	if err != nil {
		helper.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	helper.WriteJSON(w, http.StatusOK, response)
}

func (h *QuizHandler) GetQuizById(w http.ResponseWriter, r *http.Request) {
	_, ok := r.Context().Value(middleware.UserContextKey).(*helper.JWTClaims)
	if !ok {
		helper.WriteError(w, http.StatusUnauthorized, "not token provide")
		return
	}

	params := mux.Vars(r)
	quizId, _ := strconv.Atoi(params["quizid"])

	response, err := h.quizUC.GetQuizFromId(uint(quizId))
	if err != nil {
		helper.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	helper.WriteJSON(w, http.StatusOK, response)
}

func (h *QuizHandler) CreateQuiz(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(middleware.UserContextKey).(*helper.JWTClaims)
	if !ok {
		helper.WriteError(w, http.StatusUnauthorized, "not token provide")
		return
	}

	var input dto.Quiz
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		helper.WriteError(w, http.StatusBadRequest, "invalid body")
		return
	}

	if input.Title == "" {
		helper.WriteError(w, http.StatusBadRequest, "invalid body")
		return
	}

	input.Creator = claims.UserID

	response, err := h.quizUC.CreateQuiz(&input)
	if err != nil {
		helper.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	helper.WriteJSON(w, http.StatusCreated, response)

}

func (h *QuizHandler) UpdateQuiz(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(middleware.UserContextKey).(*helper.JWTClaims)
	if !ok {
		helper.WriteError(w, http.StatusUnauthorized, "not token provide")
		return
	}

	params := mux.Vars(r)
	quizId, _ := strconv.Atoi(params["quizid"])

	var input dto.UpdatedQuiz
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		helper.WriteError(w, http.StatusBadRequest, "invalid body")
		return
	}

	if input.Title == "" {
		helper.WriteError(w, http.StatusBadRequest, "invalid body")
		return
	}

	input.ID = uint(quizId)
	response, err := h.quizUC.UpdateQuiz(&input, claims.UserID)
	if err != nil {
		switch err {
		case helper.ErrUnauhorized:
			helper.WriteError(w, http.StatusUnauthorized, err.Error())
		default:
			helper.WriteError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	helper.WriteJSON(w, http.StatusOK, response)
}

func (h *QuizHandler) DeleteQuiz(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(middleware.UserContextKey).(*helper.JWTClaims)
	if !ok {
		helper.WriteError(w, http.StatusUnauthorized, "not token provide")
		return
	}

	params := mux.Vars(r)
	quizId, _ := strconv.Atoi(params["quizid"])

	if err := h.quizUC.DeleteQuiz(claims.UserID, uint(quizId)); err != nil {
		switch err {
		case helper.ErrUnauhorized:
			helper.WriteError(w, http.StatusUnauthorized, err.Error())
		default:
			helper.WriteError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	helper.WriteJSON(w, http.StatusOK, map[string]string{
		"message": "succed delete this quiz",
	})

}

// question
func (h *QuizHandler) GetQuestionAnswerByQuizId(w http.ResponseWriter, r *http.Request) {
	_, ok := r.Context().Value(middleware.UserContextKey).(*helper.JWTClaims)
	if !ok {
		helper.WriteError(w, http.StatusUnauthorized, "not token provide")
		return
	}

	params := mux.Vars(r)
	quizId, _ := strconv.Atoi(params["quizid"])

	response, err := h.quizUC.GetQuestionAnswerByQuizId(uint(quizId))
	if err != nil {
		helper.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	helper.WriteJSON(w, http.StatusOK, response)
}

func (h *QuizHandler) GetQuestionById(w http.ResponseWriter, r *http.Request) {
	_, ok := r.Context().Value(middleware.UserContextKey).(*helper.JWTClaims)
	if !ok {
		helper.WriteError(w, http.StatusUnauthorized, "not token provide")
		return
	}

	params := mux.Vars(r)
	questionId, _ := strconv.Atoi(params["questionid"])
	quizId, _ := strconv.Atoi(params["quizid"])

	response, err := h.quizUC.GetQuestionById(uint(questionId), uint(quizId))
	if err != nil {
		helper.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	helper.WriteJSON(w, http.StatusOK, response)
}

func (h *QuizHandler) CreateQuestionAndAnswer(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(middleware.UserContextKey).(*helper.JWTClaims)
	if !ok {
		helper.WriteError(w, http.StatusUnauthorized, "not token provide")
		return
	}

	params := mux.Vars(r)
	quizId, _ := strconv.Atoi(params["quizid"])

	var input dto.Question
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		helper.WriteError(w, http.StatusBadRequest, "invalid body")
		return
	}
	if input.Text == "" {
		helper.WriteError(w, http.StatusBadRequest, "question text is required")
		return
	}

	if len(input.Answers) <= 2 {
		helper.WriteError(w, http.StatusBadRequest, "at least one answer is required")
		return
	}
	if len(input.Answers) > 5 {
		helper.WriteError(w, http.StatusBadRequest, "max is 5")
		return
	}

	correctAnswer := 0
	for _, ans := range input.Answers {
		if ans.IsCorrect {
			correctAnswer++
		}
	}

	if correctAnswer != 1 {
		helper.WriteError(w, http.StatusBadRequest, "answer must 1")
		return
	}

	input.QuizID = uint(quizId)

	response, err := h.quizUC.CreateQuestionAndAnswer(&input, claims.UserID)
	if err != nil {
		switch err {
		case helper.ErrUnauhorized:
			helper.WriteError(w, http.StatusUnauthorized, err.Error())
		default:
			helper.WriteError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	helper.WriteJSON(w, http.StatusCreated, response)

}
func (h *QuizHandler) UpdateQuestion(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(middleware.UserContextKey).(*helper.JWTClaims)
	if !ok {
		helper.WriteError(w, http.StatusUnauthorized, "not token provide")
		return
	}

	params := mux.Vars(r)
	quizId, _ := strconv.Atoi(params["quizid"])
	questionId, _ := strconv.Atoi(params["questionid"])

	var input dto.QuestionUpdate
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		helper.WriteError(w, http.StatusBadRequest, "invalid body")
		return
	}
	if input.Text == "" {
		helper.WriteError(w, http.StatusBadRequest, "question text is required")
		return
	}

	input.ID = uint(questionId)
	input.QuizID = uint(quizId)
	response, err := h.quizUC.UpdateQuestion(&input, claims.UserID)
	if err != nil {
		switch err {
		case helper.ErrUnauhorized:
			helper.WriteError(w, http.StatusUnauthorized, err.Error())
		default:
			helper.WriteError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	helper.WriteJSON(w, http.StatusOK, response)

}
func (h *QuizHandler) DeleteQuestion(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(middleware.UserContextKey).(*helper.JWTClaims)
	if !ok {
		helper.WriteError(w, http.StatusUnauthorized, "not token provide")
		return
	}
	params := mux.Vars(r)
	quizId, _ := strconv.Atoi(params["quizid"])
	questionId, _ := strconv.Atoi(params["questionid"])

	if err := h.quizUC.DeleteQuestion(uint(questionId), uint(quizId), claims.UserID); err != nil {
		switch err {
		case helper.ErrUnauhorized:
			helper.WriteError(w, http.StatusUnauthorized, err.Error())
		default:
			helper.WriteError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	helper.WriteJSON(w, http.StatusOK, map[string]string{
		"message": "succed delete this question",
	})
}

// answer
func (h *QuizHandler) GetAnswerByQuestionId(w http.ResponseWriter, r *http.Request) {
	_, ok := r.Context().Value(middleware.UserContextKey).(*helper.JWTClaims)
	if !ok {
		helper.WriteError(w, http.StatusUnauthorized, "not token provide")
		return
	}

	params := mux.Vars(r)
	questionId, _ := strconv.Atoi(params["questionid"])

	response, err := h.quizUC.GetAnswerByQuestionId(uint(questionId))
	if err != nil {
		helper.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	helper.WriteJSON(w, http.StatusOK, response)
}

func (h *QuizHandler) UpdateAnswer(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(middleware.UserContextKey).(*helper.JWTClaims)
	if !ok {
		helper.WriteError(w, http.StatusUnauthorized, "not token provide")
		return
	}
	params := mux.Vars(r)
	questionId, _ := strconv.Atoi(params["questionid"])
	answerId, _ := strconv.Atoi(params["answerid"])
	quizId, _ := strconv.Atoi(params["quizid"])

	var input dto.Answer
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		helper.WriteError(w, http.StatusBadRequest, "invalid body")
		return
	}
	if input.Text == "" {
		helper.WriteError(w, http.StatusBadRequest, "answers cannot be empty")
		return
	}

	input.ID = uint(answerId)
	input.QuestionID = uint(questionId)
	response, err := h.quizUC.UpdateAnswer(claims.UserID, uint(quizId), input)
	if err != nil {
		switch err {
		case helper.ErrUnauhorized:
			helper.WriteError(w, http.StatusUnauthorized, err.Error())
		case helper.ErrCorrectAnswer:
			helper.WriteError(w, http.StatusUnprocessableEntity, err.Error())
		default:
			helper.WriteError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}
	helper.WriteJSON(w, http.StatusOK, response)
}

func (h *QuizHandler) AddAnswer(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(middleware.UserContextKey).(*helper.JWTClaims)
	if !ok {
		helper.WriteError(w, http.StatusUnauthorized, "not token provide")
		return
	}
	params := mux.Vars(r)
	questionId, _ := strconv.Atoi(params["questionid"])
	quizId, _ := strconv.Atoi(params["quizid"])

	var input []dto.Answer
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		helper.WriteError(w, http.StatusBadRequest, "invalid body")
		return
	}
	if len(input) > 5 {
		helper.WriteError(w, http.StatusBadRequest, "max 5 answer to request")
		return
	}

	for i := range input {
		input[i].QuestionID = uint(questionId)
	}

	response, err := h.quizUC.AddAnswer(claims.UserID, uint(quizId), input)
	if err != nil {
		switch err {
		case helper.ErrUnauhorized:
			helper.WriteError(w, http.StatusUnauthorized, err.Error())
		case helper.ErrToomuchAnswer:
			helper.WriteError(w, http.StatusUnprocessableEntity, err.Error())
		case helper.ErrCorrectAnswer:
			helper.WriteError(w, http.StatusUnprocessableEntity, err.Error())
		default:
			helper.WriteError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}
	helper.WriteJSON(w, http.StatusOK, response)
}

func (h *QuizHandler) DeleteAnswer(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(middleware.UserContextKey).(*helper.JWTClaims)
	if !ok {
		helper.WriteError(w, http.StatusUnauthorized, "not token provide")
		return
	}
	params := mux.Vars(r)
	questionId, _ := strconv.Atoi(params["questionid"])
	answerId, _ := strconv.Atoi(params["answerid"])
	quizId, _ := strconv.Atoi(params["quizid"])

	if err := h.quizUC.DeleteAnswer(uint(answerId), uint(questionId), uint(quizId), claims.UserID); err != nil {
		switch err {
		case helper.ErrUnauhorized:
			helper.WriteError(w, http.StatusUnauthorized, err.Error())
		case helper.ErrAnswerNotEnough:
			helper.WriteError(w, http.StatusUnprocessableEntity, err.Error())
		default:
			helper.WriteError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	helper.WriteJSON(w, http.StatusOK, map[string]string{
		"message": "succed delete this answer",
	})
}
