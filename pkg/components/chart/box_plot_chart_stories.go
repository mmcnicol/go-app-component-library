//go:build dev
// pkg/components/chart/box_plot_chart_stories.go
package chart

import (
    "github.com/maxence-charriere/go-app/v10/pkg/app"
    "github.com/mmcnicol/go-app-component-library/pkg/storybook"
)

// Use init() to auto-register when this package is imported
func init() {

    storybook.Register("Chart", "Box Plot Chart", 
        nil,
        func(controls map[string]*storybook.Control) app.UI {
            return NewChart(ChartTypeBoxPlot).
                Title("Performance Distribution").
                Data(ChartData{
                    Data: [][]float64{
                        {12, 15, 18, 22, 25, 28, 30, 32, 35},
                        {8, 12, 16, 20, 24, 28, 32, 36, 40},
                        {20, 22, 24, 26, 28, 30, 32, 34, 36},
                    },
                    Labels: []string{"Team A", "Team B", "Team C"},
                }).
                Class("dashboard-card", "chart-boxplot")
        },
    )

}
