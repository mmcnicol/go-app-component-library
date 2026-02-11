//go:build dev
// pkg/components/viz/bar_chart_stories.go
package viz

import (
	"fmt"
	"strings"
	"github.com/maxence-charriere/go-app/v10/pkg/app"
	"github.com/mmcnicol/go-app-component-library/pkg/storybook"
)

func init() {
	// Sample datasets for the story
	datasets := map[string][]float64{
		"Monthly Sales":    {45000, 52000, 48000, 61000, 58000, 63000, 72000, 68000, 74000, 82000, 79000, 85000},
		"Browser Usage":    {65, 18, 12, 5},
		"Device Types":     {45, 30, 25},
		"Customer Ratings": {22, 45, 68, 92, 78, 34, 12},
		"Population Growth": {2.4, 2.7, 3.1, 3.5, 3.9, 4.2, 4.8, 5.1},
	}

	// Sample labels for each dataset
	datasetLabels := map[string][]string{
		"Monthly Sales":     {"Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"},
		"Browser Usage":     {"Chrome", "Firefox", "Safari", "Edge"},
		"Device Types":      {"Desktop", "Mobile", "Tablet"},
		"Customer Ratings":  {"1★", "2★", "3★", "4★", "5★", "6★", "7★"},
		"Population Growth": {"2016", "2017", "2018", "2019", "2020", "2021", "2022", "2023"},
	}

	// Custom color palettes for specific datasets
	customColors := map[string][]string{
		"Browser Usage": {
			"#4285F4", // Chrome blue
			"#FF7139", // Firefox orange
			"#1D1D1F", // Safari black
			"#0078D7", // Edge blue
		},
		"Device Types": {
			"#4f46e5", // Desktop - indigo
			"#10b981", // Mobile - emerald
			"#f59e0b", // Tablet - amber
		},
		"Customer Ratings": {
			"#ef4444", // 1★ - red
			"#f97316", // 2★ - orange
			"#eab308", // 3★ - yellow
			"#84cc16", // 4★ - lime
			"#22c55e", // 5★ - green
			"#14b8a6", // 6★ - teal
			"#3b82f6", // 7★ - blue
		},
	}

	storybook.Register("Visualization", "Bar Chart",
		map[string]*storybook.Control{
			"Dataset":           storybook.NewSelectControl([]string{
				"Monthly Sales", 
				"Browser Usage", 
				"Device Types", 
				"Customer Ratings",
				"Population Growth",
			}, "Monthly Sales"),
			"Chart Title":       storybook.NewTextControl("Sales Performance 2024"),
			"Bar Color":         storybook.NewColorControl("#10b981"),
			"Show Labels":       storybook.NewBoolControl(true),
			"Show Grid":         storybook.NewBoolControl(true),
			"Animated":          storybook.NewBoolControl(true),
			"Interactive":       storybook.NewBoolControl(true),
			"Bar Width":         storybook.NewRangeControl(30, 100, 5, 70),
			"Border Radius":     storybook.NewRangeControl(0, 20, 1, 4),
		},
		func(controls map[string]*storybook.Control) app.UI {
			selectedKey := controls["Dataset"].Value.(string)
			data := datasets[selectedKey]
			labels := datasetLabels[selectedKey]
			showLabels := controls["Show Labels"].Value.(bool)
			showGrid := controls["Show Grid"].Value.(bool)
			animated := controls["Animated"].Value.(bool)
			interactive := controls["Interactive"].Value.(bool)
			barWidth := float64(controls["Bar Width"].Value.(int))
			borderRadius := float64(controls["Border Radius"].Value.(int))

			// Convert data to viz.DataSet
			dataset := DataSet{
				Labels: labels,
				Series: []Series{
					{
						Label:  selectedKey,
						Points: Values(data),
						Fill:   true,
						Stroke: Stroke{
							Width: 1,
							Color: "#333333",
						},
						PointSize: 0, // No points for bar chart
					},
				},
			}

			// Build chart options
			spec := Spec{
				Type:  ChartTypeBar,
				Title: controls["Chart Title"].Value.(string),
				Data:  dataset,
				Theme: &CustomTheme{
					BaseTheme: DefaultTheme(),
					Colors: func() []string {
						// Use custom colors if defined for this dataset
						if colors, ok := customColors[selectedKey]; ok {
							return colors
						}
						// Otherwise use single color
						return []string{controls["Bar Color"].Value.(string)}
					}(),
				},
				Width:  800,
				Height: 400,
				Interactive: InteractiveConfig{
					Enabled: interactive,
					Tooltip: TooltipConfig{
						Enabled:   true,
						Mode:      TooltipModeNearest,
						Intersect: true,
						Format: func(p Point, s Series) string {
							if p.Label != "" {
								return p.Label
							}
							return s.Label
						},
						Background:  "rgba(0, 0, 0, 0.8)",
						TextColor:   "#ffffff",
						BorderColor: "#ffffff",
						BorderWidth: 1,
						Padding:     8,
					},
					Zoom: ZoomConfig{
						Enabled: false, // Disable zoom for bar charts
					},
					Pan: PanConfig{
						Enabled: false, // Disable pan for bar charts
					},
				},
				Accessible: AccessibilityConfig{
					Enabled: true,
					Description: func() string {
						return generateBarChartDescription(selectedKey, data, labels)
					}(),
				},
				Animated: animated,
			}

			// Configure axes
			spec.Axes = AxesConfig{
				X: AxisConfig{
					Visible:     true,
					Title:       "Category",
					TitleColor:  "#666666",
					LabelColor:  "#666666",
					Grid:        GridConfig{Visible: showGrid},
				},
				Y: AxisConfig{
					Visible:     true,
					Title:       "Value",
					TitleColor:  "#666666",
					LabelColor:  "#666666",
					Grid:        GridConfig{Visible: showGrid},
					BeginAtZero: true,
				},
			}

			// Configure bar appearance
			spec.Bar = BarConfig{
				Width:         barWidth,              // Percentage of available space
				BorderRadius:  borderRadius,          // Rounded corners
				Grouped:       false,                 // Single series
				Stacked:       false,
				Horizontal:    false,
			}

			// Configure labels
			if showLabels {
				spec.Labels = LabelsConfig{
					Visible:      true,
					Position:     "top",
					FontSize:     11,
					Color:        "#333333",
					Format:       func(value float64) string {
						if value >= 1000 {
							return fmt.Sprintf("%.0fk", value/1000)
						}
						return fmt.Sprintf("%.0f", value)
					},
				}
			}

			// Create the chart
			chart := New(spec)

			// Add event handlers if interactive
			if interactive {
				chart = chart.WithClickHandler(func(ctx app.Context, e app.Event, points []Point) {
					if len(points) > 0 {
						p := points[0]
						app.Log("Bar clicked:", p.Label, p.Y)
						
						// Show a temporary notification
						storybook.ShowNotification(ctx, 
							fmt.Sprintf("Selected: %s = %.0f", p.Label, p.Y),
							"info",
						)
					}
				}).WithHoverHandler(func(ctx app.Context, e app.Event, point *Point) {
					// Optional: handle hover events
					if point != nil {
						// Update status bar or other UI elements
						ctx.Dispatch(func(ctx app.Context) {
							// You could update a status component here
						})
					}
				})
			}

			// Return the chart with wrapper for styling
			return app.Div().ID("viz-bar-chart-container").Body(
				app.Div().Class("viz-chart-wrapper").Body(
					chart,
				),
				app.Div().Class("viz-chart-footer").Body(
					app.Small().Class("text-muted").Body(
						app.Text(fmt.Sprintf("%d data points • ", len(data))),
						app.If(interactive, app.Text("Interactive • ")),
						app.If(animated, app.Text("Animated • ")),
						app.Text("Accessible"),
					),
				),
			)
		},
	)
}

