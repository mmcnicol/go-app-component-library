// pkg/components/viz/line_chart_stories.go
//go:build dev

package viz

import (
    "fmt"
    "math"
    mathRand "math/rand"
    "time"
    
    "github.com/maxence-charriere/go-app/v10/pkg/app"
    "github.com/mmcnicol/go-app-component-library/pkg/storybook"
)

// Generate shared datasets once
var (
    sinData   []Point
    cosData   []Point
    sinLabels []string
)

func init() {
    // Initialize shared datasets
    sinData = generateSineWave(50, 0, 4*math.Pi, 1.0)
    cosData = make([]Point, len(sinData))
    for i := range sinData {
        cosData[i] = Point{
            X: sinData[i].X,
            Y: math.Cos(sinData[i].X),
            Label: sinData[i].Label,
        }
    }
    
    sinLabels = make([]string, len(sinData))
    for i := range sinLabels {
        sinLabels[i] = fmt.Sprintf("%.1f", sinData[i].X)
    }
    
    // 1. Basic Line Chart Story
    storybook.Register("Visualization", "Line Chart - Basic",
        map[string]*storybook.Control{
            "Title":      storybook.NewTextControl("Sine Wave"),
            "Show Grid":  storybook.NewBoolControl(true),
            "Show Points": storybook.NewBoolControl(false),
            "Animated":   storybook.NewBoolControl(true),
            "Interactive": storybook.NewBoolControl(true),
        },
        func(controls map[string]*storybook.Control) app.UI {
            showGrid := controls["Show Grid"].Value.(bool)
            showPoints := controls["Show Points"].Value.(bool)
            
            pointSize := 0
            if showPoints {
                pointSize = 4
            }
            
            dataset := DataSet{
                Labels: sinLabels,
                Series: []Series{
                    {
                        Label:     "sin(x)",
                        Points:    sinData,
                        Color:     "#4f46e5",
                        Stroke:    Stroke{Width: 2},
                        PointSize: pointSize,
                        Fill:      false,
                    },
                },
            }
            
            spec := Spec{
                Type:   ChartTypeLine,
                Title:  controls["Title"].Value.(string),
                Data:   dataset,
                Width:  800,
                Height: 400,
                Theme:  &CustomTheme{BaseTheme: DefaultTheme()},
                Axes: AxesConfig{
                    X: AxisConfig{
                        Visible: true,
                        Title:   "x",
                        Grid:    GridConfig{Visible: showGrid},
                    },
                    Y: AxisConfig{
                        Visible:     true,
                        Title:       "sin(x)",
                        Grid:        GridConfig{Visible: showGrid},
                        BeginAtZero: false,
                    },
                },
                Interactive: InteractiveConfig{
                    Enabled: controls["Interactive"].Value.(bool),
                    Tooltip: TooltipConfig{
                        Enabled: true,
                        Mode:    TooltipModeNearest,
                        Format: func(p Point, s Series) string {
                            return fmt.Sprintf("%s: (%.2f, %.2f)", s.Label, p.X, p.Y)
                        },
                    },
                },
                Animated: controls["Animated"].Value.(bool),
            }
            
            chart := New(spec)
            
            return app.Div().ID("viz-line-basic-container").Body(
                app.Div().Class("viz-chart-wrapper").Body(chart),
                app.Div().Class("viz-chart-footer").Body(
                    app.Small().Class("text-muted").Body(
                        app.Text(fmt.Sprintf("%d data points", len(sinData))),
                    ),
                ),
            )
        },
    )
    
    // 2. Multi-Series Line Chart
    storybook.Register("Visualization", "Line Chart - Multi Series",
        map[string]*storybook.Control{
            "Title":       storybook.NewTextControl("Sine and Cosine"),
            "Show Legend": storybook.NewBoolControl(true),
            "Smooth Curves": storybook.NewBoolControl(true),
            "Tension":     storybook.NewRangeControl(0, 100, 10, 30), // 0-100% -> 0-0.3
        },
        func(controls map[string]*storybook.Control) app.UI {
            tension := float64(controls["Tension"].Value.(int)) / 100 * 0.3
            
            dataset := DataSet{
                Labels: sinLabels,
                Series: []Series{
                    {
                        Label:     "sin(x)",
                        Points:    sinData,
                        Color:     "#4f46e5",
                        Stroke:    Stroke{Width: 2},
                        Tension:   tension,
                        PointSize: 0,
                    },
                    {
                        Label:     "cos(x)",
                        Points:    cosData,
                        Color:     "#10b981",
                        Stroke:    Stroke{Width: 2},
                        Tension:   tension,
                        PointSize: 0,
                    },
                },
            }
            
            spec := Spec{
                Type:   ChartTypeLine,
                Title:  controls["Title"].Value.(string),
                Data:   dataset,
                Width:  800,
                Height: 400,
                Theme:  &CustomTheme{BaseTheme: DefaultTheme()},
                Legend: LegendConfig{
                    Visible:  controls["Show Legend"].Value.(bool),
                    Position: "top",
                },
            }
            
            chart := New(spec)
            
            return app.Div().ID("viz-line-multi-container").Body(
                app.Div().Class("viz-chart-wrapper").Body(chart),
            )
        },
    )
    
    // 3. Filled Area Chart
    storybook.Register("Visualization", "Area Chart",
        map[string]*storybook.Control{
            "Title":      storybook.NewTextControl("Area Under Curve"),
            "Opacity":    storybook.NewRangeControl(10, 100, 10, 25),
        },
        func(controls map[string]*storybook.Control) app.UI {
            opacity := float64(controls["Opacity"].Value.(int)) / 100
            
            // Generate data with baseline
            data := generateSineWave(50, 0, 4*math.Pi, 1.0)
            
            // Create color with opacity
            baseColor := "#4f46e5"
            // Convert opacity to hex (approximate)
            opacityHex := fmt.Sprintf("%02x", int(opacity*255))
            fillColor := baseColor + opacityHex
            
            dataset := DataSet{
                Labels: sinLabels,
                Series: []Series{
                    {
                        Label:  "sin(x)",
                        Points: data,
                        Color:  baseColor,
                        Stroke: Stroke{Width: 2, Color: baseColor},
                        Fill:   true,
                    },
                },
            }
            
            // Override theme colors to include opacity
            theme := &CustomTheme{
                BaseTheme: DefaultTheme(),
                Colors:    []string{fillColor},
            }
            
            spec := Spec{
                Type:  ChartTypeLine,
                Title: controls["Title"].Value.(string),
                Data:  dataset,
                Theme: theme,
                Width: 800,
                Height: 400,
                Axes: AxesConfig{
                    Y: AxisConfig{
                        BeginAtZero: true,
                    },
                },
            }
            
            chart := New(spec)
            
            return app.Div().ID("viz-area-container").Body(
                app.Div().Class("viz-chart-wrapper").Body(chart),
                app.Div().Class("viz-chart-footer").Body(
                    app.Small().Class("text-muted").Text(
                        fmt.Sprintf("Fill opacity: %.0f%%", opacity*100),
                    ),
                ),
            )
        },
    )
    
    // 4. Point Styles Chart
    storybook.Register("Visualization", "Line Chart - Point Styles",
        map[string]*storybook.Control{
            "Title":       storybook.NewTextControl("Point Markers"),
            "Point Size":  storybook.NewRangeControl(2, 12, 1, 6),
            "Point Style": storybook.NewEnumControl("style", 
                []string{"circle", "square", "triangle", "diamond", "cross"}, 
                "circle"),
        },
        func(controls map[string]*storybook.Control) app.UI {
            pointSize := controls["Point Size"].Value.(int)
            styleStr := controls["Point Style"].Value.(string)
            
            var pointStyle PointStyle
            switch styleStr {
            case "square":
                pointStyle = PointStyleSquare
            case "triangle":
                pointStyle = PointStyleTriangle
            case "diamond":
                pointStyle = PointStyleDiamond
            case "cross":
                pointStyle = PointStyleCross
            default:
                pointStyle = PointStyleCircle
            }
            
            // Generate sparse data for clear point visibility
            data := generateSineWave(15, 0, 2*math.Pi, 1.0)
            
            dataset := DataSet{
                Series: []Series{
                    {
                        Label:      "Sample Data",
                        Points:     data,
                        Color:      "#8b5cf6",
                        Stroke:     Stroke{Width: 1.5, Color: "#8b5cf6", Dash: []float64{5, 3}},
                        PointSize:  pointSize,
                        PointStyle: pointStyle,
                    },
                },
            }
            
            spec := Spec{
                Type:  ChartTypeLine,
                Title: controls["Title"].Value.(string),
                Data:  dataset,
                Width: 800,
                Height: 400,
            }
            
            chart := New(spec)
            
            return app.Div().ID("viz-points-container").Body(
                app.Div().Class("viz-chart-wrapper").Body(chart),
            )
        },
    )
    
    // 5. Large Dataset Performance
    storybook.Register("Visualization", "Line Chart - Large Dataset",
        map[string]*storybook.Control{
            "Title":        storybook.NewTextControl("50,000 Points"),
            "Sampling":     storybook.NewEnumControl("strategy",
                []string{"LTTB", "EveryNth", "MinMax", "Average"},
                "LTTB"),
            "Max Points":   storybook.NewRangeControl(100, 5000, 100, 1000),
        },
        func(controls map[string]*storybook.Control) app.UI {
            // Generate 50,000 noisy sine wave points
            data := generateNoisySineWave(50000, 0, 100*math.Pi, 1.0, 0.2)
            
            sampling := controls["Sampling"].Value.(string)
            var strategy SamplingStrategy
            switch sampling {
            case "EveryNth":
                strategy = SamplingStrategyEveryNth
            case "MinMax":
                strategy = SamplingStrategyMinMax
            case "Average":
                strategy = SamplingStrategyAverage
            default:
                strategy = SamplingStrategyLTTB
            }
            
            dataset := DataSet{
                Series: []Series{
                    {
                        Label:  "Noisy Signal",
                        Points: data,
                        Color:  "#ef4444",
                        Stroke: Stroke{Width: 1},
                    },
                },
            }
            
            spec := Spec{
                Type:      ChartTypeLine,
                Title:     controls["Title"].Value.(string),
                Data:      dataset,
                Width:     800,
                Height:    400,
                MaxPoints: controls["Max Points"].Value.(int),
                Sampling:  strategy,
            }
            
            chart := New(spec)
            
            return app.Div().ID("viz-large-container").Body(
                app.Div().Class("viz-chart-wrapper").Body(chart),
                app.Div().Class("viz-chart-footer").Body(
                    app.Small().Class("text-muted").Text(
                        fmt.Sprintf("50,000 points downsampled to %d using %s",
                            controls["Max Points"].Value.(int), sampling),
                    ),
                ),
            )
        },
    )
    
    // 6. Streaming / Real-time Chart
    storybook.Register("Visualization", "Line Chart - Streaming",
		map[string]*storybook.Control{
			"Title":      storybook.NewTextControl("Real-time Data"),
			"Max Points": storybook.NewRangeControl(20, 200, 10, 50),
			"Update Rate": storybook.NewRangeControl(10, 200, 10, 50),
		},
		func(controls map[string]*storybook.Control) app.UI {
			maxPoints := controls["Max Points"].Value.(int)
			updateRate := controls["Update Rate"].Value.(int)
			
			// Create streaming chart
			streamingChart := NewStreamingChart(ChartTypeLine).
				WithMaxPoints(maxPoints).
				WithUpdateRate(time.Duration(updateRate) * time.Millisecond)
			
			// Create a data channel
			dataChan := make(chan Point, 100)
			
			// Start streaming data in a goroutine
			go func() {
				t := 0.0
				for {
					t += 0.1
					point := Point{
						X: t,
						Y: math.Sin(t) + 0.5*math.Sin(t*3) + 0.2*math.Sin(t*10),
					}
					dataChan <- point
					time.Sleep(50 * time.Millisecond)
				}
			}()
			
			// Start the stream
			streamingChart.StreamData(dataChan)
			
			return app.Div().ID("viz-streaming-container").Body(
				app.Div().Class("viz-chart-wrapper").Body(streamingChart),
				app.Div().Class("viz-chart-footer").Body(
					app.Small().Class("text-muted").Text(
						fmt.Sprintf("Streaming %d points, update rate %dms", 
							maxPoints, updateRate),
					),
				),
			)
		},
	)
}

// Helper functions
func generateSineWave(n int, start, end, amplitude float64) []Point {
    points := make([]Point, n)
    step := (end - start) / float64(n-1)
    
    for i := 0; i < n; i++ {
        x := start + float64(i)*step
        points[i] = Point{
            X:     x,
            Y:     amplitude * math.Sin(x),
            Label: fmt.Sprintf("%.2f", x),
        }
    }
    return points
}

func generateNoisySineWave(n int, start, end, amplitude, noiseLevel float64) []Point {
    points := make([]Point, n)
    step := (end - start) / float64(n-1)
    
    for i := 0; i < n; i++ {
        x := start + float64(i)*step
        noise := (mathRand.Float64()*2 - 1) * noiseLevel * amplitude
        points[i] = Point{
            X: x,
            Y: amplitude*math.Sin(x) + noise,
        }
    }
    return points
}
