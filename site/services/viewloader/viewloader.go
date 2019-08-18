package viewloader

import (
	"html/template"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/labstack/echo"
)

const CtxKey = "view_loader"

type ViewLoader struct {
}

func (v *ViewLoader) Load(viewPath string) *template.Template {
	panic("not implemented")
}

func siteRoot() string {
	dir, _ := os.Getwd()
	for i := 0; i < 10; i++ {

		if strings.HasSuffix(filepath.Dir(dir), "/site") {
			return dir
		}

		dir = dir + "../"
	}

	panic("failed to find site root")
}

func fullViewPath(viewPath string) string {
	return siteRoot() + "/views/" + viewPath
}

func LoadView(viewPath string, c echo.Context) *template.Template {
	tmpl, err := template.ParseFiles(fullViewPath(viewPath))
	if err != nil {
		log.Println(err)
		panic("failed to load template: " + viewPath)
	}

	return tmpl
}
