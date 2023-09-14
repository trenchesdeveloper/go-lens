package main

import (
	"os"
	"html/template"
)

func main() {
	t, err := template.ParseFiles("hello.gohtml")
	if err != nil {
		panic(err)
	}

	user := struct{
		Name string
	}{
		Name: "John",
	}

	err = t.Execute(os.Stdout, user)

	if err != nil {
		panic(err)
	}
}
