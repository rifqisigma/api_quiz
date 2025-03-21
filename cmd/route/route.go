package route

import (
	"api_quiz/internal/handler"
	"api_quiz/utils/middleware"
	"net/http"

	"github.com/gorilla/mux"
)

func SetupRoutes(authHandler *handler.AuthHandler) *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/register", authHandler.Register).Methods(http.MethodPost)
	r.HandleFunc("/login", authHandler.Login).Methods(http.MethodPost)
	r.HandleFunc("/verification", authHandler.Verification).Methods(http.MethodGet)

	userRoute := r.PathPrefix("/user").Subrouter()
	userRoute.Use(middleware.JWTAuthMiddleware)

	userRoute.HandleFunc("delete", authHandler.DeleteUser).Methods(http.MethodDelete)

	return r

}
