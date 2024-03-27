package handler

import (
	"html/template"
	"net/http"
	"runtime"
)

type Data struct {
	Version string
}

func WelcomeHandle(w http.ResponseWriter, r *http.Request) {
	data := Data{
		Version: runtime.Version(),
	}

	tmpl, err := template.ParseFiles("public/views/index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
