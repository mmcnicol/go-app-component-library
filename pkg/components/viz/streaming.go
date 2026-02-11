// pkg/components/viz/streaming.go
//go:build !dev
package viz

import (
    "sync"
    "time"
    "github.com/maxence-charriere/go-app/v10/pkg/app"
)

// DataBuffer implements a circular buffer for streaming data
type DataBuffer struct {
    data      []Point
    maxSize   int
    head      int
    tail      int
    count     int
    mu        sync.RWMutex
}

// NewDataBuffer creates a new circular buffer for streaming data
func NewDataBuffer(maxSize int) *DataBuffer {
    return &DataBuffer{
        data:    make([]Point, maxSize),
        maxSize: maxSize,
    }
}

// Add adds a point to the buffer
func (db *DataBuffer) Add(point Point) {
    db.mu.Lock()
    defer db.mu.Unlock()
    
    db.data[db.head] = point
    db.head = (db.head + 1) % db.maxSize
    
    if db.count < db.maxSize {
        db.count++
    } else {
        db.tail = (db.tail + 1) % db.maxSize
    }
}

// GetSlice returns all points in the buffer
func (db *DataBuffer) GetSlice() []Point {
    db.mu.RLock()
    defer db.mu.RUnlock()
    
    if db.count == 0 {
        return []Point{}
    }
    
    result := make([]Point, db.count)
    
    if db.head > db.tail || db.count < db.maxSize {
        // Data is contiguous
        copy(result, db.data[db.tail:db.head])
    } else {
        // Data wraps around
        firstPart := db.data[db.tail:db.maxSize]
        secondPart := db.data[0:db.head]
        
        copy(result, firstPart)
        copy(result[len(firstPart):], secondPart)
    }
    
    return result
}

// StreamingChart represents a chart with streaming data capabilities
type StreamingChart struct {
    *Chart
    dataBuffer  *DataBuffer
    maxPoints   int
    updateRate  time.Duration
    streaming   bool
    lastUpdate  time.Time
}

// NewStreamingChart creates a new streaming chart
func NewStreamingChart(chartType ChartType) *StreamingChart {
    return &StreamingChart{
        Chart:      New(Spec{Type: chartType}),
        maxPoints:  100,
        updateRate: 100 * time.Millisecond,
        dataBuffer: NewDataBuffer(100),
    }
}

// WithMaxPoints sets the maximum number of points to display
func (sc *StreamingChart) WithMaxPoints(n int) *StreamingChart {
    sc.maxPoints = n
    sc.dataBuffer = NewDataBuffer(n)
    return sc
}

// WithUpdateRate sets the update rate for streaming
func (sc *StreamingChart) WithUpdateRate(rate time.Duration) *StreamingChart {
    sc.updateRate = rate
    return sc
}

// StreamData starts streaming data from a channel
func (sc *StreamingChart) StreamData(dataSource <-chan Point) {
    sc.streaming = true
    
    go func() {
        for point := range dataSource {
            sc.dataBuffer.Add(point)
            
            // Throttle updates
            if time.Since(sc.lastUpdate) >= sc.updateRate {
                sc.updateChart()
                sc.lastUpdate = time.Now()
            }
        }
    }()
}

func (sc *StreamingChart) updateChart() {
    // Get latest data
    data := sc.dataBuffer.GetSlice()
    
    // Update dataset
    if len(sc.spec.Data.Series) > 0 {
        sc.spec.Data.Series[0].Points = data
    }
    
    // Trigger re-render
    if sc.engine != nil {
        sc.engine.Update(sc.spec.Data)
    }
}

// Pause stops streaming updates
func (sc *StreamingChart) Pause() {
    sc.streaming = false
}

// Resume resumes streaming updates
func (sc *StreamingChart) Resume() {
    sc.streaming = true
}
