// pkg/components/viz/engine.go
package viz

import (
    "fmt"
    "github.com/maxence-charriere/go-app/v10/pkg/app"
)

// Engine defines the chart rendering interface
type Engine interface {
    // Lifecycle
    Init(canvas app.HTMLCanvas) error
    Render(spec *Spec) error
    Update(data DataSet) error
    Resize(width, height int) error
    Destroy() error
    
    // Interactive
    HitTest(x, y float64) ([]Point, int, error)
    GetCanvas() app.HTMLCanvas
    
    // Performance
    SetMaxPoints(n int)
    GetMetrics() Metrics
}

// Metrics contains engine performance metrics
type Metrics struct {
    FrameRate     float64
    PointCount    int
    RenderTime    float64
    MemoryUsage   uint64
}

// pkg/components/viz/engine.go - update CanvasEngine

// CanvasEngine - Pure Go canvas implementation
type CanvasEngine struct {
    canvas    app.HTMLCanvas
    ctx       app.CanvasRenderingContext2D
    width     int
    height    int
    maxPoints int
    spec      *Spec
}

func NewCanvasEngine() *CanvasEngine {
    return &CanvasEngine{
        maxPoints: 10000,
    }
}

func (e *CanvasEngine) Init(canvas app.HTMLCanvas) error {
    e.canvas = canvas
    
    // Get 2D rendering context
    ctx, err := canvas.GetContext2D()
    if err != nil {
        return err
    }
    e.ctx = ctx
    
    // Get dimensions
    e.width = canvas.Get("width").Int()
    e.height = canvas.Get("height").Int()
    
    return nil
}

func (e *CanvasEngine) Render(spec *Spec) error {
    e.spec = spec
    
    // Clear canvas
    e.ctx.ClearRect(0, 0, float64(e.width), float64(e.height))
    
    // Draw based on chart type
    switch spec.Type {
    case ChartTypeBar:
        return e.renderBarChart(spec)
    case ChartTypeLine:
        return e.renderLineChart(spec)
    case ChartTypeScatter:
        return e.renderScatterChart(spec)
    case ChartTypePie:
        return e.renderPieChart(spec)
    default:
        return e.renderPlaceholder(spec)
    }
}

func (e *CanvasEngine) renderBarChart(spec *Spec) error {
    data := spec.Data
    if len(data.Series) == 0 || len(data.Series[0].Points) == 0 {
        return e.renderPlaceholder(spec)
    }
    
    // Set background
    e.ctx.SetFillStyle(spec.Theme.GetBackgroundColor())
    e.ctx.FillRect(0, 0, float64(e.width), float64(e.height))
    
    // Calculate margins
    margin := struct {
        top, right, bottom, left float64
    }{
        top:    40,
        right:  40,
        bottom: 60,
        left:   60,
    }
    
    chartWidth := float64(e.width) - margin.left - margin.right
    chartHeight := float64(e.height) - margin.top - margin.bottom
    
    // Get series data
    series := data.Series[0]
    points := series.Points
    
    if len(points) == 0 {
        return nil
    }
    
    // Find max value for scaling
    maxValue := 0.0
    for _, p := range points {
        if p.Y > maxValue {
            maxValue = p.Y
        }
    }
    if maxValue == 0 {
        maxValue = 1
    }
    
    // Calculate bar dimensions
    barCount := len(points)
    barWidth := (chartWidth / float64(barCount)) * (spec.Bar.Width / 100)
    barSpacing := (chartWidth / float64(barCount)) * ((100 - spec.Bar.Width) / 100)
    
    // Draw grid if enabled
    if spec.Axes.Y.Grid.Visible {
        e.drawGrid(margin, chartWidth, chartHeight, maxValue)
    }
    
    // Draw axes
    e.drawAxes(margin, chartWidth, chartHeight, maxValue, spec)
    
    // Draw bars
    for i, point := range points {
        x := margin.left + float64(i)*(barWidth+barSpacing) + barSpacing/2
        barHeight := (point.Y / maxValue) * chartHeight
        
        // Determine bar color
        barColor := series.Color
        if barColor == "" && spec.Theme != nil {
            colors := spec.Theme.GetColors()
            if len(colors) > 0 {
                barColor = colors[i%len(colors)]
            }
        }
        if barColor == "" {
            barColor = "#4f46e5" // Default indigo
        }
        
        // Draw bar
        e.ctx.SetFillStyle(barColor)
        
        // Apply border radius if specified
        if spec.Bar.BorderRadius > 0 {
            e.drawRoundedRect(
                x,
                margin.top+chartHeight-barHeight,
                barWidth,
                barHeight,
                spec.Bar.BorderRadius,
            )
            e.ctx.Fill()
        } else {
            e.ctx.FillRect(
                x,
                margin.top+chartHeight-barHeight,
                barWidth,
                barHeight,
            )
        }
        
        // Draw label if enabled
        if spec.Labels.Visible {
            e.ctx.SetFillStyle(spec.Labels.Color)
            e.ctx.SetFont(fmt.Sprintf("%dpx %s", spec.Labels.FontSize, spec.Theme.GetFontFamily()))
            e.ctx.SetTextAlign("center")
            
            label := fmt.Sprintf("%.0f", point.Y)
            if spec.Labels.Format != nil {
                label = spec.Labels.Format(point.Y)
            }
            
            e.ctx.FillText(
                label,
                x+barWidth/2,
                margin.top+chartHeight-barHeight-5,
            )
        }
    }
    
    // Draw x-axis labels
    e.ctx.SetFillStyle(spec.Theme.GetTextColor())
    e.ctx.SetFont(fmt.Sprintf("12px %s", spec.Theme.GetFontFamily()))
    e.ctx.SetTextAlign("center")
    
    for i, point := range points {
        x := margin.left + float64(i)*(barWidth+barSpacing) + barSpacing/2 + barWidth/2
        
        label := point.Label
        if label == "" && i < len(data.Labels) {
            label = data.Labels[i]
        }
        
        e.ctx.FillText(label, x, margin.top+chartHeight+20)
    }
    
    // Draw title if specified
    if spec.Title != "" {
        e.ctx.SetFillStyle(spec.Theme.GetTextColor())
        e.ctx.SetFont(fmt.Sprintf("16px %s", spec.Theme.GetFontFamily()))
        e.ctx.SetTextAlign("center")
        e.ctx.FillText(spec.Title, float64(e.width)/2, 25)
    }
    
    return nil
}

