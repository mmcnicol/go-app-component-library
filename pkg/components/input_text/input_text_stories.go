//go:build dev
// pkg/components/input_text/input_text_stories.go
package input_text

import (
	"github.com/maxence-charriere/go-app/v10/pkg/app"
	"github.com/mmcnicol/go-app-component-library/pkg/storybook"
)

// Use init() to auto-register when this package is imported
func init() {

	storybook.Register("Input Text", "Default", func() app.UI {
		return &InputText{}
	})

}
