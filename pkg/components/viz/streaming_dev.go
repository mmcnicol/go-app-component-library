// pkg/components/viz/streaming_dev.go
//go:build dev

package viz

import (
    "sync"
    "time"
    "github.com/maxence-charriere/go-app/v10/pkg/app"
)

// DataBuffer for dev builds
type DataBuffer struct {
    data      []Point
    maxSize   int
    head      int
    tail      int
    count     int
    mu        sync.RWMutex
}

func NewDataBuffer(maxSize int) *DataBuffer {
    return &DataBuffer{
        data:    make([]Point, maxSize),
        maxSize: maxSize,
    }
}

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

func (db *DataBuffer) GetSlice() []Point {
    db.mu.RLock()
    defer db.mu.RUnlock()
    
    if db.count == 0 {
        return []Point{}
    }
    
    result := make([]Point, db.count)
    
    if db.head > db.tail || db.count < db.maxSize {
        copy(result, db.data[db.tail:db.head])
    } else {
        firstPart := db.data[db.tail:db.maxSize]
        secondPart := db.data[0:db.head]
        copy(result, firstPart)
        copy(result[len(firstPart):], secondPart)
    }
    
    return result
}

// StreamingChart for dev builds
type StreamingChart struct {
    *Chart
    dataBuffer  *DataBuffer
    maxPoints   int
    updateRate  time.Duration
    streaming   bool
    lastUpdate  time.Time
}

func NewStreamingChart(chartType ChartType) *StreamingChart {
    return &StreamingChart{
        Chart:      New(Spec{Type: chartType}),
        maxPoints:  100,
        updateRate: 100 * time.Millisecond,
        dataBuffer: NewDataBuffer(100),
    }
}

func (sc *StreamingChart) WithMaxPoints(n int) *StreamingChart {
    sc.maxPoints = n
    sc.dataBuffer = NewDataBuffer(n)
    return sc
}

func (sc *StreamingChart) WithUpdateRate(rate time.Duration) *StreamingChart {
    sc.updateRate = rate
    return sc
}

func (sc *StreamingChart) StreamData(dataSource <-chan Point) {
    sc.streaming = true
    
    go func() {
        for point := range dataSource {
            sc.dataBuffer.Add(point)
            
            if time.Since(sc.lastUpdate) >= sc.updateRate {
                sc.updateChart()
                sc.lastUpdate = time.Now()
            }
        }
    }()
}

func (sc *StreamingChart) updateChart() {
    data := sc.dataBuffer.GetSlice()
    
    if len(sc.spec.Data.Series) > 0 {
        sc.spec.Data.Series[0].Points = data
    }
    
    if sc.engine != nil {
        sc.engine.Update(sc.spec.Data)
    }
}

func (sc *StreamingChart) Pause() {
    sc.streaming = false
}

func (sc *StreamingChart) Resume() {
    sc.streaming = true
}

// Ensure StreamingChart implements app.UI
func (sc *StreamingChart) Render() app.UI {
    return sc.Chart.Render()
}
