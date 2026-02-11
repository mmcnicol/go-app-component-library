// pkg/components/chart/regression_chart.go
package chart

import (
	"fmt"
	"math"
	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

// RegressionChartComponent extends CanvasChart to add equation display
type RegressionChartComponent struct {
	app.Compo
	*CanvasChart
	pointColor   string
	lineColor    string
	showEquation bool
	pointSize    float64
	datasetName  string
	data         []Point
}

/*
// Override OnMount to customize the drawing
func (c *RegressionChartComponent) OnMount(ctx app.Context) {
	if c.CanvasChart == nil {
        return
    }
	// Call parent OnMount first
	c.CanvasChart.OnMount(ctx)
	
	// After canvas is mounted, we need to override the drawing
	ctx.Defer(func(ctx app.Context) {
		if c.ctx.Truthy() {
			c.drawRegressionWithEquation()
		}
	})
}
*/

func (c *RegressionChartComponent) OnMount(ctx app.Context) {
    if c.CanvasChart != nil {
        c.DrawRegression() // Capitalize this
    }
}

/*
// Override OnUpdate to handle updates
func (c *RegressionChartComponent) OnUpdate(ctx app.Context) bool {
	if c.CanvasChart != nil {
        c.drawRegression()
    }
    return true
}
*/

func (c *RegressionChartComponent) OnUpdate(ctx app.Context) bool {
    if c.CanvasChart != nil {
        c.DrawRegression()
    }
    return true // Required return
}

// OnDismount - delegate to embedded CanvasChart
func (c *RegressionChartComponent) OnDismount() {
	if c.CanvasChart != nil {
		c.CanvasChart.OnDismount()
	}
}

// ShouldUpdate - delegate to embedded CanvasChart
func (c *RegressionChartComponent) ShouldUpdate(next app.Compo) bool {
	if c.CanvasChart != nil {
		return c.CanvasChart.ShouldUpdate(next)
	}
	return true
}

/*
// Render implements app.Compo
func (c *RegressionChartComponent) Render() app.UI {
	if c.CanvasChart == nil {
		return app.Div().Text("Chart not initialized")
	}
	return c.CanvasChart.Render()
}
*/

func (c *RegressionChartComponent) Render() app.UI {
    return c.CanvasChart // Required return
}

// Custom drawing method
func (c *RegressionChartComponent) drawRegressionWithEquation() {
	if !c.ctx.Truthy() || len(c.data) == 0 {
		return
	}

	// Clear and draw the regression
	c.ctx.Set("fillStyle", "#ffffff")
	c.ctx.Call("fillRect", 0, 0, c.width, c.height)

	// Update currentPoints and DataRange
	c.currentPoints = c.data
	c.calculateLineChartRange()
	
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
	c.ctx.Set("fillStyle", "white")
	c.ctx.Set("strokeStyle", color)
	c.ctx.Set("lineWidth", 2)

	for _, pt := range data {
		px, py := c.ToPixels(pt.X, pt.Y)
		c.ctx.Call("beginPath")
		c.ctx.Call("arc", px, py, size, 0, 6.28)
		c.ctx.Call("fill")
		c.ctx.Call("stroke")
	}
}

// Draw regression equation and R-squared
func (c *RegressionChartComponent) drawRegressionStats(data []Point) {
	if len(data) < 2 {
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
	c.ctx.Set("font", "14px monospace")
	c.ctx.Set("fillStyle", "#333")
	c.ctx.Set("textAlign", "left")
	c.ctx.Set("textBaseline", "top")
	
	// Background for text
	c.ctx.Set("fillStyle", "rgba(255, 255, 255, 0.9)")
	c.ctx.Call("fillRect", c.Padding.Left, c.Padding.Top, 220, 65)
	
	c.ctx.Set("strokeStyle", "#ccc")
	c.ctx.Set("lineWidth", 1)
	c.ctx.Call("strokeRect", c.Padding.Left, c.Padding.Top, 220, 65)
	
	// Draw text
	c.ctx.Set("fillStyle", "#1e293b")
	c.ctx.Call("fillText", "Regression Statistics:", c.Padding.Left+10, c.Padding.Top+10)
	c.ctx.Set("fillStyle", "#2563eb")
	c.ctx.Call("fillText", equation, c.Padding.Left+10, c.Padding.Top+30)
	c.ctx.Set("fillStyle", "#0f172a")
	c.ctx.Call("fillText", fmt.Sprintf("R² = %.3f", r2), c.Padding.Left+10, c.Padding.Top+50)
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

// Ensure RegressionChartComponent implements app.Compo
var _ app.Compo = (*RegressionChartComponent)(nil)
