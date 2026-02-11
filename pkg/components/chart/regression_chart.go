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

func (c *RegressionChartComponent) OnMount(ctx app.Context) {
    if c.CanvasChart == nil {
        return
    }
    
    // Call parent OnMount
    c.CanvasChart.OnMount(ctx)
    
    ctx.Defer(func(ctx app.Context) {
        // Double check: is the pointer valid AND is the JS context initialized?
        if c.CanvasChart != nil && c.CanvasChart.ctx.Truthy() {
            c.drawRegressionWithEquation()
        }
    })
}

func (c *RegressionChartComponent) OnUpdate(ctx app.Context) bool {
    if c.CanvasChart == nil {
        return false
    }
    
    c.CanvasChart.OnUpdate(ctx)
    
    ctx.Defer(func(ctx app.Context) {
        if c.CanvasChart != nil && c.CanvasChart.ctx.Truthy() {
            c.drawRegressionWithEquation()
        }
    })
    return true
}

func (c *RegressionChartComponent) Render() app.UI {
    if c.CanvasChart == nil {
        return app.Div().Text("Chart not initialized")
    }
    return c.CanvasChart
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

// Custom drawing method
func (c *RegressionChartComponent) drawRegressionWithEquation() {
	if !c.CanvasChart.ctx.Truthy() || len(c.data) == 0 {
		return
	}

	// Clear and draw the regression
	c.CanvasChart.ctx.Set("fillStyle", "#ffffff")
	c.CanvasChart.ctx.Call("fillRect", 0, 0, c.width, c.height)

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

// Ensure RegressionChartComponent implements app.Compo
//var _ app.Compo = (*RegressionChartComponent)(nil)
