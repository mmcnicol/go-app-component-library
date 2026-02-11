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

// CanvasEngine - Pure Go canvas implementation
type CanvasEngine struct {
    canvas    app.HTMLCanvas
    width     int
    height    int
    maxPoints int
    spec      *Spec
    jsCanvas  app.Value // Store the JS canvas value for direct manipulation
}

func NewCanvasEngine() *CanvasEngine {
    return &CanvasEngine{
        maxPoints: 10000,
    }
}

func (e *CanvasEngine) Init(canvas app.HTMLCanvas) error {
    e.canvas = canvas
    e.jsCanvas = canvas.JSValue()
    
    // Get dimensions
    e.width = e.jsCanvas.Get("width").Int()
    e.height = e.jsCanvas.Get("height").Int()
    
    return nil
}

// InitWithValue initializes the engine with a JS value (from GetElementByID)
func (e *CanvasEngine) InitWithValue(canvasValue app.Value) error {
    e.jsCanvas = canvasValue
    
    // Get dimensions
    e.width = e.jsCanvas.Get("width").Int()
    e.height = e.jsCanvas.Get("height").Int()
    
    return nil
}

func (e *CanvasEngine) Render(spec *Spec) error {
    e.spec = spec
    
    // Get canvas context
    ctx := e.jsCanvas.Call("getContext", "2d")
    if ctx.IsNull() {
        return fmt.Errorf("failed to get 2d context")
    }
    
    // Clear canvas
    ctx.Call("clearRect", 0, 0, e.width, e.height)
    
    // Draw based on chart type
    switch spec.Type {
    case ChartTypeBar:
        return e.renderBarChart(spec, ctx)
    case ChartTypeLine:
        return e.renderLineChart(spec, ctx)
    case ChartTypeScatter:
        return e.renderScatterChart(spec, ctx)
    case ChartTypePie:
        return e.renderPieChart(spec, ctx)
    default:
        return e.renderPlaceholder(spec, ctx)
    }
}

func (e *CanvasEngine) renderBarChart(spec *Spec, ctx app.Value) error {
    data := spec.Data
    if len(data.Series) == 0 || len(data.Series[0].Points) == 0 {
        return e.renderPlaceholder(spec, ctx)
    }
    
    // Set background
    ctx.Set("fillStyle", spec.Theme.GetBackgroundColor())
    ctx.Call("fillRect", 0, 0, e.width, e.height)
    
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
        e.drawGrid(ctx, margin, chartWidth, chartHeight, maxValue)
    }
    
    // Draw axes
    e.drawAxes(ctx, margin, chartWidth, chartHeight, maxValue, spec)
    
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
        ctx.Set("fillStyle", barColor)
        
        // Apply border radius if specified
        if spec.Bar.BorderRadius > 0 {
            e.drawRoundedRect(ctx,
                x,
                margin.top+chartHeight-barHeight,
                barWidth,
                barHeight,
                spec.Bar.BorderRadius,
            )
            ctx.Call("fill")
        } else {
            ctx.Call("fillRect",
                x,
                margin.top+chartHeight-barHeight,
                barWidth,
                barHeight,
            )
        }
        
        // Draw label if enabled
        if spec.Labels.Visible {
            ctx.Set("fillStyle", spec.Labels.Color)
            ctx.Set("font", fmt.Sprintf("%dpx %s", spec.Labels.FontSize, spec.Theme.GetFontFamily()))
            ctx.Set("textAlign", "center")
            
            label := fmt.Sprintf("%.0f", point.Y)
            if spec.Labels.Format != nil {
                label = spec.Labels.Format(point.Y)
            }
            
            ctx.Call("fillText",
                label,
                x+barWidth/2,
                margin.top+chartHeight-barHeight-5,
            )
        }
    }
    
    // Draw x-axis labels
    ctx.Set("fillStyle", spec.Theme.GetTextColor())
    ctx.Set("font", fmt.Sprintf("12px %s", spec.Theme.GetFontFamily()))
    ctx.Set("textAlign", "center")
    
    for i, point := range points {
        x := margin.left + float64(i)*(barWidth+barSpacing) + barSpacing/2 + barWidth/2
        
        label := point.Label
        if label == "" && i < len(data.Labels) {
            label = data.Labels[i]
        }
        
        ctx.Call("fillText", label, x, margin.top+chartHeight+20)
    }
    
    // Draw title if specified
    if spec.Title != "" {
        ctx.Set("fillStyle", spec.Theme.GetTextColor())
        ctx.Set("font", fmt.Sprintf("16px %s", spec.Theme.GetFontFamily()))
        ctx.Set("textAlign", "center")
        ctx.Call("fillText", spec.Title, float64(e.width)/2, 25)
    }
    
    return nil
}

