// pkg/components/viz/engine.go
package viz

import (
    "fmt"
    "math"
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
    xScale    func(float64) float64
    yScale    func(float64) float64
    margin    Margin
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
    data := spec.Data
    if len(data.Series) == 0 || len(data.Series[0].Points) == 0 {
        return e.renderPlaceholder(spec, ctx)
    }
    
    // Set background
    ctx.Set("fillStyle", spec.Theme.GetBackgroundColor())
    ctx.Call("fillRect", 0, 0, e.width, e.height)
    
    // Calculate scales using all points from all series
    allPoints := []Point{}
    for _, series := range data.Series {
        allPoints = append(allPoints, series.Points...)
    }
    e.calculateScales(spec, allPoints)
    
    chartWidth := float64(e.width) - e.margin.Left - e.margin.Right
    chartHeight := float64(e.height) - e.margin.Top - e.margin.Bottom
    
    // Draw grid
    if spec.Axes.Y.Grid.Visible {
        e.drawLineChartGrid(ctx, chartWidth, chartHeight)
    }
    
    // Draw axes
    e.drawLineChartAxes(ctx, chartWidth, chartHeight, spec)
    
    // Draw each series
    for seriesIdx, series := range data.Series {
        points := series.Points
        
        if len(points) < 2 {
            continue
        }
        
        // Apply downsampling for large datasets
        if len(points) > e.maxPoints {
            points = e.downsampleLTTB(points, e.maxPoints)
        }
        
        // Determine series color
        seriesColor := series.Color
        if seriesColor == "" {
            colors := spec.Theme.GetColors()
            seriesColor = colors[seriesIdx%len(colors)]
        }
        
        // Draw the line
        ctx.Set("strokeStyle", seriesColor)
        ctx.Set("lineWidth", series.Stroke.Width)
        if series.Stroke.Width == 0 {
            ctx.Set("lineWidth", 2)
        }
        
        /*
        // Set line dash if specified
        if len(series.Stroke.Dash) > 0 {
            ctx.Set("lineDash", series.Stroke.Dash)
        } else {
            ctx.Set("lineDash", []float64{})
        }
        */
        
        // Draw the path
        ctx.Call("beginPath")
        
        for i, point := range points {
            x := e.xScale(point.X)
            y := e.yScale(point.Y)
            
            if i == 0 {
                ctx.Call("moveTo", x, y)
            } else {
                // Smooth curves if tension > 0
                if series.Tension > 0 && i < len(points)-1 {
                    e.drawSmoothSegment(ctx, points, i, series.Tension)
                } else {
                    ctx.Call("lineTo", x, y)
                }
            }
        }
        
        ctx.Call("stroke")
        
        // Fill area under line if enabled
        if series.Fill {
            ctx.Set("fillStyle", seriesColor+"40") // Add 25% opacity
            ctx.Call("beginPath")
            
            // Start from bottom-left
            firstPoint := points[0]
            ctx.Call("moveTo", e.xScale(firstPoint.X), e.yScale(0))
            
            // Draw line path
            for _, point := range points {
                ctx.Call("lineTo", e.xScale(point.X), e.yScale(point.Y))
            }
            
            // Close back to bottom-right
            lastPoint := points[len(points)-1]
            ctx.Call("lineTo", e.xScale(lastPoint.X), e.yScale(0))
            ctx.Call("closePath")
            ctx.Call("fill")
        }
        
        // Draw points if enabled
        if series.PointSize > 0 {
            for _, point := range points {
                e.drawPoint(ctx, 
                    e.xScale(point.X), 
                    e.yScale(point.Y), 
                    series.PointSize, 
                    seriesColor,
                    series.PointStyle,
                )
            }
        }
    }
    
    // Draw title
    if spec.Title != "" {
        ctx.Set("fillStyle", spec.Theme.GetTextColor())
        ctx.Set("font", fmt.Sprintf("16px %s", spec.Theme.GetFontFamily()))
        ctx.Set("textAlign", "center")
        ctx.Call("fillText", spec.Title, float64(e.width)/2, 25)
    }
    
    // Draw legend if enabled
    if spec.Legend.Visible && len(data.Series) > 1 {
        e.drawLegend(ctx, spec)
    }
    
    return nil
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

type Margin struct {
    Top, Right, Bottom, Left float64
}

func (e *CanvasEngine) calculateScales(spec *Spec, points []Point) {
    // Find min/max values
    if len(points) == 0 {
        return
    }
    
    minX, maxX := points[0].X, points[0].X
    minY, maxY := points[0].Y, points[0].Y
    
    for _, p := range points {
        minX = math.Min(minX, p.X)
        maxX = math.Max(maxX, p.X)
        minY = math.Min(minY, p.Y)
        maxY = math.Max(maxY, p.Y)
    }
    
    // Add padding
    if maxX == minX {
        maxX = minX + 1
    }
    if maxY == minY {
        maxY = minY + 1
    }
    
    // Add 10% padding to Y axis
    yPadding := (maxY - minY) * 0.1
    minY -= yPadding
    maxY += yPadding
    
    if spec.Axes.Y.BeginAtZero {
        minY = 0
    }
    
    // Set margins
    e.margin = Margin{
        Top:    40,
        Right:  40,
        Bottom: 60,
        Left:   60,
    }
    
    chartWidth := float64(e.width) - e.margin.Left - e.margin.Right
    chartHeight := float64(e.height) - e.margin.Top - e.margin.Bottom
    
    // Create scale functions
    e.xScale = func(x float64) float64 {
        return e.margin.Left + ((x - minX) / (maxX - minX)) * chartWidth
    }
    
    e.yScale = func(y float64) float64 {
        return e.margin.Top + chartHeight - ((y - minY) / (maxY - minY)) * chartHeight
    }
}

func (e *CanvasEngine) drawSmoothSegment(ctx app.Value, points []Point, i int, tension float64) {
    // Need to ensure we don't go out of bounds
    if i <= 0 || i >= len(points)-2 {
        return
    }
    
    p0 := points[i-1]
    p1 := points[i]
    p2 := points[i+1]
    p3 := points[i+2] // Need this for Catmull-Rom
    
    x0, y0 := e.xScale(p0.X), e.yScale(p0.Y)
    x1, y1 := e.xScale(p1.X), e.yScale(p1.Y)
    x2, y2 := e.xScale(p2.X), e.yScale(p2.Y)
    x3, y3 := e.xScale(p3.X), e.yScale(p3.Y)
    
    // Catmull-Rom spline to bezier control points
    cp1x := x1 + (x2-x0)*tension/6
    cp1y := y1 + (y2-y0)*tension/6
    cp2x := x2 - (x3-x1)*tension/6
    cp2y := y2 - (y3-y1)*tension/6
    
    ctx.Call("bezierCurveTo", cp1x, cp1y, cp2x, cp2y, x2, y2)
}

func (e *CanvasEngine) drawPoint(ctx app.Value, x, y float64, size int, color string, style PointStyle) {
    ctx.Set("fillStyle", color)
    ctx.Set("strokeStyle", "#ffffff")
    ctx.Set("lineWidth", 1)
    
    radius := float64(size) / 2
    
    switch style {
    case PointStyleCircle, "":
        ctx.Call("beginPath")
        ctx.Call("arc", x, y, radius, 0, 2*math.Pi)
        ctx.Call("fill")
        ctx.Call("stroke")
        
    case PointStyleSquare:
        ctx.Call("fillRect", x-radius, y-radius, radius*2, radius*2)
        ctx.Call("strokeRect", x-radius, y-radius, radius*2, radius*2)
        
    case PointStyleTriangle:
        ctx.Call("beginPath")
        ctx.Call("moveTo", x, y-radius)
        ctx.Call("lineTo", x+radius, y+radius)
        ctx.Call("lineTo", x-radius, y+radius)
        ctx.Call("closePath")
        ctx.Call("fill")
        ctx.Call("stroke")
        
    case PointStyleDiamond:
        ctx.Call("beginPath")
        ctx.Call("moveTo", x, y-radius)
        ctx.Call("lineTo", x+radius, y)
        ctx.Call("lineTo", x, y+radius)
        ctx.Call("lineTo", x-radius, y)
        ctx.Call("closePath")
        ctx.Call("fill")
        ctx.Call("stroke")
        
    case PointStyleCross:
        ctx.Set("strokeStyle", color)
        ctx.Set("lineWidth", 2)
        ctx.Call("beginPath")
        ctx.Call("moveTo", x-radius, y-radius)
        ctx.Call("lineTo", x+radius, y+radius)
        ctx.Call("moveTo", x+radius, y-radius)
        ctx.Call("lineTo", x-radius, y+radius)
        ctx.Call("stroke")
    }
}

func (e *CanvasEngine) downsampleLTTB(data []Point, threshold int) []Point {
    if threshold >= len(data) || threshold == 0 {
        return data
    }
    
    sampled := make([]Point, 0, threshold)
    sampled = append(sampled, data[0])
    
    bucketSize := float64(len(data)-2) / float64(threshold-2)
    a := 0
    
    for i := 0; i < threshold-2; i++ {
        avgRangeStart := int(math.Floor(float64(i+1)*bucketSize)) + 1
        avgRangeEnd := int(math.Floor(float64(i+2)*bucketSize)) + 1
        if avgRangeEnd > len(data) {
            avgRangeEnd = len(data)
        }
        
        avgRangeLength := avgRangeEnd - avgRangeStart
        
        var avgX, avgY float64
        for j := avgRangeStart; j < avgRangeEnd; j++ {
            avgX += data[j].X
            avgY += data[j].Y
        }
        avgX /= float64(avgRangeLength)
        avgY /= float64(avgRangeLength)
        
        rangeOffs := int(math.Floor(float64(i)*bucketSize)) + 1
        rangeTo := int(math.Floor(float64(i+1)*bucketSize)) + 1
        
        pointA := data[a]
        maxArea := -1.0
        maxAreaPoint := Point{}
        maxAreaIndex := a
        
        for j := rangeOffs; j < rangeTo; j++ {
            area := math.Abs(
                (pointA.X-avgX)*(data[j].Y-pointA.Y)-
                    (pointA.X-data[j].X)*(avgY-pointA.Y),
            ) * 0.5
            
            if area > maxArea {
                maxArea = area
                maxAreaPoint = data[j]
                maxAreaIndex = j
            }
        }
        
        sampled = append(sampled, maxAreaPoint)
        a = maxAreaIndex
    }
    
    sampled = append(sampled, data[len(data)-1])
    return sampled
}

func (e *CanvasEngine) drawLineChartGrid(ctx app.Value, chartWidth, chartHeight float64) {
    ctx.Set("strokeStyle", "#e5e7eb")
    ctx.Set("lineWidth", 0.5)
    
    // Vertical grid lines
    for i := 0; i <= 10; i++ {
        x := e.margin.Left + (float64(i)/10)*chartWidth
        ctx.Call("beginPath")
        ctx.Call("moveTo", x, e.margin.Top)
        ctx.Call("lineTo", x, e.margin.Top+chartHeight)
        ctx.Call("stroke")
    }
    
    // Horizontal grid lines
    for i := 0; i <= 5; i++ {
        y := e.margin.Top + (float64(i)/5)*chartHeight
        ctx.Call("beginPath")
        ctx.Call("moveTo", e.margin.Left, y)
        ctx.Call("lineTo", e.margin.Left+chartWidth, y)
        ctx.Call("stroke")
    }
}

func (e *CanvasEngine) drawLineChartAxes(ctx app.Value, chartWidth, chartHeight float64, spec *Spec) {
    ctx.Set("strokeStyle", "#9ca3af")
    ctx.Set("lineWidth", 1)
    
    // X-axis
    ctx.Call("beginPath")
    ctx.Call("moveTo", e.margin.Left, e.margin.Top+chartHeight)
    ctx.Call("lineTo", e.margin.Left+chartWidth, e.margin.Top+chartHeight)
    ctx.Call("stroke")
    
    // Y-axis
    ctx.Call("beginPath")
    ctx.Call("moveTo", e.margin.Left, e.margin.Top)
    ctx.Call("lineTo", e.margin.Left, e.margin.Top+chartHeight)
    ctx.Call("stroke")
}

func (e *CanvasEngine) drawLegend(ctx app.Value, spec *Spec) {
    data := spec.Data
    legendX := e.margin.Left
    legendY := e.margin.Top - 25
    
    ctx.Set("font", fmt.Sprintf("12px %s", spec.Theme.GetFontFamily()))
    ctx.Set("textBaseline", "middle")
    
    for i, series := range data.Series {
        color := series.Color
        if color == "" {
            colors := spec.Theme.GetColors()
            color = colors[i%len(colors)]
        }
        
        // Draw color box
        ctx.Set("fillStyle", color)
        ctx.Call("fillRect", legendX, legendY-6, 12, 12)
        
        // Draw label
        ctx.Set("fillStyle", spec.Theme.GetTextColor())
        ctx.Set("textAlign", "left")
        ctx.Call("fillText", series.Label, legendX+18, legendY)
        
        // Move to next legend item
        bounds := ctx.Call("measureText", series.Label)
        legendX += 18 + bounds.Get("width").Float() + 20
    }
}

