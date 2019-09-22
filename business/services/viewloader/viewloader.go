package viewloader

import (
	"encoding/json"
	"fmt"
	"goairmon/business/data/context"
	"goairmon/site/helper"
	vmodels "goairmon/site/models"
	"html/template"
	"log"
	"path/filepath"
	"time"

	"github.com/labstack/echo"
)

const CtxKey = "view_loader"

type ViewLoader struct {
}

func fullViewPath(viewPath string) string {
	return helper.ResourceRoot() + "/views/" + viewPath
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

	var reducedSensorPoints *vmodels.ReducedSensorPoints

	mainTemplate = mainTemplate.Funcs(template.FuncMap{
		"points48Hours": func() string {
			if reducedSensorPoints == nil {
				reducedSensorPoints = v.initReducedSensorPoints(c)
			}

			raw, err := json.Marshal(reducedSensorPoints.Last48Hours())
			if err != nil {
				log.Println(err)
			}

			return string(raw)
		},
		"points2Hours": func() string {
			if reducedSensorPoints == nil {
				reducedSensorPoints = v.initReducedSensorPoints(c)
			}

			raw, err := json.Marshal(reducedSensorPoints.Last2Hours())
			if err != nil {
				log.Println(err)
			}

			return string(raw)
		},
		"points7Days": func() string {
			if reducedSensorPoints == nil {
				reducedSensorPoints = v.initReducedSensorPoints(c)
			}

			raw, err := json.Marshal(reducedSensorPoints.Last7Days())
			if err != nil {
				log.Println(err)
			}

			return string(raw)
		},
	})

	parsed, err := mainTemplate.ParseFiles(files...)
	if err != nil {
		panic(fmt.Sprintf("failed to parse template: %s", err))
	}

	return parsed
}

func (v *ViewLoader) initReducedSensorPoints(c echo.Context) *vmodels.ReducedSensorPoints {
	points, err := c.Get(helper.CtxDbContext).(context.DbContext).GetSensorPoints(60 * 24 * 8)
	if err != nil {
		log.Println(err)
	}

	return vmodels.NewReducedSensorPoints(points, time.Now())
}
