//go:build dev
// Create a test file: pkg/components/chart/test_chart.go
package chart

import (
    "github.com/maxence-charriere/go-app/v10/pkg/app"
    "github.com/mmcnicol/go-app-component-library/pkg/storybook"
)

func init() {
    storybook.Register("Chart", "Canvas Test", 
        nil,
        func(controls map[string]*storybook.Control) app.UI {
            return app.Div().Body(
                app.H2().Text("Canvas Element Test"),
                app.Canvas().
                    ID("test-canvas").
                    Style("width", "100%").
                    Style("height", "300px").
                    Style("border", "2px solid red").
                    Style("background", "#f0f0f0"),
                app.Button().
                    Text("Draw on Canvas").
                    OnClick(func(ctx app.Context, e app.Event) {
                        app.Window().Call("eval", `
                            const canvas = document.getElementById('test-canvas');
                            if (canvas) {
                                canvas.width = canvas.clientWidth;
                                canvas.height = canvas.clientHeight;
                                const ctx = canvas.getContext('2d');
                                ctx.fillStyle = '#4A90E2';
                                ctx.fillRect(50, 50, 100, 100);
                                ctx.fillStyle = '#000';
                                ctx.font = '16px Arial';
                                ctx.fillText('Canvas is working!', 60, 180);
                            }
                        `)
                    }),
            )
        },
    )
}
