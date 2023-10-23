package types

// Модель для таблицы "people".
type Person struct {
	Name        string `json:"name"`
	Surname     string `json:"surname"`
	Patronymic  string `json:"patronymic"`
	Age         int    `json:"age"`
	Gender      string `json:"gender"`
	Nationality string `json:"nationality"`
}

// Условия фильтрации записей.
type PersonFilter struct {
	Name        string
	Surname     string
	Patronymic  string
	MinAge      int
	MaxAge      int
	Gender      string
	Nationality string
}

// Параметры пагинации.
type Pagination struct {
	Page        int
	PageSize    int
	Initialized bool
}

// Ответ от api ожидаемый возраст
type AgifyResponse struct {
	Count int `json:"count"`
	Age   int `json:"age"`
}

// Ответ от api ожидаемый пол
type GenderizeResponse struct {
	Count  int    `json:"count"`
	Gender string `json:"gender"`
}

// Ответ от api ожидаемая национальность
type NationalizeResponse struct {
	Count   int           `json:"count"`
	Name    string        `json:"name"`
	Country []Nationality `json:"country"`
}

type Nationality struct {
	CountryID   string  `json:"country_id"`
	Probability float64 `json:"probability"`
}
