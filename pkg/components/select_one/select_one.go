// pkg/components/select_one/select_one.go
package select_one

import (
	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

// SelectOne defines the UI component
type SelectOne struct {
	app.Compo
	Options       []string
	SelectedValue string
	PromptText    string
	Disabled      bool
    //shouldRender  bool
	OnSelect      func(ctx app.Context, val string)
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

    //o.shouldRender = false

	selectClass := "selectOne-container-select"
    if s.Disabled {
        selectClass += " selectOne-container-select-disabled"
    }
	
	return app.Div().Class("selectOne-container").Body(
        app.Select().
            Class(selectClass).
            Disabled(s.Disabled).
			//SelectedValue(o.SelectedValue). // Keeps the UI in sync with Go state
			//OnChange(o.onSelectChange).
			//OnChange(o.ValueTo(&o.SelectedValue)).
			//OnSelect(o.ValueTo(&o.SelectedValue)).
			OnChange(func(ctx app.Context, e app.Event) {
				val := ctx.JSSrc().Get("value").String()
				o.SelectedValue = val
				if o.OnSelect != nil {
					o.OnSelect(ctx, val)
				}
				//ctx.Update()
			}).
			Body(
				app.Option().
					Disabled(true).
					Selected(o.SelectedValue == "").
					Text(o.PromptText),
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
