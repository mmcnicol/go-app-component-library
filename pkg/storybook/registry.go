// pkg/storybook/registry.go
package storybook

import (
	"sort"
	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

type ControlType string

const (
    ControlText   ControlType = "text"
    ControlBool   ControlType = "bool"
    ControlNumber ControlType = "number"
)

type Control struct {
    Label    string
    Type     ControlType
    Value    any
    Options  []string // For select inputs
	ReadOnly bool
}

// Story represents a single view of a component (e.g., "Primary", "Disabled")
type Story struct {
	Name   string
	Controls map[string]*Control
    // Render now accepts the current state of controls
    Render   func(controls map[string]*Control) app.UI
}

// ComponentContainer holds all stories for a specific component
type ComponentContainer struct {
	Name    string
	Stories []Story
}

// registry stores all registered components
var registry = make(map[string][]Story)

/*
// Register adds a story to the registry
func Register(componentName string, storyName string, render func() app.UI) {
	registry[componentName] = append(registry[componentName], Story{
		Name:   storyName,
		Render: render,
	})
}
*/

// Register adds a story to the registry
func Register(componentName string, storyName string, controls map[string]*Control, render func(map[string]*Control) app.UI) {
    // Ensure map isn't nil if none provided
    if controls == nil {
        controls = make(map[string]*Control)
    }

    registry[componentName] = append(registry[componentName], Story{
        Name:     storyName,
        Controls: controls,
        Render:   render,
    })
}

// GetRegistry returns the sorted list of components for the sidebar
func GetRegistry() []ComponentContainer {
	var components []ComponentContainer
	for name, stories := range registry {
		// Sort stories by name
		sort.Slice(stories, func(i, j int) bool {
			return stories[i].Name < stories[j].Name
		})
		components = append(components, ComponentContainer{
			Name:    name,
			Stories: stories,
		})
	}
	// Sort components by name
	sort.Slice(components, func(i, j int) bool {
		return components[i].Name < components[j].Name
	})
	return components
}
