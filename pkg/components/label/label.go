// pkg/components/label/label.go
package label

import "github.com/maxence-charriere/go-app/v10/pkg/app"

type Label struct {
	app.Compo
	Text     string
	For      string // ID of the associated input
	Required bool
}

func (l *Label) Render() app.UI {
	return app.Label().
		Class("label-component").
		For(l.For).
		Body(
			app.Text(l.Text),
			app.If(l.Required, func() app.UI {
				return app.Span().Class("label-component-required").Text("*")
			}),
		)
}
