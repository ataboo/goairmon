package viewloader

import (
	"fmt"
	"goairmon/site/helper"
	"html/template"
	"path/filepath"

	"github.com/labstack/echo"
)

const CtxKey = "view_loader"

type ViewLoader struct {
}

func fullViewPath(viewPath string) string {
	return helper.SiteRoot() + "/views/" + viewPath
}

func (v *ViewLoader) layoutFilenames() []string {
	files, err := filepath.Glob(fullViewPath("layout/*.gohtml"))

	if err != nil {
		panic(fmt.Sprintf("failed to load layout at: %s", fullViewPath("layout")))
	}

	return files
}

func (v *ViewLoader) LoadView(viewPath string, c echo.Context) *template.Template {
	mainTemplate := template.New("main")
	mainTemplate, _ = mainTemplate.Parse(`{{define "main"}} {{template "base" . }} {{end}}`)
	files := append(v.layoutFilenames(), fullViewPath(viewPath))

	parsed, err := mainTemplate.ParseFiles(files...)
	if err != nil {
		panic(fmt.Sprintf("failed to parse template: %s", err))
	}

	return parsed
}
