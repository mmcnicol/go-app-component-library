
github.com/mmcnicol/go-app-component-library/pkg/components/vizzy/

a Chart/Visualization library

goals:
* Separation of Concerns: Chart logic vs rendering logic
* Testability: Pure data models are easy to test
* Extensibility: Easy to add new renderers (SVG, WebGL, etc.)
* Performance: Can optimize renderers independently
* Maintainability: Smaller, focused files
* Reusability: Same chart model can be rendered by different engines

the primary renderer is for HTML5 canvas as a go-app v10 component.

the package has a flat file structure.

file:
core_chart.go             # Pure data models
core_registry.go          # Engine registry
core_types.go             # Shared types
core_interfaces.go        # Core interfaces
render_canvas_engine.go   # Canvas engine implementation
render_canvas_bar.go      # Bar chart renderer
render_canvas_line.go     # Line chart renderer
render_canvas_...
render_svg_...
component_chart.go        # app.UI wrapper


Storybook stories for Bar charts
Storybook stories for Grouped Bar charts
Storybook stories for Horizontal Bar charts
Storybook stories for Line charts (Basic, Multi Series, Area Chart)
Storybook stories for Pie charts
Storybook stories for Donut charts
Storybook stories for Scatter plots
Storybook stories for Box plots
Storybook stories for Heatmaps
Storybook stories for Regression analysis
Storybook stories for Radar Charts
Storybook stories for Candlestick Charts
Storybook stories for Tree maps
Storybook stories for Gantt charts
Storybook stories for Sunburst diagrams
Storybook stories for Network graphs
Storybook stories for 3D Surface plots



this chart/visualization library is inspired by this book chapter:

# **25.2 Data Visualization Chart Components**

## **Introduction**

Data visualization transforms raw data into meaningful insights through graphical representations. In modern web applications, interactive charts are essential for dashboards, analytics platforms, and data-driven decision making. While numerous JavaScript charting libraries exist, integrating them with go-app's WebAssembly environment presents unique challenges and opportunities.

In this chapter, we'll build a comprehensive charting component system for go-app that bridges the gap between Go's data processing capabilities and web-based visualization. We'll create a flexible, high-performance charting library that leverages Go's strengths while providing rich, interactive visualizations.

## **Understanding Charting Requirements in go-app**

### **Why Build Custom Chart Components?**

While JavaScript chart libraries like Chart.js, D3.js, and Plotly are popular, they have limitations in a go-app context:
- **WASM-JS Bridge Overhead**: Frequent communication between WebAssembly and JavaScript can be expensive
- **Type Safety**: JavaScript libraries lack Go's compile-time type checking
- **Data Processing**: JavaScript isn't optimized for large-scale data manipulation
- **Bundle Size**: Including multiple JavaScript libraries increases bundle size significantly
- **Go Integration**: Deep integration with Go data structures and concurrency patterns

### **Core Features of Our Charting System**

Our implementation will support:
1. **Multiple Chart Types**: Line, bar, scatter, pie, area, heatmap
2. **Interactive Features**: Tooltips, zoom, pan, selection
3. **Performance**: Efficient rendering for large datasets (10k+ points)
4. **Accessibility**: Screen reader support, keyboard navigation
5. **Theming**: Consistent with design system, dark/light mode
6. **Responsive Design**: Adapts to container size
7. **Real-time Updates**: Streaming data support
8. **Export Capabilities**: PNG, SVG, CSV export

## **Architecture Overview**

### **Dual-Engine Approach**

We'll implement a hybrid architecture that uses both pure-Go rendering (via Canvas API) and optimized JavaScript libraries for complex visualizations:

```go
// Chart engine interface
type ChartEngine interface {
    Render(chart ChartSpec) error
    Update(data ChartData) error
    Destroy() error
    GetCanvas() app.HTMLCanvas
}

// Chart specification
type ChartSpec struct {
    Type        ChartType
    Data        ChartData
    Options     ChartOptions
    ContainerID string
    Engine      EngineType // "canvas", "svg", "webgl", "hybrid"
}

// Chart data structure
type ChartData struct {
    Labels   []string
    Datasets []Dataset
    Metadata map[string]interface{}
}

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
```

### **Component Hierarchy**

```go
// Base chart component
type BaseChart struct {
    app.Compo
    
    // State
    spec        ChartSpec
    containerID string
    engine      ChartEngine
    isRendered  bool
    
    // Refs
    canvasRef   app.Ref
    tooltipRef  app.Ref
    
    // Event handlers
    onPointClick    func(point DataPoint, datasetIndex int)
    onZoom          func(domain AxisRange)
    onHover         func(point DataPoint, datasetIndex int)
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
```

## **Pure-Go Canvas Rendering Engine**

### **Canvas-Based Chart Renderer**

```go
type CanvasRenderer struct {
    canvas    app.HTMLCanvas
    ctx       app.CanvasRenderingContext2D
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
    
    // Get rendering context
    ctx, err := cr.canvas.GetContext2D()
    if err != nil {
        return nil, err
    }
    cr.ctx = ctx
    
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

func (cr *CanvasRenderer) drawLineDataset(dataset Dataset, xScale, yScale Scale, datasetIndex int) {
    cr.ctx.BeginPath()
    cr.ctx.LineWidth = float64(dataset.BorderWidth)
    cr.ctx.StrokeStyle = dataset.BorderColor
    
    // Draw line
    for i, point := range dataset.Data {
        x := xScale.Convert(point.X)
        y := yScale.Convert(point.Y)
        
        if i == 0 {
            cr.ctx.MoveTo(x, y)
        } else {
            // Apply tension for smooth curves
            if dataset.Tension > 0 {
                cr.drawBezierCurve(dataset.Data, i, xScale, yScale, dataset.Tension)
            } else {
                cr.ctx.LineTo(x, y)
            }
        }
        
        // Draw points
        if dataset.PointRadius > 0 {
            cr.drawPoint(x, y, dataset.PointRadius, dataset.BackgroundColor[i])
        }
    }
    
    cr.ctx.Stroke()
    
    // Fill area under line if needed
    if dataset.Fill {
        cr.fillAreaUnderLine(dataset, xScale, yScale)
    }
}
```

