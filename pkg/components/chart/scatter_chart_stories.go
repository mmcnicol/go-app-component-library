//go:build dev
// pkg/components/chart/scatter_chart_stories.go
package chart

import (
    "github.com/maxence-charriere/go-app/v10/pkg/app"
    "github.com/mmcnicol/go-app-component-library/pkg/storybook"
)

// Use init() to auto-register when this package is imported
func init() {

    storybook.Register("Chart", "Scatter Chart", 
        nil,
        func(controls map[string]*storybook.Control) app.UI {
            return NewChart(ChartTypeScatter).
                Title("Correlation Analysis").
                WithRegression(RegressionTypeLinear, 2).
                Data(ChartData{
                    Datasets: []Dataset{{
                        Label: "Data Points",
                        Data:  generateRandomData(100),
                        PointRadius: 3,
                    }},
                }).
                Class("dashboard-card", "chart-scatter")
        },
    )

}
