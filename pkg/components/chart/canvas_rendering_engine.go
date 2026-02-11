// pkg/components/chart/canvas_rendering_engine.go
package chart

import (
    "fmt"
    "github.com/maxence-charriere/go-app/v10/pkg/app"
)

// CanvasRenderer implements ChartEngine
type CanvasRenderer struct {
    app.Compo // Add this to make it a component
    canvasID   string
    width      int
    height     int
    pixelRatio float64
    mounted    bool
    chartSpec  ChartSpec // Store the chart spec
    ctx        app.Context // Store context for updates
    ctxValid   bool // Add a flag to track if context is valid
}

func NewCanvasRenderer(containerID string) (*CanvasRenderer, error) {
    return &CanvasRenderer{
        canvasID:   containerID + "-canvas",
        pixelRatio: 1.0,
        width:      800,
        height:     400,
        ctxValid:   false,
    }, nil
}

// RenderChart renders the chart spec
func (cr *CanvasRenderer) RenderChart(chart ChartSpec) error {
    // Store the chart spec for rendering
    cr.chartSpec = chart
    
    // If already mounted and context is valid, update immediately
    if cr.mounted && cr.ctxValid {
        cr.setupCanvas(cr.ctx)
    }
    
    return nil
}

func (cr *CanvasRenderer) setupCanvas(ctx app.Context) {
    if cr.chartSpec.Type == "" {
        return // No chart to render
    }
    
    // Set canvas dimensions after mount
    app.Window().Call("eval", fmt.Sprintf(`
        (function() {
            const canvas = document.getElementById('%s');
            if (!canvas) {
                console.log('Canvas not found: %s');
                return;
            }
            
            const container = canvas.parentElement;
            if (!container) {
                console.log('Container not found for canvas');
                return;
            }
            
            // Get container dimensions
            const width = container.clientWidth || 800;
            const height = container.clientHeight || 400;
            
            // Set canvas dimensions
            canvas.width = width;
            canvas.height = height;
            canvas.style.width = width + 'px';
            canvas.style.height = height + 'px';
            
            // Draw the chart
            const ctx = canvas.getContext('2d');
            %s
        })();
    `, cr.canvasID, cr.canvasID, cr.getDrawScript()))
}

func (cr *CanvasRenderer) getDrawScript() string {
    if cr.chartSpec.Type == "" {
        return "console.log('No chart spec available');"
    }
    
    switch cr.chartSpec.Type {
    case ChartTypeBar:
        return cr.getBarChartScript()
    case ChartTypeLine:
        return cr.getLineChartScript()
    default:
        return cr.getBarChartScript()
    }
}

