//go:build dev
// pkg/components/phase_banner/phase_banner_stories.go
package phase_banner

import (
	"github.com/maxence-charriere/go-app/v10/pkg/app"
	"github.com/mmcnicol/go-app-component-library/pkg/storybook"
)

func init() {
	storybook.Register("Messages", "Phase Banner", 
		map[string]*storybook.Control{
			"Phase": {
				Label: "Phase", 
				Type: storybook.ControlSelect, 
				Value: "Alpha",
				Options: []string{"Alpha", "Beta", "Gamma", "Stable", "Deprecated"},
			},
			"Message": {
				Label: "Message", 
				Type: storybook.ControlText, 
				Value: "This is a brand new service in development.",
			},
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
	
	// Optional: Keep the individual stories as examples if needed
	storybook.Register("Messages", "Phase Banner - Alpha", 
		map[string]*storybook.Control{}, // Empty controls map
		func(controls map[string]*storybook.Control) app.UI {
			return &PhaseBanner{
				Phase: "Alpha",
				Message: app.Span().Body(
					app.Text("This is a new service in development – your "),
					app.A().Href("/feedback").Text("feedback"),
					app.Text(" will help us to improve it."),
				),
			}
		},
	)
	
	storybook.Register("Messages", "Phase Banner - Beta", 
		map[string]*storybook.Control{}, // Empty controls map
		func(controls map[string]*storybook.Control) app.UI {
			return &PhaseBanner{
				Phase: "Beta",
				Message: app.Span().Body(
					app.Text("This service is currently in beta testing. "),
					app.A().Href("/report-issue").Text("Report any issues"),
					app.Text(" you encounter."),
				),
			}
		},
	)
}
