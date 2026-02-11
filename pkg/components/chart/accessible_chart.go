// pkg/components/chart/accessible_chart.go
package chart

import (
	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

type AccessibleChart struct {
    BaseChart
    ariaLabel     string
    ariaDescribedBy string
    longDescURL   string
    dataTableID   string
}

func (ac *AccessibleChart) Render() app.UI {
    return app.Div().
        Role("img").
        Aria("label", ac.ariaLabel).
        Aria("describedby", ac.ariaDescribedBy).
        Body(
            // Visual chart
            ac.BaseChart.Render(),
            
            // Hidden data table for screen readers
            ac.renderDataTable(),
            
            // Hidden descriptive text
            app.Div().
                ID(ac.ariaDescribedBy).
                Class("sr-only").
                Body(
                    ac.generateChartDescription(),
                ),
        )
}

func (ac *AccessibleChart) generateChartDescription() string {
    var desc strings.Builder
    
    desc.WriteString(fmt.Sprintf("%s chart showing ", ac.spec.Type))
    desc.WriteString(fmt.Sprintf("%d datasets with %d data points each. ",
        len(ac.spec.Data.Datasets), len(ac.spec.Data.Datasets[0].Data)))
    
    // Describe trends
    if ac.spec.Type == ChartTypeLine {
        trend := ac.calculateTrend()
        desc.WriteString(fmt.Sprintf("The data shows a %s trend. ", trend))
    }
    
    // Describe key statistics
    stats := ac.calculateStatistics()
    desc.WriteString(fmt.Sprintf("Minimum value: %.2f. Maximum value: %.2f. ",
        stats.Min, stats.Max))
    desc.WriteString(fmt.Sprintf("Average value: %.2f.", stats.Mean))
    
    return desc.String()
}

func (ac *AccessibleChart) renderDataTable() app.UI {
    return app.Table().
        ID(ac.dataTableID).
        Class("sr-only").
        Aria("hidden", "true").
        Body(
            app.Caption().Text(ac.ariaLabel),
            app.THead().Body(
                app.Tr().Body(
                    app.Th().Text("Dataset"),
                    app.Range(ac.spec.Data.Labels).Slice(func(i int) app.UI {
                        return app.Th().Text(ac.spec.Data.Labels[i])
                    }),
                ),
            ),
            app.TBody().Body(
                app.Range(ac.spec.Data.Datasets).Slice(func(i int) app.UI {
                    dataset := ac.spec.Data.Datasets[i]
                    return app.Tr().Body(
                        app.Th().Text(dataset.Label),
                        app.Range(dataset.Data).Slice(func(j int) app.UI {
                            return app.Td().Text(fmt.Sprintf("%.2f", dataset.Data[j].Y))
                        }),
                    )
                }),
            ),
        )
}
