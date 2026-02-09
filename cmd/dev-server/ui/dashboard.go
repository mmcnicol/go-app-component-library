// cmd/dev-server/ui/dashboard.go
package ui

import (
    "github.com/maxence-charriere/go-app/v10/pkg/app"
)

type DevDashboard struct {
    app.Compo
    
    buildStatus    string
    lastBuildTime  string
    connectedClients int
    fileChanges    []string
    compileErrors  []string
}

func (d *DevDashboard) Render() app.UI {
	return app.Div().Class("dev-dashboard").Body(
		app.Header().Class("dashboard-header").Body(
			app.H1().Text("Go UI Library Dev Server"),
			// Fixed: added getStatusClass() logic below
			app.Div().Class("status-indicator").Class(d.getStatusClass()).Body(
				app.Span().Text(d.buildStatus),
			),
		),

		app.Main().Class("dashboard-main").Body(
			app.Section().Class("build-info").Body(
				app.H2().Text("Build Information"),
				app.Ul().Body(
					app.Li().Body(
						app.Strong().Text("Last build: "),
						app.Span().Text(d.lastBuildTime),
					),
					app.Li().Body(
						app.Strong().Text("Connected clients: "),
						app.Span().Text(d.connectedClients),
					),
				),
			),

			app.Section().Class("file-changes").Body(
				app.H2().Text("Recent Changes"),
				d.renderFileChanges(), // Fixed: added implementation below
			),

			app.Section().Class("errors").Body(
				app.H2().Text("Errors"),
				d.renderErrors(), // Fixed: added implementation below
			),

			app.Section().Class("actions").Body(
				app.H2().Text("Actions"),
				app.Button().
					Class("btn btn-primary").
					Text("Force Rebuild").
					OnClick(d.onForceRebuild), // Fixed: added handler below
				app.Button().
					Class("btn btn-secondary").
					Text("Clear Cache").
					OnClick(d.onClearCache), // Fixed: added handler below
			),
		),
	)
}

func (d *DevDashboard) getStatusClass() string {
	switch d.buildStatus {
	case "success":
		return "status-success"
	case "failed":
		return "status-error"
	default:
		return "status-pending"
	}
}

func (d *DevDashboard) renderFileChanges() app.UI {
	if len(d.fileChanges) == 0 {
		return app.P().Text("No recent changes detected.")
	}
	return app.Ul().Body(
		app.Range(d.fileChanges).Slice(func(i int) app.UI {
			return app.Li().Text(d.fileChanges[i])
		}),
	)
}

func (d *DevDashboard) renderErrors() app.UI {
	if len(d.compileErrors) == 0 {
		return app.P().Class("no-errors").Text("Clean build - no errors!")
	}
	return app.Div().Class("error-log").Body(
		app.Range(d.compileErrors).Slice(func(i int) app.UI {
			return app.Pre().Class("error-item").Text(d.compileErrors[i])
		}),
	)
}

// Event handlers in go-app v10 require this specific signature
func (d *DevDashboard) onForceRebuild(ctx app.Context, e app.Event) {
	app.Log("Force rebuild triggered")
	// Logic to trigger rebuild via server websocket/api would go here
}

func (d *DevDashboard) onClearCache(ctx app.Context, e app.Event) {
	app.Log("Clear cache triggered")
}
