// pkg/storybook/ui.go
package storybook

import (
	"github.com/maxence-charriere/go-app/v10/pkg/app"
	"strconv"
	"strings"
)

type Shell struct {
	app.Compo

	activeComponent string
	activeStory     string
	searchQuery     string
	shouldRender    bool
	showControls    bool
	IsDark          bool
}

func (s *Shell) OnInit() {
    // Load the preference from local storage
	app.Window().GetState("storybook-theme-dark", &s.IsDark)
}

func (s *Shell) OnMount(ctx app.Context) {
	if app.IsClient {
		app.Log("Shell OnMount()")
	}
	s.showControls = true // Default to open
	url := ctx.Page().URL()
	s.activeComponent = url.Query().Get("component")
	s.activeStory = url.Query().Get("story")
}

func (s *Shell) Render() app.UI {
	if app.IsClient {
		app.Log("Shell Render()")
	}

	s.shouldRender = false

	// Dynamically apply the dark-theme class
	layoutClass := "storybook-layout"
	if s.IsDark {
		layoutClass += " dark-theme"
	}

	allComponents := GetRegistry()
	query := strings.ToLower(s.searchQuery)

	filteredComponents := make([]ComponentContainer, 0)
	for _, c := range allComponents {
		if query == "" || strings.Contains(strings.ToLower(c.Name), query) {
			filteredComponents = append(filteredComponents, c)
		}
	}

	return app.Div().Class("storybook-layout").Body(
		// LEFT SIDEBAR
		app.Aside().Class("storybook-sidebar").Body(
			app.H2().Text("Components"),

			&ThemeSwitcher{
				IsDark: s.IsDark,
				OnChange: func(isDark bool) {
					s.IsDark = isDark
					// Persist the choice
					app.Window().SetState("storybook-theme-dark", isDark)
					//s.Update()
				},
			},

			app.Div().Class("search-container").Body(
				app.Input().
					ID("sidebar-search-input").
					Class("sidebar-search").
					Placeholder("Filter components...").
					Value(s.searchQuery).
					//AutoFocus(true).
					//OnChange(s.ValueTo(&s.searchQuery)),
					OnInput(s.onSearch),

				app.If(s.searchQuery != "", func() app.UI {
					return app.Span().
						Class("search-clear").
						Text("✕").
						OnClick(s.onClearSearch)
				}),
			),

			app.Ul().Body(
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

		// MAIN CONTENT AREA / MAIN PREVIEW
		app.Main().Class("storybook-main").Body(
            app.Div().Class("canvas-header").Body(
                app.Button().
                    Class("toggle-controls-btn").
                    Text("⚙ Controls").
                    OnClick(func(ctx app.Context, e app.Event) {
                        s.showControls = !s.showControls
                        s.shouldRender = true
                    }),
            ),
            
            app.Div().Class("canvas-content").Body(
                app.If(s.activeComponent != "", func() app.UI {
                    story := s.getActiveStory()
                    return app.Div().Class("story-container").Body(
                        story.Render(story.Controls),
                    )
                }).Else(func() app.UI {
                    return app.Div().Class("empty-state").Text("Select a story")
                }),
            ),
        ),

		// RIGHT CONTROLS PANEL (Conditional)
		app.If(s.showControls, func() app.UI {
            return app.Aside().Class("storybook-controls-panel").Body(
                s.renderControls(),
            )
        }),

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
	//s.shouldRender = true
	//ctx.Update()
}

func (s *Shell) getActiveStory() *Story {
    for _, comp := range GetRegistry() {
        if comp.Name == s.activeComponent {
            for i := range comp.Stories {
                if comp.Stories[i].Name == s.activeStory {
                    // Return the address of the story in the registry
                    return &comp.Stories[i]
                }
            }
        }
    }
    return nil
}

func (s *Shell) renderActiveStory() app.UI {
    story := s.getActiveStory()
    if story == nil {
        return app.Div().Text("Story not found")
    }
    
    return app.Div().Class("story-container").Body(
        story.Render(story.Controls), // Pass the map here
    )
}

func (s *Shell) onSearch(ctx app.Context, e app.Event) {
	if app.IsClient {
		app.Log("Shell onSearch()")
	}
	s.searchQuery = ctx.JSSrc().Get("value").String()
	s.shouldRender = true
}

func (s *Shell) onClearSearch(ctx app.Context, e app.Event) {
	if app.IsClient {
		app.Log("Shell onClearSearch()")
	}
	s.searchQuery = ""
	s.shouldRender = true
	//ctx.Update()

	elem := app.Window().Get("document").Call("getElementById", "sidebar-search-input")
	elem.Set("value", "")

	// Optional: Put the cursor back in the search box after clearing
    app.Window().GetElementByID("sidebar-search-input").Call("focus")
}

func (s *Shell) Update(ctx app.Context) bool {
	if app.IsClient {
		app.Log("Shell Update()")
	}
	return s.shouldRender
}

func (s *Shell) renderControls() app.UI {
    story := s.getActiveStory()
    if story == nil || len(story.Controls) == 0 {
        return app.Div().Text("No controls available")
    }

    return app.Div().Class("storybook-controls").Body(
        app.H3().Text("Properties"),
        app.Table().Body(
            app.Range(story.Controls).Map(func(k string) app.UI {
                ctrl := story.Controls[k]
                return app.Tr().Body(
                    app.Td().Text(ctrl.Label),
                    app.Td().Body(
                        s.renderControlInput(k, ctrl),
                    ),
                )
            }),
        ),
    )
}

func (s *Shell) renderControlInput(key string, ctrl *Control) app.UI {
    switch ctrl.Type {
    case ControlBool:
        return app.Input().
			Type("checkbox").
			Checked(ctrl.Value.(bool)).
			Disabled(ctrl.ReadOnly).
            OnChange(func(ctx app.Context, e app.Event) {
                ctrl.Value = ctx.JSSrc().Get("checked").Bool()
                s.shouldRender = true // Trigger Shell update
            })
    case ControlText, ControlNumber:
        return app.Input().
			Type("text").
			Value(ctrl.Value).
			Disabled(ctrl.ReadOnly).
            OnInput(func(ctx app.Context, e app.Event) {
                val := ctx.JSSrc().Get("value").String()
                
                // If it's a number, you might want to convert it back to int/float
                if ctrl.Type == ControlNumber {
                    i, _ := strconv.Atoi(val)
                    ctrl.Value = i
                } else {
                    ctrl.Value = val
                }
                
                s.shouldRender = true
                ctx.Update() // Essential to notify go-app to diff the DOM
            })
    default:
        return app.Text("Unsupported control")
    }
}
