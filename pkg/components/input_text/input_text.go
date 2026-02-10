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
    // Determine the class for the input element itself
    inputClass := "inputText-container-input"
    if t.Disabled {
        inputClass += " inputText-container-input-disabled"
    }

    return app.Div().Class("inputText-container").Body(
        app.Input().
            Class(inputClass). // Applied to the input
            Type("text").
            Value(t.Value).
            Disabled(t.Disabled).
            Placeholder(t.Placeholder).
            OnInput(func(ctx app.Context, e app.Event) {
                val := ctx.JSSrc().Get("value").String()
                t.Value = val
                if t.OnInput != nil {
                    t.OnInput(ctx, val)
                }
                ctx.Update() // Important for real-time feedback
            }),
    )
}
