package handler

import (
	"net/http"

	"github.com/Ayobami0/personal_blog/src/templ"
	"github.com/Ayobami0/personal_blog/src/utils"
)

func BlogsHandler(w http.ResponseWriter, r *http.Request, t *templ.Templ) {
  blogs := utils.GetBlogs()

  t.Get("blogs").ExecuteTemplate(w, "base", blogs)
}
