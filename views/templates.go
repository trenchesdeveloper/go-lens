package views

import (
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net/http"
)

type Template struct {
	htmlTpl *template.Template
}

func ParseFS(fs fs.FS, patterns ...string) (Template, error) {
	tpl, err := template.ParseFS(fs, patterns...)

	if(err != nil){
		return Template{}, fmt.Errorf("Error parsing files: %w", err)
	}

	return Template{htmlTpl: tpl}, nil
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

func Must(tpl Template, err error) Template {
	if(err != nil){
		panic(err)
	}

	return tpl
}