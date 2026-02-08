// cmd/dev-server/watch/watcher.go
package watch

import (
    "log"
    "path/filepath"
    "strings"
    
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

