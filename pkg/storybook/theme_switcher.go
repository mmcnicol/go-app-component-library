// pkg/storybook/theme_switcher.go
package storybook

import "github.com/maxence-charriere/go-app/v10/pkg/app"

type ThemeSwitcher struct {
	app.Compo
	IsDark   bool
	OnChange func(isDark bool)
}

func (s *ThemeSwitcher) Render() app.UI {
	icon := "☀️"
	label := "Light Mode"
	if s.IsDark {
		icon = "🌙"
		label = "Dark Mode"
	}

	return app.Div().
		Style("margin-bottom", "20px").
		Style("padding", "8px 12px").
		Style("border", "1px solid var(--theme-border)").
		Style("border-radius", "6px").
		Style("cursor", "pointer").
		Style("display", "flex").
		Style("align-items", "center").
		Style("justify-content", "space-between").
		Style("font-size", "14px").
		OnClick(func(ctx app.Context, e app.Event) {
			s.IsDark = !s.IsDark
			if s.OnChange != nil {
				s.OnChange(s.IsDark)
			}
		}).
		Body(
			app.Span().Text(label),
			app.Span().Style("font-size", "18px").Text(icon),
		)
}
