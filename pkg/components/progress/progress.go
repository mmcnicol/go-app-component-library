// pkg/components/progress/progress.go
package progress

import "github.com/maxence-charriere/go-app/v10/pkg/app"

type Progress struct {
    app.Compo
    Value float64
}

func (p *Progress) Render() app.UI {
    return app.Div().Class("progress-component-container").Body(
        app.Progress().
            Class("progress-component-bar"). // Custom styling hook
            Value(p.Value).
            Max(100),
    )
}
