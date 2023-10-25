package db

import (
	"fmt"
	"junior-test/db"
	"junior-test/pkg/types"
	"log"
	"testing"

	"github.com/joho/godotenv"
)

func TestMain(t *testing.T) {

	if err := godotenv.Load("../../.env"); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	datebase, err := db.InitDB()
	if err != nil {
		t.Fatalf("Error initializing test database: %v", err)
	}
	defer db.CloseDB()

	repo := NewSQLPersonRepository(datebase)
	person := &types.Person{
		Name:    "Oleg",
		Surname: "Samsonov",
	}

	id, err := repo.CreatePerson(person)

	if err != nil {
		t.Fatalf("Expected no error, but got %v", err)
	}

	fmt.Println("!!!!!!!!!!!!!!!! Test CreatePerson succes !!!!!!!!!!!!!!!! ")

	person2, err := repo.GetPersonByID(id)

	if err != nil {
		t.Fatalf("Error retrieving person with ID %d: %v", id, err)
	}

	if person.Name != person2.Name && person.Surname != person2.Surname {
		t.Fatalf("%s != %s, %s != %s,", person.Name, person2.Name, person.Surname, person2.Surname)
	}

	fmt.Println("!!!!!!!!!!!!!!!! Test GetPersonByID succes !!!!!!!!!!!!!!!!")

	person3 := &types.Person{
		Name: "Olegsandr",
	}

	err = repo.UpdatePerson(id, person3)

	if err != nil {
		t.Fatalf("Error update: %v", err)
	}

	person2, _ = repo.GetPersonByID(id)

	if person3.Name != person2.Name {
		t.Fatalf("%s != %s, %s != %s,", person.Name, person2.Name, person.Surname, person2.Surname)
	}

	fmt.Println("!!!!!!!!!!!!!!!!  Test UpdatePerson succes !!!!!!!!!!!!!!!! ")

	err = repo.DeletePerson(id)
	if err != nil {
		t.Fatalf("Error deleting %v", err)
	}

	person2, err = repo.GetPersonByID(id)

	if person2 != nil {
		t.Fatalf("Error deleting %v", err)
	}

	fmt.Println("!!!!!!!!!!!!!!!! Test DeletePerson succes !!!!!!!!!!!!!!!!")
}

func TestFiltr(t *testing.T) {
	var persons = []*types.Person{
		{
			Name:    "Oleg",
			Surname: "Samsonov",
			Age:     29,
			Gender:  "Male",
		},
		{
			Name:    "Katerina",
			Surname: "Samsonova",
			Age:     17,
			Gender:  "Female",
		},
		{
			Name:    "Vasiliy",
			Surname: "Juk",
			Age:     23,
			Gender:  "Male",
		},
		{
			Name:    "Iliya",
			Surname: "Pashkovskiy",
			Age:     30,
			Gender:  "Male",
		},
		{
			Name:    "Vladimir",
			Surname: "Chzen",
			Age:     21,
			Gender:  "Male",
		},
	}

	if err := godotenv.Load("../../.env"); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	datebase, err := db.InitDB()
	if err != nil {
		t.Fatalf("Error initializing test database: %v", err)
	}
	defer db.CloseDB()

	repo := NewSQLPersonRepository(datebase)

	var ids []int

	for _, pers := range persons {
		id, err := repo.CreatePerson(pers)
		ids = append(ids, id)
		if err != nil {
			t.Fatalf("Expected no error, but got %v", err)
		}
	}

	defer func() {
		for _, id := range ids {
			err := repo.DeletePerson(id)
			if err != nil {
				t.Fatalf("Error %v", err)
			}
		}
	}()

	var personFilter1 = types.PersonFilter{
		Gender: "Female",
	}

	filter1, _ := repo.FilterListPeople(personFilter1)

	if len(filter1) != 1 {
		t.Fatalf("The answer does not satisfy the request: %d != 1", len(filter1))
	}

	var personFilter2 = types.PersonFilter{
		Gender: "Male",
		MinAge: 20,
		MaxAge: 29,
	}

	filter2, _ := repo.FilterListPeople(personFilter2)

	if len(filter2) != 3 {
		t.Fatalf("The answer does not satisfy the request: %d != 3", len(filter2))
	}

	var personFilter3 = types.PersonFilter{
		Page:     2,
		PageSize: 2,
	}

	filter3, _ := repo.FilterListPeople(personFilter3)

	if len(filter3) != 2 {
		t.Fatalf("The answer does not satisfy the request: %d != 2", len(filter3))
	}

}
