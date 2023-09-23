package controllers

import (
	"html/template"
	"net/http"
)

func StaticHandler(tpl Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tpl.Execute(w, r, nil)
	}
}

func FAQ(tmpl Template) http.HandlerFunc {
	questions := []struct {
		Question string
		Answer   template.HTML
	}{
		{
			Question: "What is the meaning of life?",
			Answer:   "42",
		},
		{
			Question: "What is the meaning of life?",
			Answer:   "42",
		},
		{
			Question: "How do I contact support?",
			Answer: `Email us - <a href="mailto:support@lenslocked.com">
				support@lenslocked.com</a>`,
		},
	}

	return func(w http.ResponseWriter, r *http.Request) {
		tmpl.Execute(w, r, questions)
	}

}
