//go:build dev
// pkg/components/viz/grouped_bar_chart_stories.go
package viz

import (
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

// Helper function to generate accessible description
func generateBarChartDescription(title string, data []float64, labels []string) string {
	var desc strings.Builder
	
	desc.WriteString(fmt.Sprintf("Bar chart showing %s. ", title))
	
	if len(data) > 0 {
		// Find min, max, average
		min, max := data[0], data[0]
		sum := 0.0
		for _, v := range data {
			if v < min { min = v }
			if v > max { max = v }
			sum += v
		}
		avg := sum / float64(len(data))
		
		desc.WriteString(fmt.Sprintf("Total of %d categories. ", len(data)))
		desc.WriteString(fmt.Sprintf("Highest value is %.0f, lowest is %.0f. ", max, min))
		desc.WriteString(fmt.Sprintf("Average value is %.1f. ", avg))
		
		// Add top 3 categories
		if len(data) >= 3 && len(labels) >= 3 {
			// Create slice of indices and sort by value
			indices := make([]int, len(data))
			for i := range indices { indices[i] = i }
			sort.Slice(indices, func(i, j int) bool {
				return data[indices[i]] > data[indices[j]]
			})
			
			desc.WriteString("Top categories: ")
			for i := 0; i < 3; i++ {
				idx := indices[i]
				desc.WriteString(fmt.Sprintf("%s at %.0f", labels[idx], data[idx]))
				if i < 2 { desc.WriteString(", ") }
			}
			desc.WriteString(".")
		}
	}
	
	return desc.String()
}
