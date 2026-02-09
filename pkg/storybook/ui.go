// pkg/storybook/ui.go
package storybook

import (
	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

type Shell struct {
	app.Compo
	
	activeComponent string
	activeStory     string
	searchQuery     string
}

func (s *Shell) OnMount(ctx app.Context) {
	// Read URL params to restore state on reload
	url := ctx.Page().URL()
	s.activeComponent = url.Query().Get("component")
	s.activeStory = url.Query().Get("story")
	// In v10, you don't need to call Update here; 
	// modifying the struct and letting the lifecycle finish is enough.
}

func (s *Shell) Render() app.UI {
	allComponents := GetRegistry()
	query := strings.ToLower(s.searchQuery)
	
	// Filter components based on search query
	filteredComponents := make([]Component, 0)
	for _, c := range allComponents {
		
		// Simple case-insensitive match
		if query == "" || strings.Contains(strings.ToLower(c.Name), query) {
        filteredComponents = append(filteredComponents, c)
    }

	return app.Div().Class("storybook-layout").Body(
		app.Aside().Class("storybook-sidebar").Body(
			app.H2().Text("Components"),

			// Search Container
            app.Div().Class("search-container").Body(
                app.Input().
                    Class("sidebar-search").
                    Placeholder("Filter components...").
                    Value(s.searchQuery).
                    OnInput(s.onSearch),
                
                // Only show the clear button if there is text in the search box
                app.If(s.searchQuery != "", func() app.UI {
                    return app.Span().
                        Class("search-clear").
                        Text("✕"). // Multiplications X symbol
                        OnClick(s.onClearSearch)
                }),
            ),

			app.Ul().Body(
				app.Range(components).Slice(func(i int) app.UI {
					comp := components[i]
					return app.Li().Body(
						app.Strong().Text(comp.Name),
						app.Ul().Body(
							app.Range(comp.Stories).Slice(func(j int) app.UI {
								story := comp.Stories[j]
								isActive := s.activeComponent == comp.Name && s.activeStory == story.Name
								
								linkClass := "story-link"
								if isActive {
									linkClass += " active"
								}

								return app.Li().Body(
									app.A().
										Class(linkClass).
										Text(story.Name).
										OnClick(func(ctx app.Context, e app.Event) {
											s.selectStory(ctx, comp.Name, story.Name)
										}),
								)
							}),
						),
					)
				}),
			),
		),

		app.Main().Class("storybook-preview").Body(
			app.If(s.activeComponent != "", func() app.UI {
				return s.renderActiveStory()
			}).Else(func() app.UI {
				return app.Div().Class("empty-state").Text("Select a component from the sidebar")
			}),
		),
	)
}

func (s *Shell) selectStory(ctx app.Context, compName, storyName string) {
	s.activeComponent = compName
	s.activeStory = storyName
	
	// Update URL query params
	u := ctx.Page().URL()
	q := u.Query()
	q.Set("component", compName)
	q.Set("story", storyName)
	u.RawQuery = q.Encode()
	
	ctx.Navigate(u.String())
	
	// In v10, the UI refreshes automatically after an event handler.
	// No s.Update() or ctx.Update() needed here.
}

func (s *Shell) renderActiveStory() app.UI {
	for _, comp := range GetRegistry() {
		if comp.Name == s.activeComponent {
			for _, story := range comp.Stories {
				if story.Name == s.activeStory {
					return app.Div().Class("story-container").Body(
						story.Render(), 
					)
				}
			}
		}
	}
	return app.Div().Text("Story not found")
}

func (s *Shell) onSearch(ctx app.Context, e app.Event) {
	s.searchQuery = ctx.JSSrc().Get("value").String()
	// No s.Update() needed in v10; the UI refreshes after the event
}

// Simple helper for string containment
func contains(s, substr string) bool {
    // If the search box is empty, we usually want to show everything
    if substr == "" {
        return true
    }
    // Standard case-insensitive check
    return strings.Contains(strings.ToLower(s), strings.ToLower(substr))
}

// Add this handler to the Shell struct
func (s *Shell) onClearSearch(ctx app.Context, e app.Event) {
    s.searchQuery = ""
    // In v10, updating the state variable triggers the re-render automatically
}
