package main

import (
	"api_quiz/cmd/database"
	"api_quiz/cmd/route"
	"api_quiz/entity"
	"api_quiz/internal/handler"
	"api_quiz/internal/repository"
	"api_quiz/internal/usecase"
	"log"
)

func main() {
	database.ConnectDB()

	//auth
	authRepo := repository.NewAuthRepository(database.DB)
	authUsecase := usecase.NewAuthUseCase(authRepo)
	authHandler := handler.NewAuthHandler(authUsecase)

	route.SetupRoutes(authHandler)

	if database.DB == nil {
		log.Fatal("❌ Database belum diinisialisasi")
	}

	err := database.DB.AutoMigrate(&entity.User{}, &entity.Quiz{}, &entity.Question{}, &entity.Answer{}, &entity.Submission{})
	if err != nil {
		log.Fatalf("gagal migrasi boy %v", err)
	}

	log.Println("berhasil migrasi")
}
