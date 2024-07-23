package main

import (
	"contactapp/contact"
	"contactapp/counter"
	"contactapp/serve"
	"log"
	"net/http"
)

func main() {
	contact.Load()
	counter.Load()
	ptrDB := contact.Ptr()
	log.Printf("DB loaded, contacts: %d\n", len(*ptrDB))

	static := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", static))

	http.HandleFunc("/", serve.Root)
	http.HandleFunc("/contacts", serve.Contacts)
	http.HandleFunc("/contacts/new", serve.ContactsNew)
	http.HandleFunc("/contacts/new/save", serve.ContactsNewSave)
	http.HandleFunc("/contacts/", serve.ContactsShowEdit)

	port := ":3000"
	log.Printf("http://localhost%s\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
