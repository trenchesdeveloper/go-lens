package views

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

type Template struct {
	htmlTpl *template.Template
}

func Parse(filePath string) (Template, error) {
	tpl, err := template.ParseFiles(filePath)

	if(err != nil){
		return Template{}, fmt.Errorf("Error parsing template: %w", err)
	}

	return Template{htmlTpl: tpl}, nil
}

func (t Template) Execute(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	err := t.htmlTpl.Execute(w, data)

	if(err != nil){
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Error Executing Template", http.StatusInternalServerError)
		return
	}
}