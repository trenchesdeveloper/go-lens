package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/trenchesdeveloper/lenslocked/controllers"
	"github.com/trenchesdeveloper/lenslocked/models"
	"github.com/trenchesdeveloper/lenslocked/templates"
	"github.com/trenchesdeveloper/lenslocked/views"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", controllers.StaticHandler(
		views.Must(
			views.ParseFS(
				templates.FS, "home.gohtml", "tailwind.gohtml"))))

	r.Get("/contact", controllers.StaticHandler(
		views.Must(views.ParseFS(templates.FS, "contact.gohtml", "tailwind.gohtml"))))

	r.Get("/faq", controllers.FAQ(
		views.Must(views.ParseFS(templates.FS, "faq.gohtml", "tailwind.gohtml"))))

	// get the default db parameters
	cfg := models.DefaultPostgresConfig()

	// open a connection to the db
	db, err := models.NewPostgresDB(cfg)

	defer db.Close()

	if err != nil {
		panic(err)
	}

	// create a new UserService
	userService := models.UserService{
		DB: db,
	}

	users := controllers.Users{
		UserService: &userService,
	}

	users.Templates.New = views.Must(views.ParseFS(templates.FS, "signup.gohtml", "tailwind.gohtml"))
	users.Templates.SignIn = views.Must(views.ParseFS(templates.FS, "signin.gohtml", "tailwind.gohtml"))

	r.Get("/signup", users.New)
	r.Get("/signin", users.SignIn)
	r.Post("/users", users.Create)
	r.Post("/signin", users.ProcessSignIn)
	r.Get("/users/me", users.CurrentUser)

	fmt.Println("Server is running on port 3000")

	http.ListenAndServe(":3000", r)
}
