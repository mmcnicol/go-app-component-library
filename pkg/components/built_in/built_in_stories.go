//go:build dev
// pkg/components/built_in/built_in_stories.go
package built_in

import (
	"github.com/maxence-charriere/go-app/v10/pkg/app"
	"github.com/mmcnicol/go-app-component-library/pkg/storybook"
)

// Use init() to auto-register when this package is imported
func init() {

	storybook.Register("Built In", "Div", func() app.UI {
		return app.Div().
			Style("background", "#f0f0f0").
			Style("padding", "20px").
			Body(
				app.Text("hello"),
			)
	})

	selectOptions := []string{"Go", "Python", "Rust", "JavaScript"}

	storybook.Register("Built In", "Select", func() app.UI {
		return app.Select().
			Body(
				app.Option().
					Disabled(true).
					Text("Choose an option..."),
				app.Range(selectOptions).Slice(func(i int) app.UI {
					opt := selectOptions[i]
					return app.Option().
						Value(opt).
						Text(opt)
				}),
			)
	})
}
