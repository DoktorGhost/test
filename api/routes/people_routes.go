package routes

import (
	"junior-test/api/handlers"
	"net/http"
)

func SetupPeopleRoutes(handler *handlers.PeopleHandler) {
	http.HandleFunc("/api/add", handler.EnrichPerson)
}