### **Optimized Path Rendering for Large Datasets**

```go
type OptimizedPathRenderer struct {
    ctx          app.CanvasRenderingContext2D
    pathCache    map[string]*PathCacheEntry
    minSegment   float64 // minimum pixel distance to draw segment
    simplificationTolerance float64
}

func (opr *OptimizedPathRenderer) drawOptimizedLine(points []DataPoint, xScale, yScale Scale) {
    if len(points) < 2 {
        return
    }
    
    // Apply Ramer-Douglas-Peucker simplification
    simplified := opr.simplifyLine(points, opr.simplificationTolerance)
    
    opr.ctx.BeginPath()
    
    // Use line dash for very dense data to improve performance
    if len(simplified) > 1000 {
        opr.ctx.SetLineDash([]float64{1, 0}) // Solid line but optimized
    }
    
    // Draw optimized path
    for i, point := range simplified {
        x := xScale.Convert(point.X)
        y := yScale.Convert(point.Y)
        
        if i == 0 {
            opr.ctx.MoveTo(x, y)
        } else {
            // Skip segments that are too small
            prevX := xScale.Convert(simplified[i-1].X)
            prevY := yScale.Convert(simplified[i-1].Y)
            
            if math.Abs(x-prevX) > opr.minSegment || math.Abs(y-prevY) > opr.minSegment {
                opr.ctx.LineTo(x, y)
            }
        }
    }
    
    opr.ctx.Stroke()
}

// Ramer-Douglas-Peucker line simplification
func (opr *OptimizedPathRenderer) simplifyLine(points []DataPoint, epsilon float64) []DataPoint {
    if len(points) < 3 {
        return points
    }
    
    // Find the point with the maximum distance
    maxDist := 0.0
    maxIndex := 0
    
    for i := 1; i < len(points)-1; i++ {
        dist := opr.perpendicularDistance(points[i], points[0], points[len(points)-1])
        if dist > maxDist {
            maxDist = dist
            maxIndex = i
        }
    }
    
    // If max distance is greater than epsilon, recursively simplify
    if maxDist > epsilon {
        left := opr.simplifyLine(points[:maxIndex+1], epsilon)
        right := opr.simplifyLine(points[maxIndex:], epsilon)
        
        // Concatenate results, removing duplicate point
        return append(left[:len(left)-1], right...)
    }
    
    // Base case: return endpoints
    return []DataPoint{points[0], points[len(points)-1]}
}
```

## **Interactive Features**

### **Tooltip System**

```go
type TooltipManager struct {
    chart         *BaseChart
    tooltipElem   app.HTMLDiv
    currentPoint  *DataPoint
    currentDataset int
    position      TooltipPosition
    formatter     TooltipFormatter
}

func (tm *TooltipManager) setupEventHandlers() {
    // Mouse move handler
    tm.chart.canvasRef.AddEventListener("mousemove", tm.handleMouseMove)
    
    // Mouse leave handler
    tm.chart.canvasRef.AddEventListener("mouseleave", tm.handleMouseLeave)
    
    // Touch handlers for mobile
    tm.chart.canvasRef.AddEventListener("touchstart", tm.handleTouchStart)
    tm.chart.canvasRef.AddEventListener("touchmove", tm.handleTouchMove)
}

func (tm *TooltipManager) handleMouseMove(ctx app.Context, e app.Event) {
    rect := tm.chart.canvasRef.GetBoundingClientRect()
    mouseX := e.Get("clientX").Float() - rect.Get("left").Float()
    mouseY := e.Get("clientY").Float() - rect.Get("top").Float()
    
    // Find nearest data point
    point, datasetIndex, distance := tm.findNearestPoint(mouseX, mouseY)
    
    // Show tooltip if close enough
    if point != nil && distance < tm.chart.spec.Options.Tooltips.IntersectDistance {
        tm.showTooltip(point, datasetIndex, mouseX, mouseY)
    } else {
        tm.hideTooltip()
    }
}

func (tm *TooltipManager) findNearestPoint(canvasX, canvasY float64) (*DataPoint, int, float64) {
    var nearestPoint *DataPoint
    var nearestDataset int
    minDistance := math.MaxFloat64
    
    for i, dataset := range tm.chart.spec.Data.Datasets {
        for _, point := range dataset.Data {
            // Convert data point to canvas coordinates
            pointX := tm.chart.xScale.Convert(point.X)
            pointY := tm.chart.yScale.Convert(point.Y)
            
            // Calculate distance
            distance := math.Sqrt(
                math.Pow(pointX-canvasX, 2) + math.Pow(pointY-canvasY, 2),
            )
            
            // Check if this point is closer
            if distance < minDistance {
                minDistance = distance
                nearestPoint = &point
                nearestDataset = i
            }
        }
    }
    
    return nearestPoint, nearestDataset, minDistance
}

func (tm *TooltipManager) showTooltip(point *DataPoint, datasetIndex int, x, y float64) {
    // Format tooltip content
    content := tm.formatter.Format(point, datasetIndex, tm.chart.spec.Data)
    
    // Position tooltip
    tm.positionTooltip(x, y)
    
    // Update content
    tm.tooltipElem.SetInnerHTML(content)
    tm.tooltipElem.SetStyle("display", "block")
    
    // Trigger custom event
    tm.chart.onHover(*point, datasetIndex)
}
```

### **Zoom and Pan System**

