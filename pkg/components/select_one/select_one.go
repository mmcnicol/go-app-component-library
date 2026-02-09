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
	PromptText    string
	Disabled      bool
    //shouldRender  bool
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

	containerClass := ""
	if o.Disabled {
        containerClass += "selectOne-container-select-disabled"
    }

	return app.Div().Class("selectOne-container").Body(
		
		app.Select().
			Class("selectOne-container-select").
			Class(containerClass).
			Disabled(o.Disabled).
			//SelectedValue(o.selectedValue). // Keeps the UI in sync with Go state
			OnChange(o.onSelectChange).
			Body(
				app.Option().
					Disabled(true).
					Selected(o.selectedValue == "").
					Text(o.PromptText),
				app.Range(o.Options).Slice(func(i int) app.UI {
					opt := o.Options[i]
					return app.Option().
						Value(opt).
						Text(opt).
						// Set the "selected" attribute explicitly if it matches
						Selected(opt == o.selectedValue)
				}),
			),
	)
}

func (o *SelectOne) onSelectChange(ctx app.Context, e app.Event) {
	if app.IsClient {
		app.Log("SelectOne onSelectChange()")
	}
	o.selectedValue = ctx.JSSrc().Get("value").String()
	//o.shouldRender = true
	if app.IsClient {
		app.Logf("SelectOne state is now: %v", o.selectedValue)
	}
}

/*
func (o *SelectOne) Update(ctx app.Context) bool {
	if app.IsClient {
		app.Log("SelectOne Update()")
	}
	return o.shouldRender
}
*/
