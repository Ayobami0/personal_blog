package handler

import (
	"net/http"

	"github.com/Ayobami0/personal_blog/src/model"
	"github.com/Ayobami0/personal_blog/src/templ"
	"github.com/Ayobami0/personal_blog/src/utils"
)

func BlogHandler(w http.ResponseWriter, r *http.Request, t *templ.Templ) {
  name := r.PathValue("name")

  
  blog, err := utils.ParseBlogAsHTML(name)

  if err != nil {
    ErrorHandler(w, r, t, model.ErrorResponse{Message: "Page not found", Code: 404})
    return
  }
  
  t.Get("blog").ExecuteTemplate(w, "base", blog)
}
