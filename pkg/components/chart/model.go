// pkg/components/chart/model.go
package chart

import (
	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

// Chart engine interface
type ChartEngine interface {
    Render(chart ChartSpec) error
    Update(data ChartData) error
    Destroy() error
    GetCanvas() app.HTMLCanvas
}

// Chart specification
type ChartSpec struct {
    Type        ChartType
    Data        ChartData
    Options     ChartOptions
    ContainerID string
    Engine      EngineType // "canvas", "svg", "webgl", "hybrid"
}

// Chart data structure
type ChartData struct {
    Labels   []string
    Datasets []Dataset
    Metadata map[string]interface{}
}

// Dataset definition
type Dataset struct {
    Label           string
    Data            []DataPoint
    BackgroundColor []string
    BorderColor     string
    BorderWidth     int
    Fill            bool
    Tension         float64 // for line smoothing
    PointRadius     int
}

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

// Specialized chart components
type LineChart struct {
    BaseChart
    showArea      bool
    stepped       bool
    showPoints    bool
}

type BarChart struct {
    BaseChart
    horizontal    bool
    stacked       bool
    barPercentage float64
}

type PieChart struct {
    BaseChart
    donut         bool
    cutout        string // percentage or pixel
}

type ScatterChart struct {
    BaseChart
    showLine      bool
    showRegression bool
}

type HeatmapChart struct {
    BaseChart
    colorScale    ColorScale
    showValues    bool
}
