package phase-banner

import (
	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

// PhaseBanner indicates service status
type PhaseBanner struct {
	app.Compo
	Phase   string // e.g., "Alpha" or "Beta"
	Message app.UI // The content/link to display
}

func (p *PhaseBanner) Render() app.UI {
	return app.Div().Class("phase-banner").Body(
		app.P().Class("phase-banner__content").Body(
			app.Strong().
				Class("phase-banner__content__tag").
				Text(p.Phase),
			app.Span().
				Class("phase-banner__text").
				Body(p.Message),
		),
	)
}