```go
type ZoomPanManager struct {
    chart       *BaseChart
    isDragging  bool
    dragStartX  float64
    dragStartY  float64
    currentZoom ZoomState
    zoomHistory []ZoomState
}

type ZoomState struct {
    xMin, xMax float64
    yMin, yMax float64
    transform   TransformMatrix
}

func (zpm *ZoomPanManager) setupZoomPan() {
    // Mouse wheel for zoom
    zpm.chart.canvasRef.AddEventListener("wheel", zpm.handleWheel)
    
    // Drag for pan
    zpm.chart.canvasRef.AddEventListener("mousedown", zpm.handleMouseDown)
    zpm.chart.canvasRef.AddEventListener("mousemove", zpm.handleMouseMove)
    zpm.chart.canvasRef.AddEventListener("mouseup", zpm.handleMouseUp)
    
    // Touch events for mobile
    zpm.chart.canvasRef.AddEventListener("touchstart", zpm.handleTouchStart)
    zpm.chart.canvasRef.AddEventListener("touchmove", zpm.handleTouchMove)
    zpm.chart.canvasRef.AddEventListener("touchend", zpm.handleTouchEnd)
    
    // Double click to reset zoom
    zpm.chart.canvasRef.AddEventListener("dblclick", zpm.handleDoubleClick)
}

func (zpm *ZoomPanManager) handleWheel(ctx app.Context, e app.Event) {
    e.PreventDefault()
    
    delta := e.Get("deltaY").Float()
    zoomFactor := 1.1
    
    // Get mouse position relative to canvas
    rect := zpm.chart.canvasRef.GetBoundingClientRect()
    mouseX := e.Get("clientX").Float() - rect.Get("left").Float()
    mouseY := e.Get("clientY").Float() - rect.Get("top").Float()
    
    // Convert to data coordinates
    dataX := zpm.chart.xScale.Invert(mouseX)
    dataY := zpm.chart.yScale.Invert(mouseY)
    
    // Apply zoom
    if delta < 0 {
        // Zoom in
        zpm.zoomAt(dataX, dataY, zoomFactor)
    } else {
        // Zoom out
        zpm.zoomAt(dataX, dataY, 1/zoomFactor)
    }
}

func (zpm *ZoomPanManager) zoomAt(centerX, centerY, factor float64) {
    // Calculate new domain
    xRange := zpm.currentZoom.xMax - zpm.currentZoom.xMin
    yRange := zpm.currentZoom.yMax - zpm.currentZoom.yMin
    
    newXRange := xRange / factor
    newYRange := yRange / factor
    
    // Center zoom on specified point
    newXMin := centerX - (centerX - zpm.currentZoom.xMin) / factor
    newXMax := newXMin + newXRange
    
    newYMin := centerY - (centerY - zpm.currentZoom.yMin) / factor
    newYMax := newYMin + newYRange
    
    // Update zoom state
    zpm.pushZoomState()
    zpm.currentZoom = ZoomState{
        xMin: newXMin,
        xMax: newXMax,
        yMin: newYMin,
        yMax: newYMax,
    }
    
    // Update chart scales and redraw
    zpm.chart.updateScales(newXMin, newXMax, newYMin, newYMax)
    zpm.chart.Render()
    
    // Trigger zoom event
    zpm.chart.onZoom(AxisRange{X: [2]float64{newXMin, newXMax}, Y: [2]float64{newYMin, newYMax}})
}
```

## **High-Performance WebGL Rendering Engine**

### **WebGL-Based Renderer for Large Datasets**

```go
type WebGLRenderer struct {
    canvas      app.HTMLCanvas
    gl          app.WebGLRenderingContext
    programs    map[string]*WebGLProgram
    buffers     map[string]*WebGLBuffer
    textures    map[string]*WebGLTexture
    maxPoints   int
    instanceCount int
}

func NewWebGLRenderer(containerID string) (*WebGLRenderer, error) {
    wr := &WebGLRenderer{
        programs: make(map[string]*WebGLProgram),
        buffers: make(map[string]*WebGLBuffer),
        textures: make(map[string]*WebGLTexture),
        maxPoints: 100000, // Support up to 100k points
    }
    
    // Create canvas with WebGL context
    wr.canvas = app.Canvas().
        ID(containerID + "-webgl")
    
    // Get WebGL context
    gl, err := wr.canvas.GetWebGLContext()
    if err != nil {
        return nil, err
    }
    wr.gl = gl
    
    // Initialize WebGL
    wr.initWebGL()
    
    return wr, nil
}

func (wr *WebGLRenderer) initWebGL() {
    // Set clear color
    wr.gl.ClearColor(0, 0, 0, 0)
    
    // Enable blending for transparency
    wr.gl.Enable(wr.gl.BLEND)
    wr.gl.BlendFunc(wr.gl.SRC_ALPHA, wr.gl.ONE_MINUS_SRC_ALPHA)
    
    // Create shader programs
    wr.createLineShaderProgram()
    wr.createPointShaderProgram()
    wr.createBarShaderProgram()
    
    // Create vertex buffers
    wr.createVertexBuffers()
}

func (wr *WebGLRenderer) renderLineChart(chart ChartSpec) error {
    // Use instanced rendering for multiple datasets
    for i, dataset := range chart.Data.Datasets {
        wr.renderLineDataset(dataset, i, len(chart.Data.Datasets))
    }
    
    return nil
}

func (wr *WebGLRenderer) renderLineDataset(dataset Dataset, datasetIndex, totalDatasets int) {
    program := wr.programs["line"]
    wr.gl.UseProgram(program)
    
    // Set uniform values
    wr.setUniforms(program, dataset, datasetIndex, totalDatasets)
    
    // Bind vertex data
    wr.bindLineVertices(dataset.Data)
    
    // Draw line using line strip
    wr.gl.DrawArrays(wr.gl.LINE_STRIP, 0, len(dataset.Data))
    
    // Draw points if needed
    if dataset.PointRadius > 0 {
        wr.renderPoints(dataset, datasetIndex)
    }
}

func (wr *WebGLRenderer) bindLineVertices(points []DataPoint) {
    // Convert points to Float32Array for WebGL
    vertices := make([]float32, len(points)*2)
    for i, point := range points {
        vertices[i*2] = float32(point.X)
        vertices[i*2+1] = float32(point.Y)
    }
    
    // Bind to buffer
    buffer := wr.buffers["lineVertices"]
    wr.gl.BindBuffer(wr.gl.ARRAY_BUFFER, buffer)
    wr.gl.BufferData(wr.gl.ARRAY_BUFFER, vertices, wr.gl.STATIC_DRAW)
    
    // Set vertex attribute pointer
    wr.gl.VertexAttribPointer(0, 2, wr.gl.FLOAT, false, 0, 0)
    wr.gl.EnableVertexAttribArray(0)
}
```

