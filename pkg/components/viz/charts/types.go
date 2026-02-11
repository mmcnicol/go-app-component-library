// pkg/components/viz/charts/types.go
package viz

// ChartType defines available chart types
type ChartType string

const (
    ChartTypeLine       ChartType = "line"
    ChartTypeBar        ChartType = "bar"
    ChartTypeGroupedBar ChartType = "grouped-bar"
    ChartTypeStackedBar ChartType = "stacked-bar"
    ChartTypePie        ChartType = "pie"
    ChartTypeDonut      ChartType = "donut"
    ChartTypeScatter    ChartType = "scatter"
    ChartTypeBubble     ChartType = "bubble"
    ChartTypeArea       ChartType = "area"
    ChartTypeHeatmap    ChartType = "heatmap"
    ChartTypeBoxPlot    ChartType = "box-plot"
    ChartTypeViolin     ChartType = "violin"
    ChartTypeHistogram  ChartType = "histogram"
    ChartTypeRadar      ChartType = "radar"
    ChartTypeCandlestick ChartType = "candlestick"
    ChartTypeGantt      ChartType = "gantt"
    ChartTypeTreeMap    ChartType = "treemap"
    ChartTypeNetwork    ChartType = "network"
    ChartType3DSurface  ChartType = "3d-surface"
)

// ChartType shortcuts for common use cases
func LineChart(data DataSet) *Chart {
    return New(Spec{Type: ChartTypeLine, Data: data})
}

func BarChart(data DataSet) *Chart {
    return New(Spec{Type: ChartTypeBar, Data: data})
}

func PieChart(data DataSet) *Chart {
    return New(Spec{Type: ChartTypePie, Data: data})
}

// ... etc
