// pkg/components/chart/canvas_rendering_engine.go
package chart

import (
	"fmt"
    "math"
	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

// Custom type for canvas context since go-app doesn't expose CanvasRenderingContext2D directly
type CanvasRenderingContext2D struct {
    app.Value
}

type CanvasRenderer struct {
    canvas    app.HTMLCanvas
    ctx       CanvasRenderingContext2D
    width     int
    height    int
    pixelRatio float64
    animations []Animation
}

func NewCanvasRenderer(containerID string) (*CanvasRenderer, error) {
    cr := &CanvasRenderer{
        pixelRatio: 1.0,
    }
    
    // Create canvas element
    cr.canvas = app.Canvas().
        ID(containerID + "-canvas").
        Style("display", "block")
    
    // Get rendering context - in go-app, we use JS to get the context
    cr.ctx = CanvasRenderingContext2D{Value: app.Window().Get("document").Call("createElement", "canvas").Call("getContext", "2d")}
    
    return cr, nil
}

func (cr *CanvasRenderer) Render(chart ChartSpec) error {
    // Set up canvas dimensions
    cr.setupCanvas(chart.Options.Responsive)
    
    // Clear canvas
    cr.clearCanvas()
    
    // Draw chart based on type
    switch chart.Type {
    case ChartTypeLine:
        return cr.renderLineChart(chart)
    case ChartTypeBar:
        return cr.renderBarChart(chart)
    case ChartTypePie:
        return cr.renderPieChart(chart)
    case ChartTypeScatter:
        return cr.renderScatterChart(chart)
    case ChartTypeHeatmap:
        return cr.renderHeatmapChart(chart)
    default:
        return fmt.Errorf("unsupported chart type: %v", chart.Type)
    }
}

func (cr *CanvasRenderer) setupCanvas(responsive bool) {
    if responsive {
        // Get parent container dimensions
        // Implementation depends on your layout
        cr.width = 800
        cr.height = 400
    } else {
        cr.width = 800
        cr.height = 400
    }
    
    // Set canvas dimensions
    cr.canvas.SetAttr("width", cr.width)
    cr.canvas.SetAttr("height", cr.height)
    cr.canvas.SetStyle("width", fmt.Sprintf("%dpx", cr.width))
    cr.canvas.SetStyle("height", fmt.Sprintf("%dpx", cr.height))
}

func (cr *CanvasRenderer) clearCanvas() {
    cr.ctx.Call("clearRect", 0, 0, cr.width, cr.height)
}

func (cr *CanvasRenderer) renderLineChart(chart ChartSpec) error {
    // Calculate scales
    xScale := cr.calculateXScale(chart.Data)
    yScale := cr.calculateYScale(chart.Data)
    
    // Draw grid
    cr.drawGrid(xScale, yScale, chart.Options.Grid)
    
    // Draw axes
    cr.drawAxes(xScale, yScale, chart.Options.Axes)
    
    // Draw datasets
    for i, dataset := range chart.Data.Datasets {
        cr.drawLineDataset(dataset, xScale, yScale, i)
    }
    
    // Draw legend
    if chart.Options.Legend.Display {
        cr.drawLegend(chart.Data.Datasets, chart.Options.Legend)
    }
    
    return nil
}

func (cr *CanvasRenderer) calculateXScale(data ChartData) Scale {
    // Simplified implementation
    return Scale{
        Min: 0,
        Max: float64(len(data.Labels) - 1),
        Convert: func(value float64) float64 {
            // Linear mapping
            rangeSize := float64(cr.width - 100) // Account for margins
            return 50 + (value / float64(len(data.Labels)-1)) * rangeSize
        },
        Invert: func(pixel float64) float64 {
            rangeSize := float64(cr.width - 100)
            return ((pixel - 50) / rangeSize) * float64(len(data.Labels)-1)
        },
    }
}

func (cr *CanvasRenderer) calculateYScale(data ChartData) Scale {
    // Find min and max values
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
    
    // Add some padding
    rangePadding := (max - min) * 0.1
    min -= rangePadding
    max += rangePadding
    
    if min == max {
        min = 0
        max = min + 1
    }
    
    return Scale{
        Min: min,
        Max: max,
        Convert: func(value float64) float64 {
            // Y axis is inverted (0 at top)
            rangeSize := float64(cr.height - 100)
            return float64(cr.height) - 50 - ((value - min) / (max - min)) * rangeSize
        },
        Invert: func(pixel float64) float64 {
            rangeSize := float64(cr.height - 100)
            return min + ((float64(cr.height) - 50 - pixel) / rangeSize) * (max - min)
        },
    }
}

func (cr *CanvasRenderer) drawLineDataset(dataset Dataset, xScale, yScale Scale, datasetIndex int) {
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
            cr.drawPoint(x, y, dataset.PointRadius, dataset.BackgroundColor[i])
        }
    }
    
    cr.ctx.Call("stroke")
    
    // Fill area under line if needed
    if dataset.Fill {
        cr.fillAreaUnderLine(dataset, xScale, yScale)
    }
}