### **GPU-Accelerated Heatmap Rendering**

```go
type HeatmapRenderer struct {
    WebGLRenderer
    colorTexture *WebGLTexture
    densityTexture *WebGLTexture
    colorScale   ColorScale
}

func (hr *HeatmapRenderer) renderHeatmap(data [][]float64, xLabels, yLabels []string) error {
    // Create density matrix on GPU
    hr.computeDensityTexture(data)
    
    // Apply color scale
    hr.applyColorScale()
    
    // Render heatmap quad
    hr.renderHeatmapQuad()
    
    // Draw labels if needed
    if hr.showValues {
        hr.drawValueLabels(data, xLabels, yLabels)
    }
    
    return nil
}

func (hr *HeatmapRenderer) computeDensityTexture(data [][]float64) {
    // Use compute shader or fragment shader to compute densities
    program := hr.programs["heatmapDensity"]
    hr.gl.UseProgram(program)
    
    // Bind data as texture
    hr.bindDataAsTexture(data)
    
    // Dispatch compute shader
    hr.gl.DispatchCompute(
        uint32(len(data[0])/16+1),
        uint32(len(data)/16+1),
        1,
    )
    
    hr.gl.MemoryBarrier(hr.gl.SHADER_IMAGE_ACCESS_BARRIER_BIT)
}
```

## **Statistical Chart Components**

### **Statistical Visualization Components**

```go
type BoxPlotChart struct {
    BaseChart
    showOutliers bool
    showMean     bool
    whiskerType  WhiskerType // "tukey", "minmax", "percentile"
}

func (bpc *BoxPlotChart) calculateStatistics(data [][]float64) []BoxPlotStats {
    stats := make([]BoxPlotStats, len(data))
    
    for i, dataset := range data {
        // Sort data
        sorted := make([]float64, len(dataset))
        copy(sorted, dataset)
        sort.Float64s(sorted)
        
        // Calculate quartiles
        q1 := bpc.calculatePercentile(sorted, 25)
        median := bpc.calculatePercentile(sorted, 50)
        q3 := bpc.calculatePercentile(sorted, 75)
        
        // Calculate IQR
        iqr := q3 - q1
        
        // Determine whiskers
        var lowerWhisker, upperWhisker float64
        switch bpc.whiskerType {
        case WhiskerTypeTukey:
            lowerWhisker = q1 - 1.5*iqr
            upperWhisker = q3 + 1.5*iqr
        case WhiskerTypeMinMax:
            lowerWhisker = sorted[0]
            upperWhisker = sorted[len(sorted)-1]
        case WhiskerTypePercentile:
            lowerWhisker = bpc.calculatePercentile(sorted, 5)
            upperWhisker = bpc.calculatePercentile(sorted, 95)
        }
        
        // Identify outliers
        var outliers []float64
        if bpc.showOutliers {
            for _, value := range sorted {
                if value < lowerWhisker || value > upperWhisker {
                    outliers = append(outliers, value)
                }
            }
        }
        
        stats[i] = BoxPlotStats{
            Min:          sorted[0],
            Q1:           q1,
            Median:       median,
            Q3:           q3,
            Max:          sorted[len(sorted)-1],
            LowerWhisker: lowerWhisker,
            UpperWhisker: upperWhisker,
            Outliers:     outliers,
            Mean:         bpc.calculateMean(dataset),
        }
    }
    
    return stats
}
```

### **Regression Analysis Components**

```go
type RegressionChart struct {
    ScatterChart
    regressionType RegressionType // "linear", "polynomial", "exponential", "logarithmic"
    degree         int            // for polynomial regression
    showEquation   bool
    showRSquared   bool
}

func (rc *RegressionChart) calculateRegression(points []DataPoint) RegressionResult {
    switch rc.regressionType {
    case RegressionTypeLinear:
        return rc.calculateLinearRegression(points)
    case RegressionTypePolynomial:
        return rc.calculatePolynomialRegression(points, rc.degree)
    case RegressionTypeExponential:
        return rc.calculateExponentialRegression(points)
    case RegressionTypeLogarithmic:
        return rc.calculateLogarithmicRegression(points)
    default:
        return rc.calculateLinearRegression(points)
    }
}

func (rc *RegressionChart) calculateLinearRegression(points []DataPoint) RegressionResult {
    n := float64(len(points))
    
    // Calculate sums
    var sumX, sumY, sumXY, sumX2 float64
    for _, p := range points {
        sumX += p.X
        sumY += p.Y
        sumXY += p.X * p.Y
        sumX2 += p.X * p.X
    }
    
    // Calculate slope (m) and intercept (b)
    m := (n*sumXY - sumX*sumY) / (n*sumX2 - sumX*sumX)
    b := (sumY - m*sumX) / n
    
    // Calculate R-squared
    var ssTotal, ssResidual float64
    meanY := sumY / n
    
    for _, p := range points {
        yPred := m*p.X + b
        ssTotal += math.Pow(p.Y-meanY, 2)
        ssResidual += math.Pow(p.Y-yPred, 2)
    }
    
    rSquared := 1 - (ssResidual / ssTotal)
    
    return RegressionResult{
        Coefficients: []float64{b, m},
        Equation:     fmt.Sprintf("y = %.4fx + %.4f", m, b),
        RSquared:     rSquared,
        Predict: func(x float64) float64 {
            return m*x + b
        },
    }
}
```

## **Real-Time Streaming Charts**

### **Streaming Data Support**

