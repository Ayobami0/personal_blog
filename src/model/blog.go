package model

import (
	"html/template"
)


type Blog struct {
  Title string
  DateModified string
  Content template.HTML
}
