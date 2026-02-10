//go:build dev
// pkg/components/icon/icon_stories.go
package icon

import (
	"strings"
	"github.com/maxence-charriere/go-app/v10/pkg/app"
	"github.com/mmcnicol/go-app-component-library/pkg/storybook"
)

// Use init() to auto-register when this package is imported
func init() {

	storybook.Register("Misc", "Icon", 
		map[string]*storybook.Control{
			"Search": {Label: "Search Icons", Type: storybook.ControlText, Value: "error"},
			"Size":   {Label: "Icon Size", Type: storybook.ControlNumber, Value: 48},
		},
		func(controls map[string]*storybook.Control) app.UI {
			searchQuery := strings.ToLower(controls["Search"].Value.(string))
			iconSize := controls["Size"].Value.(int)

			// List of all available icon names in your GetIcon switch
			allIcons := []string{"success", "info", "warn", "error"}
			
			// Create an instance of your icon component
			i := &Icon{}

			return app.Div().Style("padding", "20px").Body(
				app.H2().Text("Icon Gallery"),
				
				// The Grid Container
				app.Div().
					Style("display", "flex").
					Style("flex-wrap", "wrap").
					Style("gap", "20px").
					Body(
						app.Range(allIcons).Slice(func(idx int) app.UI {
							name := allIcons[idx]
							
							// Filtering logic
							if searchQuery != "" && !strings.Contains(name, searchQuery) {
								return nil
							}

							// Individual Icon Card
							return app.Div().
								Style("display", "flex").
								Style("flex-direction", "column").
								Style("align-items", "center").
								Style("width", "100px").
								Style("padding", "10px").
								Style("border", "1px solid #eee").
								Style("border-radius", "8px").
								Body(
									i.GetIcon(name, iconSize),
									app.Span().
										Style("margin-top", "10px").
										Style("font-size", "12px").
										Style("color", "#666").
										Text(name),
								)
						}),
					),
			)
		},
	)

}
