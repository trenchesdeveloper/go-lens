package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func executeTemplate(w http.ResponseWriter, r *http.Request, templateName string, data interface{}) {
	w.Header().Set("Content-Type", "text/html")
	filePath := filepath.Join("templates", templateName)
	tpl, err := template.ParseFiles(filePath)

	if(err != nil){
		log.Printf("Error parsing template: %v", err)
		http.Error(w, "Error Parsing Template", http.StatusInternalServerError)
		return
	}

	err = tpl.Execute(w, data)

	if(err != nil){
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Error Executing Template", http.StatusInternalServerError)
		return
	}
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	executeTemplate(w, r, "home.gohtml", nil)
}

func contactHandler(w http.ResponseWriter, r *http.Request) {
	executeTemplate(w, r, "contact.gohtml", nil)
}

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", homeHandler)
	r.Get("/contact", contactHandler)

	fmt.Println("Server is running on port 3000")

	http.ListenAndServe(":3000", r)
}
