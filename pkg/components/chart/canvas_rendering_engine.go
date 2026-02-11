// pkg/components/chart/canvas_rendering_engine.go
package chart

import (
    "fmt"
    "math"
    "github.com/maxence-charriere/go-app/v10/pkg/app"
)

// CanvasRenderer implements ChartEngine
type CanvasRenderer struct {
    canvas    app.UI
    width     int
    height    int
    pixelRatio float64
    ctx       app.Value
}

func NewCanvasRenderer(containerID string) (*CanvasRenderer, error) {
    cr := &CanvasRenderer{
        pixelRatio: 1.0,
        width:      800,
        height:     400,
    }
    
    // Create canvas element
    cr.canvas = app.Canvas().
        ID(containerID + "-canvas").
        Class("chart-canvas").
        Style("width", "100%").
        Style("height", "100%").
        Style("display", "block").
        OnMount(func(ctx app.Context) {
            // Get the actual canvas element after mount
            ctx.JSSrc().Get("document").Call("getElementById", containerID+"-canvas")
        })
    
    return cr, nil
}

func (cr *CanvasRenderer) Render(chart ChartSpec) error {
    // Set up canvas dimensions
    cr.setupCanvas(chart.Options.Responsive)
    
    // Get the canvas context in JavaScript
    app.Window().Call("setTimeout", app.FuncOf(func(this app.Value, args []app.Value) interface{} {
        canvas := app.Window().Get("document").Call("getElementById", chart.ContainerID+"-canvas")
        if !canvas.IsUndefined() {
            ctx := canvas.Call("getContext", "2d")
            cr.ctx = ctx
            
            // Clear canvas
            ctx.Call("clearRect", 0, 0, cr.width, cr.height)
            
            // Draw chart based on type
            switch chart.Type {
            case ChartTypeLine:
                cr.renderLineChart(chart, ctx)
            case ChartTypeBar:
                cr.renderBarChart(chart, ctx)
            case ChartTypePie:
                cr.renderPieChart(chart, ctx)
            default:
                cr.renderBarChart(chart, ctx) // Default to bar chart
            }
        }
        return nil
    }), 100)
    
    return nil
}

func (cr *CanvasRenderer) setupCanvas(responsive bool) {
    if responsive {
        cr.width = 800
        cr.height = 400
    } else {
        cr.width = 800
        cr.height = 400
    }
}

func (cr *CanvasRenderer) renderBarChart(chart ChartSpec, ctx app.Value) error {
    if ctx.IsUndefined() {
        return fmt.Errorf("canvas context not available")
    }
    
    data := chart.Data
    if len(data.Datasets) == 0 {
        return fmt.Errorf("no datasets available")
    }
    
    // Calculate scales
    xScale := cr.calculateXScale(data)
    yScale := cr.calculateYScale(data)
    
    // Clear and draw background
    ctx.Call("clearRect", 0, 0, cr.width, cr.height)
    ctx.Set("fillStyle", "#ffffff")
    ctx.Call("fillRect", 0, 0, cr.width, cr.height)
    
    // Draw grid
    cr.drawGrid(ctx, xScale, yScale, chart.Options.Grid)
    
    // Draw axes
    cr.drawAxes(ctx, xScale, yScale, chart.Options.Axes, data.Labels)
    
    // Draw bars for each dataset
    numDatasets := len(data.Datasets)
    numPoints := len(data.Datasets[0].Data)
    barWidth := (xScale.Convert(1) - xScale.Convert(0)) / float64(numDatasets) * 0.8
    
    for datasetIdx, dataset := range data.Datasets {
        // Set color
        color := "#4A90E2"
        if len(dataset.BackgroundColor) > 0 {
            color = dataset.BackgroundColor[0]
        }
        ctx.Set("fillStyle", color)
        
        for pointIdx, point := range dataset.Data {
            xPos := xScale.Convert(float64(pointIdx)) + 
                   (float64(datasetIdx) * barWidth) - 
                   (float64(numDatasets) * barWidth / 2) + 
                   (barWidth / 2)
            yPos := yScale.Convert(point.Y)
            barHeight := yScale.Convert(0) - yPos
            
            // Draw bar
            ctx.Call("fillRect", 
                xPos - barWidth/2, 
                yPos, 
                barWidth, 
                barHeight)
            
            // Draw border
            ctx.Set("strokeStyle", "#333")
            ctx.Set("lineWidth", 1)
            ctx.Call("strokeRect", 
                xPos - barWidth/2, 
                yPos, 
                barWidth, 
                barHeight)
        }
    }
    
    // Draw legend
    if chart.Options.Plugins.Legend.Display {
        cr.drawLegend(ctx, data.Datasets, chart.Options.Plugins.Legend)
    }
    
    return nil
}

