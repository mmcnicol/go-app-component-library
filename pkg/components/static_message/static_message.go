// pkg/components/static_message/static_message.go
package static_message

import (
	"github.com/maxence-charriere/go-app/v10/pkg/app"
	"github.com/mmcnicol/go-app-component-library/pkg/components/icon"
)

type StaticMessage struct {
	app.Compo
	Severity string // info, warn, error, success
	Summary  string
	Detail   string
}

func (m *StaticMessage) Render() app.UI {
    i := &icon.Icon{}
    return app.Div().Class("staticMessage-container").Body(
        app.Div().Class("staticMessage-icon").Body(
            i.GetIcon(m.Severity, 20),
        ),
        app.Div().Class("staticMessage-content-summary").Text(m.Summary),
		app.Div().Class("staticMessage-content-detail").Text(m.Detail),
    )
}

/*
func (m *StaticMessage) Render() app.UI {

	i := icon.Icon{}

	// Map severity to style colors
	//icon := "pi-info-circle"
	color := "#3B82F6"      // info blue
	bgColor := "#EFF6FF"
	
	switch m.Severity {
	case "success":
		//icon = "pi-check"
		color = "#22C55E"
		bgColor = "#F0FDF4"
	case "warn":
		//icon = "pi-exclamation-triangle"
		color = "#F59E0B"
		bgColor = "#FFFBEB"
	case "error":
		//icon = "pi-times-circle"
		color = "#EF4444"
		bgColor = "#FEF2F2"
	}

	return app.Div().
		Style("display", "flex").
		Style("align-items", "center").
		Style("padding", "1rem").
		Style("background-color", bgColor).
		Style("border-radius", "6px").
		Style("border-left", "6px solid "+color).
		Body(
			// Icon placeholder (assuming FontAwesome or PrimeIcons CSS is loaded)
			//app.I().Class("pi " + icon).Style("color", color).Style("margin-right", "0.5rem"),
			//icon.GetSVGIcon(m.Severity, 24, color), // Injected SVG
			i.GetIcon(m.Severity, 72),
			app.Span().
				Style("font-weight", "bold").
				Style("margin-left", "0.5rem").
				Style("margin-right", "0.5rem").
				Text(m.Summary),
			app.Span().Text(m.Detail),
		)
}
*/