func (cr *CanvasRenderer) drawPoint(x, y float64, radius int, color string) {
    cr.ctx.Set("fillStyle", color)
    cr.ctx.Call("beginPath")
    cr.ctx.Call("arc", x, y, float64(radius), 0, 2*math.Pi)
    cr.ctx.Call("fill")
}

func (cr *CanvasRenderer) drawBezierCurve(points []DataPoint, i int, xScale, yScale Scale, tension float64) {
    // Simplified bezier curve implementation
    p0 := points[i-1]
    p1 := points[i]
    
    x0 := xScale.Convert(p0.X)
    y0 := yScale.Convert(p0.Y)
    x1 := xScale.Convert(p1.X)
    y1 := yScale.Convert(p1.Y)
    
    // Calculate control points
    cp1x := x0 + (x1-x0)*tension
    cp1y := y0
    cp2x := x1 - (x1-x0)*tension
    cp2y := y1
    
    cr.ctx.Call("bezierCurveTo", cp1x, cp1y, cp2x, cp2y, x1, y1)
}

func (cr *CanvasRenderer) fillAreaUnderLine(dataset Dataset, xScale, yScale Scale) {
    cr.ctx.Set("fillStyle", dataset.BorderColor+"33") // Add transparency
    cr.ctx.Call("beginPath")
    
    // Start at first point
    firstPoint := dataset.Data[0]
    cr.ctx.Call("moveTo", xScale.Convert(firstPoint.X), yScale.Convert(firstPoint.Y))
    
    // Draw line
    for i := 1; i < len(dataset.Data); i++ {
        point := dataset.Data[i]
        if dataset.Tension > 0 && i < len(dataset.Data)-1 {
            // Need proper bezier implementation for filling
            cr.ctx.Call("lineTo", xScale.Convert(point.X), yScale.Convert(point.Y))
        } else {
            cr.ctx.Call("lineTo", xScale.Convert(point.X), yScale.Convert(point.Y))
        }
    }
    
    // Close path to baseline
    lastPoint := dataset.Data[len(dataset.Data)-1]
    cr.ctx.Call("lineTo", xScale.Convert(lastPoint.X), yScale.Convert(yScale.Max))
    cr.ctx.Call("lineTo", xScale.Convert(firstPoint.X), yScale.Convert(yScale.Max))
    cr.ctx.Call("closePath")
    cr.ctx.Call("fill")
}

func (cr *CanvasRenderer) drawGrid(xScale, yScale Scale, grid GridOptions) {
    if !grid.Display {
        return
    }
    
    cr.ctx.Set("strokeStyle", "#e0e0e0")
    cr.ctx.Set("lineWidth", 0.5)
    
    // Draw horizontal grid lines
    for y := yScale.Min; y <= yScale.Max; y += (yScale.Max - yScale.Min) / 5 {
        yPos := yScale.Convert(y)
        cr.ctx.Call("beginPath")
        cr.ctx.Call("moveTo", 50, yPos)
        cr.ctx.Call("lineTo", float64(cr.width)-50, yPos)
        cr.ctx.Call("stroke")
    }
    
    // Draw vertical grid lines
    for x := xScale.Min; x <= xScale.Max; x += (xScale.Max - xScale.Min) / float64(len(xScale)) {
        xPos := xScale.Convert(x)
        cr.ctx.Call("beginPath")
        cr.ctx.Call("moveTo", xPos, 50)
        cr.ctx.Call("lineTo", xPos, float64(cr.height)-50)
        cr.ctx.Call("stroke")
    }
}

func (cr *CanvasRenderer) drawAxes(xScale, yScale Scale, axes AxesOptions) {
    if !axes.Display {
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
    if !legend.Display {
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
