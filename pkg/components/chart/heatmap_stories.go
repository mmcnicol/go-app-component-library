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
	userActivity := [][]float64{
        {8, 12, 15, 9, 14, 11, 13, 10},
        {5, 8, 10, 7, 9, 6, 11, 8},
        {12, 15, 18, 14, 16, 13, 17, 15},
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
