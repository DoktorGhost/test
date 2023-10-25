package db

import (
	"database/sql"
	"fmt"
	"junior-test/pkg/types"
	"log"
)

type PersonRepository interface {
	CreatePerson(person *types.Person) (int, error)
	GetPersonByID(id int) (*types.Person, error)
	UpdatePerson(id int, person *types.Person) error
	DeletePerson(id int) error
	ListPeople(filter types.PersonFilter) ([]*types.Person, error)
}

// SQLPersonRepository представляет реализацию интерфейса PersonRepository для работы с PostgreSQL.
type SQLPersonRepository struct {
	DB *sql.DB
}

// NewSQLPersonRepository создает новый экземпляр SQLPersonRepository.
func NewSQLPersonRepository(db *sql.DB) *SQLPersonRepository {
	return &SQLPersonRepository{DB: db}
}

// Реализация методов интерфейса PersonRepository
func (r *SQLPersonRepository) CreatePerson(person *types.Person) (int, error) {
	query := "INSERT INTO people (name, surname, patronymic, age, gender, nationality) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id"
	var id int
	err := r.DB.QueryRow(query, person.Name, person.Surname, person.Patronymic, person.Age, person.Gender, person.Nationality).Scan(&id)
	if err != nil {
		log.Printf("Error while inserting a new person: %v", err)
		return 0, err
	}

	log.Printf("Successfully inserted a new person with ID %d", id)
	return id, nil
}

func (r *SQLPersonRepository) GetPersonByID(id int) (*types.Person, error) {
	query := "SELECT name, surname, patronymic, age, gender, nationality FROM people WHERE id = $1"

	row := r.DB.QueryRow(query, id)

	person := &types.Person{}
	err := row.Scan(&person.Name, &person.Surname, &person.Patronymic, &person.Age, &person.Gender, &person.Nationality)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("Person with ID %d not found", id)
			return nil, nil // Запись не найдена
		}
		log.Printf("Error retrieving person with ID %d: %v", id, err)
		return nil, err // Возникла ошибка при выполнении запроса
	}

	log.Printf("Retrieved person with ID %d", id)
	return person, nil
}

func (r *SQLPersonRepository) UpdatePerson(id int, person *types.Person) error {

	var count int
	err := r.DB.QueryRow("SELECT COUNT(*) FROM people WHERE id = $1", id).Scan(&count)
	if err != nil {
		log.Printf("Error checking existence of person with ID %d: %v", id, err)
		return err
	}

	if count == 0 {
		// Записи с указанным ID не существует
		log.Printf("Person with ID %d not found", id)
		return fmt.Errorf("Person with ID %d not found", id)
	}

	// Собираем SQL-запрос динамически на основе измененных полей
	query := "UPDATE people SET"
	var params []interface{}
	paramCount := 1

	if person.Name != "" {
		query += fmt.Sprintf(" name = $%d,", paramCount)
		params = append(params, person.Name)
		paramCount++
	}
	if person.Surname != "" {
		query += fmt.Sprintf(" surname = $%d,", paramCount)
		params = append(params, person.Surname)
		paramCount++
	}
	if person.Patronymic != "" {
		query += fmt.Sprintf(" patronymic = $%d,", paramCount)
		params = append(params, person.Patronymic)
		paramCount++
	}
	if person.Age != 0 {
		query += fmt.Sprintf(" age = $%d,", paramCount)
		params = append(params, person.Age)
		paramCount++
	}
	if person.Gender != "" {
		query += fmt.Sprintf(" gender = $%d,", paramCount)
		params = append(params, person.Gender)
		paramCount++
	}
	if person.Nationality != "" {
		query += fmt.Sprintf(" nationality = $%d,", paramCount)
		params = append(params, person.Nationality)
		paramCount++
	}

	// Удаляем последнюю запятую, если она есть
	if len(query) > 0 && query[len(query)-1] == ',' {
		query = query[:len(query)-1]
	}

	// Завершаем запрос
	query += fmt.Sprintf(" WHERE id = $%d", paramCount)
	params = append(params, id)

	// Выполняем SQL-запрос
	_, err = r.DB.Exec(query, params...)
	if err != nil {
		log.Printf("Error updating person with ID %d: %v", id, err)
		return err
	}

	log.Printf("Updated person with ID %d", id)
	return nil
}

func (r *SQLPersonRepository) DeletePerson(id int) error {
	query := "DELETE FROM people WHERE id = $1"

	result, err := r.DB.Exec(query, id)
	if err != nil {
		log.Printf("Error deleting person with ID %d: %v", id, err)
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		log.Printf("Not ID %d", id)
		return sql.ErrNoRows // Возвращаем ошибку sql.ErrNoRows, если ни одна запись не была удалена
	}

	log.Printf("Deleted person with ID %d", id)
	return nil
}

func (r *SQLPersonRepository) FilterListPeople(filter types.PersonFilter) ([]*types.Person, error) {
	query := "SELECT name, surname, patronymic, age, gender, nationality FROM people WHERE 1=1"
	var args []interface{}
	paramCount := 1

	if filter.Name != "" {
		query += fmt.Sprintf(" AND LOWER(name) = LOWER($%d)", paramCount)
		args = append(args, filter.Name)
		paramCount++
	}
	if filter.Surname != "" {
		query += fmt.Sprintf(" AND LOWER(surname) = LOWER($%d)", paramCount)
		args = append(args, filter.Surname)
		paramCount++
	}
	if filter.Patronymic != "" {
		query += fmt.Sprintf(" AND LOWER(patronymic) = LOWER($%d)", paramCount)
		args = append(args, filter.Patronymic)
		paramCount++
	}
	if filter.MinAge != 0 {
		query += fmt.Sprintf(" AND age >= $%d", paramCount)
		args = append(args, filter.MinAge)
		paramCount++
	}
	if filter.MaxAge != 0 {
		query += fmt.Sprintf(" AND age <= $%d", paramCount)
		args = append(args, filter.MaxAge)
		paramCount++
	}
	if filter.Gender != "" {
		query += fmt.Sprintf(" AND gender = $%d", paramCount)
		args = append(args, filter.Gender)
		paramCount++
	}
	if filter.Nationality != "" {
		query += fmt.Sprintf(" AND LOWER(nationality) = LOWER($%d)", paramCount)
		args = append(args, filter.Nationality)
		paramCount++
	}

	if filter.Page > 0 && filter.PageSize > 0 {
		query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", paramCount, (paramCount + 1))
		args = append(args, filter.PageSize, (filter.Page-1)*filter.PageSize)
	}

	log.Printf("Executing query: %s, args: %v", query, args)

	rows, err := r.DB.Query(query, args...)
	if err != nil {
		log.Printf("Error querying people: %v", err)
		return nil, err
	}
	defer rows.Close()

	var people []*types.Person
	for rows.Next() {
		var p types.Person
		err := rows.Scan(&p.Name, &p.Surname, &p.Patronymic, &p.Age, &p.Gender, &p.Nationality)
		if err != nil {
			log.Printf("Error scanning person: %v", err)
			return nil, err
		}
		people = append(people, &p)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Error iterating through rows: %v", err)
		return nil, err
	}

	log.Printf("Query executed successfully. Retrieved %d records.", len(people))

	return people, nil
}
