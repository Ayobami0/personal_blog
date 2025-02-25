package handler

import (
	"net/http"

	"github.com/Ayobami0/personal_blog/src/model"
	"github.com/Ayobami0/personal_blog/src/templ"
)

func ErrorHandler(w http.ResponseWriter, r *http.Request, t *templ.Templ, error model.ErrorResponse) {
	t.Get("error").ExecuteTemplate(w, "base", error)
}
