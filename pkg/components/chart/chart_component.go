// pkg/components/chart/chart_component.go
package chart

import (
	"fmt"
	"math"
	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

type CanvasChart struct {
	app.Compo

	// Internal state
	canvas       app.Value
	ctx          app.Value
	width        int
	height       int
	dpr          float64
	streamData   StreamingData
	isRunning    bool
	hoverX       float64
	hoverY       float64
	showTooltip  bool
	activePoint  Point
	streamTicker app.Value
	shouldRender bool 

	// Configuration and Data
	config        ChartConfig
	currentPoints []Point
	Padding       Padding
	DataRange     DataRange
}

func (c *CanvasChart) OnMount(ctx app.Context) {
    ctx.Defer(func(ctx app.Context) {
        canvasJS := app.Window().GetElementByID("main-chart")
        
        if !canvasJS.Truthy() {
            app.Log("Error: Canvas element not found")
            return
        }

        c.canvas = canvasJS
        c.ctx = c.canvas.Call("getContext", "2d")

        c.dpr = app.Window().Get("devicePixelRatio").Float()
        if c.dpr == 0 {
            c.dpr = 1.0
        }

        c.ctx.Set("imageSmoothingEnabled", true)
        c.ctx.Set("textAlign", "center")
        c.ctx.Set("textBaseline", "middle")

        // Initialize default padding
        c.Padding = Padding{
            Top:    20,
            Right:  20,
            Bottom: 40,
            Left:   50,
        }

        c.resize()
        
        // Draw based on what type of chart we have
        c.drawAll()
        
        // Start streaming if configured
		if c.config.IsStream {
			c.startStreaming(ctx)
		}
    })
}

func (c *CanvasChart) drawPoints(data []Point, color string) {
	c.ctx.Set("fillStyle", "white")
	c.ctx.Set("strokeStyle", color)
	c.ctx.Set("lineWidth", 2)

	for _, pt := range data {
		px, py := c.ToPixels(pt.X, pt.Y)
		c.ctx.Call("beginPath")
		c.ctx.Call("arc", px, py, 4, 0, 6.28) // 6.28 radians approx 2*PI
		c.ctx.Call("fill")
		c.ctx.Call("stroke")
	}
}

func (c *CanvasChart) resize() {
	// Set logical size (CSS)
	c.width = 800
	c.height = 400

	// Set physical size (Actual pixels on the backing store)
	c.canvas.Set("width", float64(c.width)*c.dpr)
	c.canvas.Set("height", float64(c.height)*c.dpr)

	// Scale the context so drawing commands use logical pixels
	c.ctx.Call("scale", c.dpr, c.dpr)
}

func (c *CanvasChart) drawPlaceholder() {
	// Clean slate
	c.ctx.Set("fillStyle", "white")
	c.ctx.Call("fillRect", 0, 0, c.width, c.height)

	// Draw a simple border to prove it works
	c.ctx.Set("strokeStyle", "#4a90e2")
	c.ctx.Set("lineWidth", 2)
	c.ctx.Call("strokeRect", 10, 10, c.width-20, c.height-20)
    
	// Text test
	c.ctx.Set("font", "16px sans-serif")
	c.ctx.Set("fillStyle", "#333")
	c.ctx.Call("fillText", "Canvas Initialized at v10", 30, 40)
}

func (c *CanvasChart) Render() app.UI {
    // Force an update if we need to render
    if c.shouldRender {
        c.shouldRender = false
        // This will trigger OnUpdate
    }

    px, py := c.ToPixels(c.activePoint.X, c.activePoint.Y)

    return app.Div().Class("chart-wrapper").Body(
        app.Canvas().
            Class("chart-canvas").
            ID("main-chart").
            Width(c.width).
            Height(c.height).
            OnMouseMove(c.OnMouseMove),
            // REMOVED: .OnMount(c.OnMount), // This was causing the error

        app.If(c.showTooltip, func() app.UI {
            return app.Div().
                Class("chart-tooltip").
                Style("left", fmt.Sprintf("%fpx", px+15)).
                Style("top", fmt.Sprintf("%fpx", py-50)).
                Body(
                    app.Div().Class("tooltip-header").Text("Data Point"),
                    app.Div().Class("tooltip-value").Text(
                        fmt.Sprintf("X: %.1f | Y: %.1f", c.activePoint.X, c.activePoint.Y),
                    ),
                )
        }),
    )
}

// Padding defines the drawing area boundaries
type Padding struct {
	Top, Right, Bottom, Left float64
}

// DataRange defines the min/max of your dataset
type DataRange struct {
	MinX, MaxX float64
	MinY, MaxY float64
}

// ToPixels converts a data point (x, y) to canvas coordinates
func (c *CanvasChart) ToPixels(x, y float64) (px, py float64) {
	chartWidth := float64(c.width) - c.Padding.Left - c.Padding.Right
	chartHeight := float64(c.height) - c.Padding.Top - c.Padding.Bottom

	// Calculate X
	xRatio := (x - c.DataRange.MinX) / (c.DataRange.MaxX - c.DataRange.MinX)
	px = c.Padding.Left + (xRatio * chartWidth)

	// Calculate Y (Inverted for Canvas)
	yRatio := (y - c.DataRange.MinY) / (c.DataRange.MaxY - c.DataRange.MinY)
	py = (float64(c.height) - c.Padding.Bottom) - (yRatio * chartHeight)

	return px, py
}

func (c *CanvasChart) drawAxes() {
	c.ctx.Set("strokeStyle", "#ccc")
	c.ctx.Set("lineWidth", 1)
	c.ctx.Call("beginPath")

	// Y-Axis Line
	px, pyTop := c.ToPixels(c.DataRange.MinX, c.DataRange.MaxY)
	_, pyBottom := c.ToPixels(c.DataRange.MinX, c.DataRange.MinY)
	c.ctx.Call("moveTo", px, pyTop)
	c.ctx.Call("lineTo", px, pyBottom)

	// X-Axis Line
	pxLeft, py := c.ToPixels(c.DataRange.MinX, c.DataRange.MinY)
	pxRight, _ := c.ToPixels(c.DataRange.MaxX, c.DataRange.MinY)
	c.ctx.Call("moveTo", pxLeft, py)
	c.ctx.Call("lineTo", pxRight, py)

	c.ctx.Call("stroke")
}

func (c *CanvasChart) DrawLine(data []Point, color string, thickness float64) {
	if len(data) < 2 {
		return
	}

	c.ctx.Set("strokeStyle", color)
	c.ctx.Set("lineWidth", thickness)
	c.ctx.Set("lineJoin", "round")
	c.ctx.Set("lineCap", "round")
	
	c.ctx.Call("beginPath")

	for i, pt := range data {
		px, py := c.ToPixels(pt.X, pt.Y)
		
		if i == 0 {
			c.ctx.Call("moveTo", px, py)
		} else {
			c.ctx.Call("lineTo", px, py)
		}
	}

	c.ctx.Call("stroke")
}

func (c *CanvasChart) DrawBarChart(values []float64, color string) {
	if len(values) == 0 {
		return
	}

	chartWidth := float64(c.width) - c.Padding.Left - c.Padding.Right
	// Calculate width per bar including a 20% gap
	barWidth := (chartWidth / float64(len(values))) * 0.8
	gap := (chartWidth / float64(len(values))) * 0.2

	c.ctx.Set("fillStyle", color)

	for i, val := range values {
		// Calculate X based on index
		xPos := c.Padding.Left + (float64(i) * (barWidth + gap)) + (gap / 2)
		
		// Map the Y value to pixels
		_, pyValue := c.ToPixels(0, val)
		_, pyBase := c.ToPixels(0, c.DataRange.MinY)

		barHeight := pyBase - pyValue

		// Draw the rectangle: x, y, width, height
		c.ctx.Call("fillRect", xPos, pyValue, barWidth, barHeight)
	}
}

// DrawHeatmap expects a 2D slice [rows][cols]
func (c *CanvasChart) DrawHeatmap(grid [][]float64) {
	rows := len(grid)
	cols := len(grid[0])

	cellWidth := (float64(c.width) - c.Padding.Left - c.Padding.Right) / float64(cols)
	cellHeight := (float64(c.height) - c.Padding.Top - c.Padding.Bottom) / float64(rows)

	for r := 0; r < rows; r++ {
		for l := 0; l < cols; l++ {
			val := grid[r][l]
			
			// Simple color mapping: higher value = more intense blue
			// In a real app, you'd use a proper lerp function for RGB
			intensity := int((val - c.DataRange.MinY) / (c.DataRange.MaxY - c.DataRange.MinY) * 255)
			color := fmt.Sprintf("rgb(0, %d, %d)", intensity, 255-intensity)

			x := c.Padding.Left + (float64(l) * cellWidth)
			y := c.Padding.Top + (float64(r) * cellHeight)

			c.ctx.Set("fillStyle", color)
			// +1/-1 to create a tiny "grid line" effect between cells
			c.ctx.Call("fillRect", x+0.5, y+0.5, cellWidth-1, cellHeight-1)
		}
	}
}

func (c *CanvasChart) OnMouseMove(ctx app.Context, e app.Event) {
    rect := e.Get("target").Call("getBoundingClientRect")
    mouseX := e.Get("clientX").Float() - rect.Get("left").Float()
    mouseY := e.Get("clientY").Float() - rect.Get("top").Float()

    closest := Point{} // Declared here
    minDist := 20.0 
    found := false

    for _, p := range c.currentPoints {
        px, py := c.ToPixels(p.X, p.Y)
        dx := mouseX - px
        dy := mouseY - py
        dist := math.Sqrt(dx*dx + dy*dy)

        if dist < minDist {
            minDist = dist
            closest = p // Assigned here
            found = true
        }
    }

    if found {
        c.showTooltip = true
        c.activePoint = closest // USED HERE - This fixes the compiler error
    } else {
        c.showTooltip = false
    }
    
    ctx.Update()
}

func (c *CanvasChart) drawAll() {
    if !c.ctx.Truthy() {
        return
    }

    // 1. Clear the canvas with a background
    c.ctx.Set("fillStyle", "#ffffff")
    c.ctx.Call("fillRect", 0, 0, c.width, c.height)

    // 2. Reset tooltip state
    c.showTooltip = false

    // 3. Draw specific content based on config
    if len(c.config.BoxData) > 0 {
        // Draw Box Plots
        c.calculateBoxPlotRange()
        chartWidth := float64(c.width) - c.Padding.Left - c.Padding.Right
        boxSpacing := chartWidth / float64(len(c.config.BoxData)+1)
        
        for i, stats := range c.config.BoxData {
            xPos := c.Padding.Left + (float64(i+1) * boxSpacing)
            c.DrawBoxPlot(stats, xPos, c.config.BoxWidth)
        }
        c.drawAxes()
        
    } else if len(c.config.HeatmapMatrix) > 0 {
        // Draw Heatmap
        c.calculateHeatmapRange()
        c.DrawHeatmap(c.config.HeatmapMatrix)
        c.drawAxes()
        
    } else if len(c.config.PieData) > 0 {
        // Draw Pie Chart - no axes for pie charts
        c.DrawPieChart(c.config.PieData, c.getPieColors())
        
    } else if len(c.config.BarData) > 0 {
        // Draw Bar Chart
        c.calculateBarChartRange()
        c.drawAxes()
        c.DrawBarChart(c.config.BarData, c.config.LineColor, c.config.BarColors)
        // Draw labels if provided
        if len(c.config.BarLabels) == len(c.config.BarData) {
            c.drawBarLabels(c.config.BarLabels)
        }
        
    } else if c.config.IsStream {
        // Draw Streaming Chart
        if len(c.currentPoints) > 0 {
            c.calculateLineChartRange()
            c.drawAxes()
            c.DrawLine(c.currentPoints, c.config.LineColor, c.config.Thickness)
        }
        
    } else if len(c.currentPoints) > 0 {
        // Draw Line Chart
        c.calculateLineChartRange()
        c.drawAxes()
        c.DrawLine(c.currentPoints, c.config.LineColor, c.config.Thickness)
        c.drawPoints(c.currentPoints, c.config.LineColor)
    }
}

func (c *CanvasChart) DrawRegression(data []Point, pointColor, lineColor string) {
	// 1. Draw the Scatter Points
	c.drawPoints(data, pointColor)

	// 2. Calculate the Line
	m, b := CalculateLinearRegression(data)

	// 3. Define the start and end points of the line based on DataRange
	xStart := c.DataRange.MinX
	yStart := m*xStart + b

	xEnd := c.DataRange.MaxX
	yEnd := m*xEnd + b

	// 4. Draw the Trend Line
	regressionPoints := []Point{
		{X: xStart, Y: yStart},
		{X: xEnd, Y: yEnd},
	}
	
	// Set a dashed line style for the regression trend
	c.ctx.Call("setLineDash", []any{10, 5}) // 10px dash, 5px gap
	c.DrawLine(regressionPoints, lineColor, 2.0)
	c.ctx.Call("setLineDash", []any{}) // Reset to solid line
}

func (c *CanvasChart) getHeatmapColor(val float64, scheme string) string {
	// Map 0-100 value to 0-255 intensity
	intensity := int((val / 100.0) * 255)

	switch scheme {
	case "Inferno":
		// Yellow to Red to Black
		return fmt.Sprintf("rgb(255, %d, 0)", 255-intensity)
	case "Grayscale":
		return fmt.Sprintf("rgb(%d, %d, %d)", intensity, intensity, intensity)
	default: // Blue-Green
		return fmt.Sprintf("rgb(0, %d, %d)", intensity, 255-intensity)
	}
}

func (c *CanvasChart) DrawPieChart(data []float64, colors []string) {
    total := 0.0
    for _, v := range data {
        total += v
    }

    centerX := float64(c.width) / 2
    centerY := float64(c.height) / 2
    radius := math.Min(centerX, centerY) * 0.8
    currentAngle := -math.Pi / 2 // Start at 12 o'clock

    for i, val := range data {
        sliceAngle := (val / total) * 2 * math.Pi
        
        c.ctx.Set("fillStyle", colors[i%len(colors)])
        c.ctx.Call("beginPath")
        c.ctx.Call("moveTo", centerX, centerY)
        c.ctx.Call("arc", centerX, centerY, radius, currentAngle, currentAngle+sliceAngle)
        c.ctx.Call("closePath")
        c.ctx.Call("fill")

        currentAngle += sliceAngle
    }

	if c.config.InnerRadiusRatio > 0 {
		c.ctx.Set("globalCompositeOperation", "destination-out")
		c.ctx.Call("beginPath")
		c.ctx.Call("arc", centerX, centerY, radius*c.config.InnerRadiusRatio, 0, 2*math.Pi)
		c.ctx.Call("fill")
		c.ctx.Set("globalCompositeOperation", "source-over") // Reset
	}

}

func (c *CanvasChart) DrawBoxPlot(stats BoxPlotStats, xPos float64, width float64) {
    // Map statistical values to Y pixels
    _, yMin := c.ToPixels(0, stats.Min)
    _, yQ1 := c.ToPixels(0, stats.Q1)
    _, yMed := c.ToPixels(0, stats.Median)
    _, yQ3 := c.ToPixels(0, stats.Q3)
    _, yMax := c.ToPixels(0, stats.Max)

    c.ctx.Set("strokeStyle", "#333")
    c.ctx.Set("lineWidth", 2)

    // 1. Draw Whisker (Vertical line from Min to Max)
    c.ctx.Call("beginPath")
    c.ctx.Call("moveTo", xPos, yMin)
    c.ctx.Call("lineTo", xPos, yMax)
    c.ctx.Call("stroke")

    // 2. Draw the Box (Q1 to Q3)
    c.ctx.Set("fillStyle", c.config.LineColor + "80") // Add transparency
    boxHeight := yQ1 - yQ3
    c.ctx.Call("fillRect", xPos-(width/2), yQ3, width, boxHeight)
    
    c.ctx.Set("strokeStyle", "#333")
    c.ctx.Call("strokeRect", xPos-(width/2), yQ3, width, boxHeight)

    // 3. Draw Median Line
    c.ctx.Set("strokeStyle", "white")
    c.ctx.Set("lineWidth", 2)
    c.ctx.Call("beginPath")
    c.ctx.Call("moveTo", xPos-(width/2), yMed)
    c.ctx.Call("lineTo", xPos+(width/2), yMed)
    c.ctx.Call("stroke")

    // 4. Draw Whisker caps
    c.ctx.Set("strokeStyle", "#333")
    c.ctx.Set("lineWidth", 2)
    c.ctx.Call("beginPath")
    c.ctx.Call("moveTo", xPos-(width/4), yMin)
    c.ctx.Call("lineTo", xPos+(width/4), yMin)
    c.ctx.Call("moveTo", xPos-(width/4), yMax)
    c.ctx.Call("lineTo", xPos+(width/4), yMax)
    c.ctx.Call("stroke")
}

func (c *CanvasChart) OnUpdate(ctx app.Context) {
    // Stop any existing streaming when config changes
    if !c.config.IsStream && c.isRunning {
        c.isRunning = false
        if c.streamTicker.Truthy() {
            app.Window().Call("clearTimeout", c.streamTicker)
            c.streamTicker = app.Undefined()
        }
    }
    
    // Redraw when any configuration changes
    if c.ctx.Truthy() {
        c.drawAll()
        
        // Start streaming if configured and not already running
        if c.config.IsStream && !c.isRunning {
            c.startStreaming(ctx)
        }
    }
}

func (c *CanvasChart) calculateBoxPlotRange() {
    if len(c.config.BoxData) == 0 {
        return
    }
    
    minY := c.config.BoxData[0].Min
    maxY := c.config.BoxData[0].Max
    
    for _, stats := range c.config.BoxData {
        if stats.Min < minY {
            minY = stats.Min
        }
        if stats.Max > maxY {
            maxY = stats.Max
        }
    }
    
    padding := (maxY - minY) * 0.1
    c.DataRange = DataRange{
        MinX: 0,
        MaxX: float64(len(c.config.BoxData) + 1),
        MinY: minY - padding,
        MaxY: maxY + padding,
    }
}

func (c *CanvasChart) calculateHeatmapRange() {
    if len(c.config.HeatmapMatrix) == 0 {
        return
    }
    
    minY := c.config.HeatmapMatrix[0][0]
    maxY := c.config.HeatmapMatrix[0][0]
    
    for _, row := range c.config.HeatmapMatrix {
        for _, val := range row {
            if val < minY {
                minY = val
            }
            if val > maxY {
                maxY = val
            }
        }
    }
    
    // Remove the unused rows variable:
    // rows := len(c.config.HeatmapMatrix)
    cols := len(c.config.HeatmapMatrix[0])
    
    c.DataRange = DataRange{
        MinX: 0,
        MaxX: float64(cols),
        MinY: minY,
        MaxY: maxY,
    }
    
    // Adjust padding for heatmap
    c.Padding = Padding{
        Top:    20,
        Right:  20,
        Bottom: 40,
        Left:   50,
    }
}

func (c *CanvasChart) calculateLineChartRange() {
    if len(c.currentPoints) == 0 {
        return
    }
    
    minX, maxX := c.currentPoints[0].X, c.currentPoints[0].X
    minY, maxY := c.currentPoints[0].Y, c.currentPoints[0].Y
    
    for _, p := range c.currentPoints {
        if p.X < minX {
            minX = p.X
        }
        if p.X > maxX {
            maxX = p.X
        }
        if p.Y < minY {
            minY = p.Y
        }
        if p.Y > maxY {
            maxY = p.Y
        }
    }
    
    xPadding := (maxX - minX) * 0.1
    yPadding := (maxY - minY) * 0.1
    
    c.DataRange = DataRange{
        MinX: minX - xPadding,
        MaxX: maxX + xPadding,
        MinY: minY - yPadding,
        MaxY: maxY + yPadding,
    }
}

func (c *CanvasChart) getPieColors() []string {
    // Default color palette for pie charts
    return []string{
        "#4f46e5", "#10b981", "#f59e0b", "#ef4444", 
        "#8b5cf6", "#06b6d4", "#84cc16", "#f97316",
    }
}

func (c *CanvasChart) startStreaming(ctx app.Context) {
    if c.config.IsStream && !c.isRunning {
        // Clear any existing data
        c.streamData = StreamingData{
            Points:   make([]float64, 0, c.config.Capacity),
            Capacity: c.config.Capacity,
        }
        c.currentPoints = []Point{}
        
        // Start generating data
        c.isRunning = true
        
        // Store the context for the stream loop
        c.streamLoop(ctx)
    }
}

func (c *CanvasChart) streamLoop(ctx app.Context) {
    if !c.isRunning {
        return
    }

    // 1. Update Data (Simulating a data feed)
    newValue := 50.0 + (app.Window().Get("Math").Call("random").Float() * 50.0)
    c.streamData.Push(newValue)

    // 2. Update points and redraw
    c.currentPoints = make([]Point, len(c.streamData.Points))
    for i, v := range c.streamData.Points {
        c.currentPoints[i] = Point{X: float64(i), Y: v}
    }
    
    // 3. Trigger a redraw
    ctx.Dispatch(func(ctx app.Context) {
        if c.ctx.Truthy() {
            c.drawAll()
        }
    })

    // 4. Schedule next update and store the ticker
    c.streamTicker = app.Window().Call("setTimeout", app.FuncOf(func(this app.Value, args []app.Value) any {
        if c.isRunning {
            c.streamLoop(ctx)
        }
        return nil
    }), 200) // ~5 FPS
}

func (c *CanvasChart) OnDismount() {
    // Stop any running streaming
    c.isRunning = false
    
    // Clear any pending timeouts
    if c.streamTicker.Truthy() {
        app.Window().Call("clearTimeout", c.streamTicker)
        c.streamTicker = app.Undefined()
    }
    
    // Reset state
    c.canvas = app.Undefined()
    c.ctx = app.Undefined()
    c.streamData = StreamingData{}
    c.currentPoints = nil
    c.showTooltip = false
    c.shouldRender = false
}

func (c *CanvasChart) resetCanvas() {
    // Clear any existing streaming
    c.isRunning = false
    if c.streamTicker.Truthy() {
        app.Window().Call("clearTimeout", c.streamTicker)
        c.streamTicker = app.Undefined()
    }
    
    // Reset canvas state
    c.canvas = app.Undefined()
    c.ctx = app.Undefined()
    c.shouldRender = true
}

// ShouldUpdate implements the app.Compo interface
func (c *CanvasChart) ShouldUpdate(next app.Compo) bool {
    // Always update to ensure charts render correctly
    return true
}

// DrawBarChart enhanced version with optional custom colors
func (c *CanvasChart) DrawBarChart(values []float64, defaultColor string, customColors []string) {
	if len(values) == 0 {
		return
	}

	chartWidth := float64(c.width) - c.Padding.Left - c.Padding.Right
	chartHeight := float64(c.height) - c.Padding.Top - c.Padding.Bottom
	
	// Calculate width per bar including a 20% gap
	barWidth := (chartWidth / float64(len(values))) * 0.7
	gap := (chartWidth / float64(len(values))) * 0.3

	for i, val := range values {
		// Calculate X position
		xPos := c.Padding.Left + (float64(i) * (barWidth + gap)) + (gap / 2)
		
		// Map the Y value to pixels
		_, pyValue := c.ToPixels(0, val)
		_, pyBase := c.ToPixels(0, c.DataRange.MinY)

		barHeight := pyBase - pyValue

		// Set color - use custom color if provided, otherwise default
		if customColors != nil && i < len(customColors) {
			c.ctx.Set("fillStyle", customColors[i])
		} else {
			c.ctx.Set("fillStyle", defaultColor)
		}

		// Draw the rectangle
		c.ctx.Call("fillRect", xPos, pyValue, barWidth, barHeight)
		
		// Add a subtle border
		c.ctx.Set("strokeStyle", "#333")
		c.ctx.Set("lineWidth", 1)
		c.ctx.Call("strokeRect", xPos, pyValue, barWidth, barHeight)
	}
}

// drawBarLabels adds labels under each bar
func (c *CanvasChart) drawBarLabels(labels []string) {
	if len(labels) == 0 {
		return
	}

	chartWidth := float64(c.width) - c.Padding.Left - c.Padding.Right
	barWidth := (chartWidth / float64(len(labels))) * 0.7
	gap := (chartWidth / float64(len(labels))) * 0.3

	c.ctx.Set("font", "12px sans-serif")
	c.ctx.Set("fillStyle", "#666")
	c.ctx.Set("textAlign", "center")
	c.ctx.Set("textBaseline", "top")

	for i, label := range labels {
		xPos := c.Padding.Left + (float64(i) * (barWidth + gap)) + (gap / 2) + (barWidth / 2)
		yPos := float64(c.height) - c.Padding.Bottom + 10
		c.ctx.Call("fillText", label, xPos, yPos)
	}
}

// calculateBarChartRange determines the Y-axis range for bar charts
func (c *CanvasChart) calculateBarChartRange() {
	if len(c.config.BarData) == 0 {
		return
	}
	
	maxY := c.config.BarData[0]
	for _, val := range c.config.BarData[1:] {
		if val > maxY {
			maxY = val
		}
	}
	
	// Add 10% padding above the highest bar
	yPadding := maxY * 0.1
	
	c.DataRange = DataRange{
		MinX: 0,
		MaxX: float64(len(c.config.BarData)),
		MinY: 0, // Bar charts typically start at 0
		MaxY: maxY + yPadding,
	}
}
