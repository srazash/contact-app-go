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
	log.Printf("loaded %d contacts\n", len(*ptrDB))
	log.Printf("loaded %d visitors\n", counter.Count)

	static := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", static))

	http.HandleFunc("/", serve.Root)
	http.HandleFunc("/contacts", serve.Contacts)
	http.HandleFunc("/contacts/new", serve.ContactsNew)
	http.HandleFunc("/contacts/", serve.ContactsShowEdit)

	port := ":3000"
	log.Printf("http://localhost%s\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
