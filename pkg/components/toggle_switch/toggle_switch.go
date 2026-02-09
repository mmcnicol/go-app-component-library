// pkg/components/phase_banner/toggle_switch.go
package toggle_switch

import (
	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

// ToggleSwitch defines the UI component
type ToggleSwitch struct {
	app.Compo
	IsOn         bool
	Label        string
    shouldRender bool
}

func (t *ToggleSwitch) OnMount(ctx app.Context) {
	if app.IsClient {
		app.Log("ToggleSwitch OnMount()")
	}
}

// OnClick handles the toggle logic
func (t *ToggleSwitch) OnClick(ctx app.Context, e app.Event) {
	//e.StopPropagation() // Prevents the event from reaching the Storybook Shell

	t.IsOn = !t.IsOn

	// Only logs if the app was built with "-tags dev"
	if app.IsClient {
		// Simple log
		app.Log("ToggleSwitch clicked!")
		// Formatted log
		app.Logf("ToggleSwitch state is now: %v", t.IsOn)
	}

    t.shouldRender = true
	//ctx.Update()
}

func (t *ToggleSwitch) Render() app.UI {
	if app.IsClient {
		app.Log("ToggleSwitch Render()")
	}

    t.shouldRender = false

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

func (t *ToggleSwitch) Update(ctx app.Context) bool {
	if app.IsClient {
		app.Log("ToggleSwitch Update()")
	}
	return t.shouldRender
}
