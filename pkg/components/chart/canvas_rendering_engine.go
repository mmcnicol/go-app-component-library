// pkg/components/chart/canvas_rendering_engine.go
package chart

import (
	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

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
