package contact

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

const DBFILE string = "contacts.json"

var DB []Contact = []Contact{}
var nextId int = 1

type Contact struct {
	Id    int    `json:"id"`
	First string `json:"first"`
	Last  string `json:"last"`
	Email string `json:"email"`
	Phone string `json:"phone"`
}

func All() *[]Contact {
	return &DB
}

func Search(term string) []Contact {
	results := 
	return []Contact{}
}

func returnFullName(c Contact) string {
	return fmt.Sprintf("%s %s", c.First, c.Last)
}

func LoadDB() {
	dbfile, err := os.Open(DBFILE)
	if err != nil {
		if os.IsNotExist(err) {
			return
		}
		panic(err)
	}
	defer dbfile.Close()

	file, err := io.ReadAll(dbfile)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(file, &DB)
	if err != nil {
		panic(err)
	}

	nextId += len(DB)
}

func SaveDB() {
	file, err := os.Create(DBFILE)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	data, err := json.Marshal(DB)
	if err != nil {
		panic(err)
	}

	_, err = file.Write(data)
	if err != nil {
		panic(err)
	}
}

func CreateContact(first string, last string, email string, phone string) {
	contact := Contact{
		Id:    nextId,
		First: first,
		Last:  last,
		Email: email,
		Phone: phone,
	}

	nextId++

	DB = append(DB, contact)

	SaveDB()
}

func RemoveContact(id int) {
	idx := id - 1

	copy(DB[idx:], DB[idx+1:])
	DB = DB[:len(DB)-1]

	SaveDB()
	reIdContacts()
}

func reIdContacts() {
	nextId = 1
	for i := range DB {
		if DB[i].Id == nextId {
			nextId++
			continue
		}
		DB[i].Id = nextId
		nextId++
	}

	SaveDB()
}
