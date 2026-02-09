//go:build dev
// pkg/components/phase_banner/phase_banner_stories.go
package phase_banner

import (
	"github.com/maxence-charriere/go-app/v10/pkg/app"
	"github.com/mmcnicol/go-app-component-library/pkg/storybook"
)

// Use init() to auto-register when this package is imported
func init() {

	/*
	storybook.Register("Phase Banner", "Beta", func() app.UI {
		return &PhaseBanner{
			Phase: "Beta",
			//Message: app.Text("This is a new service – your "),
			// You can nest elements inside the message
			Message: app.Span().Body(
				app.Text("This is a new service – your "),
				app.A().Href("/feedback").Text("feedback"),
				app.Text(" will help us to improve it."),
			),
		}
	})
	
	storybook.Register("Phase Banner", "Alpha", func() app.UI {
		return &PhaseBanner{
			Phase: "Alpha",
			Message: app.Text("This is a brand new service in development."),
		}
	})
	*/

	storybook.Register("Messages", "Phase Banner", 
		map[string]*storybook.Control{
			"Phase": {Label: "Phase", Type: storybook.ControlText, Value: "Alpha"},
			"Message": {Label: "Message", Type: storybook.ControlText, Value: "This is a brand new service in development."},
		},
		func(controls map[string]*storybook.Control) app.UI {
			phaseString := controls["Phase"].Value.(string)
			messageString := controls["Message"].Value.(string)

			return &PhaseBanner{
				Phase: phaseString,
				Message: app.Text(messageString),
			}
		},
	)
	
}
