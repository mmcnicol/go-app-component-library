// pkg/components/chart/tooltip_manager.go
package chart

import (
    "fmt"
    "math"
    "github.com/maxence-charriere/go-app/v10/pkg/app"
)

type TooltipManager struct {
    chart         *BaseChart
    tooltipElem   app.HTMLDiv
    currentPoint  *DataPoint
    currentDataset int
    position      TooltipPosition
    formatter     TooltipFormatter
}

func (tm *TooltipManager) setupEventHandlers() {
    // Mouse move handler - need to attach to canvas element
    // tm.chart.canvasElem.On("mousemove", tm.handleMouseMove)
}

func (tm *TooltipManager) handleMouseMove(ctx app.Context, e app.Event) {
    // Get canvas position
    // Note: In go-app v10, we need to use JavaScript interop
    // This is a simplified version
    
    // rect := app.Window().GetElementByID(tm.chart.containerID + "-canvas").Call("getBoundingClientRect")
    // mouseX := e.Get("clientX").Float() - rect.Get("left").Float()
    // mouseY := e.Get("clientY").Float() - rect.Get("top").Float()
    
    // // Find nearest data point
    // point, datasetIndex, distance := tm.findNearestPoint(mouseX, mouseY)
    
    // // Show tooltip if close enough
    // if point != nil && distance < tm.chart.spec.Options.Tooltips.IntersectDistance {
    //     tm.showTooltip(point, datasetIndex, mouseX, mouseY)
    // } else {
    //     tm.hideTooltip()
    // }
}

func (tm *TooltipManager) findNearestPoint(canvasX, canvasY float64) (*DataPoint, int, float64) {
    var nearestPoint *DataPoint
    var nearestDataset int
    minDistance := math.MaxFloat64
    
    for i, dataset := range tm.chart.spec.Data.Datasets {
        for _, point := range dataset.Data {
            // Convert data point to canvas coordinates
            // This depends on having scales available
            // pointX := tm.chart.xScale.Convert(point.X)
            // pointY := tm.chart.yScale.Convert(point.Y)
            
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
    content := tm.formatter.Format(point, datasetIndex, tm.chart.spec.Data)
    
    // Position tooltip
    tm.positionTooltip(x, y)
    
    // Update content
    tm.tooltipElem.SetInnerHTML(content)
    tm.tooltipElem.SetStyle("display", "block")
    
    // Trigger custom event
    if tm.chart.onHover != nil {
        tm.chart.onHover(*point, datasetIndex)
    }
}

func (tm *TooltipManager) positionTooltip(x, y float64) {
    // Position the tooltip element
    tm.tooltipElem.SetStyle("left", fmt.Sprintf("%fpx", x))
    tm.tooltipElem.SetStyle("top", fmt.Sprintf("%fpx", y))
}

func (tm *TooltipManager) hideTooltip() {
    tm.tooltipElem.SetStyle("display", "none")
}
