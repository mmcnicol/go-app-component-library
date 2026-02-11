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
    
    // Use JavaScript to draw the chart
    jsCode := fmt.Sprintf(`
        try {
            const canvas = document.getElementById('%s');
            if (!canvas) {
                console.error('Canvas element not found: %s');
                return;
            }
            
            // Ensure canvas parent has dimensions
            const container = canvas.parentElement;
            if (container) {
                if (container.clientWidth === 0) {
                    container.style.width = '100%%';
                    container.style.minHeight = '400px';
                }
            }
            
            // Force canvas dimensions
            canvas.width = canvas.clientWidth || 800;
            canvas.height = canvas.clientHeight || 400;
            
            // Get context and draw
            const ctx = canvas.getContext('2d');
            if (!ctx) {
                console.error('Could not get 2D context');
                return;
            }
            
            // Clear and draw background
            ctx.clearRect(0, 0, canvas.width, canvas.height);
            ctx.fillStyle = '#ffffff';
            ctx.fillRect(0, 0, canvas.width, canvas.height);
            
            // Draw a test rectangle to verify canvas is working
            ctx.fillStyle = '#f0f0f0';
            ctx.fillRect(10, 10, canvas.width - 20, canvas.height - 20);
            ctx.strokeStyle = '#ccc';
            ctx.strokeRect(10, 10, canvas.width - 20, canvas.height - 20);
            
            // Draw chart type text
            ctx.fillStyle = '#333';
            ctx.font = '16px Arial';
            ctx.textAlign = 'center';
            ctx.fillText('%s Chart - Canvas Ready', canvas.width / 2, 30);
            
            // Now draw the actual chart
            %s
            
        } catch (error) {
            console.error('Error in chart rendering:', error);
        }
    `, cr.canvasID, cr.canvasID, cr.chartSpec.Type, cr.getDrawScript())
    
    app.Window().Call("eval", jsCode)
}

// pkg/components/chart/canvas_rendering_engine.go
// Update getDrawScript to ensure something is drawn:

func (cr *CanvasRenderer) getDrawScript() string {
    // First, let's just draw a simple test pattern
    testScript := `
        // Simple test drawing
        ctx.fillStyle = '#4A90E2';
        
        // Draw some shapes based on chart type
        const width = canvas.width;
        const height = canvas.height;
        
        if ('%s' === 'bar') {
            // Draw test bars
            for (let i = 0; i < 5; i++) {
                const barWidth = 40;
                const barHeight = 50 + Math.random() * 100;
                const x = 100 + i * 80;
                const y = height - 100 - barHeight;
                
                ctx.fillRect(x, y, barWidth, barHeight);
                ctx.strokeStyle = '#333';
                ctx.strokeRect(x, y, barWidth, barHeight);
            }
            ctx.fillStyle = '#000';
            ctx.font = '14px Arial';
            ctx.fillText('Bar Chart Test', width/2, 50);
        } 
        else if ('%s' === 'line') {
            // Draw test line
            ctx.beginPath();
            ctx.moveTo(50, height - 100);
            for (let i = 1; i < 6; i++) {
                ctx.lineTo(50 + i * 80, height - 100 - Math.sin(i) * 50);
            }
            ctx.strokeStyle = '#FF6384';
            ctx.lineWidth = 3;
            ctx.stroke();
            ctx.fillStyle = '#000';
            ctx.font = '14px Arial';
            ctx.fillText('Line Chart Test', width/2, 50);
        }
        else if ('%s' === 'pie') {
            // Draw test pie
            const centerX = width / 2;
            const centerY = height / 2;
            const radius = Math.min(width, height) * 0.3;
            
            const colors = ['#FF6384', '#36A2EB', '#FFCE56', '#4BC0C0'];
            let startAngle = 0;
            
            for (let i = 0; i < 4; i++) {
                const sliceAngle = Math.PI * 2 / 4;
                ctx.beginPath();
                ctx.moveTo(centerX, centerY);
                ctx.arc(centerX, centerY, radius, startAngle, startAngle + sliceAngle);
                ctx.closePath();
                ctx.fillStyle = colors[i];
                ctx.fill();
                startAngle += sliceAngle;
            }
            ctx.fillStyle = '#000';
            ctx.font = '14px Arial';
            ctx.fillText('Pie Chart Test', width/2, 50);
        }
        else {
            // Default drawing
            ctx.fillStyle = '#333';
            ctx.font = '14px Arial';
            ctx.fillText('%s Chart - Drawing not implemented', width/2, 50);
        }
    `
    
    return fmt.Sprintf(testScript, 
        string(cr.chartSpec.Type), 
        string(cr.chartSpec.Type), 
        string(cr.chartSpec.Type),
        string(cr.chartSpec.Type))
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
