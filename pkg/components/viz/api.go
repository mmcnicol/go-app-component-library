// pkg/components/viz/api.go
package viz

import (
    "github.com/maxence-charriere/go-app/v10/pkg/app"
)

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

// EngineType defines rendering engine
type EngineType string

const (
    EngineTypeAuto   EngineType = "auto"
    EngineTypeCanvas EngineType = "canvas"
    EngineTypeWebGL  EngineType = "webgl"
    EngineTypeSVG    EngineType = "svg"
)

// TooltipMode constants
type TooltipMode string

const (
    TooltipModeSingle   TooltipMode = "single"
    TooltipModeAll      TooltipMode = "all"
    TooltipModeNearest  TooltipMode = "nearest"
)

// RegressionType constants
type RegressionType string

const (
    RegressionTypeLinear      RegressionType = "linear"
    RegressionTypePolynomial  RegressionType = "polynomial"
    RegressionTypeExponential RegressionType = "exponential"
    RegressionTypeLogarithmic RegressionType = "logarithmic"
    RegressionTypePower       RegressionType = "power"
)

// ExportFormat constants
type ExportFormat string

const (
    ExportFormatPNG  ExportFormat = "png"
    ExportFormatSVG  ExportFormat = "svg"
    ExportFormatPDF  ExportFormat = "pdf"
    ExportFormatCSV  ExportFormat = "csv"
    ExportFormatJSON ExportFormat = "json"
)

// SamplingStrategy constants
type SamplingStrategy string

const (
    SamplingStrategyLTTB     SamplingStrategy = "lttb"
    SamplingStrategyEveryNth SamplingStrategy = "every-nth"
    SamplingStrategyMinMax   SamplingStrategy = "min-max"
    SamplingStrategyAverage  SamplingStrategy = "average"
)

// WhiskerType constants
type WhiskerType string

const (
    WhiskerTypeTukey      WhiskerType = "tukey"
    WhiskerTypeMinMax     WhiskerType = "minmax"
    WhiskerTypePercentile WhiskerType = "percentile"
)

// Theme interface
type Theme interface {
    GetBackgroundColor() string
    GetTextColor() string
    GetGridColor() string
    GetFontFamily() string
    GetColors() []string
}

// DefaultTheme returns the default theme
func DefaultTheme() Theme {
    return &defaultTheme{}
}

type defaultTheme struct{}

func (t *defaultTheme) GetBackgroundColor() string { return "#ffffff" }
func (t *defaultTheme) GetTextColor() string       { return "#333333" }
func (t *defaultTheme) GetGridColor() string       { return "#e5e7eb" }
func (t *defaultTheme) GetFontFamily() string      { return "sans-serif" }
func (t *defaultTheme) GetColors() []string {
    return []string{"#4f46e5", "#10b981", "#f59e0b", "#ef4444", "#8b5cf6", "#ec4899"}
}

// CustomTheme allows overriding theme properties
type CustomTheme struct {
    BaseTheme Theme
    Colors    []string
}

func (t *CustomTheme) GetBackgroundColor() string { return t.BaseTheme.GetBackgroundColor() }
func (t *CustomTheme) GetTextColor() string       { return t.BaseTheme.GetTextColor() }
func (t *CustomTheme) GetGridColor() string       { return t.BaseTheme.GetGridColor() }
func (t *CustomTheme) GetFontFamily() string      { return t.BaseTheme.GetFontFamily() }
func (t *CustomTheme) GetColors() []string {
    if len(t.Colors) > 0 {
        return t.Colors
    }
    return t.BaseTheme.GetColors()
}

// AccessibilityConfig defines accessibility options
type AccessibilityConfig struct {
    Enabled     bool
    Description string
    AriaLabel   string
    DataTableID string
}

// InteractiveConfig defines interactivity options
type InteractiveConfig struct {
    Enabled   bool
    Tooltip   TooltipConfig
    Zoom      ZoomConfig
    Pan       PanConfig
    Selection SelectionConfig
    OnClick   func(ctx app.Context, e app.Event, points []Point)
    OnHover   func(ctx app.Context, e app.Event, point *Point)
    OnZoom    func(ctx app.Context, axisRange AxisRange)
}

// TooltipConfig defines tooltip appearance and behavior
type TooltipConfig struct {
    Enabled          bool
    Mode             TooltipMode
    Intersect        bool
    IntersectDistance float64
    Format           func(point Point, series Series) string
    Background       string
    TextColor        string
    BorderColor      string
    BorderWidth      int
    Padding          int
}

