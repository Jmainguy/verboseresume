package main

import (
	"embed"
	"html/template"
	"io/fs"
	"net/http"
)

//go:embed docs/VERBOSE-RESUME.md
var verboseResumeFormatMarkdown string

//go:embed docs/VERBOSE-RESUME-QUESTIONS.md
var verboseResumeQuestionsMarkdown string

//go:embed docs/example-verbose-resume.md
var exampleVerboseResumeMarkdown string

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

func faviconHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodHead {
		w.Header().Set("Allow", "GET, HEAD")
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	data, err := embeddedStatic.ReadFile("static/brand/favicon-32.png")
	if err != nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Cache-Control", "public, max-age=86400")
	if r.Method == http.MethodHead {
		return
	}
	_, _ = w.Write(data)
}
