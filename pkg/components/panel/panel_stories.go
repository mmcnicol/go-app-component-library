//go:build dev
// pkg/components/panel/panel_stories.go
package panel

import (
	"github.com/maxence-charriere/go-app/v10/pkg/app"
	"github.com/mmcnicol/go-app-component-library/pkg/storybook"
)

func init() {
	storybook.Register("Misc", "Panel", 
		map[string]*storybook.Control{
			"Title":   {Label: "Panel Title", Type: storybook.ControlText, Value: "My Panel"},
			"Padding": {Label: "Content Padding", Type: storybook.ControlText, Value: "20px"},
			"BodyText": {Label: "Inside Text", Type: storybook.ControlText, Value: "This is a panel container."},
		},
		func(controls map[string]*storybook.Control) app.UI {
			title := controls["Title"].Value.(string)
			padding := controls["Padding"].Value.(string)
			bodyText := controls["BodyText"].Value.(string)

			return app.Div().Style("padding", "40px").Body(
				&Panel{
					Title:   title,
					Padding: padding,
					Content: app.Div().Body(
						app.P().Text(bodyText),
					),
				},
			)
		},
	)
}
