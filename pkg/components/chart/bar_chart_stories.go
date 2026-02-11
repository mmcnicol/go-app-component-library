//go:build dev
package chart

import (
	"github.com/maxence-charriere/go-app/v10/pkg/app"
	"github.com/mmcnicol/go-app-component-library/pkg/storybook"
)

func init() {
	// Sample datasets for the story
	datasets := map[string][]float64{
		"Monthly Sales":    {45000, 52000, 48000, 61000, 58000, 63000, 72000, 68000, 74000, 82000, 79000, 85000},
		"Browser Usage":    {65, 18, 12, 5},
		"Device Types":     {45, 30, 25},
		"Customer Ratings": {22, 45, 68, 92, 78, 34, 12},
	}

	// Sample labels for each dataset
	datasetLabels := map[string][]string{
		"Monthly Sales":    {"Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"},
		"Browser Usage":    {"Chrome", "Firefox", "Safari", "Edge"},
		"Device Types":     {"Desktop", "Mobile", "Tablet"},
		"Customer Ratings": {"1★", "2★", "3★", "4★", "5★", "6★", "7★"},
	}

	storybook.Register("Charts", "Bar Chart",
		map[string]*storybook.Control{
			"Dataset":      storybook.NewSelectControl([]string{"Monthly Sales", "Browser Usage", "Device Types", "Customer Ratings"}, "Monthly Sales"),
			"Bar Color":    storybook.NewColorControl("#10b981"),
			"Show Labels":  storybook.NewBoolControl(true),
			"Title":        storybook.NewTextControl("Bar Chart Example"),
		},
		func(controls map[string]*storybook.Control) app.UI {
			selectedKey := controls["Dataset"].Value.(string)
			data := datasets[selectedKey]
			labels := datasetLabels[selectedKey]
			showLabels := controls["Show Labels"].Value.(bool)

			// Create bar chart with options
			chart := New(nil,
				WithTitle(controls["Title"].Value.(string)),
				WithColor(controls["Bar Color"].Value.(string)),
				func(c *ChartConfig) {
					// Clear other chart types
					c.BoxData = nil
					c.PieData = nil
					c.HeatmapMatrix = nil
					c.IsStream = false
					
					// Set bar chart data
					c.BarData = data
					
					// Set labels if enabled
					if showLabels {
						c.BarLabels = labels
					} else {
						c.BarLabels = nil
					}
					
					// Optional: Add custom colors for different bars
					if selectedKey == "Browser Usage" {
						c.BarColors = []string{
							"#4285F4", // Chrome blue
							"#FF7139", // Firefox orange
							"#1D1D1F", // Safari black
							"#0078D7", // Edge blue
						}
					}
				},
			)
			
			return app.Div().ID("bar-chart-container").Body(
				app.Div().Class("chart-wrapper").Body(chart),
			)
		},
	)
}