// ZoomConfig defines zoom behavior
type ZoomConfig struct {
    Enabled bool
    Factor  float64
}

// PanConfig defines pan behavior
type PanConfig struct {
    Enabled bool
}

// SelectionConfig defines selection behavior
type SelectionConfig struct {
    Enabled bool
}

// AxisRange defines axis boundaries
type AxisRange struct {
    X [2]float64
    Y [2]float64
}

// AxesConfig defines axis configuration
type AxesConfig struct {
    X AxisConfig
    Y AxisConfig
}

// AxisConfig defines a single axis configuration
type AxisConfig struct {
    Visible     bool
    Title       string
    TitleColor  string
    LabelColor  string
    Grid        GridConfig
    BeginAtZero bool
    Stacked     bool
}

// GridConfig defines grid appearance
type GridConfig struct {
    Visible bool
    Color   string
    Width   float64
}

// BarConfig defines bar chart specific configuration
type BarConfig struct {
    Width        float64
    BorderRadius float64
    Grouped      bool
    Stacked      bool
    Horizontal   bool
}

// LabelsConfig defines data label configuration
type LabelsConfig struct {
    Visible  bool
    Position string
    FontSize int
    Color    string
    Format   func(value float64) string
}

// LegendConfig defines legend configuration
type LegendConfig struct {
    Visible  bool
    Position string
}

// Chart is the main component
type Chart struct {
    app.Compo
    spec *Spec
    engine Engine
}

// Spec defines the complete chart specification
type Spec struct {
    // Identity
    ID      string
    Title   string
    
    // Content
    Type    ChartType
    Data    DataSet
    
    // Appearance
    Theme   Theme
    Width   int
    Height  int
    
    // Features
    Interactive InteractiveConfig
    Accessible  AccessibilityConfig
    Animated    bool
    
    // Performance
    Engine      EngineType
    MaxPoints   int
    Sampling    SamplingStrategy
    
    // Chart specific configs
    Axes        AxesConfig
    Bar         BarConfig
    Labels      LabelsConfig
    Legend      LegendConfig
}

// New creates a new chart with the given spec
func New(spec Spec) *Chart {
    // Apply defaults
    if spec.Engine == "" {
        spec.Engine = EngineTypeAuto
    }
    if spec.Theme == nil {
        spec.Theme = DefaultTheme()
    }
    if spec.Width == 0 {
        spec.Width = 800
    }
    if spec.Height == 0 {
        spec.Height = 400
    }
    if spec.MaxPoints == 0 {
        spec.MaxPoints = 10000
    }
    if spec.Sampling == "" {
        spec.Sampling = SamplingStrategyLTTB
    }
    if spec.Interactive.Tooltip.IntersectDistance == 0 {
        spec.Interactive.Tooltip.IntersectDistance = 20
    }
    
    return &Chart{
        spec: &spec,
    }
}

// Chainable methods for fluent API
func (c *Chart) WithTitle(title string) *Chart {
    c.spec.Title = title
    return c
}

func (c *Chart) WithData(data DataSet) *Chart {
    c.spec.Data = data
    return c
}

func (c *Chart) WithTheme(theme Theme) *Chart {
    c.spec.Theme = theme
    return c
}

func (c *Chart) Interactive() *Chart {
    c.spec.Interactive.Enabled = true
    return c
}

func (c *Chart) Accessible() *Chart {
    c.spec.Accessible.Enabled = true
    return c
}

func (c *Chart) Animated() *Chart {
    c.spec.Animated = true
    return c
}

func (c *Chart) WithClickHandler(handler func(ctx app.Context, e app.Event, points []Point)) *Chart {
    c.spec.Interactive.OnClick = handler
    return c
}

func (c *Chart) WithHoverHandler(handler func(ctx app.Context, e app.Event, point *Point)) *Chart {
    c.spec.Interactive.OnHover = handler
    return c
}

func (c *Chart) WithZoomHandler(handler func(ctx app.Context, axisRange AxisRange)) *Chart {
    c.spec.Interactive.OnZoom = handler
    return c
}

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

func ScatterChart(data DataSet) *Chart {
    return New(Spec{Type: ChartTypeScatter, Data: data})
}

/*
// StreamingChart creates a chart optimized for streaming data
func StreamingChart(chartType ChartType) *StreamingChart {
    return &StreamingChart{
        Chart:       New(Spec{Type: chartType}),
        maxPoints:   100,
        updateRate:  100,
        dataBuffer:  NewDataBuffer(100),
    }
}
*/
