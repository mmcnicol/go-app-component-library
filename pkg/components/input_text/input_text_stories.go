//go:build dev
// pkg/components/input_text/input_text_stories.go
package input_text

import (
	"github.com/maxence-charriere/go-app/v10/pkg/app"
	"github.com/mmcnicol/go-app-component-library/pkg/storybook"
)

// Use init() to auto-register when this package is imported
func init() {

	inputTextDefault := &InputText{
		Label: "Label text.",
	}

	storybook.Register("Input Text", "Default", func() app.UI {
		return inputTextDefault
	})

	inputTextPlaceholder := &InputText{
		Label: "Label text.",
		Placeholder: "sample placeholder text.",
	}

	storybook.Register("Input Text", "Placeholder", func() app.UI {
		return inputTextPlaceholder
	})

	inputTextReadOnly := &InputText{
		Label: "Label text.",
		Value: "this is a test",
		Disabled: true,
	}

	storybook.Register("Input Text", "ReadOnly", func() app.UI {
		return inputTextReadOnly
	})

}
