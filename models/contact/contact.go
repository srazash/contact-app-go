package contact

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
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

func Ptr() *[]Contact {
	return &DB
}

func Load() {
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

func Save() {
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

func Create(first string, last string, email string, phone string) int {
	contact := Contact{
		Id:    nextId,
		First: first,
		Last:  last,
		Email: email,
		Phone: phone,
	}

	nextId++
	DB = append(DB, contact)
	Save()

	return contact.Id
}

func Update(id int, first, last, email, phone string) {
	c := &DB[id-1]
	ops := 0

	if first != c.First {
		c.First = first
		ops++
	}
	if last != c.Last {
		c.Last = last
		ops++
	}
	if email != c.Email {
		c.Email = email
		ops++
	}
	if phone != c.Phone {
		c.Phone = phone
		ops++
	}

	if ops > 0 {
		Save()
	}
}

func Delete(id int) {
	idx := id - 1

	copy(DB[idx:], DB[idx+1:])
	DB = DB[:len(DB)-1]

	updateIds()
	Save()
}

func updateIds() {
	nextId = 1
	for i := range DB {
		if DB[i].Id == nextId {
			nextId++
			continue
		}
		DB[i].Id = nextId
		nextId++
	}
}

func Find(id int) (Contact, error) {
	if id <= 0 || id >= nextId {
		return Contact{}, errors.New("invalid id")
	}
	return DB[id-1], nil
}

func Search(term string) []Contact {
	results := []Contact{}

	for i := range DB {
		c := contactString(DB[i])
		if strings.Contains(c, term) {
			results = append(results, DB[i])
		}
	}

	return results
}

func contactString(c Contact) string {
	return fmt.Sprintf("%s %s %s %s", c.First, c.Last, c.Email, c.Phone)
}

func GetTailId() int {
	return nextId - 1
}
