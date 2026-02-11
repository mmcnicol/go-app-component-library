//go:build dev

package chart

import (
	"github.com/maxence-charriere/go-app/v10/pkg/app"
	"github.com/mmcnicol/go-app-component-library/pkg/storybook"
)

func init() {
	// Sample data
	monthlySales := []Point{
		{X: 1, Y: 12000}, {X: 2, Y: 19000}, {X: 3, Y: 15000},
		{X: 4, Y: 21000}, {X: 5, Y: 18000}, {X: 6, Y: 24000},
	}

	storybook.Register("Charts", "Line Chart",
		map[string]*storybook.Control{
			"Title":      storybook.NewTextControl("Monthly Sales"),
			"Line Color": storybook.NewColorControl("#4f46e5"),
			"Thickness":  storybook.NewRangeControl(1, 10, 1, 3),
		},
		func(controls map[string]*storybook.Control) app.UI {
			return New(monthlySales,
				WithTitle(controls["Title"].Value.(string)),
				WithColor(controls["Line Color"].Value.(string)),
				// Assuming you add a WithThickness option
				func(c *ChartConfig) { 
					c.Thickness = float64(controls["Thickness"].Value.(int)) 
				},
			)
		},
	)
}
