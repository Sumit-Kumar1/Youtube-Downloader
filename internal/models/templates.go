package models

import (
	"html/template"
	"io"

	"github.com/labstack/echo/v5"
)

type Template struct {
	templates *template.Template
}

func (t *Template) Render(_ *echo.Context, w io.Writer, name string, data any) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

// NewTemplate this function return a template struct with parsing all html files in the path
func NewTemplate(path string) *Template {
	return &Template{
		templates: template.Must(template.ParseGlob(path + "/*.html")),
	}
}
