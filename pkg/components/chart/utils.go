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
