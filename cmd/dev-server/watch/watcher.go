// cmd/dev-server/watch/watcher.go
package watch

import (
    "log"
    "os"
    "path/filepath"
    "strings"
    "time"
    
    "github.com/fsnotify/fsnotify"
)

type Watcher struct {
    watcher      *fsnotify.Watcher
    onChange     func(changedFiles []string)
    ignoredDirs  map[string]bool
    extensions   map[string]bool
    debounceTime time.Duration
    changes      chan []string
}

func NewWatcher(sourceDir string, onChange func([]string)) (*Watcher, error) {
    fswatcher, err := fsnotify.NewWatcher()
    if err != nil {
        return nil, err
    }
    
    w := &Watcher{
        watcher:      fswatcher,
        onChange:     onChange,
        ignoredDirs:  map[string]bool{".git": true, "node_modules": true},
        extensions:   map[string]bool{".go": true, ".css": true, ".html": true},
        debounceTime: 100 * time.Millisecond,
        changes:      make(chan []string, 100),
    }
    
    // Recursively watch directories
    w.watchRecursive(sourceDir)
    
    // Start processing changes
    go w.processChanges()
    
    return w, nil
}

func (w *Watcher) watchRecursive(dir string) error {
    return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }
        
        // Skip ignored directories
        if info.IsDir() {
            dirName := filepath.Base(path)
            if w.ignoredDirs[dirName] {
                return filepath.SkipDir
            }
            return w.watcher.Add(path)
        }
        
        // Only watch relevant file extensions
        ext := filepath.Ext(path)
        if w.extensions[ext] {
            // Track this file
        }
        
        return nil
    })
}

// shouldProcessEvent determines whether a filesystem event should trigger a rebuild
func (w *Watcher) shouldProcessEvent(event fsnotify.Event) bool {
    // Filter by file operation type
    if !isRelevantFileOperation(event.Op) {
        return false
    }
    
    // Check if the file/directory is in ignored paths
    if w.isIgnoredPath(event.Name) {
        return false
    }
    
    // Check for file extension (skip non-source files)
    ext := filepath.Ext(event.Name)
    if !w.extensions[ext] && !isDirEvent(event.Name) {
        return false
    }
    
    // Skip temporary/editor backup files
    if isTemporaryFile(event.Name) {
        return false
    }
    
    return true
}

// isRelevantFileOperation checks if the event operation is relevant for rebuilding
func isRelevantFileOperation(op fsnotify.Op) bool {
    // We care about writes, creates, renames, and removes
    // Chmod events are usually not relevant for rebuilding
    return op&(fsnotify.Write|fsnotify.Create|fsnotify.Remove|fsnotify.Rename) != 0
}

// isIgnoredPath checks if a path should be ignored based on directory patterns
func (w *Watcher) isIgnoredPath(path string) bool {
    // Convert to absolute path for consistent checking
    absPath, err := filepath.Abs(path)
    if err != nil {
        // If we can't get absolute path, use relative
        absPath = path
    }
    
    // Split path into components
    components := strings.Split(filepath.ToSlash(absPath), "/")
    
    // Check each component against ignored directories
    for _, component := range components {
        if w.ignoredDirs[component] {
            return true
        }
    }
    
    // Check for hidden directories (starting with .)
    for _, component := range components {
        if strings.HasPrefix(component, ".") && component != "." && component != ".." {
            // Allow .gitignore, .env, etc. but skip .git, .idea, etc.
            if w.ignoredDirs[component] {
                return true
            }
            // Consider ignoring all dot directories except explicitly allowed ones
            if isDir, _ := isDirectory(filepath.Join(components[:len(components)-1]...)); isDir {
                return true
            }
        }
    }
    
    return false
}

// isTemporaryFile checks for common temporary/backup file patterns
func isTemporaryFile(path string) bool {
    base := filepath.Base(path)
    
    // Common editor backup/temporary file patterns
    tempPatterns := []string{
        "~",                     // Vim/Emacs backup
        ".swp", ".swo", ".swn", // Vim swap files
        ".tmp", ".temp",        // General temporary files
        "#",                    // Some editor backup
        ".sync-conflict-",      // Syncing conflicts
        ".DS_Store",           // macOS metadata
        "Thumbs.db",           // Windows thumbnail cache
    }
    
    for _, pattern := range tempPatterns {
        if strings.HasSuffix(base, pattern) || strings.Contains(base, pattern) {
            return true
        }
    }
    
    // Check for files in system temporary directories
    if strings.Contains(path, "/tmp/") || strings.Contains(path, "/Temp/") {
        return true
    }
    
    return false
}

// isDirEvent checks if the event is for a directory (not a file)
func isDirEvent(path string) bool {
    info, err := os.Stat(path)
    if err != nil {
        // If we can't stat the file, assume it's not a directory
        return false
    }
    return info.IsDir()
}

// isDirectory checks if a path is a directory (helper function)
func isDirectory(path string) (bool, error) {
    info, err := os.Stat(path)
    if err != nil {
        return false, err
    }
    return info.IsDir(), nil
}

func (w *Watcher) processChanges() {
    var timer *time.Timer
    var changedFiles []string
    
    for {
        select {
        case event, ok := <-w.watcher.Events:
            if !ok {
                return
            }
            
            // Filter relevant changes
            if !w.shouldProcessEvent(event) {
                continue
            }
            
            ext := filepath.Ext(event.Name)
            if w.extensions[ext] {
                changedFiles = append(changedFiles, event.Name)
            }
            
            // Debounce multiple changes
            if timer != nil {
                timer.Stop()
            }
            timer = time.AfterFunc(w.debounceTime, func() {
                w.onChange(changedFiles)
                changedFiles = nil
            })
            
        case err, ok := <-w.watcher.Errors:
            if !ok {
                return
            }
            log.Printf("Watcher error: %v", err)
        }
    }
}
