// pkg/components/chart/base_chart.go
package chart

import (
    "math"
    "fmt"
    //"time"
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
    styles      map[string]string
    
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
        styles: make(map[string]string), // Initialize styles map
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
    
    // Simple approach - draw directly without engine
    ctx.Defer(func(ctx app.Context) {
        if len(bc.spec.Data.Datasets) == 0 {
            fmt.Println("No datasets to render")
            return
        }
        
        fmt.Printf("Rendering %s chart with %d datasets (simple approach)\n", 
            bc.spec.Type, len(bc.spec.Data.Datasets))
        
        // Draw directly using JavaScript
        bc.drawSimpleChart(ctx)
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
    div := app.Div().
        ID(bc.containerID).
        Class(append([]string{"chart-container"}, bc.classes...)...)
    
    // Apply custom styles
    for name, value := range bc.styles {
        div = div.Style(name, value)
    }
    
    // Add default styles if not overridden
    if _, hasWidth := bc.styles["width"]; !hasWidth {
        div = div.Style("width", "100%")
    }
    if _, hasHeight := bc.styles["height"]; !hasHeight {
        div = div.Style("height", "400px")
    }
    if _, hasPosition := bc.styles["position"]; !hasPosition {
        div = div.Style("position", "relative")
    }
    
    return div.Body(
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

/*
// Style sets inline CSS styles for the chart container
func (bc *BaseChart) Style(name, value string) *BaseChart {
    // We'll need to store styles and apply them in Render()
    // For now, we can add them to classes or handle differently
    // Let's add a style map to BaseChart
    if bc.spec.Metadata == nil {
        bc.spec.Metadata = make(map[string]interface{})
    }
    if styles, ok := bc.spec.Metadata["styles"].(map[string]string); ok {
        styles[name] = value
    } else {
        bc.spec.Metadata["styles"] = map[string]string{name: value}
    }
    return bc
}
*/

func (bc *BaseChart) Style(name, value string) *BaseChart {
    bc.styles[name] = value
    return bc
}

func (bc *BaseChart) drawSimpleChart(ctx app.Context) {
    canvasID := bc.containerID + "-canvas"
    
    jsCode := fmt.Sprintf(`
        console.log('Drawing simple chart for:', '%s');
        
        const canvas = document.getElementById('%s');
        if (!canvas) {
            console.error('Canvas not found');
            return;
        }
        
        // Ensure dimensions
        if (canvas.width === 0 || canvas.height === 0) {
            canvas.width = canvas.clientWidth || 800;
            canvas.height = canvas.clientHeight || 400;
        }
        
        const ctx = canvas.getContext('2d');
        if (!ctx) {
            console.error('No 2D context');
            return;
        }
        
        // Clear and draw background
        ctx.clearRect(0, 0, canvas.width, canvas.height);
        ctx.fillStyle = '#ffffff';
        ctx.fillRect(0, 0, canvas.width, canvas.height);
        
        // Draw title
        ctx.fillStyle = '#333';
        ctx.font = 'bold 18px Arial';
        ctx.textAlign = 'center';
        ctx.fillText('%s Chart - Simple Render', canvas.width / 2, 30);
        
        // Draw a test pattern
        ctx.fillStyle = '#4A90E2';
        for (let i = 0; i < 5; i++) {
            const x = 100 + i * 100;
            const y = 200;
            const width = 60;
            const height = 50 + i * 20;
            
            ctx.fillRect(x, y - height, width, height);
            
            // Draw label
            ctx.fillStyle = '#000';
            ctx.font = '12px Arial';
            ctx.fillText('Bar ' + (i+1), x + width/2, y + 20);
            ctx.fillStyle = '#4A90E2';
        }
        
        console.log('Simple chart drawn successfully');
    `, canvasID, canvasID, bc.spec.Type)
    
    app.Window().Call("eval", jsCode)
}
