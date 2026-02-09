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
	Disabled     bool
    //shouldRender bool
	OnClick      func(ctx app.Context, val bool)
}

func (t *ToggleSwitch) OnMount(ctx app.Context) {
	if app.IsClient {
		app.Log("ToggleSwitch OnMount()")
	}
}

func (t *ToggleSwitch) Render() app.UI {
	if app.IsClient {
		app.Log("ToggleSwitch Render()")
	}

    //t.shouldRender = false

    containerClass := ""
    if t.IsOn {
        containerClass = "toggleSwitch-container-active"
    }
	if t.Disabled {
        containerClass += " toggleSwitch-container-disabled"
    }

    return app.Div().
        Class("toggleSwitch-container").
        Class(containerClass).
        OnClick(func(ctx app.Context, e app.Event) {
            t.onClick(ctx, e)
			if t.OnClick != nil {
				t.OnClick(ctx, t.IsOn)
			}
        }).
        Body(
            app.Div().Class("toggleSwitch-container-switch").Body(
                app.Span().Class("toggleSwitch-container-switch-inner"),
            ),
            app.If(t.Label != "", func() app.UI {
                return app.Div().Class("toggleSwitch-container-label").Body(
					app.Text(t.Label),
				)
            }),
        )
}

// OnClick handles the toggle logic
func (t *ToggleSwitch) onClick(ctx app.Context, e app.Event) {
	e.PreventDefault() // Prevents the event from reaching the Storybook Shell
	
	// Only logs if the app was built with "-tags dev"
	if app.IsClient {
		// Simple log
		app.Log("ToggleSwitch clicked!")
	}

	// 1. Guard clause: Stop execution if disabled
    if t.Disabled {
        return 
    }

	t.IsOn = !t.IsOn
	//t.shouldRender = true
	//ctx.Update()

	// Only logs if the app was built with "-tags dev"
	if app.IsClient {
		// Formatted log
		app.Logf("ToggleSwitch state is now: %v", t.IsOn)
	}

}

/*
func (t *ToggleSwitch) Update(ctx app.Context) bool {
	if app.IsClient {
		app.Log("ToggleSwitch Update()")
	}
	return t.shouldRender
}
*/
