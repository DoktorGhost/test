package handlers

import (
	"encoding/json"
	"fmt"
	db "junior-test/db/models"
	"junior-test/pkg/types"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type PeopleHandler struct {
	Repository *db.SQLPersonRepository
}

func (h *PeopleHandler) EnrichPerson(c *gin.Context) {
	var person types.Person

	if err := c.BindJSON(&person); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to decode JSON request"})
		return
	}

	// Создаем каналы для получения результатов
	ageCh := make(chan int)
	genderCh := make(chan string)
	nationalityCh := make(chan string)
	errCh := make(chan error)

	// Запускаем горутины для каждого запроса
	go func() {
		age, err := h.getAge(person.Name)
		if err != nil {
			errCh <- err
			return
		}
		ageCh <- age
	}()

	go func() {
		gender, err := h.getGender(person.Name)
		if err != nil {
			errCh <- err
			return
		}
		genderCh <- gender
	}()

	go func() {
		nationality, err := h.getNationality(person.Name)
		if err != nil {
			errCh <- err
			return
		}
		nationalityCh <- nationality
	}()

	// Ожидание результатов из горутины
	var Err error

	for i := 0; i < 3; i++ {
		select {
		case age := <-ageCh:
			person.Age = age
		case gender := <-genderCh:
			person.Gender = gender
		case nationality := <-nationalityCh:
			person.Nationality = nationality
		case Err = <-errCh:
			break
		}
	}

	if Err == nil {
		_, err := h.Repository.CreatePerson(&person)
		if err != nil {
			log.Printf("Error CreatePerson: %v", err)
		}
	} else {
		log.Printf("Error: %v", Err)
	}

	c.JSON(http.StatusOK, gin.H{"message": "Person enriched and created"})
}

func (h *PeopleHandler) getAge(name string) (int, error) {
	// Формирование URL для запроса к API Agify
	url := "https://api.agify.io/?name=" + name

	// Отправка HTTP GET-запроса
	resp, err := http.Get(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	// Обработка ответа
	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("API request failed with status: %s", resp.Status)
	}

	// Распарсинг JSON-ответа
	var data types.AgifyResponse

	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return 0, err
	}

	// Проверка наличия возраста и корректного имени
	if data.Count == 0 {
		return 0, fmt.Errorf("Name is incorrect")
	}

	// Возврат возраста
	return data.Age, nil
}

func (h *PeopleHandler) getGender(name string) (string, error) {
	// Формирование URL для запроса к API Genderize
	url := "https://api.genderize.io/?name=" + name

	// Отправка HTTP GET-запроса
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Обработка ответа
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API request failed with status: %s", resp.Status)
	}

	// Распарсинг JSON-ответа
	var data types.GenderizeResponse

	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return "", err
	}

	// Проверка наличия пола и корректного имени
	if data.Count == 0 {
		return "", fmt.Errorf("Name is incorrect")
	}

	// Возврат пола
	return data.Gender, nil
}

func (h *PeopleHandler) getNationality(name string) (string, error) {
	// Формирование URL для запроса к API Nationalize
	url := "https://api.nationalize.io/?name=" + name

	// Отправка HTTP GET-запроса
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Обработка ответа
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API request failed with status: %s", resp.Status)
	}

	// Распарсинг JSON-ответа
	var data types.NationalizeResponse

	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return "", err
	}

	// Проверка наличия национальности и корректного имени
	if data.Count == 0 {
		return "", fmt.Errorf("Name is incorrect")
	}

	// Возврат национальности (берем национальность с наибольшей вероятностью)
	maxProbability := 0.0
	var nationality string
	for _, country := range data.Country {
		if country.Probability > maxProbability {
			maxProbability = country.Probability
			nationality = country.CountryID
		}
	}

	return nationality, nil
}

func (h *PeopleHandler) GetPersonByID(c *gin.Context) {
	id := c.Param("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	person, err := h.Repository.GetPersonByID(idInt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get person"})
		return
	}

	c.JSON(http.StatusOK, person)
}

func (h *PeopleHandler) UpdatePerson(c *gin.Context) {
	// Извлекаем ID из URL.
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	// Извлекаем данные, которые необходимо обновить, из запроса.
	var updatedPerson types.Person
	if err := c.BindJSON(&updatedPerson); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to decode JSON request"})
		return
	}

	// Вызываем функцию репозитория для обновления информации о человеке по ID.
	err = h.Repository.UpdatePerson(id, &updatedPerson)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update person"})
		return
	}

	// Отправляем успешный ответ.
	c.JSON(http.StatusOK, gin.H{"message": "Person updated successfully"})
}

func (h *PeopleHandler) DeletePerson(c *gin.Context) {
	// Извлекаем ID из URL.
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	// Вызываем функцию репозитория для удаления информации о человеке по ID.
	err = h.Repository.DeletePerson(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete person"})
		return
	}

	// Отправляем успешный ответ.
	c.JSON(http.StatusOK, gin.H{"message": "Person deleted successfully"})
}

func (h *PeopleHandler) FilterListPeople(c *gin.Context) {
	// Парсинг параметров запроса, таких как фильтр и параметры пагинации
	var filter types.PersonFilter

	if err := c.BindJSON(&filter); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to decode JSON request"})
		return
	}

	// Вызов функции репозитория для получения списка людей
	people, err := h.Repository.FilterListPeople(filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve people"})
		return
	}

	// Отправка списка людей в формате JSON
	c.JSON(http.StatusOK, people)
}
