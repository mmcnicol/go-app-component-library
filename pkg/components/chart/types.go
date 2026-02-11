// pkg/components/chart/types.go
package chart

import (
	"github.com/maxence-charriere/go-app/v10/pkg/app"
	"time"
	//"sync"
	//"math"
	"fmt"
	//"strings"
	//"sort"
)

// Chart types
type ChartType string

const (
	ChartTypeLine     ChartType = "line"
	ChartTypeBar      ChartType = "bar"
	ChartTypePie      ChartType = "pie"
	ChartTypeScatter  ChartType = "scatter"
	ChartTypeHeatmap  ChartType = "heatmap"
	ChartTypeBoxPlot  ChartType = "boxplot"
)

// Engine types
type EngineType string

const (
	EngineTypeCanvas EngineType = "canvas"
	EngineTypeSVG    EngineType = "svg"
	EngineTypeWebGL  EngineType = "webgl"
	EngineTypeHybrid EngineType = "hybrid"
)

// Data structures
type DataPoint struct {
	X     float64
	Y     float64
	Label string
	Value float64
}

type AxisRange struct {
	X [2]float64
	Y [2]float64
}

// Chart options
type ChartOptions struct {
	Responsive        bool
	MaintainAspectRatio bool
	Plugins           ChartPlugins
	Scales            ChartScales
	Grid              GridOptions
	Axes              AxesOptions
	Interaction       ChartInteraction
	Tooltips          TooltipOptions
	Legend            LegendOptions
	ColorScale        ColorScale
}

type ChartPlugins struct {
	Legend  LegendOptions
	Tooltip TooltipOptions
}

type ChartScales struct {
	X Axis
	Y Axis
}

type Axis struct {
	BeginAtZero bool
	Stacked     bool
	Title       AxisTitle
}

type AxisTitle struct {
	Display bool
	Text    string
}

type GridOptions struct {
	Display bool
}

type AxesOptions struct {
	Display bool
}

type ChartInteraction struct {
	Mode      string
	Intersect bool
}

type TooltipOptions struct {
	Enabled          bool
	IntersectDistance float64
	Callbacks        TooltipCallbacks
}

type TooltipCallbacks struct {
	Label func(context TooltipContext) string
}

type TooltipContext struct {
	label  string
	value  float64
	total  float64
}

func (tc TooltipContext) Label() string { return tc.label }
func (tc TooltipContext) Value() float64 { return tc.value }
func (tc TooltipContext) Total() float64 { return tc.total }

type LegendOptions struct {
	Display  bool
	Position string
}

type ColorScale struct {
	Min    float64
	Max    float64
	Colors []string
}

// Scale types
type Scale struct {
	Min     float64
	Max     float64
	Convert func(value float64) float64
	Invert  func(pixel float64) float64
}

// Regression types
type RegressionType string

const (
	RegressionTypeLinear      RegressionType = "linear"
	RegressionTypePolynomial  RegressionType = "polynomial"
	RegressionTypeExponential RegressionType = "exponential"
	RegressionTypeLogarithmic RegressionType = "logarithmic"
)

type RegressionResult struct {
	Coefficients []float64
	Equation     string
	RSquared     float64
	Predict      func(x float64) float64
}

// Box plot types
type WhiskerType string

const (
	WhiskerTypeTukey      WhiskerType = "tukey"
	WhiskerTypeMinMax     WhiskerType = "minmax"
	WhiskerTypePercentile WhiskerType = "percentile"
)

type BoxPlotStats struct {
	Min          float64
	Q1           float64
	Median       float64
	Q3           float64
	Max          float64
	LowerWhisker float64
	UpperWhisker float64
	Outliers     []float64
	Mean         float64
}

// Tooltip types
type TooltipPosition struct {
	X float64
	Y float64
}

type TooltipFormatter interface {
	Format(point *DataPoint, datasetIndex int, data ChartData) string
}

type DefaultTooltipFormatter struct{}

func (dtf DefaultTooltipFormatter) Format(point *DataPoint, datasetIndex int, data ChartData) string {
	return fmt.Sprintf("Value: %.2f", point.Y)
}

// Zoom types
type TransformMatrix struct {
	A, B, C, D, E, F float64
}

// Animation types
type Animation struct {
	Duration time.Duration
	Easing   string
	Callback func()
}

// Statistics for accessible chart
type ChartStatistics struct {
	Min  float64
	Max  float64
	Mean float64
}

type ChartData struct {
    Labels   []string
    Datasets []Dataset
    Data     [][]float64  // For box plot
    XLabels  []string     // For heatmap
    YLabels  []string     // For heatmap
    Metadata map[string]interface{}
}


// Chart engine interface
type ChartEngine interface {
    Render(chart ChartSpec) error
    Update(data ChartData) error
    Destroy() error
    GetCanvas() app.UI
}

// Chart specification
type ChartSpec struct {
    Type        ChartType
    Data        ChartData
    Options     ChartOptions
    ContainerID string
    Engine      EngineType // "canvas", "svg", "webgl", "hybrid"
}

/*
// Chart data structure
type ChartData struct {
    Labels   []string
    Datasets []Dataset
    Metadata map[string]interface{}
}
*/

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
