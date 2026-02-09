//go:build dev
// pkg/components/hello/hello_stories.go
package hello

import (
	"github.com/maxence-charriere/go-app/v10/pkg/app"
	"github.com/mmcnicol/go-app-component-library/pkg/storybook"
)

// Use init() to auto-register when this package is imported
func init() {
	storybook.Register("Hello", "Default", func() app.UI {
		return &Hello{}
	})

    // You can easily create variants
	storybook.Register("Hello", "In Container", func() app.UI {
		return app.Div().
			Style("background", "#f0f0f0").
			Style("padding", "20px").
			Body(
				&Hello{},
			)
	})
}
