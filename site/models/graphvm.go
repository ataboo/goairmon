package models

import (
	"fmt"
	"goairmon/business/data/models"
	"html/template"
	"strings"
)

type GraphVm struct {
	SensorPoints []*models.SensorPoint
}

/*
var data = {
            labels: [1, 2, 3, 4, 5, 6],
            datasets: [{
                data: [
                    {x: 1, y: 12},
                    {x: 2, y: 19},
                    {x: 3, y: 3},
                    {x: 4, y: 5},
                    {x: 5, y: 2},
                    {x: 6, y: 3},
                ]
            }]
        }
*/

func (g *GraphVm) GraphJsData() template.HTML {
	labels := make([]string, 0)
	points := make([]string, 0)
	for i := 0; i < len(g.SensorPoints); i++ {
		if g.SensorPoints[i] == nil {
			continue
		}
		timeStr := g.SensorPoints[i].Time.Format("20060102T1504")
		timeStr = fmt.Sprintf(`moment("%s")`, timeStr)

		labels = append(labels, timeStr)
		points = append(points, fmt.Sprintf("{x: %s, y: %.2f}", timeStr, g.SensorPoints[i].Co2Value))
	}

	return template.HTML(fmt.Sprintf("{labels: [%s], datasets: [{ data: [%s] }]}", strings.Join(labels, ", "), strings.Join(points, ", ")))
}
