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

// Theme interface
type Theme interface {
    GetBackgroundColor() string
    GetTextColor() string
    GetGridColor() string
    GetFontFamily() string
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

// CustomTheme allows overriding theme properties
type CustomTheme struct {
    BaseTheme Theme
    Colors    []string
}

func (t *CustomTheme) GetBackgroundColor() string { return t.BaseTheme.GetBackgroundColor() }
func (t *CustomTheme) GetTextColor() string       { return t.BaseTheme.GetTextColor() }
func (t *CustomTheme) GetGridColor() string       { return t.BaseTheme.GetGridColor() }
func (t *CustomTheme) GetFontFamily() string      { return t.BaseTheme.GetFontFamily() }

// AccessibilityConfig defines accessibility options
type AccessibilityConfig struct {
    Enabled     bool
    Description string
}

// InteractiveConfig defines interactivity options
type InteractiveConfig struct {
    Enabled  bool
    Tooltip  TooltipConfig
    Zoom     ZoomConfig
    Pan      PanConfig
    OnClick  func(ctx app.Context, e app.Event, points []Point)
    OnHover  func(ctx app.Context, e app.Event, point *Point)
}

// TooltipConfig defines tooltip appearance and behavior
type TooltipConfig struct {
    Enabled     bool
    Mode        TooltipMode
    Intersect   bool
    Format      func(point Point, series Series) string
    Background  string
    TextColor   string
    BorderColor string
    BorderWidth int
    Padding     int
}

// TooltipMode constants
type TooltipMode string

const (
    TooltipModeSingle   TooltipMode = "single"
    TooltipModeAll      TooltipMode = "all"
    TooltipModeNearest  TooltipMode = "nearest"
)

// ZoomConfig defines zoom behavior
type ZoomConfig struct {
    Enabled bool
}

// PanConfig defines pan behavior
type PanConfig struct {
    Enabled bool
}

// SelectionConfig defines selection behavior
type SelectionConfig struct {
    Enabled bool
}

// Event represents a chart event
type Event struct {
    Type string
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
