package main

import (
	"contactapp/contact"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

func main() {
	ptrDB := contact.All()
	fmt.Printf("DB len: %d\n", len(*ptrDB))

	contact.LoadDB()

	// fmt.Printf("BEFORE: len: %d\n", len(*ptrDB))
	// fmt.Printf("DB: %v\n", *ptrDB)

	// contact.CreateContact("Ryan", "Shaw-Harrison", "ryan@mail.local", "+44 (0) 1234 567890")
	// contact.CreateContact("John", "Smith", "john@mail.local", "+44 (0) 1234 567999")
	// contact.CreateContact("David", "Jones", "david@mail.local", "+44 (0) 1234 567000")
	// contact.CreateContact("Sally", "Brown", "david@mail.local", "+44 (0) 1234 567000")

	// contact.RemoveContact(2)
	// contact.ReIdContacts()

	// fmt.Printf("AFTER: len: %d\n", len(*ptrDB))
	// fmt.Printf("DB: %v\n", *ptrDB)

	search := contact.Search("ryan")
	fmt.Println(search)
	search = contact.Search("mail.local")
	fmt.Println(search)
	search = contact.Search("567000")
	fmt.Println(search)

	result, err := contact.Find(1)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(result)
	}

	for _, c := range *ptrDB {
		fmt.Printf("\tID: %d\n", c.Id)
		fmt.Printf("\tName: %s, %s\n", c.Last, c.First)
		fmt.Printf("\tEmail: %s\n", c.Email)
		fmt.Printf("\tPhone: %s\n", c.Phone)
	}

	http.HandleFunc("/", serveRoot)
	http.HandleFunc("/contacts", serveContacts)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func serveRoot(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/contacts", http.StatusFound)
}

type serveContactsData struct {
	Contacts []contact.Contact
}

func serveContacts(w http.ResponseWriter, r *http.Request) {
	layout := filepath.Join("templates", "layout.html")
	index := filepath.Join("templates", "index.html")

	tmpl, err := template.ParseFiles(layout, index)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := serveContactsData{
		Contacts: *contact.All(),
	}

	err = tmpl.ExecuteTemplate(w, "layout", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
