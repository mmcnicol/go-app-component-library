// cmd/dev-server/main.go
package main

import (
    "embed"
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "path/filepath"
    "sync"
    "time"
    
    "github.com/mmcnicol/go-app-component-library/cmd/dev-server/build"
    "github.com/mmcnicol/go-app-component-library/cmd/dev-server/handlers"
    "github.com/mmcnicol/go-app-component-library/cmd/dev-server/watch"
)

//go:embed static/*
var staticFiles embed.FS

type Server struct {
    port          int
    workDir       string
    compiler      *build.Compiler
    watcher       *watch.Watcher
    liveReload    *handlers.LiveReloadServer
    currentWasm   string
    wasmMu        sync.RWMutex
    dashboardData *DashboardData
}

func NewServer(port int, workDir string) (*Server, error) {
    outputDir := filepath.Join(workDir, "tmp", "wasm")
    os.MkdirAll(outputDir, 0755)
    
    s := &Server{
        port:        port,
        workDir:     workDir,
        compiler:    build.NewCompiler(workDir, outputDir),
        liveReload:  handlers.NewLiveReloadServer(),
        dashboardData: &DashboardData{},
    }
    
    // Initialize watcher
    watcher, err := watch.NewWatcher(workDir, s.onFileChange)
    if err != nil {
        return nil, err
    }
    s.watcher = watcher
    
    // Initial build
    initialWasm, err := s.compiler.BuildWasm(context.Background(),
        filepath.Join(workDir, "cmd/dev-server/main_wasm.go"),
        nil)
    if err != nil {
        return nil, err
    }
    s.currentWasm = initialWasm
    
    return s, nil
}

func (s *Server) onFileChange(changedFiles []string) {
    log.Printf("Files changed: %v", changedFiles)
    
    s.dashboardData.AddFileChanges(changedFiles)
    s.dashboardData.SetBuildStatus("building")
    
    // Try incremental build first
    wasmPath, err := s.compiler.BuildOnlyChanged(context.Background(), changedFiles)
    if err != nil || wasmPath == "" {
        // Fall back to full build
        wasmPath, err = s.compiler.BuildWasm(context.Background(),
            filepath.Join(s.workDir, "cmd/dev-server/main_wasm.go"),
            changedFiles)
    }
    
    if err != nil {
        log.Printf("Build failed: %v", err)
        s.dashboardData.AddError(err.Error())
        s.dashboardData.SetBuildStatus("error")
    } else {
        s.wasmMu.Lock()
        s.currentWasm = wasmPath
        s.wasmMu.Unlock()
        
        s.dashboardData.SetBuildStatus("success")
        s.dashboardData.SetLastBuildTime(time.Now())
        
        // Notify clients to reload
        s.liveReload.BroadcastReload("file change")
    }
}

func (s *Server) Start() error {
    mux := http.NewServeMux()
    
    // Static files
    mux.Handle("/static/", http.FileServer(http.FS(staticFiles)))
    
    // WASM file with cache busting
    mux.HandleFunc("/app.wasm", s.serveWasm)
    
    // Live reload WebSocket
    mux.Handle("/ws", s.liveReload)
    
    // Development dashboard
    mux.HandleFunc("/_dashboard", s.serveDashboard)
    
    // Main application handler
    mux.HandleFunc("/", s.serveApp)
    
    // API for dashboard data
    mux.HandleFunc("/api/dashboard", s.serveDashboardData)
    
    addr := fmt.Sprintf(":%d", s.port)
    log.Printf("Development server starting on http://localhost%s", addr)
    log.Printf("Dashboard: http://localhost%s/_dashboard", addr)
    
    return http.ListenAndServe(addr, mux)
}

func (s *Server) serveWasm(w http.ResponseWriter, r *http.Request) {
    s.wasmMu.RLock()
    wasmPath := s.currentWasm
    s.wasmMu.RUnlock()
    
    // Add cache busting headers
    w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
    w.Header().Set("Pragma", "no-cache")
    w.Header().Set("Expires", "0")
    
    http.ServeFile(w, r, wasmPath)
}
