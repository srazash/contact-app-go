package contact

import (
	"encoding/json"
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

func (c *Contact) GetId() int {
	return c.Id
}

func (c *Contact) GetFirst() string {
	return c.First
}

func (c *Contact) GetLast() string {
	return c.Last
}

func (c *Contact) GetEmail() string {
	return c.Email
}

func (c *Contact) GetPhone() string {
	return c.Phone
}

func (c *Contact) SetId(id int) {
	c.Id = id
}

func (c *Contact) SetFirst(first string) {
	c.First = first
}

func (c *Contact) SetLast(last string) {
	c.Last = last
}

func (c *Contact) SetEmail(email string) {
	c.Email = email
}

func (c *Contact) SetPhone(phone string) {
	c.Phone = phone
}

func All() *[]Contact {
	return &DB
}

func Search(term string) []Contact {
	return []Contact{}
}

func LoadDB() {
	file, err := os.ReadFile(DBFILE)
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
	file, err := os.OpenFile(DBFILE, os.O_APPEND|os.O_CREATE, 0666)
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
