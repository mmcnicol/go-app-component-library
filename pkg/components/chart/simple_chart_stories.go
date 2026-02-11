//go:build dev
// pkg/components/chart/simple_chart_stories.go
package chart

import (
    "github.com/maxence-charriere/go-app/v10/pkg/app"
    "github.com/mmcnicol/go-app-component-library/pkg/storybook"
)

func init() {
    // Test with SimpleChart first
    storybook.Register("Chart", "Simple Bar Chart Test", 
        nil,
        func(controls map[string]*storybook.Control) app.UI {
            return app.Div().Body(
                app.H3().Text("Simple Bar Chart Test"),
                app.P().Text("This should definitely show a bar chart"),
                NewSimpleChart(ChartTypeBar).
                    Title("Simple Test Chart").
                    Data(ChartData{
                        Labels: []string{"Jan", "Feb", "Mar", "Apr", "May", "Jun"},
                        Datasets: []Dataset{{
                            Label: "Sales",
                            Data: []DataPoint{
                                {X: 0, Y: 12000},
                                {X: 1, Y: 19000},
                                {X: 2, Y: 15000},
                                {X: 3, Y: 21000},
                                {X: 4, Y: 18000},
                                {X: 5, Y: 24000},
                            },
                        }},
                    }).
                    Width(800).
                    Height(400),
            )
        },
    )
    
    // Also test the original but fixed
    storybook.Register("Chart", "Fixed Bar Chart", 
        nil,
        func(controls map[string]*storybook.Control) app.UI {
            userActivity := [][]float64{
                {8, 12, 15, 9, 14, 11, 13, 10},
                {5, 8, 10, 7, 9, 6, 11, 8},
                {12, 15, 18, 14, 16, 13, 17, 15},
            }
            
            return NewChart(ChartTypeBar).
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
                Class("dashboard-card", "chart-bar")
        },
    )
}
