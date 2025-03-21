package main

import (
	"api_quiz/cmd/database"
	"fmt"
	"log"
	"os"

	"github.com/lpernett/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("⚠ No .env file found, using system environment variables")
	}

	database.ConnectDB()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Println("🚀 Server running on http://localhost:" + port)
	// log.Fatal(http.ListenAndServe(":"+port, r))
}
