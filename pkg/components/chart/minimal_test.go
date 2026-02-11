//go:build dev
// Create a new test file: pkg/components/chart/minimal_test.go
package chart

import (
    "github.com/maxence-charriere/go-app/v10/pkg/app"
    "github.com/mmcnicol/go-app-component-library/pkg/storybook"
)

func init() {
    storybook.Register("Chart", "Minimal Working Chart", 
        nil,
        func(controls map[string]*storybook.Control) app.UI {
            return app.Div().Body(
                app.H3().Text("Minimal Working Chart Test"),
                app.P().Text("This uses pure JavaScript to draw on canvas"),
                app.Canvas().
                    ID("minimal-canvas").
                    Style("width", "100%").
                    Style("height", "300px").
                    Style("border", "2px solid blue").
                    Style("background", "#f9f9f9").
                    OnMount(func(ctx app.Context) {
                        // Draw directly when mounted
                        jsCode := `
                            const canvas = document.getElementById('minimal-canvas');
                            canvas.width = canvas.clientWidth;
                            canvas.height = canvas.clientHeight;
                            const ctx = canvas.getContext('2d');
                            
                            // Clear
                            ctx.clearRect(0, 0, canvas.width, canvas.height);
                            
                            // Draw background
                            ctx.fillStyle = '#ffffff';
                            ctx.fillRect(0, 0, canvas.width, canvas.height);
                            
                            // Draw a simple bar chart
                            const barData = [30, 50, 80, 40, 60];
                            const barWidth = 50;
                            const barSpacing = 20;
                            const maxBarHeight = 150;
                            
                            for (let i = 0; i < barData.length; i++) {
                                const barHeight = (barData[i] / 100) * maxBarHeight;
                                const x = 50 + i * (barWidth + barSpacing);
                                const y = 200 - barHeight;
                                
                                // Draw bar
                                ctx.fillStyle = '#4A90E2';
                                ctx.fillRect(x, y, barWidth, barHeight);
                                
                                // Draw label
                                ctx.fillStyle = '#333';
                                ctx.font = '14px Arial';
                                ctx.textAlign = 'center';
                                ctx.fillText('Bar ' + (i+1), x + barWidth/2, 220);
                                ctx.fillText(barData[i].toString(), x + barWidth/2, y - 10);
                            }
                            
                            // Draw title
                            ctx.fillStyle = '#000';
                            ctx.font = 'bold 18px Arial';
                            ctx.textAlign = 'center';
                            ctx.fillText('Minimal Working Chart', canvas.width/2, 30);
                            
                            console.log('Minimal chart drawn successfully');
                        `
                        
                        app.Window().Call("eval", jsCode)
                    }),
            )
        },
    )
}
