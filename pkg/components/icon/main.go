// pkg/components/icon/main.go
package icon

import (
	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

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
