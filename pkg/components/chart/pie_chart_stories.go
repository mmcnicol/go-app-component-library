//go:build dev
// pkg/components/chart/pie_chart_stories.go
package chart

import (
    //"fmt"
    "github.com/maxence-charriere/go-app/v10/pkg/app"
    "github.com/mmcnicol/go-app-component-library/pkg/storybook"
)

// Use init() to auto-register when this package is imported
// pkg/components/chart/pie_chart_stories.go
func init() {
    categoryRevenue := []DataPoint{
        {Label: "Electronics", Value: 45000},
        {Label: "Clothing", Value: 32000},
        {Label: "Home Goods", Value: 28000},
        {Label: "Books", Value: 15000},
    }

    storybook.Register("Chart", "Pie Chart", 
        nil,
        func(controls map[string]*storybook.Control) app.UI {
            return NewAccessibleChart(ChartTypePie).  // Use AccessibleChart
                Title("Revenue by Category").
                Data(ChartData{
                    Datasets: []Dataset{{
                        Label: "Revenue",
                        Data: categoryRevenue,
                        BackgroundColor: []string{
                            "#FF6384", "#36A2EB", "#FFCE56", "#4BC0C0",
                        },
                    }},
                }).
                Options(ChartOptions{
                    Plugins: ChartPlugins{
                        Legend: LegendOptions{
                            Position: "right",
                        },
                    },
                }).
                Class("dashboard-card", "chart-pie")
        },
    )
}
