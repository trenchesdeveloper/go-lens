package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/trenchesdeveloper/lenslocked/controllers"
	"github.com/trenchesdeveloper/lenslocked/templates"
	"github.com/trenchesdeveloper/lenslocked/views"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", controllers.StaticHandler(
		views.Must(views.ParseFS(templates.FS, "home.gohtml"))))

	r.Get("/contact", controllers.StaticHandler(
		views.Must(views.ParseFS(templates.FS, "contact.gohtml"))))

	r.Get("/faq", controllers.FAQ(
		views.Must(views.ParseFS(templates.FS, "faq.gohtml"))))

	fmt.Println("Server is running on port 3000")

	http.ListenAndServe(":3000", r)
}
