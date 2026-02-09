// pkg/components/phase_banner/toggle_switch.go
package toggle_switch

import (
	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

// ToggleSwitch defines the UI component
type ToggleSwitch struct {
	app.Compo
	IsOn  bool
	Label string
}

// OnClick handles the toggle logic
func (t *ToggleSwitch) OnClick(ctx app.Context, e app.Event) {
	t.IsOn = !t.IsOn
	//t.Update()
}

func (t *ToggleSwitch) Render() app.UI {
	return app.Div().
		Class("toggle-container").
		Class(app.If(t.IsOn, "active").Else("")). // Dynamic class binding
		OnClick(t.OnClick).
		Body(
			app.Div().Class("switch").Body(
				app.Span().Class("switch-inner"),
			),
			app.If(t.Label != "", app.Span().Text(t.Label)),
		)
}
