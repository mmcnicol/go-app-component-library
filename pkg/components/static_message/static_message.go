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
