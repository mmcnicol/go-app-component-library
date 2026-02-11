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
