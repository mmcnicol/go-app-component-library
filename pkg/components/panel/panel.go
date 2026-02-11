// pkg/components/panel/panel.go
package panel

import (
	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

type Panel struct {
	app.Compo
	Content app.UI
	Padding string
	Title   string
}

func (p *Panel) Render() app.UI {
	return app.Div().
		Style("border", "1px solid #eee").
		Style("border-radius", "8px").
		Style("background-color", "inherit").
		Style("display", "flex").
		Style("flex-direction", "column").
		Body(
			app.If(p.Title != "",
				app.Div().
					Style("padding", "10px 15px").
					Style("border-bottom", "1px solid #eee").
					Style("font-weight", "bold").
					Text(p.Title),
			),
			app.Div().
				Style("padding", p.getPadding()).
				Body(p.Content),
		)
}

func (p *Panel) getPadding() string {
	if p.Padding == "" {
		return "20px"
	}
	return p.Padding
}
