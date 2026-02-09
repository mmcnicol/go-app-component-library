// pkg/components/select_one/select_one.go
package select_one

import (
	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

// SelectOne defines the UI component
type SelectOne struct {
	app.Compo
	Options       []string
	selectedValue string
	Label         string
    shouldRender  bool
}

func (o *SelectOne) OnMount(ctx app.Context) {
	if app.IsClient {
		app.Log("SelectOne OnMount()")
	}
}

func (o *SelectOne) Render() app.UI {
	if app.IsClient {
		app.Log("SelectOne Render()")
	}

    o.shouldRender = false

	return app.Div().Class("picklist-container").Body(
		app.If(o.Label != "", app.Label().Text(o.Label)),
		
		app.Select().
			Class("picklist-select").
			Value(o.SelectedValue). // Keeps the UI in sync with Go state
			OnChange(o.onSelectChange).
			Body(
				app.Option().
					Disabled(true).
					Selected(p.SelectedValue == "").
					Text("Choose an option..."),
				app.Range(o.Options).Slice(func(i int) app.UI {
					opt := o.Options[i]
					return app.Option().
						Value(opt).
						Text(opt).
						// Set the "selected" attribute explicitly if it matches
						Selected(opt == o.SelectedValue)
				}),
			),
	)
}

func (o *SelectOne) onSelectChange(ctx app.Context, e app.Event) {
	if app.IsClient {
		app.Log("SelectOne onSelectChange()")
	}
	o.selectedValue = ctx.JSSrc().Get("value").String()
	o.shouldRender = true
}

func (o *SelectOne) Update(ctx app.Context) bool {
	if app.IsClient {
		app.Log("SelectOne Update()")
	}
	return o.shouldRender
}
