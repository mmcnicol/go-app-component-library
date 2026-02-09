//go:build dev
// pkg/components/phase-banner/phase-banner_stories.go
package phase-banner

import (
	"github.com/maxence-charriere/go-app/v10/pkg/app"
	"github.com/mmcnicol/go-app-component-library/pkg/storybook"
)

// Use init() to auto-register when this package is imported
func init() {
	storybook.Register("Phase Banner", "Beta", func() app.UI {
		return &PhaseBanner{
			Phase: "Beta",
			Message: app.Text("This is a new service – your "),
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
}
