// pkg/components/vizzy/core_registry.go
package core

import (
    "fmt"
    "sync"
)

// EngineFactory creates a chart engine instance
type EngineFactory func() ChartEngine

var (
    engines     = make(map[ChartType]map[EngineType]EngineFactory)
    enginesMu   sync.RWMutex
)

// RegisterEngine registers a rendering engine for a chart type
func RegisterEngine(chartType ChartType, engineType EngineType, factory EngineFactory) {
    enginesMu.Lock()
    defer enginesMu.Unlock()
    
    if _, ok := engines[chartType]; !ok {
        engines[chartType] = make(map[EngineType]EngineFactory)
    }
    engines[chartType][engineType] = factory
}

// GetEngine returns an appropriate engine for the chart
func GetEngine(chartType ChartType, preferred EngineType) (ChartEngine, error) {
    enginesMu.RLock()
    defer enginesMu.RUnlock()
    
    // Try preferred engine
    if chartEngines, ok := engines[chartType]; ok {
        if factory, ok := chartEngines[preferred]; ok {
            return factory(), nil
        }
        
        // Fallback to canvas
        if factory, ok := chartEngines[EngineTypeCanvas]; ok {
            return factory(), nil
        }
    }
    
    return nil, fmt.Errorf("no engine found for chart type: %s", chartType)
}
