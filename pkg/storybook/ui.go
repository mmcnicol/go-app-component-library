// pkg/storybook/ui.go
package storybook

import (
	"github.com/maxence-charriere/go-app/v10/pkg/app"
	"strings"
)

type Shell struct {
	app.Compo

	activeComponent string
	activeStory     string
	searchQuery     string
}

func (s *Shell) OnMount(ctx app.Context) {
	if app.IsClient {
		app.Log("Shell OnMount()")
	}
	url := ctx.Page().URL()
	s.activeComponent = url.Query().Get("component")
	s.activeStory = url.Query().Get("story")
}

func (s *Shell) Render() app.UI {
	if app.IsClient {
		app.Log("Shell Render()")
	}
	allComponents := GetRegistry()
	query := strings.ToLower(s.searchQuery)

	// 1. Correctly filter components and close the loop/if blocks
	filteredComponents := make([]ComponentContainer, 0)
	for _, c := range allComponents {
		if query == "" || strings.Contains(strings.ToLower(c.Name), query) {
			filteredComponents = append(filteredComponents, c)
		}
	} // This closing brace was missing

	return app.Div().Class("storybook-layout").Body(
		app.Aside().Class("storybook-sidebar").Body(
			app.H2().Text("Components"),

			app.Div().Class("search-container").Body(
				app.Input().
					Class("sidebar-search").
					Placeholder("Filter components...").
					Value(s.searchQuery).
					OnInput(s.onSearch),

				app.If(s.searchQuery != "", func() app.UI {
					return app.Span().
						Class("search-clear").
						Text("✕").
						OnClick(s.onClearSearch)
				}),
			),

			app.Ul().Body(
				// 2. Use filteredComponents here instead of 'components'
				app.Range(filteredComponents).Slice(func(i int) app.UI {
					comp := filteredComponents[i]
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
	if app.IsClient {
		app.Log("Shell selectStory()")
	}
	s.activeComponent = compName
	s.activeStory = storyName

	u := ctx.Page().URL()
	q := u.Query()
	q.Set("component", compName)
	q.Set("story", storyName)
	u.RawQuery = q.Encode()

	ctx.Navigate(u.String())
}

func (s *Shell) renderActiveStory() app.UI {
	if app.IsClient {
		app.Log("Shell renderActiveStory()")
	}
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
}

func (s *Shell) onClearSearch(ctx app.Context, e app.Event) {
	s.searchQuery = ""
	ctx.Update()
}
