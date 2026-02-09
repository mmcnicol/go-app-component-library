// pkg/components/input_text/input_text.go
package input_text

import (
	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

// InputText indicates service status
type InputText struct {
	app.Compo
	Value       string
	Placeholder string
	Disabled    bool
	OnInput     func(ctx app.Context, val string)
}

func (t *InputText) Render() app.UI {
	if app.IsClient {
		app.Log("InputText Render()")
	}

	containerClass := ""
	if t.Disabled {
        containerClass += "inputText-container-input-disabled"
    }

	return app.Div().Class("inputText-container").Body(

		/*
		app.If(t.Label != "", func() app.UI {
			return app.Label().
				Class("inputText-container-label").
				Text(t.Label)
		}),
		*/

		app.Input().
			Class("inputText-container-input").
			Class(containerClass).
			Type("text").
			Value(t.Value).
			Disabled(t.Disabled).
			Placeholder(t.Placeholder).
			//AutoFocus(true).
			//OnChange(t.ValueTo(&t.Value)),
			OnInput(func(ctx app.Context, e app.Event) {
				val := ctx.JSSrc().Get("value").String()
				t.Value = val
				if t.OnInput != nil {
					t.OnInput(ctx, val)
				}
				//ctx.Update()
			},
		),
	)
}
