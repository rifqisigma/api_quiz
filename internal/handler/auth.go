package handler

import (
	"api_quiz/dto"
	"api_quiz/internal/usecase"
	"api_quiz/utils/helper"
	"api_quiz/utils/middleware"
	"encoding/json"
	"net/http"
)

type AuthHandler struct {
	authUC usecase.AuthUseCase
}

func NewAuthHandler(authUC usecase.AuthUseCase) *AuthHandler {
	return &AuthHandler{authUC}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var input dto.Register

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		helper.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}
	if input.Email == "" || input.Username == "" || input.Password == "" {
		helper.WriteError(w, http.StatusBadRequest, "invalid body")
		return
	}

	if err := h.authUC.Register(&input); err != nil {
		switch err {
		case helper.ErrInvalidEmail:
			helper.WriteError(w, http.StatusBadRequest, "invalid type email")
			return
		default:
			helper.WriteError(w, http.StatusInternalServerError, err.Error())
			return
		}
	}

	helper.WriteJSON(w, http.StatusOK, map[string]string{
		"message": "link verification has been send or your email",
	})

}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var input dto.Login
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		helper.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}
	if input.Email == "" || input.Password == "" {
		helper.WriteError(w, http.StatusBadRequest, "invalid body")
		return
	}
	token, err := h.authUC.Login(&input)
	if err != nil {
		switch err {
		case helper.ErrInvalidEmail:
			helper.WriteError(w, http.StatusBadRequest, "invalid type email")
			return
		default:
			helper.WriteError(w, http.StatusInternalServerError, err.Error())
			return
		}
	}

	helper.WriteJSON(w, http.StatusOK, map[string]string{
		"token_jwt": token,
	})

}

func (h *AuthHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(middleware.UserContextKey).(*helper.JWTClaims)
	if !ok {
		helper.WriteError(w, http.StatusUnauthorized, "")
		return
	}

	if err := h.authUC.DeleteUser(claims.UserID); err != nil {
		helper.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	helper.WriteJSON(w, http.StatusOK, map[string]string{
		"message": "succed delete this account",
	})
}

func (h *AuthHandler) Verification(w http.ResponseWriter, r *http.Request) {
	tokenstring := r.URL.Query().Get("token")

	claims, err := helper.ParseJWT(tokenstring)
	if err != nil {
		helper.WriteError(w, http.StatusInternalServerError, "failed verified jwt user")
		return
	}
	if err := h.authUC.ValidateUser(claims.UserID); err != nil {
		helper.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	helper.WriteJSON(w, http.StatusOK, map[string]string{
		"message": "link berhasil terkirim",
	})
}
