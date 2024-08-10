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
	if r.URL.Path == "/contacts/" {
		http.Redirect(w, r, "/contacts", http.StatusPermanentRedirect)
		return
	}

	url := strings.Split(r.URL.Path, "/")

	layout := filepath.Join("templates", "layout.html")
	var body string = ""
	var id int = 0
	var err error = errors.ErrUnsupported

	action := func() string {
		if len(url) == 3 {
			id, err = strconv.Atoi(url[2])
			if err != nil {
				return url[2]
			}
			return ""
		} else if len(url) == 4 {
			id, err = strconv.Atoi(url[2])
			if err != nil {
				return url[2]
			}
			return url[3]
		} else {
			return url[2]
		}
	}()

	if action == "edit" {
		if r.Method == http.MethodGet {
			body = filepath.Join("templates", "edit.html")
			id, err = strconv.Atoi(url[len(url)-2])
			if err != nil {
				log.Fatal(err)
			}
		} else if r.Method == http.MethodPost {
			contact.Update(id,
				r.FormValue("first"),
				r.FormValue("last"),
				r.FormValue("email"),
				r.FormValue("phone"))
			http.Redirect(w, r, "/contacts", http.StatusFound)
			return
		}
	} else if action == "" {
		if r.Method == http.MethodGet {
			body = filepath.Join("templates", "show.html")
		} else if r.Method == http.MethodDelete {
			contact.Delete(id)
			http.Redirect(w, r, "/contacts", http.StatusSeeOther)
			return
		}
	} else {
		http.Redirect(w, r, "/contacts", http.StatusPermanentRedirect)
		return
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
