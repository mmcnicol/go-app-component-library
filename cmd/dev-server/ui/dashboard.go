// cmd/dev-server/ui/dashboard.go
package ui

import (
    "github.com/maxenglander/go-app/v9/pkg/app"
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
                d.renderFileChanges(),
            ),
            
            app.Section().Class("errors").Body(
                app.H2().Text("Errors"),
                d.renderErrors(),
            ),
            
            app.Section().Class("actions").Body(
                app.H2().Text("Actions"),
                app.Button().
                    Class("btn btn-primary").
                    Text("Force Rebuild").
                    OnClick(d.onForceRebuild),
                app.Button().
                    Class("btn btn-secondary").
                    Text("Clear Cache").
                    OnClick(d.onClearCache),
            ),
        ),
    )
}