func (e *CanvasEngine) renderLineChart(spec *Spec, ctx app.Value) error {
    // Placeholder implementation
    return e.renderPlaceholder(spec, ctx)
}

func (e *CanvasEngine) renderScatterChart(spec *Spec, ctx app.Value) error {
    // Placeholder implementation
    return e.renderPlaceholder(spec, ctx)
}

func (e *CanvasEngine) renderPieChart(spec *Spec, ctx app.Value) error {
    // Placeholder implementation
    return e.renderPlaceholder(spec, ctx)
}

func (e *CanvasEngine) renderPlaceholder(spec *Spec, ctx app.Value) error {
    // Clear canvas
    ctx.Call("clearRect", 0, 0, e.width, e.height)
    
    // Set background
    ctx.Set("fillStyle", "#f9fafb")
    ctx.Call("fillRect", 0, 0, e.width, e.height)
    
    // Draw placeholder text
    ctx.Set("fillStyle", "#6b7280")
    ctx.Set("font", "14px sans-serif")
    ctx.Set("textAlign", "center")
    ctx.Set("textBaseline", "middle")
    
    chartType := "Chart"
    if spec != nil {
        chartType = string(spec.Type)
    }
    
    ctx.Call("fillText",
        fmt.Sprintf("%s Rendering Not Yet Implemented", chartType),
        float64(e.width)/2,
        float64(e.height)/2,
    )
    
    return nil
}

func (e *CanvasEngine) drawGrid(ctx app.Value, margin struct{ top, right, bottom, left float64 }, chartWidth, chartHeight, maxValue float64) {
    // Draw horizontal grid lines
    ctx.Set("strokeStyle", "#e5e7eb")
    ctx.Set("lineWidth", 1)
    
    for i := 0; i <= 5; i++ {
        y := margin.top + (float64(i)/5)*chartHeight
        ctx.Call("beginPath")
        ctx.Call("moveTo", margin.left, y)
        ctx.Call("lineTo", margin.left+chartWidth, y)
        ctx.Call("stroke")
        
        // Draw y-axis labels
        value := maxValue * (1 - float64(i)/5)
        ctx.Set("fillStyle", "#6b7280")
        ctx.Set("font", "11px sans-serif")
        ctx.Set("textAlign", "right")
        ctx.Call("fillText", fmt.Sprintf("%.0f", value), margin.left-10, y+4)
    }
}

func (e *CanvasEngine) drawAxes(ctx app.Value, margin struct{ top, right, bottom, left float64 }, chartWidth, chartHeight, maxValue float64, spec *Spec) {
    // Draw x-axis
    ctx.Call("beginPath")
    ctx.Set("strokeStyle", "#9ca3af")
    ctx.Set("lineWidth", 1)
    ctx.Call("moveTo", margin.left, margin.top+chartHeight)
    ctx.Call("lineTo", margin.left+chartWidth, margin.top+chartHeight)
    ctx.Call("stroke")
    
    // Draw y-axis
    ctx.Call("beginPath")
    ctx.Call("moveTo", margin.left, margin.top)
    ctx.Call("lineTo", margin.left, margin.top+chartHeight)
    ctx.Call("stroke")
}

func (e *CanvasEngine) drawRoundedRect(ctx app.Value, x, y, width, height, radius float64) {
    ctx.Call("beginPath")
    ctx.Call("moveTo", x+radius, y)
    ctx.Call("lineTo", x+width-radius, y)
    ctx.Call("quadraticCurveTo", x+width, y, x+width, y+radius)
    ctx.Call("lineTo", x+width, y+height-radius)
    ctx.Call("quadraticCurveTo", x+width, y+height, x+width-radius, y+height)
    ctx.Call("lineTo", x+radius, y+height)
    ctx.Call("quadraticCurveTo", x, y+height, x, y+height-radius)
    ctx.Call("lineTo", x, y+radius)
    ctx.Call("quadraticCurveTo", x, y, x+radius, y)
    ctx.Call("closePath")
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
    e.jsCanvas.Set("width", width)
    e.jsCanvas.Set("height", height)
    
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
