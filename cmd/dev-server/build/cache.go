// cmd/dev-server/build/cache.go
package build

import (
    "crypto/sha256"
    "encoding/hex"
    "encoding/json"
    "fmt"
    "io"
    "os"
    "path/filepath"
    "sort"
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

// fileExists is a helper function
func fileExists(path string) bool {
    _, err := os.Stat(path)
    return err == nil
}

// calculateHash computes a unique hash for the build based on:
// 1. The content of the main file
// 2. The content of all changed files (dependencies)
// 3. The list of dependencies themselves
func (c *BuildCache) calculateHash(mainFile string, changedFiles []string) string {
    hasher := sha256.New()
    
    // 1. Add main file content to hash
    if err := addFileContentToHash(hasher, mainFile); err != nil {
        // If we can't read the main file, return empty hash
        return ""
    }
    
    // 2. Add all changed files content to hash
    // Sort to ensure consistent hash regardless of order
    sort.Strings(changedFiles)
    
    for _, file := range changedFiles {
        // Add the filename itself (dependency relationship matters)
        hasher.Write([]byte(file))
        
        // Add the file content
        if err := addFileContentToHash(hasher, file); err != nil {
            // If a dependency is missing, we should rebuild
            continue
        }
    }
    
    // 3. Add current timestamp (optional, for freshness)
    // This ensures builds aren't cached indefinitely
    hasher.Write([]byte(time.Now().Format("2006-01-02-15"))) // Hourly granularity
    
    // Return hex-encoded hash
    return hex.EncodeToString(hasher.Sum(nil))
}

// addFileContentToHash reads a file and adds its content to the hash
func addFileContentToHash(hasher io.Writer, filePath string) error {
    file, err := os.Open(filePath)
    if err != nil {
        return fmt.Errorf("failed to open file %s: %v", filePath, err)
    }
    defer file.Close()
    
    // Get file info for size (optional)
    info, err := file.Stat()
    if err != nil {
        return fmt.Errorf("failed to stat file %s: %v", filePath, err)
    }
    
    // Add file size to hash (quick check)
    fmt.Fprintf(hasher, "%d:", info.Size())
    
    // Copy file content to hasher
    _, err = io.Copy(hasher, file)
    if err != nil {
        return fmt.Errorf("failed to read file %s: %v", filePath, err)
    }
    
    return nil
}
