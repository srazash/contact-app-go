package serve

import (
	"contactapp/contact"
	"contactapp/counter"
	"errors"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
)

func Root(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		counter.Increment()
	}
	http.Redirect(w, r, "/contacts", http.StatusFound)
}

func Contacts(w http.ResponseWriter, r *http.Request) {
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

	data := struct {
		Contacts []contact.Contact
		Term     string
		Counter  string
	}{
		Contacts: func() []contact.Contact {
			if term != "" {
				found := contact.Search(term)
				if len(found) > 0 {
					return found
				}
			}
			return *contact.Ptr()
		}(),
		Term:    term,
		Counter: counter.PaddedCount(),
	}

	err = tmpl.ExecuteTemplate(w, "layout", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func ContactsNew(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		layout := filepath.Join("templates", "layout.html")
		contactnew := filepath.Join("templates", "new.html")

		tmpl, err := template.ParseFiles(layout, contactnew)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		data := struct {
			Counter string
		}{
			Counter: counter.PaddedCount(),
		}

		err = tmpl.ExecuteTemplate(w, "layout", data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	} else if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Failed to parse form", http.StatusBadRequest)
			return
		}

		contact.Create(r.FormValue("first"),
			r.FormValue("last"),
			r.FormValue("email"),
			r.FormValue("phone"))

		http.Redirect(w, r, "/contacts", http.StatusFound)
		return
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
}

func ContactsShowEdit(w http.ResponseWriter, r *http.Request) {
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
	} else if url[len(url)-1] == "delete" {
		id, err = strconv.Atoi(url[len(url)-2])
		if err != nil {
			log.Fatal(err)
		}
		contact.Delete(id)
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

	con, err := contact.Find(id)
	if err != nil {
		log.Fatal(err)
	}

	data := struct {
		Contact contact.Contact
		Counter string
	}{
		Contact: con,
		Counter: counter.PaddedCount(),
	}

	err = tmpl.ExecuteTemplate(w, "layout", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
