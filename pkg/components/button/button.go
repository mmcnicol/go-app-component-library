// pkg/components/button/button.go
package button

import "github.com/maxence-charriere/go-app/v10/pkg/app"

type ButtonLook string

const (
    LookPrimary   ButtonLook = "primary"
    LookSecondary ButtonLook = "secondary"
    LookDanger    ButtonLook = "danger"
)

type Button struct {
    app.Compo
    Label    string
    Look     ButtonLook
    Disabled bool
    OnClick  func(ctx app.Context, e app.Event)
}

func (b *Button) Render() app.UI {
    class := "button-component button-component-" + string(b.Look)
    if b.Disabled {
        class += "button-component-disabled"
    }

    return app.Button().
        Class(class).
        Disabled(b.Disabled).
        OnClick(b.OnClick).
        Text(b.Label)
}
