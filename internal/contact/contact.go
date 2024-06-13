package contact

var DB []Contact = []Contact{}
var nextId int = 1

type Contact struct {
	id    int
	first string
	last  string
	email string
	phone string
}

func All() *[]Contact {
	return &DB
}

func CreateContact(first string, last string, email string, phone string) {
	contact := Contact{
		id:    nextId,
		first: first,
		last:  last,
		email: email,
		phone: phone,
	}

	nextId++

	DB = append(DB, contact)
}
