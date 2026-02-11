// pkg/components/chart/base_chart.go
package chart

import (
    "math"
    "fmt"
    "time"
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
    classes     []string
    
    // Managers
    tooltipManager *TooltipManager
    zoomPanManager *ZoomPanManager
    
    // Event handlers
    onPointClick    func(point DataPoint, datasetIndex int)
    onZoom          func(domain AxisRange)
    onHover         func(point DataPoint, datasetIndex int)
}

// NewChart creates a new chart of the specified type
func NewChart(chartType ChartType) *BaseChart {
    id := fmt.Sprintf("chart-%s", GenerateID())
    return &BaseChart{
        spec: ChartSpec{
            Type: chartType,
            Data: ChartData{},
            Options: ChartOptions{
                Responsive: true,
                MaintainAspectRatio: false,
                Plugins: ChartPlugins{
                    Title: TitleOptions{
                        Display: true,
                        Text:    string(chartType) + " Chart",
                    },
                    Legend: LegendOptions{
                        Display: true,
                        Position: "top",
                    },
                    Tooltip: TooltipOptions{
                        Enabled: true,
                        IntersectDistance: 10,
                    },
                },
            },
            ContainerID: id,
            Engine: EngineTypeCanvas,
        },
        containerID: id,
    }
}

// Title sets the chart title
func (bc *BaseChart) Title(title string) *BaseChart {
    bc.spec.Options.Plugins.Title = TitleOptions{
        Display: true,
        Text:    title,
    }
    return bc
}

// Data sets the chart data
func (bc *BaseChart) Data(data ChartData) *BaseChart {
    bc.spec.Data = data
    return bc
}

// Options sets the chart options
func (bc *BaseChart) Options(options ChartOptions) *BaseChart {
    bc.spec.Options = options
    return bc
}

// Initialize managers in constructor or setup method
func (bc *BaseChart) setupManagers() {
    bc.tooltipManager = NewTooltipManager(bc)
    bc.zoomPanManager = NewZoomPanManager(bc)
}

// pkg/components/chart/base_chart.go
// Update the OnMount method completely:

func (bc *BaseChart) OnMount(ctx app.Context) {
    bc.isRendered = false
    bc.setupManagers()
    
    // Initialize the render engine using the wrapper
    engine, err := NewChartEngineWrapper(bc.containerID)
    if err != nil {
        fmt.Printf("Error creating chart engine: %v\n", err)
        return
    }
    
    bc.engine = engine
    
    // Use Defer to ensure rendering happens after component is fully mounted
    ctx.Defer(func(ctx app.Context) {
        // Double-check we have data
        if len(bc.spec.Data.Datasets) == 0 {
            fmt.Println("No datasets to render")
            return
        }
        
        // Wait a bit to ensure DOM is ready
        ctx.After(100*time.Millisecond, func(ctx app.Context) {
            fmt.Printf("Rendering %s chart with %d datasets\n", 
                bc.spec.Type, len(bc.spec.Data.Datasets))
            
            err := bc.engine.Render(bc.spec)
            if err != nil {
                fmt.Printf("Error rendering chart: %v\n", err)
            } else {
                bc.isRendered = true
                fmt.Println("Chart rendered successfully")
            }
        })
    })
}

// OnNav is called when the component is navigated to
func (bc *BaseChart) OnNav(ctx app.Context) {
    if !bc.isRendered && bc.engine != nil && len(bc.spec.Data.Datasets) > 0 {
        err := bc.engine.Render(bc.spec)
        if err == nil {
            bc.isRendered = true
        }
    }
}

// Render renders the chart
func (bc *BaseChart) Render() app.UI {
    return app.Div().
        ID(bc.containerID).
        Class(append([]string{"chart-container"}, bc.classes...)...).
        Style("position", "relative").
        Style("width", "100%").
        Style("height", "400px").
        Body(
            // Canvas element for drawing
            func() app.UI {
                if bc.engine != nil {
                    return bc.engine.GetCanvas()
                }
                return app.Canvas().
                    ID(bc.containerID + "-canvas").
                    Class("chart-canvas").
                    Style("width", "100%").
                    Style("height", "100%").
                    Style("display", "block")
            }(),
            // Tooltip
            func() app.UI {
                if bc.tooltipManager != nil {
                    return bc.tooltipManager.GetTooltipUI()
                }
                return app.Div()
            }(),
        )
}

// Add this method to render the chart
func (bc *BaseChart) renderChart() {
    if bc.engine != nil && !bc.isRendered && len(bc.spec.Data.Datasets) > 0 {
        err := bc.engine.Render(bc.spec)
        if err == nil {
            bc.isRendered = true
        }
    }
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

// Class sets CSS classes
func (bc *BaseChart) Class(classes ...string) *BaseChart {
    bc.classes = append(bc.classes, classes...)
    return bc
}

func (bc *BaseChart) WithRegression(regType RegressionType, degree int) *ScatterChart {
    // Return a scatter chart with regression
    return &ScatterChart{
        BaseChart: *bc,
        showRegression: true,
    }
}

// Helper methods for box plot chart
func (bc *BaseChart) calculatePercentile(sorted []float64, percentile float64) float64 {
    if len(sorted) == 0 {
        return 0
    }
    
    if len(sorted) == 1 {
        return sorted[0]
    }
    
    index := (percentile / 100) * float64(len(sorted)-1)
    lower := int(index)
    upper := lower + 1
    
    if upper >= len(sorted) {
        return sorted[lower]
    }
    
    weight := index - float64(lower)
    return sorted[lower]*(1-weight) + sorted[upper]*weight
}

func (bc *BaseChart) calculateMean(data []float64) float64 {
    if len(data) == 0 {
        return 0
    }
    
    var sum float64
    for _, v := range data {
        sum += v
    }
    return sum / float64(len(data))
}
