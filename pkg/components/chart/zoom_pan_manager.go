// pkg/components/chart/zoom_pan_manager.go
package chart

import (
    "fmt"
    "github.com/maxence-charriere/go-app/v10/pkg/app"
)

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

func (zpm *ZoomPanManager) pushZoomState() {
    zpm.zoomHistory = append(zpm.zoomHistory, zpm.currentZoom)
    // Limit history size
    if len(zpm.zoomHistory) > 50 {
        zpm.zoomHistory = zpm.zoomHistory[1:]
    }
}
