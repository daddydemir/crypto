package assets

import (
	"embed"
	"html/template"
	"log/slog"
)

//go:embed templates/*.html
var templateFs embed.FS

func GetTemplate(name string) *template.Template {
	tmpl, err := template.ParseFS(templateFs, name)
	if err != nil {
		slog.Error("ParseFs", "error", err)
	}
	return tmpl
}
