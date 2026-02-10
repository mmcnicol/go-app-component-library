// pkg/components/icon/main.go
package icon

import (
	"fmt"
	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

type Icon struct {
	app.Compo
}

// GetIcon returns an SVG icon
func (i *Icon) GetIcon(name string, size int) app.UI {
	var d string
	switch name {
	case "success":
		d = "M9 16.17L4.83 12l-1.42 1.41L9 19 21 7l-1.41-1.41z"
	case "info":
		d = "M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm1 15h-2v-6h2v6zm0-8h-2V7h2v2z"
	case "warn":
		d = "M1 21h22L12 2 1 21zm12-3h-2v-2h2v2zm0-4h-2v-4h2v4z"
	case "error":
		d = "M12 2C6.47 2 2 6.47 2 12s4.47 10 10 10 10-4.47 10-10S17.53 2 12 2zm5 13.59L15.59 17 12 13.41 8.41 17 7 15.59 10.59 12 7 8.41 8.41 7 12 10.59 15.59 7 17 8.41 13.41 12 17 15.59z"
	}

	// Using app.Raw avoids the "undefined: app.Svg" compilation error
	return app.Raw(fmt.Sprintf(
		`<svg class="icon icon-%s" viewBox="0 0 24 24" width="%d" height="%d"><path d="%s" /></svg>`,
		name, size, size, d,
	))
}

/*
// Icon returns an SVG icon based on the provided name and size.
func Icon(name string, size int) app.UI {
	var pathData string

	switch name {
	case "success":
		pathData = "M9 16.17L4.83 12l-1.42 1.41L9 19 21 7l-1.41-1.41z"
	case "info":
		pathData = "M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm1 15h-2v-6h2v6zm0-8h-2V7h2v2z"
	case "warn":
		pathData = "M1 21h22L12 2 1 21zm12-3h-2v-2h2v2zm0-4h-2v-4h2v4z"
	case "error":
		pathData = "M12 2C6.47 2 2 6.47 2 12s4.47 10 10 10 10-4.47 10-10S17.53 2 12 2zm5 13.59L15.59 17 12 13.41 8.41 17 7 15.59 10.59 12 7 8.41 8.41 7 12 10.59 15.59 7 17 8.41 13.41 12 17 15.59z"
	}

	// Use app.Svg() and app.Path() with the correct capitalization
	return app.Svg().
		Attr("viewBox", "0 0 24 24").
		Attr("width", size).
		Attr("height", size).
		Attr("fill", "currentColor").
		Body(
			app.Path().D(pathData),
		)
}
*/

/*
// Icon returns an SVG icon based on the provided name and size.
func Icon(name string, size int) app.UI {
	var paths []app.UI

	switch name {
	case "success":
		paths = []app.UI{app.Path().D("M9 16.17L4.83 12l-1.42 1.41L9 19 21 7l-1.41-1.41z")}
	case "info":
		paths = []app.UI{app.Path().D("M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm1 15h-2v-6h2v6zm0-8h-2V7h2v2z")}
	case "warn":
		paths = []app.UI{app.Path().D("M1 21h22L12 2 1 21zm12-3h-2v-2h2v2zm0-4h-2v-4h2v4z")}
	case "error":
		paths = []app.UI{app.Path().D("M12 2C6.47 2 2 6.47 2 12s4.47 10 10 10 10-4.47 10-10S17.53 2 12 2zm5 13.59L15.59 17 12 13.41 8.41 17 7 15.59 10.59 12 7 8.41 8.41 7 12 10.59 15.59 7 17 8.41 13.41 12 17 15.59z")}
	}

	return app.Raw("<svg viewBox='0 0 24 24' width='" + fmt.Sprint(size) + "' height='" + fmt.Sprint(size) + "' fill='currentColor'>" +
		app.HTMLString(app.Range(paths).Slice(func(i int) app.UI { return paths[i] })) +
		"</svg>")
}
*/

/*
func GetSVGIcon(name string, size int, color string) app.UI {
	iconBody := func() []app.UI {
		switch name {
		case "success":
			return []app.UI{app.Path().D("M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z")}
		case "info":
			return []app.UI{app.Path().D("M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z")}
		case "warn":
			return []app.UI{app.Path().D("M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z")}
		case "error":
			return []app.UI{app.Path().D("M10 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2m7-2a9 9 0 11-18 0 9 9 0 0118 0z")}
		default:
			return nil
		}
	}

	return app.Svg().
		Attribute("viewBox", "0 0 24 24").
		Attribute("width", size).
		Attribute("height", size).
		Attribute("stroke", color).
		Attribute("fill", "none").
		Attribute("stroke-width", "2").
		Body(iconBody()...)
}
*/
