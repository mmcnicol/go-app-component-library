//go:build dev
// pkg/components/built_in/built_in_stories.go
package built_in

import (
	"github.com/maxence-charriere/go-app/v10/pkg/app"
	"github.com/mmcnicol/go-app-component-library/pkg/storybook"
	"fmt"
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
			"Caption": {Label: "Caption", Type: storybook.ControlText, Value: "Employee Directory"},
			//"Footer": {Label: "Footer", Type: storybook.ControlText, Value: "table footer text."},
		},
		func(controls map[string]*storybook.Control) app.UI {
			caption := controls["Caption"].Value.(string)
			//footerCaption := controls["Footer"].Value.(string)

			// Process CSV headers into a slice
			headers := []string{"Name", "Role", "Location"}

			// Dummy data for the table body
			rows := [][]string{
				{"Alice", "Engineer", "New York"},
				{"Bob", "Designer", "London"},
				{"Charlie", "Manager", "Tokyo"},
			}

			return app.Table().
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
					/*
					app.TFoot().Body(
						app.Tr().Body(
							app.Td().ColSpan(3).Text(footerCaption),
						),
					),
					*/
				)
		},
	)

	storybook.Register("Built In", "InputTextArea", 
		map[string]*storybook.Control{
			"Value":       {Label: "Value", Type: storybook.ControlText, Value: "Hello\nWorld!"}, 
			"Disabled":    {Label: "Disabled", Type: storybook.ControlBool, Value: false},
			"Placeholder": {Label: "Placeholder", Type: storybook.ControlText, Value: "Enter your message..."},
			"Rows":        {Label: "Rows", Type: storybook.ControlNumber, Value: 5},
		},
		func(controls map[string]*storybook.Control) app.UI {
			valueString := controls["Value"].Value.(string)
			isDisabled := controls["Disabled"].Value.(bool)
			placeholderString := controls["Placeholder"].Value.(string)
			rows := controls["Rows"].Value.(int)

			return app.Textarea().
				Style("width", "100%").
				Style("padding", "8px").
				Rows(rows).
				Placeholder(placeholderString).
				Disabled(isDisabled).
				OnInput(func(ctx app.Context, e app.Event) {
					// Sync the change back to the registry
					newVal := ctx.JSSrc().Get("value").String()
					controls["Value"].Value = newVal
					ctx.Update()
				}).
				Body(
					// This sets the content of the textarea
					app.Text(valueString),
				)
		},
	)

	storybook.Register("Built In", "Button", 
		map[string]*storybook.Control{
			"Label":    {Label: "Button Text", Type: storybook.ControlText, Value: "Click Me"}, 
			"Disabled": {Label: "Disabled", Type: storybook.ControlBool, Value: false},
			"Title":    {Label: "Tooltip (Title)", Type: storybook.ControlText, Value: "Submit Form"}, 
		},
		func(controls map[string]*storybook.Control) app.UI {
			label := controls["Label"].Value.(string)
			isDisabled := controls["Disabled"].Value.(bool)
			title := controls["Title"].Value.(string)

			return app.Button().
				Text(label).
				Title(title).
				Disabled(isDisabled).
				//Style("padding", "10px 20px").
				//Style("cursor", "pointer").
				OnClick(func(ctx app.Context, e app.Event) {
					// For a storybook, you might want to log the click
					app.Log("Button clicked: " + label)
				})
		},
	)

	storybook.Register("Built In", "Meter", 
		map[string]*storybook.Control{
			"Value":   {Label: "Current Value", Type: storybook.ControlNumber, Value: 60}, 
			"Min":     {Label: "Min Value", Type: storybook.ControlNumber, Value: 0}, 
			"Max":     {Label: "Max Value", Type: storybook.ControlNumber, Value: 100}, 
			"Low":     {Label: "Low Threshold", Type: storybook.ControlNumber, Value: 30}, 
			"High":    {Label: "High Threshold", Type: storybook.ControlNumber, Value: 80}, 
			"Optimum": {Label: "Optimum Value", Type: storybook.ControlNumber, Value: 90}, 
		},
		func(controls map[string]*storybook.Control) app.UI {
			// Asserting as int, then converting to float64 as required by meter methods
			val := float64(controls["Value"].Value.(int))
			min := float64(controls["Min"].Value.(int))
			max := float64(controls["Max"].Value.(int))
			low := float64(controls["Low"].Value.(int))
			high := float64(controls["High"].Value.(int))
			optimum := float64(controls["Optimum"].Value.(int))

			return app.Div().Body(
				//app.P().Text(fmt.Sprintf("Status: %.0f%%", val)),
				app.Meter().
					Style("width", "100%").
					Style("height", "25px").
					Min(min).
					Max(max).
					Value(val).
					Low(low).
					High(high).
					Optimum(optimum).
					Body(
						// Fallback text for older browsers
						app.Text(fmt.Sprintf("%.0f", val)),
					),
			)
		},
	)

	storybook.Register("Built In", "Time", 
		map[string]*storybook.Control{
			"Value":    {Label: "Time", Type: storybook.ControlText, Value: "12:00"}, 
			"Disabled": {Label: "Disabled", Type: storybook.ControlBool, Value: false},
			"Min":      {Label: "Min (08:00)", Type: storybook.ControlText, Value: "08:00"},
			"Max":      {Label: "Max (18:00)", Type: storybook.ControlText, Value: "18:00"},
		},
		func(controls map[string]*storybook.Control) app.UI {
			val := controls["Value"].Value.(string)
			dis := controls["Disabled"].Value.(bool)
			min := controls["Min"].Value.(string)
			max := controls["Max"].Value.(string)

			return app.Input().
				Type("time").
				Value(val).
				// Use the built-in methods for go-app v10
				Disabled(dis).
				Min(min). 
				Max(max).
				// Add a style to visualize invalid states
				Style("border", "2px solid").
				Style("border-color", "initial"). 
				OnInput(func(ctx app.Context, e app.Event) {
					newVal := ctx.JSSrc().Get("value").String()
					controls["Value"].Value = newVal
					
					// Trigger a browser check for min/max validity
					isInvalid := ctx.JSSrc().Get("validity").Get("valid").Bool()
					if !isInvalid {
						app.Log("Time is outside restricted range!")
					}
					
					ctx.Update()
				})
		},
	)

	storybook.Register("Built In", "Dialog", 
		map[string]*storybook.Control{
			"Open":    {Label: "Open", Type: storybook.ControlBool, Value: true}, 
			"Title":   {Label: "Dialog Title", Type: storybook.ControlText, Value: "Confirmation"}, 
			"Message": {Label: "Message Body", Type: storybook.ControlText, Value: "Are you sure you want to proceed?"}, 
		},
		func(controls map[string]*storybook.Control) app.UI {
			isOpen := controls["Open"].Value.(bool)
			title := controls["Title"].Value.(string)
			message := controls["Message"].Value.(string)

			return app.Div().Body(
				app.Dialog().
					// The 'open' attribute determines visibility in go-app
					Open(isOpen).
					//Style("border", "1px solid #ccc").
					//Style("border-radius", "8px").
					//Style("padding", "20px").
					//Style("box-shadow", "0 4px 6px rgba(0,0,0,0.1)").
					Body(
						app.H3().Text(title),
						app.P().Text(message),
						app.Div().Style("text-align", "right").Body(
							app.Button().Text("Cancel").OnClick(func(ctx app.Context, e app.Event) {
								controls["Open"].Value = false
								app.Log("Dialog Cancel button clicked")
								ctx.Update()
							}),
							app.Button().
								Style("margin-left", "10px").
								Text("Confirm").
								OnClick(func(ctx app.Context, e app.Event) {
									app.Log("Dialog Confirm button clicked")
									controls["Open"].Value = false
									ctx.Update()
								}),
						),
					),
			)
		},
	)

	storybook.Register("Built In", "Progress", 
        map[string]*storybook.Control{
            "Value": {Label: "Value", Type: storybook.ControlNumber, Value: 50},
        },
        func(controls map[string]*storybook.Control) app.UI {
            val := float64(controls["Value"].Value.(int))
            return app.Div().Style("padding", "20px").Body(
                app.Progress().Value(val).Max(100),
            )
        },
    )

}
