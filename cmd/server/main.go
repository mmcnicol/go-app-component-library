// cmd/server/main.go
package main

import (
	"github.com/maxence-charriere/go-app/v10/pkg/app"
	"github.com/mmcnicol/go-app-component-library/pkg/storybook"
	
    // Import your components so their init() functions run and register stories
	_ "github.com/mmcnicol/go-app-component-library/pkg/components/built_in"
	//_ "github.com/mmcnicol/go-app-component-library/pkg/components/hello"
	_ "github.com/mmcnicol/go-app-component-library/pkg/components/phase_banner"
	_ "github.com/mmcnicol/go-app-component-library/pkg/components/toggle_switch"
	_ "github.com/mmcnicol/go-app-component-library/pkg/components/select_one"
	_ "github.com/mmcnicol/go-app-component-library/pkg/components/input_text"
	_ "github.com/mmcnicol/go-app-component-library/pkg/components/static_message"
	_ "github.com/mmcnicol/go-app-component-library/pkg/components/icon"
	_ "github.com/mmcnicol/go-app-component-library/pkg/components/input_text_area"
	_ "github.com/mmcnicol/go-app-component-library/pkg/components/tree"
	_ "github.com/mmcnicol/go-app-component-library/pkg/components/progress"
	_ "github.com/mmcnicol/go-app-component-library/pkg/components/button"
	_ "github.com/mmcnicol/go-app-component-library/pkg/components/label"
	_ "github.com/mmcnicol/go-app-component-library/pkg/components/table"
	_ "github.com/mmcnicol/go-app-component-library/pkg/components/panel"
	_ "github.com/mmcnicol/go-app-component-library/pkg/components/chart"

	"log"
	"net/http"
)

func main() {
	// CRITICAL: The server MUST know the route exists
	// Route "/" to the Storybook shell
	app.Route("/", func() app.Composer { return &storybook.Shell{} })

	h := &app.Handler{
		Name:      "go-app component library",
		Description: "A go-app UI library using Go and WebAssembly",
		Author:      "mmcnicol",
		Styles: []string{
			"/web/style/variables.css", // Load this FIRST
			"/web/style/main.css",
			"/web/style/phase_banner.css",
			"/web/style/toggle_switch.css",
			"/web/style/select_one.css",
			"/web/style/input_text.css",
			"/web/style/input_text_area.css",
			"/web/style/tree.css",
			"/web/style/static_message.css",
			"/web/style/progress.css",
			"/web/style/button.css",
			"/web/style/icon.css",
			"/web/style/label.css",
			"/web/style/table.css",
			"/web/style/sortable_table.css",
			"/web/style/data_grid.css",
			"/web/style/chart.css",
		},
		Icon: app.Icon{
			//Default: "/web/images/logo.png",
		},
		//Resources: app.LocalDir("web"),
	}

	http.Handle("/", h)

	// Example API endpoint
	http.HandleFunc("/api/data", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"status": "ok"}`))
	})

	log.Println("Serving at http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
