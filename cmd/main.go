package main

import (
	"junior-test/api/handlers"
	"junior-test/api/routes"
	"junior-test/db"
	dbModels "junior-test/db/models"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	var envFilePath string

	if _, err := os.Stat("/app/.env"); err == nil {
		// Файл .env найден внутри контейнера
		envFilePath = "/app/.env"
	} else if _, err := os.Stat("../.env"); err == nil {
		// Файл .env найден в корне проекта (для локальной разработки)
		envFilePath = "../.env"
	} else {
		log.Fatal("Error: .env file not found")
	}

	if err := godotenv.Load(envFilePath); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	httpPort := os.Getenv("HTTP_PORT")
	// Инициализация базы данных
	database, err := db.InitDB()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.CloseDB()

	// Создание экземпляра репозитория
	repository := dbModels.NewSQLPersonRepository(database)

	// Создание экземпляра PeopleHandler с передачей репозитория
	handler := &handlers.PeopleHandler{
		Repository: repository,
	}

	// Запуск сервера
	routes.RunServer(handler)
	log.Fatal(http.ListenAndServe(":"+httpPort, nil))

}
