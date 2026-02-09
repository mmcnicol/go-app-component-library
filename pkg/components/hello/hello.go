// pkg/components/hello/hello.go
package hello

import "github.com/maxence-charriere/go-app/v10/pkg/app"

type Hello struct {
    app.Compo
}

func (n *Hello) Render() app.UI {
    return app.Div().Body(
        app.P().Text("Hello World!"),
    )
}
