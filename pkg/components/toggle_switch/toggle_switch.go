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
	ctx.Update()
}

func (t *ToggleSwitch) Render() app.UI {
    activeClass := ""
    if t.IsOn {
        activeClass = "active"
    }

    return app.Div().
        Class("toggle-container").
        Class(activeClass).
        OnClick(func(ctx app.Context, e app.Event) {
            t.OnClick(ctx, e)
        }).
        Body(
            app.Div().Class("switch").Body(
                app.Span().Class("switch-inner"),
            ),
            app.If(t.Label != "", func() app.UI {
                return app.Span().Text(t.Label)
            }),
        )
}