func (e *CanvasEngine) renderLineChart(spec *Spec) error {
    // Placeholder implementation
    return e.renderPlaceholder(spec)
}

func (e *CanvasEngine) renderScatterChart(spec *Spec) error {
    // Placeholder implementation
    return e.renderPlaceholder(spec)
}

func (e *CanvasEngine) renderPieChart(spec *Spec) error {
    // Placeholder implementation
    return e.renderPlaceholder(spec)
}

func (e *CanvasEngine) renderPlaceholder(spec *Spec) error {
    // Clear canvas
    e.ctx.ClearRect(0, 0, float64(e.width), float64(e.height))
    
    // Set background
    e.ctx.SetFillStyle("#f9fafb")
    e.ctx.FillRect(0, 0, float64(e.width), float64(e.height))
    
    // Draw placeholder text
    e.ctx.SetFillStyle("#6b7280")
    e.ctx.SetFont("14px sans-serif")
    e.ctx.SetTextAlign("center")
    e.ctx.SetTextBaseline("middle")
    
    chartType := "Chart"
    if spec != nil {
        chartType = string(spec.Type)
    }
    
    e.ctx.FillText(
        fmt.Sprintf("%s Rendering Not Yet Implemented", chartType),
        float64(e.width)/2,
        float64(e.height)/2,
    )
    
    return nil
}

func (e *CanvasEngine) drawGrid(margin struct{ top, right, bottom, left float64 }, chartWidth, chartHeight, maxValue float64) {
    // Draw horizontal grid lines
    e.ctx.SetStrokeStyle("#e5e7eb")
    e.ctx.SetLineWidth(1)
    
    for i := 0; i <= 5; i++ {
        y := margin.top + (float64(i)/5)*chartHeight
        e.ctx.BeginPath()
        e.ctx.MoveTo(margin.left, y)
        e.ctx.LineTo(margin.left+chartWidth, y)
        e.ctx.Stroke()
        
        // Draw y-axis labels
        value := maxValue * (1 - float64(i)/5)
        e.ctx.SetFillStyle("#6b7280")
        e.ctx.SetFont("11px sans-serif")
        e.ctx.SetTextAlign("right")
        e.ctx.FillText(fmt.Sprintf("%.0f", value), margin.left-10, y+4)
    }
}

func (e *CanvasEngine) drawAxes(margin struct{ top, right, bottom, left float64 }, chartWidth, chartHeight, maxValue float64, spec *Spec) {
    // Draw x-axis
    e.ctx.BeginPath()
    e.ctx.SetStrokeStyle("#9ca3af")
    e.ctx.SetLineWidth(1)
    e.ctx.MoveTo(margin.left, margin.top+chartHeight)
    e.ctx.LineTo(margin.left+chartWidth, margin.top+chartHeight)
    e.ctx.Stroke()
    
    // Draw y-axis
    e.ctx.BeginPath()
    e.ctx.MoveTo(margin.left, margin.top)
    e.ctx.LineTo(margin.left, margin.top+chartHeight)
    e.ctx.Stroke()
}

func (e *CanvasEngine) drawRoundedRect(x, y, width, height, radius float64) {
    e.ctx.BeginPath()
    e.ctx.MoveTo(x+radius, y)
    e.ctx.LineTo(x+width-radius, y)
    e.ctx.QuadraticCurveTo(x+width, y, x+width, y+radius)
    e.ctx.LineTo(x+width, y+height-radius)
    e.ctx.QuadraticCurveTo(x+width, y+height, x+width-radius, y+height)
    e.ctx.LineTo(x+radius, y+height)
    e.ctx.QuadraticCurveTo(x, y+height, x, y+height-radius)
    e.ctx.LineTo(x, y+radius)
    e.ctx.QuadraticCurveTo(x, y, x+radius, y)
    e.ctx.ClosePath()
}

func (e *CanvasEngine) Update(data DataSet) error {
    // Update data and re-render
    if e.spec != nil {
        e.spec.Data = data
        return e.Render(e.spec)
    }
    return nil
}

func (e *CanvasEngine) Resize(width, height int) error {
    e.width = width
    e.height = height
    e.canvas.Set("width", width)
    e.canvas.Set("height", height)
    
    if e.spec != nil {
        return e.Render(e.spec)
    }
    return nil
}

func (e *CanvasEngine) Destroy() error {
    // Clean up resources
    return nil
}

func (e *CanvasEngine) HitTest(x, y float64) ([]Point, int, error) {
    // Simple hit testing implementation
    // This would need to be expanded for actual hit testing
    return nil, 0, nil
}

func (e *CanvasEngine) GetCanvas() app.HTMLCanvas {
    return e.canvas
}

func (e *CanvasEngine) SetMaxPoints(n int) {
    e.maxPoints = n
}

func (e *CanvasEngine) GetMetrics() Metrics {
    return Metrics{
        PointCount: e.maxPoints,
    }
}
