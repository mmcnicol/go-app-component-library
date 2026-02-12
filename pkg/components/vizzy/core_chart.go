// pkg/components/vizzy/core_chart.go
package core

import (
    "github.com/maxence-charriere/go-app/v10/pkg/app"
)

// ChartEngine defines the rendering contract
type ChartEngine interface {
    // Lifecycle
    Initialize(width, height int, canvas app.Value) error
    Render(chart *Chart) error
    Update(chart *Chart) error
    Resize(width, height int) error
    Destroy() error
    
    // Interactive
    HitTest(x, y float64) ([]Point, int, error)
    GetCanvas() app.Value
    
    // Performance
    SetMaxPoints(n int)
    GetMetrics() Metrics
}

// Chart represents a pure data model with no rendering logic
type Chart struct {
    // Identity
    ID    string
    Title string
    
    // Content
    Type ChartType
    Data DataSet
    
    // Appearance
    Theme  Theme
    Width  int
    Height int
    
    // Features
    Interactive InteractiveConfig
    Accessible  AccessibilityConfig
    Animated    bool
    
    // Performance
    MaxPoints int
    Sampling  SamplingStrategy
    
    // Chart specific configs
    Axes   AxesConfig
    Bar    BarConfig
    Labels LabelsConfig
    Legend LegendConfig
    
    // Rendering engine (optional - can be set later)
    engine ChartEngine
}

// ChartBuilder provides a fluent API for building charts
type ChartBuilder struct {
    chart *Chart
}

func NewChartBuilder() *ChartBuilder {
    return &ChartBuilder{
        chart: &Chart{
            Width:      800,
            Height:     400,
            MaxPoints:  10000,
            Sampling:   SamplingStrategyLTTB,
            Theme:      DefaultTheme(),
        },
    }
}
