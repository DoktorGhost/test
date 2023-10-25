## Тестовое задание
Реализовать сервис, который будет получать по апи ФИО, из открытых апи обогащать ответ наиболее вероятными возрастом, полом и национальностью и сохранять данные в БД. По запросу выдавать инфу о найденных людях. Необходимо реализовать следующее:
### 1. Выставить rest методы
1) Для получения данных с различными фильтрами и пагинацией
2) Для удаления по идентификатору
3) Для изменения сущности
4) Для добавления новых людей в формате:
``` json
{
    "name": "Dmitriy",
    "surname": "Ushakov",
    "patronymic": "Vasilevich" // необязательно
}

```
### 2 Корректное сообщение обогатить
1. Возрастом - https://api.agify.io/?name=Dmitriy
2. Полом - https://api.genderize.io/?name=Dmitriy
3. Национальностью - https://api.nationalize.io/?name=Dmitriy
3. Обогащенное сообщение положить в БД postgres (структура БД должна быть создана
путем миграций)
4. Покрыть код debug- и info-логами
5. Вынести конфигурационные данные в .env


# Запуск приложения
В репозитории есть докер-компос файл, поэтому с помощью него поднимаем необходимые контейнеры:
``` golang
docker-compose up
```
Автоматически поднимается база данных postgres (миграция прописана в докерфайле) и само приложение.
Идем по порядку задач, тестируем через Postman
1) метод для получения данных с различными фильтрами и пагинацией:
Выбираем метод GET, роут выставлен по адресу http://localhost:8080 **/filter**. Для фильтрации есть отдельный тип:
``` golang
type PersonFilter struct {
	Name        string `json:"name"`
	Surname     string `json:"surname"`
	Patronymic  string `json:"patronymic"`
	MinAge      int    `json:"minage"`
	MaxAge      int    `json:"maxage"`
	Gender      string `json:"gender"`
	Nationality string `json:"nationality"`
	Page        int    `json:"page"`
	PageSize    int    `json:"pagesize"`
}
```
Заполняем необходимые поля в теле запроса, пагинация так же прописана в PersonFilter: Page и PageSize.  

2. метод для удаления по идентификатору:
Выбираем метод DELETE, роут http://localhost:8080 **/delete/:id** . В теле запроса ничего передавать не нужно, необходимый ID передаем напрямую в url

3. метод для изменения сущности:
Выбираем метод PUT, роут http://localhost:8080 **/update/:id** . В теле запроса нужно передать структуру, с полями, которые нужно отредактировать, пустые поля из структуры не будут обновляться в бд.
``` golang
type Person struct {
	Name        string `json:"name"`
	Surname     string `json:"surname"`
	Patronymic  string `json:"patronymic"`
	Age         int    `json:"age"`
	Gender      string `json:"gender"`
	Nationality string `json:"nationality"`
}
```
Структура для обновления та же самая, что и для записи в бд.

4. метод для добавления новых людей:
Выбираем метод POST, роут http://localhost:8080 **/add**. В теле запроса передаем структуру с необходимыми данными (ФИО), тут же происходит дополнение информации.

Файл .env для нашего примера содержит следующие данные:
``` golang
DB_USER=admin
DB_PASSWORD=admin
DB_NAME=people_database
DB_HOST=composepostgres //меняем на "localhost" при запуске на локальной машине, вне контейнера
DB_PORT=5432
HTTP_PORT=8080
```




