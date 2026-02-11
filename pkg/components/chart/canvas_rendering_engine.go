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
        fmt.Println("No chart spec to render")
        return
    }
    
    fmt.Printf("Setting up canvas for %s chart\n", cr.chartSpec.Type)
    
    // Create a VERY simple JavaScript that won't have syntax errors
    jsCode := fmt.Sprintf(`
        console.log('Setting up canvas: %s');
        
        try {
            const canvas = document.getElementById('%s');
            if (!canvas) {
                console.error('Canvas not found');
                return;
            }
            
            // Set dimensions
            canvas.width = 800;
            canvas.height = 400;
            
            const ctx = canvas.getContext('2d');
            if (!ctx) {
                console.error('No 2D context');
                return;
            }
            
            // Draw something simple to verify it works
            ctx.fillStyle = '#ffffff';
            ctx.fillRect(0, 0, canvas.width, canvas.height);
            
            ctx.fillStyle = '#4A90E2';
            ctx.fillRect(50, 50, 100, 200);
            
            ctx.fillStyle = '#000000';
            ctx.font = '16px Arial';
            ctx.fillText('Simple Test - Chart Type: %s', 50, 30);
            
            console.log('Simple test drawn');
        } catch(e) {
            console.error('Error in simple test:', e);
        }
    `, cr.canvasID, cr.canvasID, cr.chartSpec.Type)
    
    // Use Dispatch to execute JavaScript in the main thread
    ctx.Dispatch(func(ctx app.Context) {
        app.Window().Call("eval", jsCode)
    })
}

