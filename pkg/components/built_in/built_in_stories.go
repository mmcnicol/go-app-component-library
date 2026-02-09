//go:build dev
// pkg/components/built_in/built_in_stories.go
package built_in

import (
	"github.com/maxence-charriere/go-app/v10/pkg/app"
	"github.com/mmcnicol/go-app-component-library/pkg/storybook"
)

// Use init() to auto-register when this package is imported
func init() {

	/*
	storybook.Register("Built In", "Div", func() app.UI {
		return app.Div().
			Style("background", "#f0f0f0").
			Style("padding", "20px").
			Body(
				app.Text("hello"),
			)
	})
	*/

	storybook.Register("Built In", "Div", 
		map[string]*storybook.Control{
			"Text": {Label: "Text", Type: storybook.ControlText, Value: "Hello"},
		},
		func(controls map[string]*storybook.Control) app.UI {
			text := controls["Text"].Value.(string)

			return app.Div().
				Style("background", "#f0f0f0").
				Style("padding", "20px").
				Body(
					app.Text(text),
				)

		},
	)

	selectOptions := []string{"Go", "Python", "Rust", "JavaScript"}

	/*
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
	*/

	storybook.Register("Built In", "Select", 
		map[string]*storybook.Control{
			"PromptText": {Label: "Prompt Text", Type: storybook.ControlText, Value: "Choose an option..."},
			"Disabled": {Label: "Disabled", Type: storybook.ControlBool, Value: false},
			"Options": {Label: "Options", Type: storybook.ControlText, Options: selectOptions},
		},
		func(controls map[string]*storybook.Control) app.UI {
			// Retrieve values from the controls map
			promptText := controls["PromptText"].Value.(string)
			isDisabled := controls["Disabled"].Value.(bool)
			selectOptions := controls["Options"].Value.([]string)

			return app.Select().
			Body(
				app.Option().
					Disabled(isDisabled).
					Text(promptText),
				app.Range(selectOptions).Slice(func(i int) app.UI {
					opt := selectOptions[i]
					return app.Option().
						Value(opt).
						Text(opt)
				}),
			)
		},
	)


}
