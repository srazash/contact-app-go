package contact

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

const DBFILE string = "contacts.json"

var DB []Contact = []Contact{}
var EMAIL = make(map[string]bool)
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

	for _, c := range DB {
		EMAIL[c.Email] = true
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
	EMAIL[email] = true
	Save()

	return contact.Id
}

func ValidateForm(values *map[string]string) map[string]string {
	v := *values
	errors := make(map[string]string)

	if v["First"] == "" {
		errors["First"] = "First name is required"
	}
	if v["Last"] == "" {
		errors["Last"] = "Last name is required"
	}
	if v["Email"] == "" {
		errors["Email"] = "Email is required"
	}
	if EMAIL[v["Email"]] {
		errors["Email"] = "Email must be unique"
	}
	if v["Phone"] == "" {
		errors["Phone"] = "Phone number is required"
	}

	return errors
}

func ValidateEmail(email string) string {
	if email == "" {
		return "Email is required"
	}
	if EMAIL[email] {
		return "Email must be unique"
	}
	return ""
}

func ValidateContactId(contact_id int) bool {
	return contact_id > 0 && contact_id <= len(DB)+1
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

func MultiDelete(ids []int) {

	del := []int{}
	for i := len(ids) - 1; i >= 0; i-- {
		del = append(del, ids[i])
	}

	for _, id := range del {
		idx := id - 1

		copy(DB[idx:], DB[idx+1:])
		DB = DB[:len(DB)-1]
	}

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
		t := strings.ToLower(term)
		c := strings.ToLower(contactString(DB[i]))
		if strings.Contains(c, t) {
			results = append(results, DB[i])
		}
	}

	time.Sleep(time.Millisecond * 1000)

	return results
}

func contactString(c Contact) string {
	return fmt.Sprintf("%s %s %s %s", c.First, c.Last, c.Email, c.Phone)
}

func GetTailId() int {
	return nextId - 1
}

func ContactsCount() int {
	return len(DB)
}

func PaginatedContacts(page int, items int) []Contact {
	start := (page - 1) * items
	end := start + items

	if start >= len(DB) {
		return []Contact{}
	}

	if end > len(DB) {
		end = len(DB)
	}

	return DB[start:end]
}

func NextPage(page int, items int) bool {
	end := (page-1)*items + items

	return end <= len(DB)
}

func PrevPage(page int) bool {
	return page > 1
}

func TotalPages(items int) int {
	total := len(DB) / items
	if len(DB)%items > 0 {
		total += len(DB) % items
	}
	return total
}

func LastPage(page int, items int) bool {
	return page == TotalPages(items)
}

func JsonAllContacts() []byte {
	jsonData, err := json.Marshal(DB)
	if err != nil {
		panic(err)
	}
	return jsonData
}

func JsonContactById(contact_id int) ([]byte, error) {
	contact, err := Find(contact_id)
	if err != nil {
		return nil, err
	}

	jsonData, err := json.Marshal(contact)
	if err != nil {
		return nil, err
	}
	return jsonData, nil
}