func (cr *CanvasRenderer) getDrawScript() string {
    data := cr.chartSpec.Data
    if len(data.Datasets) == 0 {
        return "console.log('No data to draw');"
    }
    
    // For now, just draw a simple test chart
    return `
        // Draw a simple test chart
        const width = canvas.width;
        const height = canvas.height;
        
        // Draw title
        ctx.fillStyle = '#333';
        ctx.font = 'bold 16px Arial';
        ctx.textAlign = 'center';
        ctx.fillText('Test Chart', width / 2, 30);
        
        // Draw some test bars
        const barCount = 6;
        const barWidth = 40;
        const barSpacing = 20;
        const totalWidth = (barCount * barWidth) + ((barCount - 1) * barSpacing);
        const startX = (width - totalWidth) / 2;
        const maxBarHeight = height - 100;
        
        // Test data
        const values = [30, 50, 80, 40, 60, 90];
        
        for (let i = 0; i < barCount; i++) {
            const barHeight = (values[i] / 100) * maxBarHeight;
            const x = startX + i * (barWidth + barSpacing);
            const y = height - 50 - barHeight;
            
            // Draw bar
            ctx.fillStyle = '#4A90E2';
            ctx.fillRect(x, y, barWidth, barHeight);
            
            // Draw bar border
            ctx.strokeStyle = '#2c6fb3';
            ctx.lineWidth = 2;
            ctx.strokeRect(x, y, barWidth, barHeight);
            
            // Draw value label
            ctx.fillStyle = '#333';
            ctx.font = '12px Arial';
            ctx.textAlign = 'center';
            ctx.fillText(values[i].toString(), x + barWidth/2, y - 10);
            
            // Draw x-axis label
            ctx.fillText('Bar ' + (i+1), x + barWidth/2, height - 30);
        }
        
        // Draw y-axis
        ctx.beginPath();
        ctx.moveTo(startX - 10, height - 50);
        ctx.lineTo(startX - 10, height - 50 - maxBarHeight);
        ctx.strokeStyle = '#333';
        ctx.lineWidth = 2;
        ctx.stroke();
        
        // Draw y-axis labels
        ctx.textAlign = 'right';
        for (let i = 0; i <= 5; i++) {
            const y = height - 50 - (i * maxBarHeight / 5);
            const value = Math.round((i * 100) / 5);
            ctx.fillText(value.toString(), startX - 15, y + 4);
            
            // Draw grid line
            ctx.beginPath();
            ctx.moveTo(startX, y);
            ctx.lineTo(startX + totalWidth, y);
            ctx.strokeStyle = '#e0e0e0';
            ctx.lineWidth = 1;
            ctx.stroke();
        }
    `
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
            color: '%s',
            borderColor: '%s',
            borderWidth: %d
        },`, dataset.Label, points, color, dataset.BorderColor, dataset.BorderWidth)
    }
    
    labels := ""
    for _, label := range data.Labels {
        labels += fmt.Sprintf("'%s',", label)
    }
    
    return fmt.Sprintf(`
        // Draw bar chart
        function drawBarChart(ctx, width, height) {
            console.log('Drawing bar chart, width:', width, 'height:', height);
            
            // Clear canvas
            ctx.clearRect(0, 0, width, height);
            
            // Draw white background
            ctx.fillStyle = '#ffffff';
            ctx.fillRect(0, 0, width, height);
            
            const datasets = [%s];
            const labels = [%s];
            const numDatasets = datasets.length;
            const numBars = labels.length;
            
            console.log('Datasets:', numDatasets, 'Bars per dataset:', datasets[0]?.data?.length);
            
            // Calculate scales
            let minY = Infinity;
            let maxY = -Infinity;
            
            datasets.forEach(dataset => {
                dataset.data.forEach(point => {
                    if (point.y < minY) minY = point.y;
                    if (point.y > maxY) maxY = point.y;
                });
            });
            
            console.log('Data range - min:', minY, 'max:', maxY);
            
            // Add padding
            const yPadding = (maxY - minY) * 0.1;
            minY = Math.max(0, minY - yPadding);  // Ensure minY doesn't go below 0 for bar charts
            maxY = maxY + yPadding;
            
            // Ensure minY is 0 if BeginAtZero is true
            if (%t && minY < 0) {
                minY = 0;
            }
            
            // Margins
            const margin = { top: 40, right: 40, bottom: 60, left: 60 };
            const plotWidth = width - margin.left - margin.right;
            const plotHeight = height - margin.top - margin.bottom;
            
            console.log('Margins - left:', margin.left, 'right:', margin.right, 'plotWidth:', plotWidth);
            
            // X scale - fix bar positioning
            const barGroupWidth = plotWidth / numBars;
            const barSpacing = barGroupWidth * 0.1; // 10%% spacing between bar groups
            const barWidth = (barGroupWidth - barSpacing) / numDatasets;
            
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
            
            // Draw bars - FIXED POSITIONING
            datasets.forEach((dataset, datasetIdx) => {
                ctx.fillStyle = dataset.color;
                ctx.strokeStyle = dataset.borderColor || '#333';
                ctx.lineWidth = dataset.borderWidth || 1;
                
                dataset.data.forEach((point, pointIdx) => {
                    // Calculate bar position - FIXED
                    const groupCenterX = margin.left + (pointIdx * barGroupWidth) + (barGroupWidth / 2);
                    const barX = groupCenterX - (numDatasets * barWidth / 2) + (datasetIdx * barWidth);
                    
                    const yPos = yScale(point.y);
                    const zeroY = yScale(0);
                    const barHeight = zeroY - yPos;
                    
                    // Only draw if barHeight is positive
                    if (barHeight > 0) {
                        // Draw bar
                        ctx.fillRect(
                            barX,
                            yPos,
                            barWidth,
                            barHeight
                        );
                        
                        // Draw border
                        ctx.strokeRect(
                            barX,
                            yPos,
                            barWidth,
                            barHeight
                        );
                        
                        // Draw value label on top of bar
                        ctx.fillStyle = '#000000';
                        ctx.font = '10px Arial';
                        ctx.textAlign = 'center';
                        ctx.textBaseline = 'bottom';
                        ctx.fillText(
                            Math.round(point.y).toString(),
                            barX + barWidth / 2,
                            yPos - 5
                        );
                        ctx.fillStyle = dataset.color; // Reset fill color
                    }
                });
            });
            
            // Draw labels on X axis - FIXED
            ctx.fillStyle = '#000000';
            ctx.textAlign = 'center';
            ctx.textBaseline = 'top';
            labels.forEach((label, i) => {
                if (i < numBars) {
                    const x = margin.left + (i * barGroupWidth) + (barGroupWidth / 2);
                    ctx.fillText(label, x, height - margin.bottom + 10);
                }
            });
            
            // Draw Y axis labels
            ctx.textAlign = 'right';
            ctx.textBaseline = 'middle';
            for (let i = 0; i <= 5; i++) {
                const yValue = minY + ((maxY - minY) * i / 5);
                const yPos = yScale(yValue);
                ctx.fillText(Math.round(yValue).toString(), margin.left - 10, yPos);
            }
            
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
                
                // Draw border around box
                ctx.strokeStyle = '#333';
                ctx.lineWidth = 1;
                ctx.strokeRect(legendX, y, boxSize, boxSize);
                
                // Draw label
                ctx.fillStyle = '#000000';
                ctx.textAlign = 'left';
                ctx.textBaseline = 'middle';
                ctx.font = '12px Arial';
                ctx.fillText(dataset.label, legendX + boxSize + 10, y + boxSize/2);
            });
            
            // Draw title if exists
            if ('%s') {
                ctx.fillStyle = '#000000';
                ctx.font = 'bold 16px Arial';
                ctx.textAlign = 'center';
                ctx.textBaseline = 'top';
                ctx.fillText('%s', width / 2, 10);
            }
        }
        
        // Call draw function
        drawBarChart(ctx, width, height);
    `, datasets, labels, cr.chartSpec.Options.Scales.Y.BeginAtZero, 
       cr.chartSpec.Options.Plugins.Title.Display, 
       cr.chartSpec.Options.Plugins.Title.Text)
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

func (cr *CanvasRenderer) getPieChartScript() string {
    data := cr.chartSpec.Data
    if len(data.Datasets) == 0 || len(data.Datasets[0].Data) == 0 {
        return "console.log('No data for pie chart');"
    }
    
    dataset := data.Datasets[0]
    colors := dataset.BackgroundColor
    points := ""
    
    for i, point := range dataset.Data {
        color := "#4A90E2"
        if i < len(colors) {
            color = colors[i]
        }
        label := point.Label
        if label == "" {
            label = fmt.Sprintf("Item %d", i+1)
        }
        points += fmt.Sprintf(`{
            label: '%s',
            value: %f,
            color: '%s'
        },`, label, point.Y, color)
    }
    
    return fmt.Sprintf(`
        // Draw pie chart
        function drawPieChart(ctx, width, height) {
            // Clear canvas
            ctx.clearRect(0, 0, width, height);
            
            // Draw white background
            ctx.fillStyle = '#ffffff';
            ctx.fillRect(0, 0, width, height);
            
            const data = [%s];
            
            // Calculate total
            let total = 0;
            data.forEach(item => {
                total += item.value;
            });
            
            // Center and radius
            const centerX = width / 2;
            const centerY = height / 2;
            const radius = Math.min(width, height) * 0.35;
            
            // Starting angle
            let startAngle = 0;
            
            // Draw each segment
            data.forEach((item, i) => {
                const sliceAngle = (item.value / total) * 2 * Math.PI;
                
                // Draw segment
                ctx.beginPath();
                ctx.moveTo(centerX, centerY);
                ctx.arc(centerX, centerY, radius, startAngle, startAngle + sliceAngle);
                ctx.closePath();
                
                ctx.fillStyle = item.color;
                ctx.fill();
                ctx.strokeStyle = '#ffffff';
                ctx.lineWidth = 2;
                ctx.stroke();
                
                // Draw label
                const midAngle = startAngle + sliceAngle / 2;
                const labelRadius = radius * 0.7;
                const labelX = centerX + Math.cos(midAngle) * labelRadius;
                const labelY = centerY + Math.sin(midAngle) * labelRadius;
                
                ctx.fillStyle = '#000000';
                ctx.font = '12px Arial';
                ctx.textAlign = 'center';
                ctx.textBaseline = 'middle';
                
                const percentage = (item.value / total * 100).toFixed(1);
                ctx.fillText(item.label, labelX, labelY);
                
                // Update start angle for next segment
                startAngle += sliceAngle;
            });
        }
        
        // Call draw function
        drawPieChart(ctx, width, height);
    `, points)
}
