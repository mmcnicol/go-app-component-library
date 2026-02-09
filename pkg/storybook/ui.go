// pkg/storybook/ui.go
package storybook

import (
	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

type Shell struct {
	app.Compo
	
	activeComponent string
	activeStory     string
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
	components := GetRegistry()

	return app.Div().Class("storybook-layout").Body(
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