// Helper function to generate accessible description
func generateBarChartDescription(title string, data []float64, labels []string) string {
	var desc strings.Builder
	
	desc.WriteString(fmt.Sprintf("Bar chart showing %s. ", title))
	
	if len(data) > 0 {
		// Find min, max, average
		min, max := data[0], data[0]
		sum := 0.0
		for _, v := range data {
			if v < min { min = v }
			if v > max { max = v }
			sum += v
		}
		avg := sum / float64(len(data))
		
		desc.WriteString(fmt.Sprintf("Total of %d categories. ", len(data)))
		desc.WriteString(fmt.Sprintf("Highest value is %.0f, lowest is %.0f. ", max, min))
		desc.WriteString(fmt.Sprintf("Average value is %.1f. ", avg))
		
		// Add top 3 categories
		if len(data) >= 3 && len(labels) >= 3 {
			// Create slice of indices and sort by value
			indices := make([]int, len(data))
			for i := range indices { indices[i] = i }
			sort.Slice(indices, func(i, j int) bool {
				return data[indices[i]] > data[indices[j]]
			})
			
			desc.WriteString("Top categories: ")
			for i := 0; i < 3; i++ {
				idx := indices[i]
				desc.WriteString(fmt.Sprintf("%s at %.0f", labels[idx], data[idx]))
				if i < 2 { desc.WriteString(", ") }
			}
			desc.WriteString(".")
		}
	}
	
	return desc.String()
}
