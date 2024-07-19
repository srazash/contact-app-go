package main

import (
	"contactapp/contact"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
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

	static := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", static))

	http.HandleFunc("/", serveRoot)
	http.HandleFunc("/contacts", serveContacts)
	http.HandleFunc("/contacts/new", serveContactsNew)
	http.HandleFunc("/contacts/new/save", serveContactsNewSave)
	http.HandleFunc("/contacts/", serveContactsShowEdit)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func serveRoot(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/contacts", http.StatusFound)
}

type serveContactsData struct {
	Contacts []contact.Contact
	Term     string
}

func serveContacts(w http.ResponseWriter, r *http.Request) {
	layout := filepath.Join("templates", "layout.html")
	index := filepath.Join("templates", "index.html")

	tmpl, err := template.ParseFiles(layout, index)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var term string = ""
	arg := r.URL.Query()
	term = arg.Get("q")

	var data serveContactsData = serveContactsData{}

	if term == "" {
		data = serveContactsData{
			Contacts: *contact.All(),
			Term:     term,
		}
	} else {
		data = serveContactsData{
			Contacts: contact.Search(term),
			Term:     term,
		}
	}

	err = tmpl.ExecuteTemplate(w, "layout", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func serveContactsNew(w http.ResponseWriter, r *http.Request) {
	layout := filepath.Join("templates", "layout.html")
	contactnew := filepath.Join("templates", "new.html")

	tmpl, err := template.ParseFiles(layout, contactnew)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.ExecuteTemplate(w, "layout", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func serveContactsNewSave(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	contact.CreateContact(r.FormValue("first"),
		r.FormValue("last"),
		r.FormValue("email"),
		r.FormValue("phone"))

	http.Redirect(w, r, "/contacts", http.StatusFound)
}

func serveContactsShowEdit(w http.ResponseWriter, r *http.Request) {
	url := strings.Split(r.URL.Path, "/")

	layout := filepath.Join("templates", "layout.html")
	var body string = ""
	var id int = 0
	var err error = errors.ErrUnsupported

	if url[len(url)-1] == "edit" {
		body = filepath.Join("templates", "edit.html")
		id, err = strconv.Atoi(url[len(url)-2])
		if err != nil {
			log.Fatal(err)
		}
	} else if url[len(url)-1] == "save" {
		id, err = strconv.Atoi(url[len(url)-2])
		if err != nil {
			log.Fatal(err)
		}
		contact.Update(id,
			r.FormValue("first"),
			r.FormValue("last"),
			r.FormValue("email"),
			r.FormValue("phone"))
		http.Redirect(w, r, "/contacts", http.StatusFound)
		return
	} else {
		body = filepath.Join("templates", "show.html")
		id, err = strconv.Atoi(url[len(url)-1])
		if err != nil {
			log.Fatal(err)
		}
	}

	tmpl, err := template.ParseFiles(layout, body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data, err := contact.Find(id)
	if err != nil {
		log.Fatal(err)
	}

	err = tmpl.ExecuteTemplate(w, "layout", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
