package routes

import (
	"junior-test/api/handlers"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
)

func RunServer(handler *handlers.PeopleHandler) {
	r := gin.Default()

	// Установка маршрутов
	SetupPeopleRoutes(r, handler)

	// Запуск сервера
	httpPort := os.Getenv("HTTP_PORT")
	r.Run(":" + httpPort)
}

func SetupPeopleRoutes(r *gin.Engine, handler *handlers.PeopleHandler) {

	// Маршрут для добавления человека
	r.POST("/add", handler.EnrichPerson)

	// Маршрут для получения информации о человеке по ID
	r.GET("/getperson/:id", func(c *gin.Context) {
		id := c.Param("id")
		idInt, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
			return
		}

		// Здесь вызываем функцию репозитория для получения информации о человеке по ID
		person, err := handler.Repository.GetPersonByID(idInt)

		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Person not found"})
			return
		}

		c.JSON(http.StatusOK, person)
	})

}
