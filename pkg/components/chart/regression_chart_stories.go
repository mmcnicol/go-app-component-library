//go:build dev
// pkg/components/chart/regression_chart_stories.go
package chart

import (
	"github.com/maxence-charriere/go-app/v10/pkg/app"
	"github.com/mmcnicol/go-app-component-library/pkg/storybook"
)

func init() {
	// Sample datasets with different correlation strengths
	datasets := map[string][]Point{
		"Strong Positive": {
			{X: 1, Y: 12}, {X: 2, Y: 18}, {X: 3, Y: 25}, {X: 4, Y: 32},
			{X: 5, Y: 38}, {X: 6, Y: 45}, {X: 7, Y: 52}, {X: 8, Y: 58},
			{X: 9, Y: 65}, {X: 10, Y: 72},
		},
		"Moderate Positive": {
			{X: 1, Y: 15}, {X: 2, Y: 22}, {X: 3, Y: 28}, {X: 4, Y: 30},
			{X: 5, Y: 42}, {X: 6, Y: 48}, {X: 7, Y: 55}, {X: 8, Y: 60},
			{X: 9, Y: 63}, {X: 10, Y: 70},
		},
		"Negative Correlation": {
			{X: 1, Y: 85}, {X: 2, Y: 78}, {X: 3, Y: 72}, {X: 4, Y: 65},
			{X: 5, Y: 58}, {X: 6, Y: 52}, {X: 7, Y: 45}, {X: 8, Y: 38},
			{X: 9, Y: 32}, {X: 10, Y: 25},
		},
		"Weak Correlation": {
			{X: 1, Y: 45}, {X: 2, Y: 52}, {X: 3, Y: 38}, {X: 4, Y: 65},
			{X: 5, Y: 42}, {X: 6, Y: 58}, {X: 7, Y: 55}, {X: 8, Y: 48},
			{X: 9, Y: 62}, {X: 10, Y: 51},
		},
		"With Outliers": {
			{X: 1, Y: 15}, {X: 2, Y: 22}, {X: 3, Y: 28}, {X: 4, Y: 32},
			{X: 5, Y: 95}, // Outlier
			{X: 6, Y: 45}, {X: 7, Y: 52}, {X: 8, Y: 58},
			{X: 9, Y: 12}, // Outlier
			{X: 10, Y: 70},
		},
	}

	storybook.Register("Charts", "Regression Chart",
		map[string]*storybook.Control{
			"Dataset":        storybook.NewSelectControl(
				[]string{"Strong Positive", "Moderate Positive", "Negative Correlation", "Weak Correlation", "With Outliers"}, 
				"Strong Positive",
			),
			"Point Color":    storybook.NewColorControl("#4f46e5"),
			"Line Color":     storybook.NewColorControl("#ef4444"),
			"Show Equation":  storybook.NewBoolControl(true),
			"Title":          storybook.NewTextControl("Linear Regression Analysis"),
			"Point Size":     storybook.NewRangeControl(2, 10, 1, 4),
		},
		func(controls map[string]*storybook.Control) app.UI {
			selectedKey := controls["Dataset"].Value.(string)
			data := datasets[selectedKey]
			
			// Use the constructor instead of manual struct creation
			regressionChart := NewRegressionChart(data,
				WithTitle(controls["Title"].Value.(string)),
				WithColor(controls["Point Color"].Value.(string)),
				func(c *ChartConfig) {
					c.BoxData = nil
					c.PieData = nil
					c.HeatmapMatrix = nil
					c.IsStream = false
					c.Thickness = 2.0
				},
			).
			WithPointColor(controls["Point Color"].Value.(string)).
			WithLineColor(controls["Line Color"].Value.(string)).
			WithShowEquation(controls["Show Equation"].Value.(bool)).
			WithPointSize(float64(controls["Point Size"].Value.(int))).
			WithDatasetName(selectedKey)
			
			return app.Div().ID("regression-chart-container").Body(
				app.Div().Class("regression-wrapper").Body(
					regressionChart,
				),
			)
		},
	)
}
