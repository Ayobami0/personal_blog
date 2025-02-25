package templ

import (
	"fmt"
	"html/template"

	"github.com/Ayobami0/personal_blog/src/constants"
)

/// Manages the registration of nested templates
type Templ struct {
	templates map[string]*template.Template
}

func New() *Templ {
  return &Templ{
    templates: make(map[string]*template.Template),
  }
}

func (t *Templ) Register(name string, files ...string) {
  var f []string

  for _, v := range files {
    f = append(f, string(fmt.Sprintf("%s/%s", constants.TEMPLATES_DIR, v)))
  }

  t.templates[name] = template.Must(template.ParseFiles(f...))
}

func (t Templ) Get(name string) *template.Template {
  tmp, ok := t.templates[name]

  if !ok {
    return nil
  }

  return tmp
}
