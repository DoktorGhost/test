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
	ListPeople(filter types.PersonFilter, pagination types.Pagination) ([]*types.Person, error)
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
	query := "UPDATE people SET name = $1, surname = $2, patronymic = $3, age = $4, gender = $5, nationality = $6 WHERE id = $7"

	_, err := r.DB.Exec(query, person.Name, person.Surname, person.Patronymic, person.Age, person.Gender, person.Nationality, id)
	if err != nil {
		log.Printf("Error updating person with ID %d: %v", id, err)
		return err
	}

	log.Printf("Updated person with ID %d", id)
	return nil
}

func (r *SQLPersonRepository) DeletePerson(id int) error {
	query := "DELETE FROM people WHERE id = $1"

	_, err := r.DB.Exec(query, id)
	if err != nil {
		log.Printf("Error deleting person with ID %d: %v", id, err)
		return err
	}

	log.Printf("Deleted person with ID %d", id)
	return nil
}

func (r *SQLPersonRepository) ListPeople(filter types.PersonFilter, pagination types.Pagination) ([]*types.Person, error) {
	query := "SELECT name, surname, patronymic, age, gender, nationality FROM people WHERE 1=1"
	var args []interface{}
	paramCount := 1

	if filter.Name != "" {
		query += fmt.Sprintf(" AND name = $%d", paramCount)
		args = append(args, filter.Name)
		paramCount++
	}
	if filter.Surname != "" {
		query += fmt.Sprintf(" AND surname = $%d", paramCount)
		args = append(args, filter.Surname)
		paramCount++
	}
	if filter.Patronymic != "" {
		query += fmt.Sprintf(" AND patronymic = $%d", paramCount)
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
		query += fmt.Sprintf(" AND nationality = $%d", paramCount)
		args = append(args, filter.Nationality)
		paramCount++
	}

	if pagination.Initialized == true {
		query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", paramCount, (paramCount + 1))
		args = append(args, pagination.PageSize, (pagination.Page-1)*pagination.PageSize)
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
