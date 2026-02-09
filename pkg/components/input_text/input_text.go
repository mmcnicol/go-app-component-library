// pkg/components/input_text/input_text.go
package input_text

import (
	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

// InputText indicates service status
type InputText struct {
	app.Compo
	value     string
}

func (t *InputText) Render() app.UI {
	if app.IsClient {
		app.Log("InputText Render()")
	}
	return app.Input().
		Type("text").
		Value(t.value).
		//Placeholder("What is your name?").
		//AutoFocus(true).
		OnChange(t.ValueTo(&t.value)),
	)
}
