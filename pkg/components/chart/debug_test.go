//go:build dev
// pkg/components/chart/debug_test.go
package chart

import (
    "github.com/maxence-charriere/go-app/v10/pkg/app"
    "github.com/mmcnicol/go-app-component-library/pkg/storybook"
)

type DebugChart struct {
    app.Compo
    mounted bool
}

func (dc *DebugChart) OnMount(ctx app.Context) {
    dc.mounted = true
    
    ctx.Defer(func(ctx app.Context) {
        // Super simple JavaScript that should work
        jsCode := `
            console.log('Debug chart: Creating canvas');
            
            const container = document.querySelector('[data-debug-chart="true"]');
            if (!container) {
                console.error('Debug container not found');
                return;
            }
            
            // Create canvas
            const canvas = document.createElement('canvas');
            canvas.width = 800;
            canvas.height = 400;
            canvas.style.width = '100%';
            canvas.style.height = '400px';
            canvas.style.border = '3px solid green';
            canvas.style.background = '#f0f8ff';
            
            container.innerHTML = '';
            container.appendChild(canvas);
            
            // Draw on it
            const ctx = canvas.getContext('2d');
            ctx.fillStyle = '#228B22';
            ctx.font = 'bold 24px Arial';
            ctx.textAlign = 'center';
            ctx.fillText('DEBUG CHART - THIS SHOULD APPEAR', 400, 50);
            
            // Draw a rectangle
            ctx.fillStyle = '#FF6347';
            ctx.fillRect(100, 100, 200, 150);
            
            console.log('Debug chart drawn successfully');
        `
        
        app.Window().Call("eval", jsCode)
    })
}

func (dc *DebugChart) Render() app.UI {
    return app.Div().
        Attr("data-debug-chart", "true").
        Style("width", "100%").
        Style("height", "400px").
        Style("margin", "20px 0").
        Style("border", "2px dashed red")
}

func init() {
    storybook.Register("Chart", "Debug Chart", 
        nil,
        func(controls map[string]*storybook.Control) app.UI {
            return app.Div().Body(
                app.H3().Text("Debug Chart Test"),
                app.P().Text("This is a minimal test to see if canvas works"),
                &DebugChart{},
            )
        },
    )
}
