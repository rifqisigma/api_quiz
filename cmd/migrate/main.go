package main

import (
	"api_quiz/cmd/database"
	"api_quiz/entity"
	"log"
)

func main() {
	database.ConnectDB()

	if database.DB == nil {
		log.Fatal("❌ Database belum diinisialisasi")
	}

	err := database.DB.AutoMigrate(&entity.User{}, &entity.Quiz{}, &entity.Question{}, &entity.Answer{}, &entity.Submission{}, &entity.SubmissionUserAnswer{})
	if err != nil {
		log.Fatalf("gagal migrasi boy %v", err)
	}

	log.Println("berhasil migrasi")
}
