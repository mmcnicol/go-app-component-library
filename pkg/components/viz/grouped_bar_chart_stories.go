// pkg/components/viz/grouped_bar_chart_stories.go
//go:build dev
package viz

import (
    "fmt"
    "sort"
    "strings"
    
    "github.com/maxence-charriere/go-app/v10/pkg/app"
    "github.com/mmcnicol/go-app-component-library/pkg/storybook"
)

func init() {
	// Grouped bar chart example - multiple series
	groupedData := DataSet{
		Labels: []string{"Q1", "Q2", "Q3", "Q4"},
		Series: []Series{
			{
				Label:  "Product A",
				Points: Values([]float64{45000, 52000, 48000, 61000}),
				Color:  "#4f46e5",
			},
			{
				Label:  "Product B",
				Points: Values([]float64{38000, 49000, 53000, 59000}),
				Color:  "#10b981",
			},
			{
				Label:  "Product C",
				Points: Values([]float64{29000, 41000, 45000, 52000}),
				Color:  "#f59e0b",
			},
		},
	}

	storybook.Register("Visualization", "Grouped Bar Chart",
		map[string]*storybook.Control{
			"Title":        storybook.NewTextControl("Quarterly Product Performance"),
			"Show Legend":  storybook.NewBoolControl(true),
			"Show Values":  storybook.NewBoolControl(true),
			"Interactive":  storybook.NewBoolControl(true),
		},
		func(controls map[string]*storybook.Control) app.UI {
			spec := Spec{
				Type:  ChartTypeBar,
				Title: controls["Title"].Value.(string),
				Data:  groupedData,
				Width:  800,
				Height: 400,
				Interactive: InteractiveConfig{
					Enabled: controls["Interactive"].Value.(bool),
					Tooltip: TooltipConfig{
						Enabled: true,
						Mode:    TooltipModeAll,
					},
				},
				Bar: BarConfig{
					Width:         70,
					Grouped:       true,
					BorderRadius:  4,
				},
				Legend: LegendConfig{
					Visible:  controls["Show Legend"].Value.(bool),
					Position: "top",
				},
				Labels: LabelsConfig{
					Visible: controls["Show Values"].Value.(bool),
					Position: "top",
					FontSize: 11,
				},
			}

			chart := New(spec)
			
			return app.Div().ID("viz-grouped-bar-chart-container").Body(
				app.Div().Class("viz-chart-wrapper").Body(chart),
			)
		},
	)
}
