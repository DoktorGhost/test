package main

import (
	"junior-test/api/handlers"
	"junior-test/api/routes"
	"junior-test/db"
	dbModels "junior-test/db/models"
	"log"
	"net/http"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load("../.env"); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
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

	// Установка маршрутов
	routes.SetupPeopleRoutes(handler)

	// Запуск сервера
	log.Fatal(http.ListenAndServe(":8080", nil))
}
