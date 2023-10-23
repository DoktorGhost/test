package handlers

import (
	"encoding/json"
	"fmt"
	db "junior-test/db/models"
	"junior-test/pkg/types"
	"log"
	"net/http"
)

type PeopleHandler struct {
	Repository *db.SQLPersonRepository
}

func (h *PeopleHandler) EnrichPerson(w http.ResponseWriter, r *http.Request) {
	var person types.Person

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&person); err != nil {
		http.Error(w, "Failed to decode JSON request", http.StatusBadRequest)
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
