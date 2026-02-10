//go:build dev
// pkg/components/built_in/built_in_stories.go
package built_in

import (
	"github.com/maxence-charriere/go-app/v10/pkg/app"
	"github.com/mmcnicol/go-app-component-library/pkg/storybook"
	"strings"
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

	/*
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
	*/

	storybook.Register("Built In", "Select", 
		map[string]*storybook.Control{
			"PromptText": {Label: "Prompt Text", Type: storybook.ControlText, Value: "Choose an option..."},
			"Disabled":   {Label: "Disabled", Type: storybook.ControlBool, Value: false},
			// Make sure to set the Value field so the type assertion doesn't panic
			"Options":    {Label: "Options", Type: storybook.ControlText, Value: selectOptions}, 
			"SelectedValue": {Label: "Selected Value", Type: storybook.ControlText, Value: ""}, 
		},
		func(controls map[string]*storybook.Control) app.UI {
			promptText := controls["PromptText"].Value.(string)
			isDisabled := controls["Disabled"].Value.(bool)
			opts := controls["Options"].Value.([]string)
			selectedValue := controls["SelectedValue"].Value.(string)

			return app.Select().
				Disabled(isDisabled).
				//Value(selectedValue).
				OnChange(func(ctx app.Context, e app.Event) {
					val := ctx.JSSrc().Get("value").String()
					controls["SelectedValue"].Value = val
					ctx.Update()
				}).
				Body(
					// Placeholder option
					app.Option().Text(promptText).Value("").Selected(selectedValue == ""),
					
					app.Range(opts).Slice(func(i int) app.UI {
						optVal := opts[i]
						return app.Option().
							Text(optVal).
							Value(optVal).
							Selected(optVal == selectedValue) 
					},
				),
			)
		},
	)

	storybook.Register("Built In", "InputText", 
		map[string]*storybook.Control{
			"Value": {Label: "Value", Type: storybook.ControlText, Value: ""}, 
			"Disabled": {Label: "Disabled", Type: storybook.ControlBool, Value: false},
			"Placeholder": {Label: "Placeholder", Type: storybook.ControlText, Value: "First Name"}, 
		},
		func(controls map[string]*storybook.Control) app.UI {
			valueString := controls["Value"].Value.(string)
			isDisabled := controls["Disabled"].Value.(bool)
			placeholderString := controls["Placeholder"].Value.(string)

			return app.Input().
				Type("text").
				Value(valueString).
				Placeholder(placeholderString).
				Disabled(isDisabled).
				OnInput(func(ctx app.Context, e app.Event) {
					newVal := ctx.JSSrc().Get("value").String()
					controls["Value"].Value = newVal
					ctx.Update()
				},
			)
		},
	)

	storybook.Register("Built In", "Table", 
		map[string]*storybook.Control{
			"Caption":  {Label: "Caption", Type: storybook.ControlText, Value: "Employee Directory"},
			"Headers":  {Label: "Headers (CSV)", Type: storybook.ControlText, Value: "Name, Role, Location"},
			"Striped":  {Label: "Striped Rows", Type: storybook.ControlBool, Value: true},
			"Bordered": {Label: "Show Borders", Type: storybook.ControlBool, Value: true},
		},
		func(controls map[string]*storybook.Control) app.UI {
			caption := controls["Caption"].Value.(string)
			headersRaw := controls["Headers"].Value.(string)
			isStriped := controls["Striped"].Value.(bool)
			isBordered := controls["Bordered"].Value.(bool)

			// Process CSV headers into a slice
			headers := strings.Split(headersRaw, ",")

			// Dummy data for the table body
			rows := [][]string{
				{"Alice", "Engineer", "New York"},
				{"Bob", "Designer", "London"},
				{"Charlie", "Manager", "Tokyo"},
			}

			// Apply dynamic styling based on controls
			tableClass := "storybook-table"
			if isStriped {
				tableClass += " striped"
			}
			if isBordered {
				tableClass += " bordered"
			}

			return app.Table().
				Class(tableClass).
				Body(
					app.Caption().Text(caption),
					app.THead().Body(
						app.Tr().Body(
							app.Range(headers).Slice(func(i int) app.UI {
								return app.Th().Text(strings.TrimSpace(headers[i]))
							}),
						),
					),
					app.TBody().Body(
						app.Range(rows).Slice(func(i int) app.UI {
							return app.Tr().Body(
								app.Range(rows[i]).Slice(func(j int) app.UI {
									return app.Td().Text(rows[i][j])
								}),
							)
						}),
					),
					app.TFoot().Body(
						app.Text("table footer text."),
					),
				)
		},
	)

}
