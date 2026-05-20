package main

import (
	"embed"
	"html/template"
	"io/fs"
	"net/http"
)

//go:embed docs/VERBOSE-RESUME.md
var verboseResumeFormatMarkdown string

//go:embed all:static
var embeddedStatic embed.FS

//go:embed all:templates
var embeddedTemplates embed.FS

func loadTemplates() *template.Template {
	return template.Must(
		template.New("").Funcs(templateFuncs()).ParseFS(
			embeddedTemplates,
			"templates/partials/*.html",
			"templates/*.html",
		),
	)
}

func staticFileServer() (http.Handler, error) {
	sub, err := fs.Sub(embeddedStatic, "static")
	if err != nil {
		return nil, err
	}
	return http.StripPrefix("/static/", http.FileServer(http.FS(sub))), nil
}
