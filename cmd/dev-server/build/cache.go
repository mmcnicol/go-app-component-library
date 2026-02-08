// cmd/dev-server/build/cache.go
package build

import (
    "crypto/sha256"
    "encoding/hex"
    "encoding/json"
    "os"
    "path/filepath"
    "time"
)

type CacheEntry struct {
    OutputPath  string    `json:"output_path"`
    Hash        string    `json:"hash"`
    Dependencies []string `json:"dependencies"`
    Timestamp   time.Time `json:"timestamp"`
}

type BuildCache struct {
    cacheDir string
    entries  map[string]CacheEntry
}

func (c *BuildCache) IsValid(mainFile string, changedFiles []string) bool {
    entry, exists := c.entries[mainFile]
    if !exists {
        return false
    }
    
    // Check if any dependency changed
    currentHash := c.calculateHash(mainFile, changedFiles)
    return entry.Hash == currentHash && 
           time.Since(entry.Timestamp) < 5*time.Minute &&
           fileExists(entry.OutputPath)
}