```go
type StreamingChart struct {
    BaseChart
    dataBuffer   *DataBuffer
    maxPoints    int
    streaming    bool
    updateRate   time.Duration
    lastUpdate   time.Time
}

type DataBuffer struct {
    data      []DataPoint
    maxSize   int
    head      int
    tail      int
    count     int
    mu        sync.RWMutex
}

func NewDataBuffer(maxSize int) *DataBuffer {
    return &DataBuffer{
        data:    make([]DataPoint, maxSize),
        maxSize: maxSize,
    }
}

func (db *DataBuffer) Add(point DataPoint) {
    db.mu.Lock()
    defer db.mu.Unlock()
    
    db.data[db.head] = point
    db.head = (db.head + 1) % db.maxSize
    
    if db.count < db.maxSize {
        db.count++
    } else {
        db.tail = (db.tail + 1) % db.maxSize
    }
}

func (db *DataBuffer) GetSlice() []DataPoint {
    db.mu.RLock()
    defer db.mu.RUnlock()
    
    if db.count == 0 {
        return []DataPoint{}
    }
    
    result := make([]DataPoint, db.count)
    
    if db.head > db.tail || db.count < db.maxSize {
        // Data is contiguous
        copy(result, db.data[db.tail:db.head])
    } else {
        // Data wraps around
        firstPart := db.data[db.tail:db.maxSize]
        secondPart := db.data[0:db.head]
        
        copy(result, firstPart)
        copy(result[len(firstPart):], secondPart)
    }
    
    return result
}

func (sc *StreamingChart) StreamData(dataSource <-chan DataPoint) {
    sc.streaming = true
    
    go func() {
        for point := range dataSource {
            sc.dataBuffer.Add(point)
            
            // Throttle updates
            if time.Since(sc.lastUpdate) >= sc.updateRate {
                sc.updateChart()
                sc.lastUpdate = time.Now()
            }
        }
    }()
}

func (sc *StreamingChart) updateChart() {
    // Get latest data
    data := sc.dataBuffer.GetSlice()
    
    // Update dataset
    sc.spec.Data.Datasets[0].Data = data
    
    // Trigger incremental update
    sc.engine.Update(sc.spec.Data)
}
```

## **Accessibility and Internationalization**

### **Accessible Chart Components**

```go
type AccessibleChart struct {
    BaseChart
    ariaLabel     string
    ariaDescribedBy string
    longDescURL   string
    dataTableID   string
}

func (ac *AccessibleChart) Render() app.UI {
    return app.Div().
        Role("img").
        Aria("label", ac.ariaLabel).
        Aria("describedby", ac.ariaDescribedBy).
        Body(
            // Visual chart
            ac.BaseChart.Render(),
            
            // Hidden data table for screen readers
            ac.renderDataTable(),
            
            // Hidden descriptive text
            app.Div().
                ID(ac.ariaDescribedBy).
                Class("sr-only").
                Body(
                    ac.generateChartDescription(),
                ),
        )
}

func (ac *AccessibleChart) generateChartDescription() string {
    var desc strings.Builder
    
    desc.WriteString(fmt.Sprintf("%s chart showing ", ac.spec.Type))
    desc.WriteString(fmt.Sprintf("%d datasets with %d data points each. ",
        len(ac.spec.Data.Datasets), len(ac.spec.Data.Datasets[0].Data)))
    
    // Describe trends
    if ac.spec.Type == ChartTypeLine {
        trend := ac.calculateTrend()
        desc.WriteString(fmt.Sprintf("The data shows a %s trend. ", trend))
    }
    
    // Describe key statistics
    stats := ac.calculateStatistics()
    desc.WriteString(fmt.Sprintf("Minimum value: %.2f. Maximum value: %.2f. ",
        stats.Min, stats.Max))
    desc.WriteString(fmt.Sprintf("Average value: %.2f.", stats.Mean))
    
    return desc.String()
}

func (ac *AccessibleChart) renderDataTable() app.UI {
    return app.Table().
        ID(ac.dataTableID).
        Class("sr-only").
        Aria("hidden", "true").
        Body(
            app.Caption().Text(ac.ariaLabel),
            app.THead().Body(
                app.Tr().Body(
                    app.Th().Text("Dataset"),
                    app.Range(ac.spec.Data.Labels).Slice(func(i int) app.UI {
                        return app.Th().Text(ac.spec.Data.Labels[i])
                    }),
                ),
            ),
            app.TBody().Body(
                app.Range(ac.spec.Data.Datasets).Slice(func(i int) app.UI {
                    dataset := ac.spec.Data.Datasets[i]
                    return app.Tr().Body(
                        app.Th().Text(dataset.Label),
                        app.Range(dataset.Data).Slice(func(j int) app.UI {
                            return app.Td().Text(fmt.Sprintf("%.2f", dataset.Data[j].Y))
                        }),
                    )
                }),
            ),
        )
}
```

## **Complete Usage Examples**

### **Dashboard with Multiple Chart Types**

