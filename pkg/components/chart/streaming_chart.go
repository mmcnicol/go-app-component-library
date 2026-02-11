// pkg/components/chart/streaming_chart.go
package chart

import (
    "sync"
	"time"
)

type StreamingChart struct {
    BaseChart
    dataBuffer   *DataBuffer
    maxPoints    int
    streaming    bool
    updateRate   time.Duration
    lastUpdate   time.Time
}

type DataBuffer struct {
    data      []DataPoint
    maxSize   int
    head      int
    tail      int
    count     int
    mu        sync.RWMutex
}

func NewDataBuffer(maxSize int) *DataBuffer {
    return &DataBuffer{
        data:    make([]DataPoint, maxSize),
        maxSize: maxSize,
    }
}

func (db *DataBuffer) Add(point DataPoint) {
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

func (db *DataBuffer) GetSlice() []DataPoint {
    db.mu.RLock()
    defer db.mu.RUnlock()
    
    if db.count == 0 {
        return []DataPoint{}
    }
    
    result := make([]DataPoint, db.count)
    
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

func (sc *StreamingChart) StreamData(dataSource <-chan DataPoint) {
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
    sc.spec.Data.Datasets[0].Data = data
    
    // Trigger incremental update
    sc.engine.Update(sc.spec.Data)
}
