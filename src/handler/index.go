package handler

import (
	"net/http"

	"github.com/Ayobami0/personal_blog/src/model"
	"github.com/Ayobami0/personal_blog/src/templ"
)

func IndexHandler(w http.ResponseWriter, r *http.Request, t *templ.Templ) {
	if r.URL.Path != "/" {
		ErrorHandler(w, r, t, model.ErrorResponse{Message: "Page not found", Code: 404})
		return
	}
	t.Get("index").ExecuteTemplate(w, "base", nil)
}