```go
func ExampleDataDashboard() app.UI {
    // Sample data
    monthlySales := []DataPoint{
        {X: 1, Y: 12000}, {X: 2, Y: 19000}, {X: 3, Y: 15000},
        {X: 4, Y: 21000}, {X: 5, Y: 18000}, {X: 6, Y: 24000},
    }
    
    categoryRevenue := []DataPoint{
        {Label: "Electronics", Value: 45000},
        {Label: "Clothing", Value: 32000},
        {Label: "Home Goods", Value: 28000},
        {Label: "Books", Value: 15000},
    }
    
    userActivity := [][]float64{
        {8, 12, 15, 9, 14, 11, 13, 10},
        {5, 8, 10, 7, 9, 6, 11, 8},
        {12, 15, 18, 14, 16, 13, 17, 15},
    }
    
    return app.Div().
        Class("dashboard").
        Body(
            app.H1().Text("Sales Dashboard"),
            
            // First row: Line chart and Pie chart
            app.Div().
                Class("dashboard-row").
                Body(
                    NewChart(ChartTypeLine).
                        Title("Monthly Sales Trend").
                        Data(ChartData{
                            Labels: []string{"Jan", "Feb", "Mar", "Apr", "May", "Jun"},
                            Datasets: []Dataset{{
                                Label: "Sales",
                                Data:  monthlySales,
                                BorderColor: "#4A90E2",
                                Fill: true,
                            }},
                        }).
                        Options(ChartOptions{
                            Responsive: true,
                            MaintainAspectRatio: false,
                            Interaction: ChartInteraction{
                                Mode: "nearest",
                                Intersect: false,
                            },
                            Scales: ChartScales{
                                Y: Axis{
                                    BeginAtZero: true,
                                    Title: AxisTitle{
                                        Display: true,
                                        Text: "Sales ($)",
                                    },
                                },
                            },
                        }).
                        Class("dashboard-card", "chart-line"),
                    
                    NewChart(ChartTypePie).
                        Title("Revenue by Category").
                        Data(ChartData{
                            Datasets: []Dataset{{
                                Data: categoryRevenue,
                                BackgroundColor: []string{
                                    "#FF6384", "#36A2EB", "#FFCE56", "#4BC0C0",
                                },
                            }},
                        }).
                        Options(ChartOptions{
                            Plugins: ChartPlugins{
                                Legend: LegendOptions{
                                    Position: "right",
                                },
                                Tooltip: TooltipOptions{
                                    Callbacks: TooltipCallbacks{
                                        Label: func(context TooltipContext) string {
                                            value := context.Value()
                                            total := context.Total()
                                            percentage := (value / total) * 100
                                            return fmt.Sprintf("%s: $%.0f (%.1f%%)",
                                                context.Label(), value, percentage)
                                        },
                                    },
                                },
                            },
                        }).
                        Class("dashboard-card", "chart-pie"),
                ),
            
            // Second row: Bar chart and Heatmap
            app.Div().
                Class("dashboard-row").
                Body(
                    NewChart(ChartTypeBar).
                        Title("User Activity by Hour").
                        Data(ChartData{
                            Labels: []string{"9AM", "10AM", "11AM", "12PM", "1PM", "2PM", "3PM", "4PM"},
                            Datasets: []Dataset{
                                {
                                    Label: "Week 1",
                                    Data:  userActivity[0],
                                    BackgroundColor: "rgba(74, 144, 226, 0.7)",
                                },
                                {
                                    Label: "Week 2",
                                    Data:  userActivity[1],
                                    BackgroundColor: "rgba(255, 99, 132, 0.7)",
                                },
                                {
                                    Label: "Week 3",
                                    Data:  userActivity[2],
                                    BackgroundColor: "rgba(75, 192, 192, 0.7)",
                                },
                            },
                        }).
                        Options(ChartOptions{
                            Scales: ChartScales{
                                Y: Axis{
                                    Stacked: true,
                                    Title: AxisTitle{
                                        Display: true,
                                        Text: "Active Users",
                                    },
                                },
                                X: Axis{
                                    Stacked: true,
                                },
                            },
                        }).
                        Class("dashboard-card", "chart-bar"),
                    
                    NewChart(ChartTypeHeatmap).
                        Title("Correlation Matrix").
                        Data(ChartData{
                            Data: userActivity,
                            XLabels: []string{"Feature A", "Feature B", "Feature C", "Feature D"},
                            YLabels: []string{"Metric 1", "Metric 2", "Metric 3"},
                        }).
                        Options(ChartOptions{
                            Plugins: ChartPlugins{
                                Legend: LegendOptions{
                                    Display: false,
                                },
                            },
                            ColorScale: ColorScale{
                                Min: 0,
                                Max: 1,
                                Colors: []string{"#FFFFFF", "#4A90E2"},
                            },
                        }).
                        Class("dashboard-card", "chart-heatmap"),
                ),
            
            // Third row: Statistical charts
            app.Div().
                Class("dashboard-row").
                Body(
                    NewChart(ChartTypeBoxPlot).
                        Title("Performance Distribution").
                        Data(ChartData{
                            Data: [][]float64{
                                {12, 15, 18, 22, 25, 28, 30, 32, 35},
                                {8, 12, 16, 20, 24, 28, 32, 36, 40},
                                {20, 22, 24, 26, 28, 30, 32, 34, 36},
                            },
                            Labels: []string{"Team A", "Team B", "Team C"},
                        }).
                        Class("dashboard-card", "chart-boxplot"),
                    
                    NewChart(ChartTypeScatter).
                        Title("Correlation Analysis").
                        WithRegression(RegressionTypeLinear, 2).
                        Data(ChartData{
                            Datasets: []Dataset{{
                                Label: "Data Points",
                                Data:  generateRandomData(100),
                                PointRadius: 3,
                            }},
                        }).
                        Class("dashboard-card", "chart-scatter"),
                ),
        )
}
```

### **Real-Time Monitoring Dashboard**

