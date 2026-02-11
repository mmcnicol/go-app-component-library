// pkg/components/chart/tooltip_manager.go
package chart

import (
    "fmt"
    "math"
    "github.com/maxence-charriere/go-app/v10/pkg/app"
)

type TooltipManager struct {
    chart         *BaseChart
    tooltipElem   app.UI // Changed from app.HTMLDiv to app.UI
    currentPoint  *DataPoint
    currentDataset int
    position      TooltipPosition
    formatter     TooltipFormatter
}

func NewTooltipManager(chart *BaseChart) *TooltipManager {
    tm := &TooltipManager{
        chart: chart,
        formatter: DefaultTooltipFormatter{},
    }
    
    // Create tooltip element
    tm.tooltipElem = app.Div().
        Class("chart-tooltip").
        Style("position", "absolute").
        Style("display", "none").
        Style("background", "white").
        Style("border", "1px solid #ccc").
        Style("padding", "5px").
        Style("border-radius", "3px").
        Style("pointer-events", "none").
        Style("z-index", "1000")
    
    return tm
}

func (tm *TooltipManager) setupEventHandlers() {
    // Mouse move handler - would be attached to canvas element
    // Note: This is a simplified implementation
}

func (tm *TooltipManager) handleMouseMove(ctx app.Context, e app.Event) {
    // Get canvas position
    // This would require JavaScript interop
}

func (tm *TooltipManager) findNearestPoint(canvasX, canvasY float64) (*DataPoint, int, float64) {
    var nearestPoint *DataPoint
    var nearestDataset int
    minDistance := math.MaxFloat64
    
    for i, dataset := range tm.chart.spec.Data.Datasets {
        for _, point := range dataset.Data {
            // Convert data point to canvas coordinates
            // This depends on having scales available
            pointX := 0.0 // placeholder
            pointY := 0.0 // placeholder
            
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
    //content := tm.formatter.Format(point, datasetIndex, tm.chart.spec.Data)
    _ = tm.formatter.Format(point, datasetIndex, tm.chart.spec.Data)
    
    // Update tooltip position and content
    // In go-app v10, we would need to update the component state
    // and re-render the tooltip with new content and position
    
    // For now, this is a placeholder implementation
    tm.position = TooltipPosition{X: x, Y: y}
    tm.currentPoint = point
    tm.currentDataset = datasetIndex
    
    // Trigger custom event
    if tm.chart.onHover != nil {
        tm.chart.onHover(*point, datasetIndex)
    }
}

func (tm *TooltipManager) positionTooltip(x, y float64) {
    tm.position = TooltipPosition{X: x, Y: y}
}

func (tm *TooltipManager) hideTooltip() {
    tm.currentPoint = nil
}

// GetTooltipUI returns the tooltip UI element
func (tm *TooltipManager) GetTooltipUI() app.UI {
    if tm.currentPoint == nil {
        return app.Div().
            Class("chart-tooltip").
            Style("display", "none")
    }
    
    content := tm.formatter.Format(tm.currentPoint, tm.currentDataset, tm.chart.spec.Data)
    
    return app.Div().
        Class("chart-tooltip").
        Style("position", "absolute").
        Style("left", fmt.Sprintf("%fpx", tm.position.X)).
        Style("top", fmt.Sprintf("%fpx", tm.position.Y)).
        Style("background", "white").
        Style("border", "1px solid #ccc").
        Style("padding", "5px").
        Style("border-radius", "3px").
        Style("pointer-events", "none").
        Style("z-index", "1000").
        Text(content)
}
