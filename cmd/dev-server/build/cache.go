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
    "sync"
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
    mu       sync.RWMutex // Add mutex for thread-safe access
}

// NewBuildCache creates a new build cache instance
// If cacheDir is empty, it uses a default directory (~/.cache/{app-name}/build)
func NewBuildCache(cacheDir string) (*BuildCache, error) {
    // If no cache directory provided, use default
    if cacheDir == "" {
        userCacheDir, err := os.UserCacheDir()
        if err != nil {
            // Fallback to temp directory
            userCacheDir = os.TempDir()
        }
        cacheDir = filepath.Join(userCacheDir, "go-app-dev-server", "build-cache")
    }
    
    // Create cache directory if it doesn't exist
    if err := os.MkdirAll(cacheDir, 0755); err != nil {
        return nil, fmt.Errorf("failed to create cache directory: %v", err)
    }
    
    cache := &BuildCache{
        cacheDir: cacheDir,
        entries:  make(map[string]CacheEntry),
    }
    
    // Load existing cache entries
    if err := cache.load(); err != nil {
        // Don't fail if cache can't be loaded, just log and start fresh
        fmt.Printf("Warning: Could not load cache: %v\n", err)
    }
    
    // Start periodic cleanup goroutine
    go cache.periodicCleanup()
    
    return cache, nil
}

// load reads cache entries from disk
func (c *BuildCache) load() error {
    cacheFile := filepath.Join(c.cacheDir, "cache.json")
    
    // Check if cache file exists
    if _, err := os.Stat(cacheFile); os.IsNotExist(err) {
        // No cache file, start with empty cache
        return nil
    }
    
    // Read cache file
    data, err := os.ReadFile(cacheFile)
    if err != nil {
        return fmt.Errorf("failed to read cache file: %v", err)
    }
    
    c.mu.Lock()
    defer c.mu.Unlock()
    
    // Unmarshal JSON
    if err := json.Unmarshal(data, &c.entries); err != nil {
        // If JSON is corrupted, start fresh
        c.entries = make(map[string]CacheEntry)
        return fmt.Errorf("cache file corrupted: %v", err)
    }
    
    // Clean up any stale entries on load
    c.cleanupStaleEntries()
    
    return nil
}

// save writes cache entries to disk
func (c *BuildCache) save() error {
    c.mu.RLock()
    defer c.mu.RUnlock()
    
    // Create cache directory if it doesn't exist
    if err := os.MkdirAll(c.cacheDir, 0755); err != nil {
        return err
    }
    
    cacheFile := filepath.Join(c.cacheDir, "cache.json")
    
    // Marshal entries to JSON
    data, err := json.MarshalIndent(c.entries, "", "  ")
    if err != nil {
        return fmt.Errorf("failed to marshal cache: %v", err)
    }
    
    // Write to temp file first, then rename (atomic operation)
    tempFile := cacheFile + ".tmp"
    if err := os.WriteFile(tempFile, data, 0644); err != nil {
        return fmt.Errorf("failed to write cache file: %v", err)
    }
    
    // Atomic rename
    if err := os.Rename(tempFile, cacheFile); err != nil {
        return fmt.Errorf("failed to rename cache file: %v", err)
    }
    
    return nil
}

// Set adds or updates a cache entry
func (c *BuildCache) Set(mainFile string, entry CacheEntry) error {
    c.mu.Lock()
    c.entries[mainFile] = entry
    c.mu.Unlock()
    
    // Save to disk (could be done async in production)
    if err := c.save(); err != nil {
        return fmt.Errorf("failed to save cache: %v", err)
    }
    
    return nil
}

// Get retrieves a cache entry
func (c *BuildCache) Get(mainFile string) (CacheEntry, bool) {
    c.mu.RLock()
    defer c.mu.RUnlock()
    
    entry, exists := c.entries[mainFile]
    return entry, exists
}

// Clear removes all cache entries
func (c *BuildCache) Clear() error {
    c.mu.Lock()
    c.entries = make(map[string]CacheEntry)
    c.mu.Unlock()
    
    // Remove cache file
    cacheFile := filepath.Join(c.cacheDir, "cache.json")
    if err := os.Remove(cacheFile); err != nil && !os.IsNotExist(err) {
        return err
    }
    
    return nil
}

// ClearEntry removes a specific cache entry
func (c *BuildCache) ClearEntry(mainFile string) {
    c.mu.Lock()
    delete(c.entries, mainFile)
    c.mu.Unlock()
    
    // Save asynchronously
    go c.save()
}

// cleanupStaleEntries removes entries older than the cache duration
func (c *BuildCache) cleanupStaleEntries() {
    cutoff := time.Now().Add(-5 * time.Minute) // Same as IsValid check
    
    for key, entry := range c.entries {
        if entry.Timestamp.Before(cutoff) || !fileExists(entry.OutputPath) {
            delete(c.entries, key)
        }
    }
}

// periodicCleanup runs periodic cleanup of stale entries
func (c *BuildCache) periodicCleanup() {
    ticker := time.NewTicker(15 * time.Minute)
    defer ticker.Stop()
    
    for range ticker.C {
        c.mu.Lock()
        c.cleanupStaleEntries()
        c.mu.Unlock()
        
        // Save after cleanup
        if err := c.save(); err != nil {
            fmt.Printf("Warning: Failed to save cache during cleanup: %v\n", err)
        }
    }
}

// GetCacheDir returns the cache directory path
func (c *BuildCache) GetCacheDir() string {
    return c.cacheDir
}

// GetEntryCount returns the number of cache entries
func (c *BuildCache) GetEntryCount() int {
    c.mu.RLock()
    defer c.mu.RUnlock()
    return len(c.entries)
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