```go
func ExampleRealTimeMonitoring() app.UI {
    var (
        cpuChart    *StreamingChart
        memoryChart *StreamingChart
        networkChart *StreamingChart
    )
    
    // Initialize charts
    cpuChart = NewStreamingChart(ChartTypeLine).
        Title("CPU Usage (%)").
        MaxPoints(100).
        UpdateRate(100 * time.Millisecond)
    
    memoryChart = NewStreamingChart(ChartTypeArea).
        Title("Memory Usage (MB)").
        MaxPoints(100).
        UpdateRate(100 * time.Millisecond)
    
    networkChart = NewStreamingChart(ChartTypeLine).
        Title("Network Throughput (MB/s)").
        MaxPoints(100).
        UpdateRate(100 * time.Millisecond)
    
    // Start streaming data
    OnMount(func() {
        // Simulate data streams
        go simulateCPUData(cpuChart)
        go simulateMemoryData(memoryChart)
        go simulateNetworkData(networkChart)
    })
    
    return app.Div().
        Class("monitoring-dashboard").
        Body(
            app.H1().Text("System Monitoring"),
            
            app.Div().
                Class("metrics-grid").
                Body(
                    cpuChart.
                        Class("metric-card", "metric-critical").
                        WarningThreshold(80).
                        DangerThreshold(90),
                    
                    memoryChart.
                        Class("metric-card", "metric-warning"),
                    
                    networkChart.
                        Class("metric-card", "metric-normal"),
                    
                    // Alert panel
                    NewAlertPanel().
                        Title("System Alerts").
                        Class("metric-card", "metric-alerts").
                        Body(
                            // Alert components would go here
                        ),
                ),
            
            // Controls
            app.Div().
                Class("dashboard-controls").
                Body(
                    app.Button().
                        Class("btn", "btn-primary").
                        Text("Pause All").
                        OnClick(func(ctx app.Context, e app.Event) {
                            cpuChart.Pause()
                            memoryChart.Pause()
                            networkChart.Pause()
                        }),
                    
                    app.Button().
                        Class("btn", "btn-secondary").
                        Text("Export Data").
                        OnClick(func(ctx app.Context, e app.Event) {
                            exportChartData(cpuChart, "cpu_usage.csv")
                        }),
                    
                    app.Select().
                        Class("time-range-selector").
                        Body(
                            app.Option().Value("1h").Text("Last Hour"),
                            app.Option().Value("6h").Text("Last 6 Hours"),
                            app.Option().Value("24h").Text("Last 24 Hours"),
                            app.Option().Value("7d").Text("Last 7 Days"),
                        ).
                        OnChange(func(ctx app.Context, e app.Event) {
                            range := ctx.JSSrc().Get("value").String()
                            updateTimeRange(range)
                        }),
                ),
        )
}
```

## **Performance Optimization Techniques**

### **Data Sampling for Large Datasets**

```go
type DataSampler struct {
    strategy SamplingStrategy
    maxPoints int
}

func (ds *DataSampler) Sample(points []DataPoint) []DataPoint {
    if len(points) <= ds.maxPoints {
        return points
    }
    
    switch ds.strategy {
    case SamplingStrategyLTTB:
        return ds.sampleLTTB(points, ds.maxPoints)
    case SamplingStrategyEveryNth:
        return ds.sampleEveryNth(points, ds.maxPoints)
    case SamplingStrategyMinMax:
        return ds.sampleMinMax(points, ds.maxPoints)
    case SamplingStrategyAverage:
        return ds.sampleAverage(points, ds.maxPoints)
    default:
        return ds.sampleLTTB(points, ds.maxPoints)
    }
}

// Largest Triangle Three Buckets (LTTB) downsampling algorithm
func (ds *DataSampler) sampleLTTB(data []DataPoint, threshold int) []DataPoint {
    if threshold >= len(data) || threshold == 0 {
        return data
    }
    
    sampled := make([]DataPoint, 0, threshold)
    
    // Bucket size. Leave room for start and end data points
    every := float64(len(data)-2) / float64(threshold-2)
    
    sampled = append(sampled, data[0]) // Always add the first point
    
    a := 0 // Initially a is the first point in the triangle
    
    for i := 0; i < threshold-2; i++ {
        // Calculate point average for next bucket (containing c)
        avgRangeStart := int(math.Floor(float64(i+1)*every)) + 1
        avgRangeEnd := int(math.Floor(float64(i+2)*every)) + 1
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
        
        // Get the range for this bucket
        rangeOffs := int(math.Floor(float64(i)*every)) + 1
        rangeTo := int(math.Floor(float64(i+1)*every)) + 1
        
        // Point a
        pointA := data[a]
        maxArea := -1.0
        var maxAreaPoint DataPoint
        
        for ; rangeOffs < rangeTo; rangeOffs++ {
            // Calculate triangle area
            area := math.Abs(
                (pointA.X-avgX)*(data[rangeOffs].Y-pointA.Y) -
                (pointA.X-data[rangeOffs].X)*(avgY-pointA.Y),
            ) * 0.5
            
            if area > maxArea {
                maxArea = area
                maxAreaPoint = data[rangeOffs]
                a = rangeOffs // Next a is this b
            }
        }
        
        sampled = append(sampled, maxAreaPoint)
    }
    
    sampled = append(sampled, data[len(data)-1]) // Always add last
    
    return sampled
}
```

### **Web Workers for Data Processing**

```go
type ChartWorker struct {
    worker     *app.Worker
    operations chan WorkerOperation
    results    chan WorkerResult
}

type WorkerOperation struct {
    ID     string
    Type   WorkerOpType
    Data   interface{}
    Params map[string]interface{}
}

func NewChartWorker() *ChartWorker {
    cw := &ChartWorker{
        operations: make(chan WorkerOperation, 100),
        results:    make(chan WorkerResult, 100),
    }
    
    // Create Web Worker
    cw.worker = app.NewWorker("/workers/chart-worker.js")
    
    // Set up message handlers
    cw.setupMessageHandlers()
    
    // Start processing loop
    go cw.processOperations()
    
    return cw
}

func (cw *ChartWorker) processOperations() {
    for op := range cw.operations {
        switch op.Type {
        case WorkerOpCalculateStatistics:
            cw.calculateStatistics(op)
        case WorkerOpDownsampleData:
            cw.downsampleData(op)
        case WorkerOpCalculateRegression:
            cw.calculateRegression(op)
        case WorkerOpClusterPoints:
            cw.clusterPoints(op)
        }
    }
}

func (cw *ChartWorker) calculateStatistics(op WorkerOperation) {
    // Perform heavy statistical calculations in worker
    data := op.Data.([]float64)
    
    // Calculate in worker
    cw.worker.PostMessage(map[string]interface{}{
        "id":   op.ID,
        "type": "calculateStatistics",
        "data": data,
    })
}
```

## **Export and Integration Features**

### **Chart Export System**

