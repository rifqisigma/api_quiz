package main

import (
	"api_quiz/cmd/database"
	"api_quiz/cmd/route"
	"api_quiz/internal/handler"
	"api_quiz/internal/repository"
	"api_quiz/internal/usecase"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/lpernett/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("⚠ No .env file found, using system environment variables")
	}

	database.ConnectDB()
	//auth
	authRepo := repository.NewAuthRepository(database.DB)
	authUsecase := usecase.NewAuthUseCase(authRepo)
	authHandler := handler.NewAuthHandler(authUsecase)

	quizRepo := repository.NewQuizRepository(database.DB)
	quizUsecase := usecase.NewQuizUseCase(quizRepo)
	quizHandler := handler.NewQuizHandler(quizUsecase)

	submissionRepo := repository.NewSubmissionRepository(database.DB)
	submissionUseCase := usecase.NewSubmissionUseCase(submissionRepo, quizRepo)
	submissionHandler := handler.NewSubmissionHandler(submissionUseCase)

	r := route.SetupRoutes(authHandler, quizHandler, submissionHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Println("🚀 Server running on http://localhost:" + port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
