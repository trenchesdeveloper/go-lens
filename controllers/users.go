package controllers

import (
	"fmt"
	"net/http"

	"github.com/trenchesdeveloper/lenslocked/models"
)

type Users struct {
	Templates struct {
		New    Template
		SignIn Template
	}

	UserService    *models.UserService
	SessionService *models.SessionService
}

func (u Users) New(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email string
	}
	data.Email = r.FormValue("email")
	u.Templates.New.Execute(w, r, data)
}

func (u Users) SignIn(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email string
	}
	data.Email = r.FormValue("email")
	u.Templates.SignIn.Execute(w, r, data)
}

func (u Users) ProcessSignIn(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	password := r.FormValue("password")

	user, err := u.UserService.Authenticate(email, password)

	if err != nil {
		fmt.Println(err)
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	session, err := u.SessionService.Create(user.ID)

	if err != nil {
		fmt.Println(err)
		http.Redirect(w, r, "/signin", http.StatusFound)
		return
	}

	SetCookie(w, CookieSession, session.Token)

	http.Redirect(w, r, "/users/me", http.StatusFound)

	fmt.Fprintf(w, "User authenticated successfully: %v", user)
}

func (u Users) Create(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	password := r.FormValue("password")

	user, err := u.UserService.Create(email, password)

	if err != nil {
		fmt.Println(err)
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	session, err := u.SessionService.Create(user.ID)

	if err != nil {
		fmt.Println(err)
		http.Redirect(w, r, "/signin", http.StatusFound)
		return
	}

	SetCookie(w, CookieSession, session.Token)

	http.Redirect(w, r, "/users/me", http.StatusFound)

	fmt.Fprintf(w, "User created successfully: %v", user)

}

func (u Users) CurrentUser(w http.ResponseWriter, r *http.Request) {
	token, err := readCookie(r, CookieSession)

	if err != nil {
		fmt.Println(err)
		http.Redirect(w, r, "/signin", http.StatusFound)
		return
	}
	user, err := u.SessionService.User(token)
	fmt.Fprintf(w, "Current user: %s\n", user.Email)
}


func (u Users) SignOut(w http.ResponseWriter, r *http.Request) {
	token, err := readCookie(r, CookieSession)

	if err != nil {
		http.Redirect(w, r, "/signin", http.StatusFound)
		return
	}

	if err := u.SessionService.Delete(token); err != nil {
		fmt.Println(err)
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	// delete the cookie
	deleteCookie(w, CookieSession)

	http.Redirect(w, r, "/signin", http.StatusFound)
}