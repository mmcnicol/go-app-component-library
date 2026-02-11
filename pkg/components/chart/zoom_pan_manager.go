// pkg/components/chart/zoom_pan_manager.go
package chart

import (
    //"fmt"
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

func NewZoomPanManager(chart *BaseChart) *ZoomPanManager {
    return &ZoomPanManager{
        chart: chart,
        currentZoom: ZoomState{
            xMin: 0,
            xMax: 100,
            yMin: 0,
            yMax: 100,
        },
    }
}

func (zpm *ZoomPanManager) setupZoomPan() {
    // Event handlers would be set up in the chart component
    // This is a simplified implementation
}

func (zpm *ZoomPanManager) handleWheel(ctx app.Context, e app.Event) {
    e.PreventDefault()
    
    delta := e.Get("deltaY").Float()
    zoomFactor := 1.1
    
    // Note: This implementation assumes we have access to mouse position
    // and canvas scales, which would require JavaScript interop
    
    // Apply zoom
    if delta < 0 {
        // Zoom in
        zpm.zoomAt(50, 50, zoomFactor) // Center on middle for now
    } else {
        // Zoom out
        zpm.zoomAt(50, 50, 1/zoomFactor)
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
    
    // Trigger zoom event
    if zpm.chart.onZoom != nil {
        zpm.chart.onZoom(AxisRange{X: [2]float64{newXMin, newXMax}, Y: [2]float64{newYMin, newYMax}})
    }
}

func (zpm *ZoomPanManager) pushZoomState() {
    zpm.zoomHistory = append(zpm.zoomHistory, zpm.currentZoom)
    // Limit history size
    if len(zpm.zoomHistory) > 50 {
        zpm.zoomHistory = zpm.zoomHistory[1:]
    }
}

func (zpm *ZoomPanManager) resetZoom() {
    zpm.pushZoomState()
    zpm.currentZoom = ZoomState{
        xMin: 0,
        xMax: 100,
        yMin: 0,
        yMax: 100,
    }
    
    if zpm.chart.onZoom != nil {
        zpm.chart.onZoom(AxisRange{X: [2]float64{0, 100}, Y: [2]float64{0, 100}})
    }
}
