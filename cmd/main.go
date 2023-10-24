package main

import (
	"database/sql"
	"junior-test/api/handlers"
	"junior-test/api/routes"
	"junior-test/db"
	dbModels "junior-test/db/models"
	"log"
	"net/http"
	"os"
	"time"

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
	var database *sql.DB
	var dbErr error

	// Логика попыток повторного подключения к базе данных
	for retries := 1; retries <= 5; retries++ {
		database, dbErr = db.InitDB()
		if dbErr == nil {
			defer db.CloseDB() // Закрываем соединение при завершении работы
			break
		}

		log.Printf("Failed to initialize database (attempt %d): %v", retries, dbErr)
		time.Sleep(5 * time.Second)
	}

	if dbErr != nil {
		log.Fatal("Failed to initialize database after 5 attempts")
	}

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
