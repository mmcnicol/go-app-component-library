// pkg/components/built_in/canvas.go
package built_in

import (
	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

type BuiltInCanvas struct {
	app.Compo
}

func (c *BuiltInCanvas) Render() app.UI {
	return app.Canvas().
		ID("builtin-canvas").
		Width(200).
		Height(200).
		Style("border", "1px solid #ccc")
}

// OnMount is called when the component is inserted into the DOM
func (c *BuiltInCanvas) OnMount(ctx app.Context) {
	// We must draw after the element is mounted
	c.drawCircle()
}

func (c *BuiltInCanvas) drawCircle() {
	// Get the canvas element from the DOM
	canvas := app.Window().GetElementByID("builtin-canvas")
	if !canvas.Truthy() {
		return
	}

	// Get the 2D context
	ctx := canvas.Call("getContext", "2d")
	
	// Clear Canvas
	ctx.Call("clearRect", 0, 0, 200, 200)

	// Draw Green Circle
	ctx.Call("beginPath")
	// arc(x, y, radius, startAngle, endAngle)
	ctx.Call("arc", 100, 100, 50, 0, 2*3.14159)
	ctx.Set("fillStyle", "green")
	ctx.Call("fill")
	ctx.Set("strokeStyle", "darkgreen")
	ctx.Set("lineWidth", 5)
	ctx.Call("stroke")
}