func (cr *CanvasRenderer) getLineChartScript() string {
    data := cr.chartSpec.Data
    if len(data.Datasets) == 0 {
        return "console.log('No data for line chart');"
    }
    
    // Prepare data for JavaScript
    datasets := ""
    for _, dataset := range data.Datasets {
        points := ""
        for _, point := range dataset.Data {
            points += fmt.Sprintf("{x: %f, y: %f},", point.X, point.Y)
        }
        color := "#4A90E2"
        if len(dataset.BackgroundColor) > 0 {
            color = dataset.BackgroundColor[0]
        }
        datasets += fmt.Sprintf(`{
            label: '%s',
            data: [%s],
            color: '%s',
            borderWidth: %d,
            pointRadius: %d
        },`, dataset.Label, points, color, dataset.BorderWidth, dataset.PointRadius)
    }
    
    labels := ""
    for _, label := range data.Labels {
        labels += fmt.Sprintf("'%s',", label)
    }
    
    return fmt.Sprintf(`
        // Draw line chart
        function drawLineChart(ctx, width, height) {
            // Clear canvas
            ctx.clearRect(0, 0, width, height);
            
            // Draw white background
            ctx.fillStyle = '#ffffff';
            ctx.fillRect(0, 0, width, height);
            
            const datasets = [%s];
            const labels = [%s];
            const numPoints = datasets[0].data.length;
            
            // Calculate scales
            let minY = Infinity;
            let maxY = -Infinity;
            
            datasets.forEach(dataset => {
                dataset.data.forEach(point => {
                    if (point.y < minY) minY = point.y;
                    if (point.y > maxY) maxY = point.y;
                });
            });
            
            // Add padding
            const yPadding = (maxY - minY) * 0.1;
            minY -= yPadding;
            maxY += yPadding;
            
            // Margins
            const margin = { top: 20, right: 80, bottom: 60, left: 60 };
            const plotWidth = width - margin.left - margin.right;
            const plotHeight = height - margin.top - margin.bottom;
            
            // X scale
            const xScale = (index) => margin.left + (index * plotWidth / (numPoints - 1));
            
            // Y scale (inverted)
            const yScale = (value) => margin.top + plotHeight - ((value - minY) / (maxY - minY)) * plotHeight;
            
            // Draw grid
            ctx.strokeStyle = '#e0e0e0';
            ctx.lineWidth = 0.5;
            
            // Horizontal grid lines
            for (let i = 0; i <= 5; i++) {
                const y = minY + ((maxY - minY) * i / 5);
                const yPos = yScale(y);
                ctx.beginPath();
                ctx.moveTo(margin.left, yPos);
                ctx.lineTo(width - margin.right, yPos);
                ctx.stroke();
            }
            
            // Draw axes
            ctx.strokeStyle = '#000000';
            ctx.lineWidth = 1;
            ctx.font = '12px Arial';
            ctx.fillStyle = '#000000';
            
            // Y axis
            ctx.beginPath();
            ctx.moveTo(margin.left, margin.top);
            ctx.lineTo(margin.left, height - margin.bottom);
            ctx.stroke();
            
            // X axis
            ctx.beginPath();
            ctx.moveTo(margin.left, height - margin.bottom);
            ctx.lineTo(width - margin.right, height - margin.bottom);
            ctx.stroke();
            
            // Draw lines
            datasets.forEach((dataset) => {
                ctx.strokeStyle = dataset.color;
                ctx.lineWidth = dataset.borderWidth || 2;
                ctx.beginPath();
                
                dataset.data.forEach((point, i) => {
                    const x = xScale(i);
                    const y = yScale(point.y);
                    
                    if (i === 0) {
                        ctx.moveTo(x, y);
                    } else {
                        ctx.lineTo(x, y);
                    }
                });
                
                ctx.stroke();
                
                // Draw points
                if (dataset.pointRadius > 0) {
                    ctx.fillStyle = dataset.color;
                    dataset.data.forEach((point, i) => {
                        const x = xScale(i);
                        const y = yScale(point.y);
                        ctx.beginPath();
                        ctx.arc(x, y, dataset.pointRadius, 0, Math.PI * 2);
                        ctx.fill();
                    });
                }
            });
            
            // Draw labels on X axis
            ctx.textAlign = 'center';
            ctx.textBaseline = 'top';
            labels.forEach((label, i) => {
                if (i < numPoints) {
                    const x = xScale(i);
                    ctx.fillText(label, x, height - margin.bottom + 10);
                }
            });
            
            // Draw legend
            const legendX = width - margin.right + 10;
            const legendY = margin.top;
            const boxSize = 15;
            const spacing = 25;
            
            datasets.forEach((dataset, i) => {
                const y = legendY + i * spacing;
                
                // Draw color box
                ctx.fillStyle = dataset.color;
                ctx.fillRect(legendX, y, boxSize, boxSize);
                
                // Draw label
                ctx.fillStyle = '#000000';
                ctx.textAlign = 'left';
                ctx.textBaseline = 'middle';
                ctx.fillText(dataset.label, legendX + boxSize + 10, y + boxSize/2);
            });
        }
        
        // Call draw function
        drawLineChart(ctx, width, height);
    `, datasets, labels)
}

