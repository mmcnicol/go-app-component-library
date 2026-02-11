// pkg/components/chart/regression_chart.go
package chart

import (
	"fmt"
	"math"
	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

// RegressionChartComponent must embed app.Compo to be a valid go-app component
type RegressionChartComponent struct {
	app.Compo
	*CanvasChart
	
	pointColor   string
	lineColor    string
	showEquation bool
	pointSize    float64
	datasetName  string
	data         []Point

	// Add a flag to track if we've initialized
	initialized bool
}

// NewRegressionChart creates a new RegressionChartComponent
func NewRegressionChart(data []Point, opts ...Option) *RegressionChartComponent {
	canvasChart := New(data, opts...)
	
	return &RegressionChartComponent{
		CanvasChart:  canvasChart,
		data:         data,
		pointColor:   canvasChart.config.LineColor,
		lineColor:    "#ef4444", // default line color
		showEquation: true,
		pointSize:    4,
		initialized:  false,
	}
}

// WithPointColor sets the point color
func (c *RegressionChartComponent) WithPointColor(color string) *RegressionChartComponent {
	c.pointColor = color
	return c
}

// WithLineColor sets the regression line color
func (c *RegressionChartComponent) WithLineColor(color string) *RegressionChartComponent {
	c.lineColor = color
	return c
}

// WithShowEquation toggles equation display
func (c *RegressionChartComponent) WithShowEquation(show bool) *RegressionChartComponent {
	c.showEquation = show
	return c
}

// WithPointSize sets the point size
func (c *RegressionChartComponent) WithPointSize(size float64) *RegressionChartComponent {
	c.pointSize = size
	return c
}

// WithDatasetName sets the dataset name
func (c *RegressionChartComponent) WithDatasetName(name string) *RegressionChartComponent {
	c.datasetName = name
	return c
}

// Render implements app.Compo
func (c *RegressionChartComponent) Render() app.UI {
    if c.CanvasChart == nil {
        return app.Div().Text("Chart not initialized")
    }
    
    // Generate a unique ID for this chart instance
    chartID := fmt.Sprintf("regression-chart-%p", c)
    
    // Initialize width/height if not set
    if c.CanvasChart.width == 0 {
        c.CanvasChart.width = 800
    }
    if c.CanvasChart.height == 0 {
        c.CanvasChart.height = 400
    }
    
    return app.Div().Class("chart-wrapper").Body(
        app.Canvas().
            Class("chart-canvas").
            ID(chartID).
            Width(c.CanvasChart.width).
            Height(c.CanvasChart.height).
            OnMouseMove(c.OnMouseMove),
        
        app.If(c.CanvasChart.showTooltip, func() app.UI {
            px, py := c.CanvasChart.ToPixels(c.CanvasChart.activePoint.X, c.CanvasChart.activePoint.Y)
            return app.Div().
                Class("chart-tooltip").
                Style("left", fmt.Sprintf("%fpx", px+15)).
                Style("top", fmt.Sprintf("%fpx", py-50)).
                Body(
                    app.Div().Class("tooltip-header").Text("Data Point"),
                    app.Div().Class("tooltip-value").Text(
                        fmt.Sprintf("X: %.1f | Y: %.1f", 
                            c.CanvasChart.activePoint.X, 
                            c.CanvasChart.activePoint.Y),
                    ),
                )
        }),
    )
}

// OnUpdate implements app.Compo
func (c *RegressionChartComponent) OnUpdate(ctx app.Context) {
	// Only draw if we're initialized and have a valid context
	if c.CanvasChart != nil && c.CanvasChart.ctx.Truthy() && c.initialized {
		c.drawRegressionWithEquation()
	}
}

/ Render implements app.Compo
func (c *RegressionChartComponent) Render() app.UI {
    if c.CanvasChart == nil {
        return app.Div().Text("Chart not initialized")
    }
    
    // Generate a unique ID for this chart instance
    chartID := fmt.Sprintf("regression-chart-%p", c)
    
    return app.Div().Class("chart-wrapper").Body(
        app.Canvas().
            Class("chart-canvas").
            ID(chartID).  // Use unique ID instead of "main-chart"
            Width(c.CanvasChart.width).
            Height(c.CanvasChart.height).
            OnMouseMove(c.OnMouseMove),  // Now handled by this component
        
        app.If(c.CanvasChart.showTooltip, func() app.UI {
            px, py := c.CanvasChart.ToPixels(c.CanvasChart.activePoint.X, c.CanvasChart.activePoint.Y)
            return app.Div().
                Class("chart-tooltip").
                Style("left", fmt.Sprintf("%fpx", px+15)).
                Style("top", fmt.Sprintf("%fpx", py-50)).
                Body(
                    app.Div().Class("tooltip-header").Text("Data Point"),
                    app.Div().Class("tooltip-value").Text(
                        fmt.Sprintf("X: %.1f | Y: %.1f", 
                            c.CanvasChart.activePoint.X, 
                            c.CanvasChart.activePoint.Y),
                    ),
                )
        }),
    )
}

// Add OnMouseMove handler
func (c *RegressionChartComponent) OnMouseMove(ctx app.Context, e app.Event) {
    if c.CanvasChart != nil {
        // Use the canvas chart's mouse handling logic
        c.CanvasChart.OnMouseMove(ctx, e)
    }
}

// OnDismount implements app.Compo
func (c *RegressionChartComponent) OnDismount() {
	// Clean up
	if c.CanvasChart != nil {
		c.CanvasChart.OnDismount()
	}
	c.initialized = false
}

// ShouldUpdate implements app.Compo
func (c *RegressionChartComponent) ShouldUpdate(next app.Compo) bool {
	// Always update to reflect changes
	return true
}

// Custom drawing method
func (c *RegressionChartComponent) drawRegressionWithEquation() {
	// Check the inner pointer AND the inner context
	if c.CanvasChart == nil || !c.CanvasChart.ctx.Truthy() || len(c.data) == 0 {
		return 
	}

	// Update currentPoints and DataRange
	c.currentPoints = c.data
	c.calculateLineChartRange()
	
	// Clear the canvas first
	c.CanvasChart.ctx.Set("fillStyle", "#ffffff")
	c.CanvasChart.ctx.Call("fillRect", 0, 0, c.width, c.height)
	
	// Draw axes
	c.drawAxes()
	
	// Draw the regression line and points
	c.DrawRegression(c.data, c.pointColor, c.lineColor)
	
	// Draw points with custom size
	c.drawPointsWithSize(c.data, c.pointColor, c.pointSize)
	
	// Draw equation and statistics if enabled
	if c.showEquation {
		c.drawRegressionStats(c.data)
	}
}

// Draw points with custom size
func (c *RegressionChartComponent) drawPointsWithSize(data []Point, color string, size float64) {
	if c.CanvasChart == nil || !c.CanvasChart.ctx.Truthy() {
		return
	}
	
	c.CanvasChart.ctx.Set("fillStyle", "white")
	c.CanvasChart.ctx.Set("strokeStyle", color)
	c.CanvasChart.ctx.Set("lineWidth", 2)

	for _, pt := range data {
		px, py := c.ToPixels(pt.X, pt.Y)
		c.CanvasChart.ctx.Call("beginPath")
		c.CanvasChart.ctx.Call("arc", px, py, size, 0, 6.28)
		c.CanvasChart.ctx.Call("fill")
		c.CanvasChart.ctx.Call("stroke")
	}
}

// Draw regression equation and R-squared
func (c *RegressionChartComponent) drawRegressionStats(data []Point) {
	if len(data) < 2 || c.CanvasChart == nil || !c.CanvasChart.ctx.Truthy() {
		return
	}

	m, b := CalculateLinearRegression(data)
	
	// Calculate R-squared
	r2 := calculateRSquared(data, m, b)
	
	// Format equation
	equation := fmt.Sprintf("y = %.2fx + %.2f", m, b)
	if b < 0 {
		equation = fmt.Sprintf("y = %.2fx - %.2f", m, math.Abs(b))
	}
	
	// Draw on canvas
	c.CanvasChart.ctx.Set("font", "14px monospace")
	c.CanvasChart.ctx.Set("fillStyle", "#333")
	c.CanvasChart.ctx.Set("textAlign", "left")
	c.CanvasChart.ctx.Set("textBaseline", "top")
	
	// Background for text
	c.CanvasChart.ctx.Set("fillStyle", "rgba(255, 255, 255, 0.9)")
	c.CanvasChart.ctx.Call("fillRect", c.Padding.Left, c.Padding.Top, 220, 65)
	
	c.CanvasChart.ctx.Set("strokeStyle", "#ccc")
	c.CanvasChart.ctx.Set("lineWidth", 1)
	c.CanvasChart.ctx.Call("strokeRect", c.Padding.Left, c.Padding.Top, 220, 65)
	
	// Draw text
	c.CanvasChart.ctx.Set("fillStyle", "#1e293b")
	c.CanvasChart.ctx.Call("fillText", "Regression Statistics:", c.Padding.Left+10, c.Padding.Top+10)
	c.CanvasChart.ctx.Set("fillStyle", "#2563eb")
	c.CanvasChart.ctx.Call("fillText", equation, c.Padding.Left+10, c.Padding.Top+30)
	c.CanvasChart.ctx.Set("fillStyle", "#0f172a")
	c.CanvasChart.ctx.Call("fillText", fmt.Sprintf("R² = %.3f", r2), c.Padding.Left+10, c.Padding.Top+50)
}

// Helper function to calculate R-squared
func calculateRSquared(data []Point, m, b float64) float64 {
	if len(data) < 2 {
		return 0
	}

	// Calculate mean of Y
	sumY := 0.0
	for _, p := range data {
		sumY += p.Y
	}
	meanY := sumY / float64(len(data))

	// Calculate SSres and SStot
	ssRes := 0.0
	ssTot := 0.0
	
	for _, p := range data {
		// Predicted Y from regression line
		yPred := m*p.X + b
		
		// Residual sum of squares
		ssRes += math.Pow(p.Y-yPred, 2)
		
		// Total sum of squares
		ssTot += math.Pow(p.Y-meanY, 2)
	}

	if ssTot == 0 {
		return 0
	}

	return 1 - (ssRes / ssTot)
}
