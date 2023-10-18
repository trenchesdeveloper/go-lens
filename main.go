package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/gorilla/csrf"
	"github.com/joho/godotenv"
	"github.com/trenchesdeveloper/lenslocked/controllers"
	"github.com/trenchesdeveloper/lenslocked/migrations"
	"github.com/trenchesdeveloper/lenslocked/models"
	"github.com/trenchesdeveloper/lenslocked/templates"
	"github.com/trenchesdeveloper/lenslocked/views"
)

func main() {
	err := godotenv.Load()
  if err != nil {
    log.Fatal("Error loading .env file")
  }
	// Setup DB
	// get the default db parameters
	cfg := models.DefaultPostgresConfig()

	fmt.Println(cfg.String())

	// open a connection to the db
	db, err := models.NewPostgresDB(cfg)

	defer db.Close()

	err = models.MigrateFS(db, migrations.FS, ".")

	if err != nil {
		panic(err)
	}

	//setup services
	// create a new UserService
	userService := models.UserService{
		DB: db,
	}

	// create session service
	sessionService := models.SessionService{
		DB: db,
	}

	//setup Middleware
	umw := controllers.UserMiddleware{
		SessionService: &sessionService,
	}

	csrfMw := csrf.Protect(
		[]byte("vj-ytid-459k-3mck-oit"),
		// TODO: change this to true in production
		csrf.Secure(false),
	)

	//setup controllers
	users := controllers.Users{
		UserService:    &userService,
		SessionService: &sessionService,
	}

	users.Templates.New = views.Must(views.ParseFS(templates.FS, "signup.gohtml", "tailwind.gohtml"))
	users.Templates.SignIn = views.Must(views.ParseFS(templates.FS, "signin.gohtml", "tailwind.gohtml"))

	// setup router
	r := chi.NewRouter()
	r.Use(csrfMw)
	r.Use(umw.SetUser)
	r.Use(middleware.Logger)
	r.Get("/", controllers.StaticHandler(
		views.Must(
			views.ParseFS(
				templates.FS, "home.gohtml", "tailwind.gohtml"))))

	r.Get("/contact", controllers.StaticHandler(
		views.Must(views.ParseFS(templates.FS, "contact.gohtml", "tailwind.gohtml"))))

	r.Get("/faq", controllers.FAQ(
		views.Must(views.ParseFS(templates.FS, "faq.gohtml", "tailwind.gohtml"))))

	if err != nil {
		panic(err)
	}

	r.Get("/signup", users.New)
	r.Get("/signin", users.SignIn)
	r.Post("/users", users.Create)
	r.Post("/signin", users.ProcessSignIn)
	r.Post("/signout", users.SignOut)
	r.Route("/users/me", func(r chi.Router) {
		r.Use(umw.RequireUser)
		r.Get("/", users.CurrentUser)
	})
	//r.Get("/users/me", users.CurrentUser)

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Page not found", http.StatusNotFound)
	})

	fmt.Println("Server is running on port 3000")

	http.ListenAndServe(":3000", (r))
}
