//go:build dev
// pkg/components/chart/line_chart_stories.go
package chart

import (
	"github.com/maxence-charriere/go-app/v10/pkg/app"
	"github.com/mmcnicol/go-app-component-library/pkg/storybook"
)

// Use init() to auto-register when this package is imported
func init() {

	// Sample data
    monthlySales := []DataPoint{
        {X: 1, Y: 12000}, {X: 2, Y: 19000}, {X: 3, Y: 15000},
        {X: 4, Y: 21000}, {X: 5, Y: 18000}, {X: 6, Y: 24000},
    }

    storybook.Register("Chart", "Line Chart", 
        nil,
        func(controls map[string]*storybook.Control) app.UI {
            return NewChart(ChartTypeLine).
				Title("Monthly Sales Trend").
				Data(ChartData{
					Labels: []string{"Jan", "Feb", "Mar", "Apr", "May", "Jun"},
					Datasets: []Dataset{{
						Label: "Sales",
						Data:  monthlySales,
						BorderColor: "#4A90E2",
						Fill: true,
					}},
				}).
				Options(ChartOptions{
					Responsive: true,
					MaintainAspectRatio: false,
					Interaction: ChartInteraction{
						Mode: "nearest",
						Intersect: false,
					},
					Scales: ChartScales{
						Y: Axis{
							BeginAtZero: true,
							Title: AxisTitle{
								Display: true,
								Text: "Sales ($)",
							},
						},
					},
				}).
				Class("dashboard-card", "chart-line")
        },
    )

}
