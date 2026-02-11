//go:build dev
package chart

import (
	"github.com/maxence-charriere/go-app/v10/pkg/app"
	"github.com/mmcnicol/go-app-component-library/pkg/storybook"
)

func init() {
	// Grouped bar chart example
	groupedData := []struct {
		Category string
		Values   []float64
		Labels   []string
	}{
		{"Q1 2024", []float64{45000, 52000, 48000}, []string{"Product A", "Product B", "Product C"}},
		{"Q2 2024", []float64{61000, 58000, 63000}, []string{"Product A", "Product B", "Product C"}},
		{"Q3 2024", []float64{72000, 68000, 74000}, []string{"Product A", "Product B", "Product C"}},
	}

	storybook.Register("Charts", "Grouped Bar Chart",
		map[string]*storybook.Control{
			"Quarter": storybook.NewSelectControl([]string{"Q1 2024", "Q2 2024", "Q3 2024"}, "Q2 2024"),
			"Bar Color": storybook.NewColorControl("#8b5cf6"),
		},
		func(controls map[string]*storybook.Control) app.UI {
			selectedQuarter := controls["Quarter"].Value.(string)
			
			var data []float64
			var labels []string
			for _, d := range groupedData {
				if d.Category == selectedQuarter {
					data = d.Values
					labels = d.Labels
					break
				}
			}

			chart := New(nil,
				WithTitle("Product Performance - " + selectedQuarter),
				WithColor(controls["Bar Color"].Value.(string)),
				func(c *ChartConfig) {
					c.BarData = data
					c.BarLabels = labels
					// Different colors for each product
					c.BarColors = []string{"#4f46e5", "#10b981", "#f59e0b"}
				},
			)
			
			return app.Div().ID("grouped-bar-chart-container").Body(chart)
		},
	)
}
