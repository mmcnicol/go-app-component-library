// pkg/components/select_one/select_one.go
package select_one

import (
	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

// SelectOne defines the UI component
type SelectOne struct {
	app.Compo
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
}
