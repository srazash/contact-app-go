package main

import (
	"contactapp/internal/contact"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

func main() {
	ptrDB := contact.All()
	fmt.Printf("DB len: %d\n", len(*ptrDB))

	//contact.CreateContact("Ryan", "Shaw-Harrison", "ryan@mail.local", "+44 (0) 1234 567890")
	contact.LoadDB()
	fmt.Printf("DB loaded, len: %d\n", len(*ptrDB))
	fmt.Printf("DB: %v\n", *ptrDB)

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

func serveContacts(w http.ResponseWriter, r *http.Request) {
	layout := filepath.Join("templates", "layout.html")
	index := filepath.Join("templates", "index.html")

	tmpl, err := template.ParseFiles(layout, index)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.ExecuteTemplate(w, "layout", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
