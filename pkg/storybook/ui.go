// pkg/storybook/ui.go
package storybook

import (
	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

// Shell is the main layout for the component library
type Shell struct {
	app.Compo // This embedding is crucial: it gives you s.Update()
	
	activeComponent string
	activeStory     string
}

func (s *Shell) OnMount(ctx app.Context) {
	// Read URL params to restore state on reload
	url := ctx.Page().URL()
	s.activeComponent = url.Query().Get("component")
	s.activeStory = url.Query().Get("story")
	s.Update()
}

func (s *Shell) Render() app.UI {
	components := GetRegistry()

	return app.Div().Class("storybook-layout").Body(
		// 1. Sidebar
		app.Aside().Class("storybook-sidebar").Body(
			app.H2().Text("Components"),
			app.Ul().Body(
				app.Range(components).Slice(func(i int) app.UI {
					comp := components[i]
					return app.Li().Body(
						app.Strong().Text(comp.Name),
						app.Ul().Body(
							app.Range(comp.Stories).Slice(func(j int) app.UI {
								story := comp.Stories[j]
								isActive := s.activeComponent == comp.Name && s.activeStory == story.Name
								
								// Calculate class string manually since ClassIf isn't standard
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

		// 2. Preview Area
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
	
	// Update URL without reloading so sharing works
	ctx.Page().URL().Query().Set("component", compName)
	ctx.Page().URL().Query().Set("story", storyName)
	ctx.Navigate("?component=" + compName + "&story=" + storyName)
	
	s.Update()
}

func (s *Shell) renderActiveStory() app.UI {
	// Find the render function in the registry
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
