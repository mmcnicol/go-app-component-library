// pkg/components/chart/simple_chart.go
package chart

import (
    "fmt"
    "github.com/maxence-charriere/go-app/v10/pkg/app"
)

type SimpleChart struct {
    app.Compo
    chartType ChartType
    data      ChartData
    title     string
    width     int
    height    int
    mounted   bool
}

func NewSimpleChart(chartType ChartType) *SimpleChart {
    return &SimpleChart{
        chartType: chartType,
        width:     800,
        height:    400,
    }
}

func (sc *SimpleChart) Title(title string) *SimpleChart {
    sc.title = title
    return sc
}

func (sc *SimpleChart) Data(data ChartData) *SimpleChart {
    sc.data = data
    return sc
}

func (sc *SimpleChart) Width(width int) *SimpleChart {
    sc.width = width
    return sc
}

func (sc *SimpleChart) Height(height int) *SimpleChart {
    sc.height = height
    return sc
}

func (sc *SimpleChart) OnMount(ctx app.Context) {
    sc.mounted = true
    
    // Draw the chart after mount
    ctx.Defer(func(ctx app.Context) {
        sc.drawChart()
    })
}

func (sc *SimpleChart) drawChart() {
    canvasID := fmt.Sprintf("simple-chart-%s", GenerateID())
    
    jsCode := ""
    
    if sc.chartType == ChartTypeBar {
        jsCode = sc.getBarChartJS(canvasID)
    } else {
        // Default to bar chart for now
        jsCode = sc.getBarChartJS(canvasID)
    }
    
    app.Window().Call("eval", fmt.Sprintf(`
        (function() {
            %s
        })();
    `, jsCode))
}

func (sc *SimpleChart) getBarChartJS(canvasID string) string {
    // Prepare data
    labelsJS := "["
    for _, label := range sc.data.Labels {
        labelsJS += fmt.Sprintf("'%s',", label)
    }
    if len(sc.data.Labels) > 0 {
        labelsJS = labelsJS[:len(labelsJS)-1] // Remove trailing comma
    }
    labelsJS += "]"
    
    valuesJS := "["
    if len(sc.data.Datasets) > 0 && len(sc.data.Datasets[0].Data) > 0 {
        for _, point := range sc.data.Datasets[0].Data {
            valuesJS += fmt.Sprintf("%f,", point.Y)
        }
        valuesJS = valuesJS[:len(valuesJS)-1] // Remove trailing comma
    }
    valuesJS += "]"
    
    // Fix: Use proper string escaping and avoid % signs
    return fmt.Sprintf(`
        (function() {
            const canvas = document.createElement('canvas');
            canvas.id = '%s';
            canvas.width = %d;
            canvas.height = %d;
            canvas.style.width = '100%%';
            canvas.style.height = '100%%';
            canvas.style.display = 'block';
            
            // Find and replace the chart container
            const container = document.querySelector('[data-chart-container="%s"]');
            if (container) {
                container.innerHTML = '';
                container.appendChild(canvas);
                
                const ctx = canvas.getContext('2d');
                
                // Clear canvas
                ctx.clearRect(0, 0, canvas.width, canvas.height);
                
                // Draw background
                ctx.fillStyle = '#ffffff';
                ctx.fillRect(0, 0, canvas.width, canvas.height);
                
                // Draw title
                if ('%s' !== '') {
                    ctx.fillStyle = '#333';
                    ctx.font = 'bold 18px Arial';
                    ctx.textAlign = 'center';
                    ctx.fillText('%s', canvas.width / 2, 30);
                }
                
                // Draw bars
                const labels = %s;
                const values = %s;
                const barCount = labels.length;
                
                if (barCount > 0 && values.length > 0) {
                    const margin = { top: 60, right: 40, bottom: 60, left: 60 };
                    const plotWidth = canvas.width - margin.left - margin.right;
                    const plotHeight = canvas.height - margin.top - margin.bottom;
                    const barWidth = plotWidth / barCount * 0.7;
                    
                    // Find max value
                    let maxValue = Math.max(...values);
                    maxValue = maxValue * 1.1; // Add 10%% padding
                    
                    // Draw bars
                    for (let i = 0; i < barCount; i++) {
                        const barHeight = (values[i] / maxValue) * plotHeight;
                        const x = margin.left + (i * plotWidth / barCount) + (plotWidth / barCount - barWidth) / 2;
                        const y = canvas.height - margin.bottom - barHeight;
                        
                        // Draw bar
                        ctx.fillStyle = '#4A90E2';
                        ctx.fillRect(x, y, barWidth, barHeight);
                        
                        // Draw bar border
                        ctx.strokeStyle = '#2c6fb3';
                        ctx.lineWidth = 1;
                        ctx.strokeRect(x, y, barWidth, barHeight);
                        
                        // Draw value on top
                        ctx.fillStyle = '#333';
                        ctx.font = '12px Arial';
                        ctx.textAlign = 'center';
                        ctx.fillText(values[i].toFixed(0), x + barWidth/2, y - 5);
                        
                        // Draw label
                        ctx.fillText(labels[i], x + barWidth/2, canvas.height - margin.bottom + 20);
                    }
                    
                    // Draw Y axis
                    ctx.beginPath();
                    ctx.moveTo(margin.left, margin.top);
                    ctx.lineTo(margin.left, canvas.height - margin.bottom);
                    ctx.strokeStyle = '#333';
                    ctx.lineWidth = 1;
                    ctx.stroke();
                    
                    // Draw grid lines
                    ctx.strokeStyle = '#e0e0e0';
                    ctx.lineWidth = 0.5;
                    for (let i = 0; i <= 5; i++) {
                        const y = canvas.height - margin.bottom - (i * plotHeight / 5);
                        ctx.beginPath();
                        ctx.moveTo(margin.left, y);
                        ctx.lineTo(canvas.width - margin.right, y);
                        ctx.stroke();
                        
                        // Draw Y axis labels
                        const value = (i * maxValue / 5).toFixed(0);
                        ctx.fillStyle = '#666';
                        ctx.font = '12px Arial';
                        ctx.textAlign = 'right';
                        ctx.fillText(value, margin.left - 10, y + 4);
                    }
                } else {
                    // Draw "no data" message
                    ctx.fillStyle = '#999';
                    ctx.font = '16px Arial';
                    ctx.textAlign = 'center';
                    ctx.fillText('No data available', canvas.width / 2, canvas.height / 2);
                }
            }
        })();
    `, canvasID, sc.width, sc.height, canvasID, sc.title, sc.title, labelsJS, valuesJS)
}

func (sc *SimpleChart) Render() app.UI {
    canvasID := fmt.Sprintf("simple-chart-%s", GenerateID())
    
    return app.Div().
        Attr("data-chart-container", canvasID).
        Style("width", "100%").
        Style("height", fmt.Sprintf("%dpx", sc.height)).
        Style("position", "relative").
        Style("border", "1px solid #e0e0e0").
        Style("border-radius", "4px").
        Style("background", "#ffffff").
        Style("margin", "20px 0")
}
