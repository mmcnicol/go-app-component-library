//go:build dev
// pkg/components/viz/horizontal_bar_chart_stories.go
package viz

import (
	"github.com/maxence-charriere/go-app/v10/pkg/app"
	"github.com/mmcnicol/go-app-component-library/pkg/storybook"
)

func init() {
	storybook.Register("Visualization", "Horizontal Bar Chart",
		map[string]*storybook.Control{
			"Title":      storybook.NewTextControl("Customer Satisfaction Survey"),
			"Bar Color":  storybook.NewColorControl("#8b5cf6"),
		},
		func(controls map[string]*storybook.Control) app.UI {
			spec := Spec{
				Type:  ChartTypeBar,
				Title: controls["Title"].Value.(string),
				Data: DataSet{
					Labels: []string{"Very Satisfied", "Satisfied", "Neutral", "Dissatisfied", "Very Dissatisfied"},
					Series: []Series{{
						Label:  "Responses",
						Points: Values([]float64{245, 180, 65, 42, 18}),
						Color:  controls["Bar Color"].Value.(string),
					}},
				},
				Width:  800,
				Height: 400,
				Bar: BarConfig{
					Horizontal: true,
					Width:      80,
				},
			}

			chart := New(spec)
			
			return app.Div().ID("viz-horizontal-bar-chart-container").Body(
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
