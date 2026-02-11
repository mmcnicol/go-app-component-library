//go:build dev
// pkg/components/chart/bar_chart_stories.go
package chart

import (
	"github.com/maxence-charriere/go-app/v10/pkg/app"
	"github.com/mmcnicol/go-app-component-library/pkg/storybook"
)

// Use init() to auto-register when this package is imported
func init() {

    storybook.Register("Chart", "Bar Chart", 
        nil,
        func(controls map[string]*storybook.Control) app.UI {
            return NewChart(ChartTypeBar).
				Title("User Activity by Hour").
				Data(ChartData{
					Labels: []string{"9AM", "10AM", "11AM", "12PM", "1PM", "2PM", "3PM", "4PM"},
					Datasets: []Dataset{
						{
							Label: "Week 1",
							Data:  userActivity[0],
							BackgroundColor: "rgba(74, 144, 226, 0.7)",
						},
						{
							Label: "Week 2",
							Data:  userActivity[1],
							BackgroundColor: "rgba(255, 99, 132, 0.7)",
						},
						{
							Label: "Week 3",
							Data:  userActivity[2],
							BackgroundColor: "rgba(75, 192, 192, 0.7)",
						},
					},
				}).
				Options(ChartOptions{
					Scales: ChartScales{
						Y: Axis{
							Stacked: true,
							Title: AxisTitle{
								Display: true,
								Text: "Active Users",
							},
						},
						X: Axis{
							Stacked: true,
						},
					},
				}).
				Class("dashboard-card", "chart-bar")
        },
    )

}
