package main

import (
	"log"
	"net/http"

	"github.com/Ayobami0/personal_blog/src/handler"
	"github.com/Ayobami0/personal_blog/src/templ"
)

const STATIC_DIR = "src/static"

func main() {
	// Templ
	t := templ.New()

	mux := http.NewServeMux()

	t.Register("index", "index.html", "base.html")
  t.Register("blogs", "blogs.html", "base.html")
  t.Register("blog", "blog.html", "base.html")
  t.Register("error", "error.html", "base.html")

  mux.HandleFunc("/blogs/{name}", func(w http.ResponseWriter, r *http.Request) { handler.BlogHandler(w, r, t) })
	mux.HandleFunc("/blogs/", func(w http.ResponseWriter, r *http.Request) { handler.BlogsHandler(w, r, t) })
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) { handler.IndexHandler(w, r, t) })
  mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir(STATIC_DIR))))

  if err := http.ListenAndServe("0.0.0.0:5000", mux); err != nil {
		log.Fatal(err)
	}
}
