//go:build dev
// Create a new file: pkg/components/chart/simple_test.go
package chart

import (
    "github.com/maxence-charriere/go-app/v10/pkg/app"
    "github.com/mmcnicol/go-app-component-library/pkg/storybook"
)

type SimpleTestChart struct {
    app.Compo
}

func (stc *SimpleTestChart) OnMount(ctx app.Context) {
    // Use the simplest possible JavaScript
    jsCode := `(function() {
        console.log('SimpleTestChart: Starting');
        
        // Find container
        const container = document.querySelector('[data-simple-test="true"]');
        if (!container) {
            console.error('Container not found');
            return;
        }
        
        // Create canvas
        const canvas = document.createElement('canvas');
        canvas.width = 800;
        canvas.height = 400;
        canvas.style.width = '100%';
        canvas.style.height = '400px';
        canvas.style.border = '3px solid blue';
        canvas.style.background = '#f0f8ff';
        
        // Clear container and add canvas
        container.innerHTML = '';
        container.appendChild(canvas);
        
        // Draw
        const ctx = canvas.getContext('2d');
        
        // Clear
        ctx.clearRect(0, 0, canvas.width, canvas.height);
        
        // Draw background
        ctx.fillStyle = '#ffffff';
        ctx.fillRect(0, 0, canvas.width, canvas.height);
        
        // Draw title
        ctx.fillStyle = '#0000ff';
        ctx.font = 'bold 24px Arial';
        ctx.textAlign = 'center';
        ctx.fillText('SIMPLE TEST CHART - SHOULD WORK', 400, 50);
        
        // Draw a rectangle
        ctx.fillStyle = '#ff0000';
        ctx.fillRect(100, 100, 200, 150);
        
        // Draw text in rectangle
        ctx.fillStyle = '#ffffff';
        ctx.font = '20px Arial';
        ctx.fillText('SUCCESS!', 200, 180);
        
        console.log('SimpleTestChart: Finished successfully');
    })();`
    
    app.Window().Call("eval", jsCode)
}

func (stc *SimpleTestChart) Render() app.UI {
    return app.Div().
        Attr("data-simple-test", "true").
        Style("width", "100%").
        Style("height", "400px").
        Style("margin", "20px 0")
}

func init() {
    storybook.Register("Chart", "Simple Test Chart", 
        nil,
        func(controls map[string]*storybook.Control) app.UI {
            return app.Div().Body(
                app.H3().Text("Simple Test Chart"),
                app.P().Text("This is the simplest possible test"),
                &SimpleTestChart{},
            )
        },
    )
}
