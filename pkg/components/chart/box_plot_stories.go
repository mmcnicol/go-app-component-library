//go:build dev

package chart

import (
	"github.com/maxence-charriere/go-app/v10/pkg/app"
	"github.com/mmcnicol/go-app-component-library/pkg/storybook"
)

func init() {
	// Sample raw data for different groups
	groupA := []float64{12, 15, 17, 22, 25, 28, 30, 45, 52}
	groupB := []float64{5, 8, 30, 32, 35, 38, 40, 42, 85}

	storybook.Register("Charts", "Box Plot",
		map[string]*storybook.Control{
			"Box Width":   storybook.NewRangeControl(20, 100, 5, 50),
			"Theme Color": storybook.NewColorControl("#4a90e2"),
		},
		// In box_plot_stories.go, modify the render function:
		func(controls map[string]*storybook.Control) app.UI {
			statsA := CalculateBoxStats(groupA)
			statsB := CalculateBoxStats(groupB)

			return New(nil,
				WithTitle("Comparison of Distributions"),
				WithColor(controls["Theme Color"].Value.(string)),
				func(c *ChartConfig) {
					c.BoxWidth = float64(controls["Box Width"].Value.(int))
					c.BoxData = []BoxPlotStats{statsA, statsB}
					// Ensure the chart knows to draw box plots
					c.PieData = nil  // Clear any pie chart data
					c.HeatmapMatrix = nil  // Clear any heatmap data
				},
			)
		},
	)
}
