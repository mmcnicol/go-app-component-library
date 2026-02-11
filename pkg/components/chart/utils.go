// pkg/components/chart/utils.go
package chart

import (
	"crypto/rand"
    "fmt"
	"math"
	"time"
)

func generateRandomData(count int) []DataPoint {
	var points []DataPoint
	for i := 0; i < count; i++ {
		points = append(points, DataPoint{
			X: float64(i),
			Y: float64(i) + (math.Sin(float64(i)*0.1) * 10),
		})
	}
	return points
}

// GenerateID creates a unique ID for chart elements
func GenerateID() string {
    b := make([]byte, 8)
    if _, err := rand.Read(b); err != nil {
        // Fallback to timestamp if crypto fails
        return fmt.Sprintf("%d", time.Now().UnixNano())
    }
    return fmt.Sprintf("%x", b)
}

// escapeJS properly escapes strings for JavaScript
func escapeJS(s string) string {
    // Simple escaping - replace backslashes and quotes
    result := ""
    for _, r := range s {
        switch r {
        case '\\':
            result += "\\\\"
        case '\'':
            result += "\\'"
        case '"':
            result += "\\\""
        case '\n':
            result += "\\n"
        case '\r':
            result += "\\r"
        case '\t':
            result += "\\t"
        default:
            result += string(r)
        }
    }
    return result
}