// Simplified scale calculations
func (cr *CanvasRenderer) calculateXScale(data ChartData) Scale {
    numLabels := len(data.Labels)
    if numLabels == 0 && len(data.Datasets) > 0 {
        numLabels = len(data.Datasets[0].Data)
    }
    
    return Scale{
        Min: 0,
        Max: float64(numLabels),
        Convert: func(value float64) float64 {
            margin := 80.0
            plotWidth := float64(cr.width) - 2*margin
            return margin + (value * plotWidth / float64(numLabels))
        },
    }
}

func (cr *CanvasRenderer) calculateYScale(data ChartData) Scale {
    min := math.MaxFloat64
    max := -math.MaxFloat64
    
    for _, dataset := range data.Datasets {
        for _, point := range dataset.Data {
            if point.Y < min {
                min = point.Y
            }
            if point.Y > max {
                max = point.Y
            }
        }
    }
    
    if min == math.MaxFloat64 {
        min = 0
        max = 100
    }
    
    // Add padding
    rangePadding := (max - min) * 0.1
    min -= rangePadding
    max += rangePadding
    
    if min < 0 && !data.Options.Scales.Y.BeginAtZero {
        min = 0
    }
    
    finalMin := min
    finalMax := max
    
    return Scale{
        Min: finalMin,
        Max: finalMax,
        Convert: func(value float64) float64 {
            margin := 60.0
            plotHeight := float64(cr.height) - 2*margin
            return float64(cr.height) - margin - ((value - finalMin) / (finalMax - finalMin)) * plotHeight
        },
    }
}

// Add this method to ChartData to fix compilation error
func (cd ChartData) GetOptions() ChartOptions {
    return ChartOptions{}
}

func (cr *CanvasRenderer) drawLineDataset(dataset Dataset, xScale, yScale Scale, datasetIndex int) {
    if cr.ctx.Value.IsUndefined() {
        return
    }
    
    cr.ctx.Set("lineWidth", float64(dataset.BorderWidth))
    cr.ctx.Set("strokeStyle", dataset.BorderColor)
    cr.ctx.Call("beginPath")
    
    // Draw line
    for i, point := range dataset.Data {
        x := xScale.Convert(point.X)
        y := yScale.Convert(point.Y)
        
        if i == 0 {
            cr.ctx.Call("moveTo", x, y)
        } else {
            // Apply tension for smooth curves
            if dataset.Tension > 0 {
                cr.drawBezierCurve(dataset.Data, i, xScale, yScale, dataset.Tension)
            } else {
                cr.ctx.Call("lineTo", x, y)
            }
        }
        
        // Draw points
        if dataset.PointRadius > 0 {
            color := dataset.BorderColor
            if len(dataset.BackgroundColor) > i {
                color = dataset.BackgroundColor[i]
            }
            cr.drawPoint(x, y, dataset.PointRadius, color)
        }
    }
    
    cr.ctx.Call("stroke")
    
    // Fill area under line if needed
    if dataset.Fill {
        cr.fillAreaUnderLine(dataset, xScale, yScale)
    }
}

func (cr *CanvasRenderer) drawPoint(x, y float64, radius int, color string) {
    if cr.ctx.Value.IsUndefined() {
        return
    }
    
    cr.ctx.Set("fillStyle", color)
    cr.ctx.Call("beginPath")
    cr.ctx.Call("arc", x, y, float64(radius), 0, 2*math.Pi)
    cr.ctx.Call("fill")
}

func (cr *CanvasRenderer) drawBezierCurve(points []DataPoint, i int, xScale, yScale Scale, tension float64) {
    if cr.ctx.Value.IsUndefined() || i <= 0 {
        return
    }
    
    p0 := points[i-1]
    p1 := points[i]
    
    x0 := xScale.Convert(p0.X)
    y0 := yScale.Convert(p0.Y)
    x1 := xScale.Convert(p1.X)
    y1 := yScale.Convert(p1.Y)
    
    // Calculate control points (simplified)
    cp1x := x0 + (x1-x0)*tension
    cp1y := y0
    cp2x := x1 - (x1-x0)*tension
    cp2y := y1
    
    cr.ctx.Call("bezierCurveTo", cp1x, cp1y, cp2x, cp2y, x1, y1)
}

func (cr *CanvasRenderer) fillAreaUnderLine(dataset Dataset, xScale, yScale Scale) {
    if cr.ctx.Value.IsUndefined() || len(dataset.Data) == 0 {
        return
    }
    
    cr.ctx.Set("fillStyle", dataset.BorderColor+"33") // Add transparency
    cr.ctx.Call("beginPath")
    
    // Start at first point
    firstPoint := dataset.Data[0]
    cr.ctx.Call("moveTo", xScale.Convert(firstPoint.X), yScale.Convert(firstPoint.Y))
    
    // Draw line through all points
    for i := 1; i < len(dataset.Data); i++ {
        point := dataset.Data[i]
        cr.ctx.Call("lineTo", xScale.Convert(point.X), yScale.Convert(point.Y))
    }
    
    // Close path to baseline
    lastPoint := dataset.Data[len(dataset.Data)-1]
    cr.ctx.Call("lineTo", xScale.Convert(lastPoint.X), yScale.Convert(yScale.Max))
    cr.ctx.Call("lineTo", xScale.Convert(firstPoint.X), yScale.Convert(yScale.Max))
    cr.ctx.Call("closePath")
    cr.ctx.Call("fill")
}

