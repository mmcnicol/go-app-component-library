// pkg/components/chart/base_chart.go
package chart

import (
	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

// Base chart component
type BaseChart struct {
    app.Compo
    
    // State
    spec        ChartSpec
    containerID string
    engine      ChartEngine
    isRendered  bool
    
    // Refs
    canvasRef   app.Ref
    tooltipRef  app.Ref
    
    // Event handlers
    onPointClick    func(point DataPoint, datasetIndex int)
    onZoom          func(domain AxisRange)
    onHover         func(point DataPoint, datasetIndex int)
}

func (bc *BaseChart) calculateTrend() string {
    // Simplified trend calculation
    if len(bc.spec.Data.Datasets) == 0 || len(bc.spec.Data.Datasets[0].Data) < 2 {
        return "insufficient data"
    }
    
    first := bc.spec.Data.Datasets[0].Data[0].Y
    last := bc.spec.Data.Datasets[0].Data[len(bc.spec.Data.Datasets[0].Data)-1].Y
    
    if last > first {
        return "positive"
    } else if last < first {
        return "negative"
    }
    return "stable"
}

func (bc *BaseChart) calculateStatistics() ChartStatistics {
    stats := ChartStatistics{
        Min:  math.MaxFloat64,
        Max:  -math.MaxFloat64,
        Mean: 0,
    }
    
    var sum float64
    var count int
    
    for _, dataset := range bc.spec.Data.Datasets {
        for _, point := range dataset.Data {
            if point.Y < stats.Min {
                stats.Min = point.Y
            }
            if point.Y > stats.Max {
                stats.Max = point.Y
            }
            sum += point.Y
            count++
        }
    }
    
    if count > 0 {
        stats.Mean = sum / float64(count)
    } else {
        stats.Min = 0
        stats.Max = 0
    }
    
    return stats
}

func (bc *BaseChart) updateScales(xMin, xMax, yMin, yMax float64) {
    // Update scale domains
    // Implementation depends on your scale system
}

func (bc *BaseChart) WithRegression(regType RegressionType, degree int) *ScatterChart {
    // Return a scatter chart with regression
    return &ScatterChart{
        BaseChart: *bc,
        showRegression: true,
    }
}
