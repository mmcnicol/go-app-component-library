// pkg/components/chart/bar_chart_stories.go
package chart

import (
    "github.com/maxence-charriere/go-app/v10/pkg/app"
    "github.com/mmcnicol/go-app-component-library/pkg/storybook"
)

// Helper function to convert float64 slice to DataPoint slice
func convertToDataPoints(values []float64) []DataPoint {
    points := make([]DataPoint, len(values))
    for i, v := range values {
        points[i] = DataPoint{
            X: float64(i),
            Y: v,
        }
    }
    return points
}

// Use init() to auto-register when this package is imported
// pkg/components/chart/bar_chart_stories.go
func init() {
    userActivity := [][]float64{
        {8, 12, 15, 9, 14, 11, 13, 10},
        {5, 8, 10, 7, 9, 6, 11, 8},
        {12, 15, 18, 14, 16, 13, 17, 15},
    }

    storybook.Register("Chart", "Bar Chart", 
        nil,
        func(controls map[string]*storybook.Control) app.UI {
            return NewAccessibleChart(ChartTypeBar).  // Use AccessibleChart
                Title("User Activity by Hour").
                Data(ChartData{
                    Labels: []string{"9AM", "10AM", "11AM", "12PM", "1PM", "2PM", "3PM", "4PM"},
                    Datasets: []Dataset{
                        {
                            Label: "Week 1",
                            Data:  convertToDataPoints(userActivity[0]),
                            BackgroundColor: []string{"rgba(74, 144, 226, 0.7)"},
                            BorderColor: "rgba(74, 144, 226, 1)",
                            BorderWidth: 1,
                        },
                        {
                            Label: "Week 2",
                            Data:  convertToDataPoints(userActivity[1]),
                            BackgroundColor: []string{"rgba(255, 99, 132, 0.7)"},
                            BorderColor: "rgba(255, 99, 132, 1)",
                            BorderWidth: 1,
                        },
                        {
                            Label: "Week 3",
                            Data:  convertToDataPoints(userActivity[2]),
                            BackgroundColor: []string{"rgba(75, 192, 192, 0.7)"},
                            BorderColor: "rgba(75, 192, 192, 1)",
                            BorderWidth: 1,
                        },
                    },
                }).
                Options(ChartOptions{
                    Scales: ChartScales{
                        Y: Axis{
                            BeginAtZero: true,
                            Title: AxisTitle{
                                Display: true,
                                Text: "Active Users",
                            },
                        },
                        X: Axis{
                            Title: AxisTitle{
                                Display: true,
                                Text: "Time of Day",
                            },
                        },
                    },
                }).
                Class("dashboard-card", "chart-bar")
        },
    )

    // Add a simple test chart
    storybook.Register("Chart", "Test Bar Chart", 
        nil,
        func(controls map[string]*storybook.Control) app.UI {
            // Very simple test data
            return app.Div().Body(
                app.H2().Text("Simple Test Bar Chart"),
                app.P().Text("This should show a basic chart"),
                NewChart(ChartTypeBar).
                    Data(ChartData{
                        Labels: []string{"A", "B", "C", "D"},
                        Datasets: []Dataset{{
                            Label: "Test Data",
                            Data: []DataPoint{
                                {X: 0, Y: 10},
                                {X: 1, Y: 20},
                                {X: 2, Y: 15},
                                {X: 3, Y: 25},
                            },
                            BackgroundColor: []string{"#4A90E2"},
                        }},
                    }).
                    Class("test-chart").
                    Style("border", "1px solid #ccc").
                    Style("padding", "20px"),
            )
        },
    )
}
