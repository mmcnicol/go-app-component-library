// pkg/components/icon/icon.go
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
	isSpinner := name == "spinner"

	switch name {
	case "success":
		d = "M9 16.17L4.83 12l-1.42 1.41L9 19 21 7l-1.41-1.41z"
	case "info":
		d = "M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm1 15h-2v-6h2v6zm0-8h-2V7h2v2z"
	case "warn":
		d = "M1 21h22L12 2 1 21zm12-3h-2v-2h2v2zm0-4h-2v-4h2v4z"
	case "error":
		d = "M12 2C6.47 2 2 6.47 2 12s4.47 10 10 10 10-4.47 10-10S17.53 2 12 2zm5 13.59L15.59 17 12 13.41 8.41 17 7 15.59 10.59 12 7 8.41 8.41 7 12 10.59 15.59 7 17 8.41 13.41 12 17 15.59z"
	case "chevron-right":
		d = "M10 6L8.59 7.41 13.17 12l-4.58 4.59L10 18l6-6z"
	case "chevron-down":
		d = "M16.59 8.59L12 13.17 7.41 8.59 6 10l6 6 6-6z"
	case "chevron-left":
		d = "M15.41 7.41L14 6l-6 6 6 6 1.41-1.41L10.83 12z"
	case "chevron-up":
		d = "M7.41 15.41L12 10.83l4.59 4.58L18 14l-6-6-6 6z"
	case "sort":
		d = "M3 18h6v-2H3v2zM3 6v2h18V6H3zm0 7h12v-2H3v2z"
	case "settings":
		d = "M19.14 12.94c.04-.3.06-.61.06-.94 0-.32-.02-.64-.07-.94l2.03-1.58c.18-.14.23-.41.12-.61l-1.92-3.32c-.12-.22-.37-.29-.59-.22l-2.39.96c-.5-.38-1.03-.7-1.62-.94l-.36-2.54c-.04-.24-.24-.41-.48-.41h-3.84c-.24 0-.43.17-.47.41l-.36 2.54c-.59.24-1.13.57-1.62.94l-2.39-.96c-.22-.08-.47 0-.59.22L2.74 8.87c-.12.21-.08.47.12.61l2.03 1.58c-.05.3-.09.63-.09.94s.02.64.07.94l-2.03 1.58c-.18.14-.23.41-.12.61l1.92 3.32c.12.22.37.29.59.22l2.39-.96c.5.38 1.03.7 1.62.94l.36 2.54c.05.24.24.41.48.41h3.84c.24 0 .44-.17.47-.41l.36-2.54c.59-.24 1.13-.56 1.62-.94l2.39.96c.22.08.47 0 .59-.22l1.92-3.32c.12-.22.07-.47-.12-.61l-2.01-1.58zM12 15.6c-1.98 0-3.6-1.62-3.6-3.6s1.62-3.6 3.6-3.6 3.6 1.62 3.6 3.6-1.62 3.6-3.6 3.6z"
	case "logout":
		d = "M17 7l-1.41 1.41L18.17 11H8v2h10.17l-2.58 2.58L17 17l5-5zM4 5h8V3H4c-1.1 0-2 .9-2 2v14c0 1.1.9 2 2 2h8v-2H4V5z"
	case "folder":
		d = "M10 4H4c-1.1 0-1.99.9-1.99 2L2 18c0 1.1.9 2 2 2h16c1.1 0 2-.9 2-2V8c0-1.1-.9-2-2-2h-8l-2-2z"
	case "file":
		d = "M14 2H6c-1.1 0-1.99.9-1.99 2L4 20c0 1.1.89 2 1.99 2H18c1.1 0 2-.9 2-2V8l-6-6zm2 16H8v-2h8v2zm0-4H8v-2h8v2zm-3-5V3.5L18.5 9H13z"
	case "hospital-inpatient":
		// Person symbol next to a medical cross (Outpatient/Ambulatory care)
    	d = "M19 3H5c-1.1 0-1.99.9-1.99 2L3 19c0 1.1.9 2 2 2h14c1.1 0 2-.9 2-2V5c0-1.1-.9-2-2-2zm-1 11h-4v4h-2v-4H8v-2h4V8h2v4h4v2zM7 7c1.1 0 2 .9 2 2s-.9 2-2 2-2-.9-2-2 .9-2 2-2zm-2 9v-1c0-1.1.9-2 2-2s2 .9 2 2v1h-4z"
	case "hospital-outpatient":
		// Simple: Person + arrow (outpatient concept)
    	d = "M12 6c-1.1 0-2 .9-2 2s.9 2 2 2 2-.9 2-2-.9-2-2-2zm-2 9h4v2h-4zm6-8l-4 4 4 4 1.4-1.4L13.8 12H16v-2h-2.2l1.6-1.6L14 7z"
	case "hospital-alert-covid":
		d = "M12 2C6.47 2 2 6.47 2 12s4.47 10 10 10 10-4.47 10-10S17.53 2 12 2zm5 13h-2v2h-2v-2h-2v-2h2v-2h2v2h2v2zm-4.88-6.42L11 8.5V6h2v2.5l-.88.08zM7 13v-2h2.5l.08.88.08.12H7zm5 6v-2.5l.88-.08.12-.08V19h-2zm5-6h-2.5l-.08-.88-.08-.12H17v2z"
	case "spinner":
		d = "M12 4V2C6.48 2 2 6.48 2 12h2c0-4.41 3.59-8 8-8z"
	}

	classNames := fmt.Sprintf("icon icon-%s", name)
	if isSpinner {
		classNames += " icon-spin"
	}

	// Using app.Raw avoids the "undefined: app.Svg" compilation error
	// ViewBox stays 0 0 24 24. Width/Height handles the zoom.
	return app.Raw(fmt.Sprintf(
		`<svg fill="currentColor" class="%s" viewBox="0 0 24 24" width="%d" height="%d"><path d="%s" /></svg>`,
		classNames, size, size, d,
	))
}
