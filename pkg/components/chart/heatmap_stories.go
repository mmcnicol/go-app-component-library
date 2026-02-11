//go:build dev
// pkg/components/chart/heatmap_stories.go
package chart

import (
	"github.com/maxence-charriere/go-app/v10/pkg/app"
	"github.com/mmcnicol/go-app-component-library/pkg/storybook"
)

// Use init() to auto-register when this package is imported
func init() {

	// Sample data
    categoryRevenue := []DataPoint{
        {Label: "Electronics", Value: 45000},
        {Label: "Clothing", Value: 32000},
        {Label: "Home Goods", Value: 28000},
        {Label: "Books", Value: 15000},
    }

    storybook.Register("Chart", "Heatmap", 
        nil,
        func(controls map[string]*storybook.Control) app.UI {
            return NewChart(ChartTypeHeatmap).
				Title("Correlation Matrix").
				Data(ChartData{
					Data: userActivity,
					XLabels: []string{"Feature A", "Feature B", "Feature C", "Feature D"},
					YLabels: []string{"Metric 1", "Metric 2", "Metric 3"},
				}).
				Options(ChartOptions{
					Plugins: ChartPlugins{
						Legend: LegendOptions{
							Display: false,
						},
					},
					ColorScale: ColorScale{
						Min: 0,
						Max: 1,
						Colors: []string{"#FFFFFF", "#4A90E2"},
					},
				}).
				Class("dashboard-card", "chart-heatmap")
        },
    )

}
