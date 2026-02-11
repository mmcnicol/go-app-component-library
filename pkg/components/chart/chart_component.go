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
	canvas      app.Value
	ctx         app.Value
	width       int
	height      int
	dpr         float64
	streamData  StreamingData
	isRunning   bool
	hoverX      float64
	hoverY      float64
	showTooltip bool
	activePoint Point

	// Configuration and Data
	config        ChartConfig
	currentPoints []Point
	Padding       Padding
	DataRange     DataRange
}

/*
func (c *CanvasChart) OnMount(ctx app.Context) {
	c.canvas = ctx.SrcElement()
	c.ctx = c.canvas.Call("getContext", "2d")
	c.dpr = app.Window().Get("devicePixelRatio").Float()
	if c.dpr == 0 { c.dpr = 1.0 }

	// Phase B: Setup defaults
	c.Padding = Padding{Top: 20, Right: 20, Bottom: 40, Left: 50}
	c.DataRange = DataRange{MinX: 0, MaxX: 100, MinY: 0, MaxY: 100}

	c.resize()
	c.drawAxes() // Instead of placeholder
}
*/

/*
func (c *CanvasChart) OnMount(ctx app.Context) {
	// 1. Get the Canvas Element
	c.canvas = ctx.JSSrc()
	c.ctx = c.canvas.Call("getContext", "2d")

	// 2. Determine Display Density
	c.dpr = app.Window().Get("devicePixelRatio").Float()
	if c.dpr == 0 {
		c.dpr = 1.0
	}

	ctx.Set("imageSmoothingEnabled", true)

	// Test Data for Bars
    barData := []float64{45, 80, 55, 92, 30, 66, 78}
    
    c.resize()
    c.drawAxes()
    c.DrawBarChart(barData, "#4a90e2")



	c.streamData = StreamingData{
        Points:   make([]float64, 0, 100),
        Capacity: 100,
    }
    
    c.resize()
    
    // Start the animation
    c.StartStreaming(ctx)
}
*/

/*
func (c *CanvasChart) OnMount(ctx app.Context) {
	// 1. Get the Canvas Element using JSSrc()
	c.canvas = ctx.JSSrc()
	c.ctx = c.canvas.Call("getContext", "2d")

	// 2. Determine Display Density
	c.dpr = app.Window().Get("devicePixelRatio").Float()
	if c.dpr == 0 {
		c.dpr = 1.0
	}

	// Set property on the JS Context, not the app.Context
	c.ctx.Set("imageSmoothingEnabled", true)

	c.resize()
	c.drawAxes()
}
*/