```go
type ChartExporter struct {
    chart *BaseChart
}

func (ce *ChartExporter) Export(format ExportFormat, options ExportOptions) error {
    switch format {
    case ExportFormatPNG:
        return ce.exportPNG(options)
    case ExportFormatSVG:
        return ce.exportSVG(options)
    case ExportFormatPDF:
        return ce.exportPDF(options)
    case ExportFormatCSV:
        return ce.exportCSV(options)
    case ExportFormatJSON:
        return ce.exportJSON(options)
    default:
        return fmt.Errorf("unsupported export format: %v", format)
    }
}

func (ce *ChartExporter) exportPNG(options ExportOptions) error {
    // Get canvas data URL
    dataURL, err := ce.chart.canvasRef.ToDataURL("image/png")
    if err != nil {
        return err
    }
    
    // Create download link
    app.Window().Call("downloadDataURL", dataURL, 
        fmt.Sprintf("chart_%s.png", time.Now().Format("20060102_150405")))
    
    return nil
}

func (ce *ChartExporter) exportCSV(options ExportOptions) error {
    var csv strings.Builder
    
    // Write headers
    headers := []string{"Label"}
    for _, dataset := range ce.chart.spec.Data.Datasets {
        headers = append(headers, dataset.Label)
    }
    csv.WriteString(strings.Join(headers, ",") + "\n")
    
    // Write data rows
    for i, label := range ce.chart.spec.Data.Labels {
        row := []string{label}
        for _, dataset := range ce.chart.spec.Data.Datasets {
            if i < len(dataset.Data) {
                row = append(row, fmt.Sprintf("%f", dataset.Data[i].Y))
            } else {
                row = append(row, "")
            }
        }
        csv.WriteString(strings.Join(row, ",") + "\n")
    }
    
    // Trigger download
    app.Window().Call("downloadText", csv.String(), 
        fmt.Sprintf("chart_data_%s.csv", time.Now().Format("20060102_150405")))
    
    return nil
}
```

## **Testing and Quality Assurance**

```go
func TestChartComponents(t *testing.T) {
    t.Run("renders line chart correctly", func(t *testing.T) {
        chart := NewLineChart(testData)
        rendered := chart.Render()
        
        assert.NotNil(t, rendered)
        assert.Contains(t, rendered.String(), "canvas")
    })
    
    t.Run("handles large datasets efficiently", func(t *testing.T) {
        largeData := generateTestData(100000) // 100k points
        
        start := time.Now()
        chart := NewLineChart(largeData)
        chart.Render()
        elapsed := time.Since(start)
        
        assert.Less(t, elapsed, 500*time.Millisecond, 
            "Rendering 100k points should take less than 500ms")
    })
    
    t.Run("calculates statistics correctly", func(t *testing.T) {
        data := []float64{1, 2, 3, 4, 5}
        stats := calculateStatistics(data)
        
        assert.InDelta(t, 3.0, stats.Mean, 0.001)
        assert.InDelta(t, 1.581, stats.StdDev, 0.001)
    })
    
    t.Run("downsamples data while preserving shape", func(t *testing.T) {
        data := generateSineWave(1000)
        sampled := downsampleLTTB(data, 100)
        
        assert.Len(t, sampled, 100)
        
        // Check that extremes are preserved
        assert.Equal(t, data[0].X, sampled[0].X)
        assert.Equal(t, data[len(data)-1].X, sampled[len(sampled)-1].X)
    })
    
    t.Run("meets accessibility requirements", func(t *testing.T) {
        chart := NewAccessibleChart(testData)
        rendered := chart.Render()
        
        // Check ARIA attributes
        assert.Contains(t, rendered.String(), "role=\"img\"")
        assert.Contains(t, rendered.String(), "aria-label")
        
        // Check for hidden data table
        assert.Contains(t, rendered.String(), "sr-only")
    })
}
```

## **Best Practices and Considerations**

### **Performance Optimization Checklist**
- ✓ Implement data sampling for >10k points
- ✓ Use WebGL for >50k points
- ✓ Debounce resize and interaction events
- ✓ Implement canvas layer caching
- ✓ Use requestAnimationFrame for animations
- ✓ Offload heavy calculations to Web Workers
- ✓ Implement virtual rendering for extremely large datasets

### **Accessibility Checklist**
- ✓ Provide ARIA labels and descriptions
- ✓ Include hidden data tables for screen readers
- ✓ Support keyboard navigation
- ✓ Ensure proper color contrast
- ✓ Provide text alternatives for color-coded information
- ✓ Support high contrast mode
- ✓ Make interactive elements keyboard focusable

### **Responsive Design Considerations**
- ✓ Maintain aspect ratio or allow flexible sizing
- ✓ Adjust font sizes based on container size
- ✓ Responsive hide/show chart elements
- ✓ Touch-optimized interactions
- ✓ Mobile-first tooltip positioning
- ✓ Adaptive sampling based on screen size

## **Conclusion**

Building data visualization chart components in go-app requires balancing performance, interactivity, and accessibility while leveraging Go's strengths in data processing. By implementing a multi-engine architecture with canvas, SVG, and WebGL renderers, we can handle datasets of varying sizes efficiently.

Key architectural insights include:
1. **Hybrid rendering approach**: Use appropriate technology for each use case
2. **Progressive enhancement**: Start with basic canvas, add WebGL for large datasets
3. **Accessibility-first design**: Build accessible charts from the ground up
4. **Performance optimization**: Implement downsampling, caching, and Web Workers
5. **Real-time capabilities**: Support streaming data with efficient updates

This implementation serves as a foundation that can be extended with:
- **Geospatial charts**: Maps with data overlays
- **3D visualizations**: WebGL-based 3D charts
- **Network graphs**: Force-directed diagrams
- **Hierarchical charts**: Tree maps, sunburst diagrams
- **Gantt charts**: Project timeline visualizations
- **Candlestick charts**: Financial data visualization
- **Radar charts**: Multi-variable comparisons

Remember that effective data visualization is both an art and a science. By providing flexible, performant charting components, you empower developers to create compelling data stories while maintaining the type safety and performance benefits of Go in the browser.