func (cr *CanvasRenderer) getBarChartScript() string {
    data := cr.chartSpec.Data
    if len(data.Datasets) == 0 {
        return "console.log('No data for bar chart');"
    }
    
    // Prepare data for JavaScript
    datasets := ""
    for _, dataset := range data.Datasets {
        points := ""
        for _, point := range dataset.Data {
            points += fmt.Sprintf("{x: %f, y: %f},", point.X, point.Y)
        }
        color := "#4A90E2"
        if len(dataset.BackgroundColor) > 0 {
            color = dataset.BackgroundColor[0]
        }
        datasets += fmt.Sprintf(`{
            label: '%s',
            data: [%s],
            color: '%s'
        },`, dataset.Label, points, color)
    }
    
    labels := ""
    for _, label := range data.Labels {
        labels += fmt.Sprintf("'%s',", label)
    }
    
    return fmt.Sprintf(`
        // Draw bar chart
        function drawBarChart(ctx, width, height) {
            // Clear canvas
            ctx.clearRect(0, 0, width, height);
            
            // Draw white background
            ctx.fillStyle = '#ffffff';
            ctx.fillRect(0, 0, width, height);
            
            const datasets = [%s];
            const labels = [%s];
            const numDatasets = datasets.length;
            const numPoints = datasets[0].data.length;
            
            // Calculate scales
            let minY = Infinity;
            let maxY = -Infinity;
            
            datasets.forEach(dataset => {
                dataset.data.forEach(point => {
                    if (point.y < minY) minY = point.y;
                    if (point.y > maxY) maxY = point.y;
                });
            });
            
            // Add padding
            const yPadding = (maxY - minY) * 0.1;
            minY -= yPadding;
            maxY += yPadding;
            
            // Ensure minY is 0 if BeginAtZero is true
            if (%t && minY < 0) {
                minY = 0;
            }
            
            // Margins
            const margin = { top: 20, right: 80, bottom: 60, left: 60 };
            const plotWidth = width - margin.left - margin.right;
            const plotHeight = height - margin.top - margin.bottom;
            
            // X scale
            const xScale = (index) => margin.left + (index * plotWidth / numPoints);
            
            // Y scale (inverted)
            const yScale = (value) => {
                if (maxY === minY) return margin.top + plotHeight / 2;
                return margin.top + plotHeight - ((value - minY) / (maxY - minY)) * plotHeight;
            };
            
            // Draw grid
            ctx.strokeStyle = '#e0e0e0';
            ctx.lineWidth = 0.5;
            
            // Horizontal grid lines
            for (let i = 0; i <= 5; i++) {
                const y = minY + ((maxY - minY) * i / 5);
                const yPos = yScale(y);
                ctx.beginPath();
                ctx.moveTo(margin.left, yPos);
                ctx.lineTo(width - margin.right, yPos);
                ctx.stroke();
            }
            
            // Draw axes
            ctx.strokeStyle = '#000000';
            ctx.lineWidth = 1;
            ctx.font = '12px Arial';
            ctx.fillStyle = '#000000';
            
            // Y axis
            ctx.beginPath();
            ctx.moveTo(margin.left, margin.top);
            ctx.lineTo(margin.left, height - margin.bottom);
            ctx.stroke();
            
            // X axis
            ctx.beginPath();
            ctx.moveTo(margin.left, height - margin.bottom);
            ctx.lineTo(width - margin.right, height - margin.bottom);
            ctx.stroke();
            
            // Draw bars
            const barWidth = (plotWidth / numPoints) / numDatasets * 0.8;
            
            datasets.forEach((dataset, datasetIdx) => {
                ctx.fillStyle = dataset.color;
                
                dataset.data.forEach((point, pointIdx) => {
                    const xPos = xScale(pointIdx) + 
                                (datasetIdx * barWidth) - 
                                (numDatasets * barWidth / 2) + 
                                (barWidth / 2);
                    const yPos = yScale(point.y);
                    const barHeight = yScale(0) - yPos;
                    
                    // Draw bar
                    ctx.fillRect(
                        xPos - barWidth/2,
                        yPos,
                        barWidth,
                        barHeight
                    );
                    
                    // Draw border
                    ctx.strokeStyle = '#333';
                    ctx.lineWidth = 1;
                    ctx.strokeRect(
                        xPos - barWidth/2,
                        yPos,
                        barWidth,
                        barHeight
                    );
                });
            });
            
            // Draw labels on X axis
            ctx.textAlign = 'center';
            ctx.textBaseline = 'top';
            labels.forEach((label, i) => {
                const x = xScale(i);
                ctx.fillText(label, x, height - margin.bottom + 10);
            });
            
            // Draw legend
            const legendX = width - margin.right + 10;
            const legendY = margin.top;
            const boxSize = 15;
            const spacing = 25;
            
            datasets.forEach((dataset, i) => {
                const y = legendY + i * spacing;
                
                // Draw color box
                ctx.fillStyle = dataset.color;
                ctx.fillRect(legendX, y, boxSize, boxSize);
                
                // Draw label
                ctx.fillStyle = '#000000';
                ctx.textAlign = 'left';
                ctx.textBaseline = 'middle';
                ctx.fillText(dataset.label, legendX + boxSize + 10, y + boxSize/2);
            });
        }
        
        // Call draw function
        drawBarChart(ctx, width, height);
    `, datasets, labels, cr.chartSpec.Options.Scales.Y.BeginAtZero)
}

// Update implements ChartEngine.Update
func (cr *CanvasRenderer) Update(data ChartData) error {
    // Update the chart spec and re-render
    cr.chartSpec.Data = data
    if cr.mounted && cr.ctxValid {
        cr.setupCanvas(cr.ctx)
    }
    return nil
}

func (cr *CanvasRenderer) Destroy() error {
    // Clean up resources
    return nil
}

// Render method for the component (implements app.UI)
func (cr *CanvasRenderer) Render() app.UI {
    return app.Canvas().
        ID(cr.canvasID).
        Class("chart-canvas").
        Style("width", "100%").
        Style("height", "100%").
        Style("display", "block")
}

// OnMount is called when the component is mounted
func (cr *CanvasRenderer) OnMount(ctx app.Context) {
    cr.mounted = true
    cr.ctx = ctx
    cr.ctxValid = true // Mark context as valid
    
    // Schedule drawing after mount
    ctx.Defer(func(ctx app.Context) {
        cr.setupCanvas(ctx)
    })
}

// GetCanvas returns the UI element
func (cr *CanvasRenderer) GetCanvas() app.UI {
    return cr
}