func (c *CanvasChart) OnMount(ctx app.Context) {
    // 1. Wait for the next frame to ensure the DOM is fully rendered
    ctx.Defer(func(ctx app.Context) {
        // 2. Explicitly find the canvas element by its ID
        // Note: In a library, you should generate a unique ID, 
        // but for now we use the ID defined in Render()
        canvasJS := app.Window().GetElementByID("main-chart")
        
        if !canvasJS.Truthy() {
            app.Log("Error: Canvas element not found")
            return
        }

        c.canvas = canvasJS
        c.ctx = c.canvas.Call("getContext", "2d")

        // 3. Setup Display Density
        c.dpr = app.Window().Get("devicePixelRatio").Float()
        if c.dpr == 0 {
            c.dpr = 1.0
        }

        c.ctx.Set("imageSmoothingEnabled", true)

        // 4. Initialize layout and draw
        c.resize()
        c.drawAxes()
        
        // Trigger specific painters based on config
        if len(c.config.BoxData) > 0 {
            // Draw your box plots here
            for i, stats := range c.config.BoxData {
                xPos := float64(i+1) * 100.0 // Example spacing
                c.DrawBoxPlot(stats, xPos, c.config.BoxWidth)
            }
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

/*
func (c *CanvasChart) Render() app.UI {
	return app.Canvas().
		ID("main-chart").
		Style("width", "800px").
		Style("height", "400px").
		Style("background", "#f9f9f9").
		OnMount(c.OnMount)
}
*/

/*
func (c *CanvasChart) Render() app.UI {
    px, py := c.ToPixels(c.activePoint.X, c.activePoint.Y)

    return app.Div().Style("position", "relative").Body(
        app.Canvas().
            ID("main-chart").
            Style("width", "800px").
            Style("height", "400px").
            OnMount(c.OnMount).
            OnMouseMove(c.OnMouseMove), // Added event

        // Tooltip Overlay
        app.If(c.showTooltip, func() app.UI {
            return app.Div().
                Style("position", "absolute").
                Style("left", fmt.Sprintf("%fpx", px+10)).
                Style("top", fmt.Sprintf("%fpx", py-40)).
                Style("background", "rgba(0,0,0,0.8)").
                Style("color", "white").
                Style("padding", "5px 10px").
                Style("border-radius", "4px").
                Style("pointer-events", "none"). // Don't block mouse events
                Body(
                    app.Text(fmt.Sprintf("X: %.2f, Y: %.2f", c.activePoint.X, c.activePoint.Y)),
                )
        }),
    )
}
*/

/*
func (c *CanvasChart) Render() app.UI {
	px, py := c.ToPixels(c.activePoint.X, c.activePoint.Y)

	return app.Div().Class("chart-wrapper").Body(
		// The Drawing Surface
		app.Canvas().
			Class("chart-canvas").
			ID("main-chart").
			Width(c.width).
			Height(c.height).
			OnMount(c.OnMount).
			OnMouseMove(c.OnMouseMove),

		// Tooltip Logic
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

		// Legend (Optional)
		app.Div().Class("chart-legend").Body(
			app.Div().Class("legend-item").Body(
				app.Div().Class("legend-color").Style("background", c.config.LineColor),
				app.Span().Text(c.config.Title),
			),
		),
	)
}
*/

func (c *CanvasChart) Render() app.UI {
	px, py := c.ToPixels(c.activePoint.X, c.activePoint.Y)

	return app.Div().Class("chart-wrapper").Body(
		app.Canvas().
			Class("chart-canvas").
			ID("main-chart"). // This must match GetElementByID
			Width(c.width).
			Height(c.height).
			OnMouseMove(c.OnMouseMove),

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

func (c *CanvasChart) StartStreaming(ctx app.Context) {
	c.isRunning = true
	c.streamLoop(ctx)
}

func (c *CanvasChart) streamLoop(ctx app.Context) {
	if !c.isRunning {
		return
	}

	// 1. Update Data (Simulating a data feed)
	// In a real app, this might come from a WebSocket or Channel
	newValue := 50.0 + (app.Window().Get("Math").Call("random").Float() * 20.0)
	c.streamData.Push(newValue)

	// 2. Clear and Redraw
	ctx.Dispatch(func(ctx app.Context) {
		// Clear Canvas
		c.ctx.Call("clearRect", 0, 0, c.width, c.height)
		
		// Draw Infrastructure
		c.drawAxes()
		
		// Draw the moving line
		points := make([]Point, len(c.streamData.Points))
		for i, v := range c.streamData.Points {
			points[i] = Point{X: float64(i), Y: v}
		}
		
		// Update DataRange MaxX to match Capacity so the line "fills" the space
		c.DataRange.MaxX = float64(c.streamData.Capacity)
		c.DrawLine(points, "#28a745", 2.0)
	})

	// 3. Request next frame
	app.Window().Call("requestAnimationFrame", app.FuncOf(func(this app.Value, args []app.Value) any {
		c.streamLoop(ctx)
		return nil
	}))
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

/*
// Centralized drawing method to prevent code duplication
func (c *CanvasChart) drawAll() {
    c.ctx.Call("clearRect", 0, 0, c.width, c.height)
    c.drawAxes()
    c.DrawLine(c.currentPoints, c.config.LineColor, c.config.Thickness)
}
*/

func (c *CanvasChart) drawAll() {
	if !c.ctx.Truthy() {
		return
	}

	// 1. Clear the canvas
	c.ctx.Call("clearRect", 0, 0, c.width, c.height)

	// 2. Draw the foundation
	c.drawAxes()

	// 3. Draw specific content based on config
	if len(c.config.BoxData) > 0 {
		// Draw Box Plots
		// Space them out based on the number of items
		for i, stats := range c.config.BoxData {
			// Calculate horizontal position
			xPos := c.Padding.Left + (float64(i+1) * 100.0) 
			c.DrawBoxPlot(stats, xPos, c.config.BoxWidth)
		}
	} else if len(c.currentPoints) > 0 {
		// Draw Line/Scatter if points exist
		c.DrawLine(c.currentPoints, c.config.LineColor, c.config.Thickness)
	}
    // ... add other types like Pie/Heatmap here ...
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
	c.ctx.Set("fillStyle", c.config.LineColor)
	boxHeight := yQ1 - yQ3
	c.ctx.Call("fillRect", xPos-(width/2), yQ3, width, boxHeight)
	c.ctx.Call("strokeRect", xPos-(width/2), yQ3, width, boxHeight)

	// 3. Draw Median Line
	c.ctx.Set("strokeStyle", "white")
	c.ctx.Call("beginPath")
	c.ctx.Call("moveTo", xPos-(width/2), yMed)
	c.ctx.Call("lineTo", xPos+(width/2), yMed)
	c.ctx.Call("stroke")
}

/*
func (c *CanvasChart) OnUpdate(ctx app.Context) {
    // Check if data changed and trigger a redraw
    ctx.Dispatch(func(ctx app.Context) {
        c.drawAll() 
    })
}
*/

func (c *CanvasChart) OnUpdate(ctx app.Context) {
    // Whenever the storybook controls change the config, redraw
	c.drawAll()
}