func (cr *CanvasRenderer) drawGrid(xScale, yScale Scale, grid GridOptions) {
    if !grid.Display || cr.ctx.Value.IsUndefined() {
        return
    }
    
    cr.ctx.Set("strokeStyle", "#e0e0e0")
    cr.ctx.Set("lineWidth", 0.5)
    
    // Draw horizontal grid lines (5 lines)
    yRange := yScale.Max - yScale.Min
    for i := 0; i <= 5; i++ {
        y := yScale.Min + (yRange * float64(i) / 5)
        yPos := yScale.Convert(y)
        cr.ctx.Call("beginPath")
        cr.ctx.Call("moveTo", 50, yPos)
        cr.ctx.Call("lineTo", float64(cr.width)-50, yPos)
        cr.ctx.Call("stroke")
    }
    
    // Draw vertical grid lines - use number of labels if available
    numLabels := 10 // default
    xRange := xScale.Max - xScale.Min
    for i := 0; i <= numLabels; i++ {
        x := xScale.Min + (xRange * float64(i) / float64(numLabels))
        xPos := xScale.Convert(x)
        cr.ctx.Call("beginPath")
        cr.ctx.Call("moveTo", xPos, 50)
        cr.ctx.Call("lineTo", xPos, float64(cr.height)-50)
        cr.ctx.Call("stroke")
    }
}

func (cr *CanvasRenderer) drawAxes(xScale, yScale Scale, axes AxesOptions) {
    if !axes.Display || cr.ctx.Value.IsUndefined() {
        return
    }
    
    cr.ctx.Set("strokeStyle", "#000000")
    cr.ctx.Set("lineWidth", 1)
    cr.ctx.Set("font", "12px Arial")
    cr.ctx.Set("fillStyle", "#000000")
    
    // Draw Y axis
    cr.ctx.Call("beginPath")
    cr.ctx.Call("moveTo", 50, 50)
    cr.ctx.Call("lineTo", 50, float64(cr.height)-50)
    cr.ctx.Call("stroke")
    
    // Draw X axis
    cr.ctx.Call("beginPath")
    cr.ctx.Call("moveTo", 50, float64(cr.height)-50)
    cr.ctx.Call("lineTo", float64(cr.width)-50, float64(cr.height)-50)
    cr.ctx.Call("stroke")
}

func (cr *CanvasRenderer) drawLegend(datasets []Dataset, legend LegendOptions) {
    if !legend.Display || cr.ctx.Value.IsUndefined() {
        return
    }
    
    startX := float64(cr.width) - 200
    startY := 20
    boxSize := 15
    spacing := 25
    
    for i, dataset := range datasets {
        y := startY + i*spacing
        
        // Draw color box
        color := dataset.BorderColor
        if color == "" && len(dataset.BackgroundColor) > 0 {
            color = dataset.BackgroundColor[0]
        }
        
        cr.ctx.Set("fillStyle", color)
        cr.ctx.Call("fillRect", startX, float64(y), float64(boxSize), float64(boxSize))
        
        // Draw label
        cr.ctx.Set("fillStyle", "#000000")
        cr.ctx.Call("fillText", dataset.Label, startX+float64(boxSize)+10, float64(y+12))
    }
}

// Stub methods for other chart types
func (cr *CanvasRenderer) renderBarChart(chart ChartSpec) error {
    return fmt.Errorf("bar chart not yet implemented")
}

func (cr *CanvasRenderer) renderPieChart(chart ChartSpec) error {
    return fmt.Errorf("pie chart not yet implemented")
}

func (cr *CanvasRenderer) renderScatterChart(chart ChartSpec) error {
    return fmt.Errorf("scatter chart not yet implemented")
}

func (cr *CanvasRenderer) renderHeatmapChart(chart ChartSpec) error {
    return fmt.Errorf("heatmap chart not yet implemented")
}

// Additional required methods for ChartEngine interface
func (cr *CanvasRenderer) Update(data ChartData) error {
    // In a real implementation, this would update the chart data
    return fmt.Errorf("update not yet implemented")
}

func (cr *CanvasRenderer) Destroy() error {
    // Clean up resources
    return nil
}

func (cr *CanvasRenderer) GetCanvas() app.UI {
    return cr.canvas
}
