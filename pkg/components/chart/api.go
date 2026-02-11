// pkg/components/chart/api.go
package chart

//import "github.com/maxence-charriere/go-app/v10/pkg/app"

// ChartConfig defines the public interface for customization
type ChartConfig struct {
	Title            string
	LineColor        string
	Thickness        float64
	IsStream         bool
	Capacity         int
	BoxWidth         float64
	BoxData          []BoxPlotStats
	HeatmapMatrix    [][]float64
	ColorScheme      string
	Opacity          float64
	PieData          []float64
	InnerRadiusRatio float64
	BarData          []float64  // Simple bar chart values
	BarLabels        []string   // Optional labels for bars
	BarColors        []string   // Optional custom colors per bar
}

type Option func(*ChartConfig)

// Utility options for the user
func WithTitle(t string) Option      { return func(c *ChartConfig) { c.Title = t } }
func WithColor(col string) Option    { return func(c *ChartConfig) { c.LineColor = col } }
func SetStreaming(cap int) Option    { return func(c *ChartConfig) { c.IsStream = true; c.Capacity = cap } }

// New creates a configured instance of your CanvasChart
func New(data []Point, opts ...Option) *CanvasChart {
	config := &ChartConfig{
		LineColor: "#4a90e2",
		Thickness: 2.0,
	}
	for _, opt := range opts {
		opt(config)
	}

	return &CanvasChart{
		currentPoints: data,
		config:        *config,
	}
}

// Bar chart specific option functions
func WithBarData(data []float64) Option {
	return func(c *ChartConfig) {
		c.BarData = data
		// Clear other chart types
		c.BoxData = nil
		c.PieData = nil
		c.HeatmapMatrix = nil
		c.IsStream = false
	}
}

func WithBarLabels(labels []string) Option {
	return func(c *ChartConfig) {
		c.BarLabels = labels
	}
}

func WithBarColors(colors []string) Option {
	return func(c *ChartConfig) {
		c.BarColors = colors
	}
}
