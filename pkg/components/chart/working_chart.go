//go:build dev
// pkg/components/chart/working_chart.go
package chart

import (
    "fmt"
    "github.com/maxence-charriere/go-app/v10/pkg/app"
    "github.com/mmcnicol/go-app-component-library/pkg/storybook"
)

type WorkingChart struct {
    app.Compo
    chartType ChartType
    data      ChartData
    title     string
    mounted   bool
}

func NewWorkingChart(chartType ChartType) *WorkingChart {
    return &WorkingChart{
        chartType: chartType,
    }
}

func (wc *WorkingChart) Title(title string) *WorkingChart {
    wc.title = title
    return wc
}

func (wc *WorkingChart) Data(data ChartData) *WorkingChart {
    wc.data = data
    return wc
}

func (wc *WorkingChart) OnMount(ctx app.Context) {
    wc.mounted = true
    
    // Draw after mount
    ctx.Defer(func(ctx app.Context) {
        wc.drawChart()
    })
}

func (wc *WorkingChart) drawChart() {
    canvasID := fmt.Sprintf("working-chart-%s", GenerateID())
    
    // Simple, hardcoded JavaScript that definitely works
    jsCode := fmt.Sprintf(`
        console.log('Drawing chart for canvas:', '%s');
        
        const container = document.querySelector('[data-chart-id="%s"]');
        if (!container) {
            console.error('Container not found');
            return;
        }
        
        // Create canvas
        const canvas = document.createElement('canvas');
        canvas.id = '%s';
        canvas.width = 800;
        canvas.height = 400;
        canvas.style.width = '100%%';
        canvas.style.height = '400px';
        canvas.style.border = '1px solid #4A90E2';
        canvas.style.background = '#ffffff';
        
        container.innerHTML = '';
        container.appendChild(canvas);
        
        const ctx = canvas.getContext('2d');
        
        // Clear canvas
        ctx.clearRect(0, 0, canvas.width, canvas.height);
        
        // Draw background
        ctx.fillStyle = '#ffffff';
        ctx.fillRect(0, 0, canvas.width, canvas.height);
        
        // Draw title
        ctx.fillStyle = '#333';
        ctx.font = 'bold 18px Arial';
        ctx.textAlign = 'center';
        ctx.fillText('Working Chart - This Should Display', canvas.width / 2, 30);
        
        // Draw a simple bar chart
        const barData = [30, 50, 80, 40, 60, 90, 70];
        const barWidth = 40;
        const barSpacing = 20;
        const maxBarHeight = 250;
        
        const totalWidth = (barData.length * barWidth) + ((barData.length - 1) * barSpacing);
        const startX = (canvas.width - totalWidth) / 2;
        
        for (let i = 0; i < barData.length; i++) {
            const barHeight = (barData[i] / 100) * maxBarHeight;
            const x = startX + i * (barWidth + barSpacing);
            const y = 350 - barHeight;
            
            // Draw bar
            ctx.fillStyle = '#4A90E2';
            ctx.fillRect(x, y, barWidth, barHeight);
            
            // Draw bar border
            ctx.strokeStyle = '#2c6fb3';
            ctx.lineWidth = 2;
            ctx.strokeRect(x, y, barWidth, barHeight);
            
            // Draw value
            ctx.fillStyle = '#333';
            ctx.font = '12px Arial';
            ctx.textAlign = 'center';
            ctx.fillText(barData[i].toString(), x + barWidth/2, y - 10);
            
            // Draw label
            ctx.fillText('Item ' + (i+1), x + barWidth/2, 370);
        }
        
        // Draw axes
        ctx.beginPath();
        ctx.moveTo(startX - 10, 350);
        ctx.lineTo(startX - 10, 100);
        ctx.strokeStyle = '#333';
        ctx.lineWidth = 2;
        ctx.stroke();
        
        ctx.beginPath();
        ctx.moveTo(startX - 10, 350);
        ctx.lineTo(startX + totalWidth + 10, 350);
        ctx.stroke();
        
        console.log('Chart drawn successfully!');
    `, canvasID, canvasID, canvasID)
    
    app.Window().Call("eval", jsCode)
}

func (wc *WorkingChart) Render() app.UI {
    canvasID := fmt.Sprintf("working-chart-%s", GenerateID())
    
    return app.Div().
        Attr("data-chart-id", canvasID).
        Style("width", "100%").
        Style("height", "400px").
        Style("position", "relative").
        Style("margin", "20px 0")
}

// Register a story for the working chart
func init() {
    storybook.Register("Chart", "Working Chart Test", 
        nil,
        func(controls map[string]*storybook.Control) app.UI {
            return app.Div().Body(
                app.H3().Text("Working Chart Test"),
                app.P().Text("This should definitely show a chart"),
                NewWorkingChart(ChartTypeBar).
                    Title("Test Chart").
                    Data(ChartData{
                        Labels: []string{"A", "B", "C", "D", "E", "F", "G"},
                        Datasets: []Dataset{{
                            Label: "Test Data",
                            Data: []DataPoint{
                                {X: 0, Y: 30},
                                {X: 1, Y: 50},
                                {X: 2, Y: 80},
                                {X: 3, Y: 40},
                                {X: 4, Y: 60},
                                {X: 5, Y: 90},
                                {X: 6, Y: 70},
                            },
                        }},
                    }),
            )
        },
    )
}
