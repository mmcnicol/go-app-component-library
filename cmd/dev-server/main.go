// cmd/dev-server/main.go
package main

import (
    "bytes"
    "embed"
    "encoding/json"
    "flag"
    "fmt"
    "log"
    "net/http"
    "os"
    "os/exec"
    "path/filepath"
    "runtime"
    "strings"
    "sync"
    "time"
    
    "github.com/mmcnicol/go-app-component-library/cmd/dev-server/build"
    "github.com/mmcnicol/go-app-component-library/cmd/dev-server/handlers"
    "github.com/mmcnicol/go-app-component-library/cmd/dev-server/watch"
)

// DashboardData holds metrics for the development dashboard
type DashboardData struct {
    BuildStatus       string    `json:"build_status"`
    LastBuildTime     time.Time `json:"last_build_time"`
    BuildCount        int       `json:"build_count"`
    ConnectedClients  int       `json:"connected_clients"`
    FileChanges       []string  `json:"file_changes"`
    CompileErrors     []string  `json:"compile_errors"`
    mu                sync.RWMutex
}

// AddFileChanges adds file changes to the dashboard data
func (d *DashboardData) AddFileChanges(files []string) {
    d.mu.Lock()
    defer d.mu.Unlock()
    d.FileChanges = append(d.FileChanges, files...)
}

// AddError adds an error to the dashboard data
func (d *DashboardData) AddError(err string) {
    d.mu.Lock()
    defer d.mu.Unlock()
    d.CompileErrors = append(d.CompileErrors, err)
}

// SetBuildStatus sets the build status
func (d *DashboardData) SetBuildStatus(status string) {
    d.mu.Lock()
    defer d.mu.Unlock()
    d.BuildStatus = status
}

// SetLastBuildTime sets the last build time and increments build count
func (d *DashboardData) SetLastBuildTime(t time.Time) {
    d.mu.Lock()
    defer d.mu.Unlock()
    d.LastBuildTime = t
    d.BuildCount++
}

// SetConnectedClients sets the number of connected clients
func (d *DashboardData) SetConnectedClients(count int) {
    d.mu.Lock()
    defer d.mu.Unlock()
    d.ConnectedClients = count
}

// Clear clears all dashboard data
func (d *DashboardData) Clear() {
    d.mu.Lock()
    defer d.mu.Unlock()
    d.FileChanges = nil
    d.CompileErrors = nil
}

// GetData returns a copy of the dashboard data
func (d *DashboardData) GetData() DashboardData {
    d.mu.RLock()
    defer d.mu.RUnlock()
    return DashboardData{
        BuildStatus:      d.BuildStatus,
        LastBuildTime:    d.LastBuildTime,
        BuildCount:       d.BuildCount,
        ConnectedClients: d.ConnectedClients,
        FileChanges:      append([]string{}, d.FileChanges...),
        CompileErrors:    append([]string{}, d.CompileErrors...),
    }
}

// Server represents the development server
type Server struct {
    port          int
    workDir       string
    outputDir     string
    compiler      *build.Compiler
    watcher       *watch.Watcher
    liveReload    *handlers.LiveReloadServer
    currentWasm   string
    wasmMu        sync.RWMutex
    dashboardData *DashboardData
    enableDashboard bool
    profile       bool
}

// NewServer creates a new development server
func NewServer(port int, workDir string, enableDashboard, profile bool) (*Server, error) {
    // Use web folder as output directory
    outputDir := filepath.Join(workDir, "web")
    if err := os.MkdirAll(outputDir, 0755); err != nil {
        return nil, fmt.Errorf("failed to create web directory: %v", err)
    }
    
    // Create compiler with error handling
    compiler, err := build.NewCompiler(workDir, outputDir)
    if err != nil {
        // Log the error but continue - the compiler might still work
        log.Printf("Warning: failed to create compiler: %v", err)
        // Don't return error here, we'll try to build directly
    }
    
    s := &Server{
        port:            port,
        workDir:         workDir,
        outputDir:       outputDir,
        compiler:        compiler,
        liveReload:      handlers.NewLiveReloadServer(),
        dashboardData:   &DashboardData{},
        enableDashboard: enableDashboard,
        profile:         profile,
    }
    
    // Initialize watcher - also watch web folder for CSS/JS changes
    watcher, err := watch.NewWatcher(workDir, s.onFileChange)
    if err != nil {
        return nil, fmt.Errorf("failed to create watcher: %v", err)
    }
    s.watcher = watcher
    
    // Initial build - output to web/app.wasm
    initialWasm, err := s.buildWasm()
    if err != nil {
        log.Printf("Initial build warning: %v", err)
        // Don't fail here - the server can still start
        // Create a dummy wasm file or continue without it
        s.currentWasm = ""
    } else {
        s.currentWasm = initialWasm
    }
    
    // Start monitoring connected clients
    go s.monitorClients()
    
    return s, nil
}

// buildWasm builds the WebAssembly binary
func (s *Server) buildWasm() (string, error) {
    outputPath := filepath.Join(s.outputDir, "app.wasm")
    
    // The wasm entry point is at cmd/wasm/main.go
    mainPackage := "./cmd/wasm"
    
    // First verify the package exists
    verifyCmd := exec.Command("go", "list", mainPackage)
    verifyCmd.Dir = s.workDir
    if output, err := verifyCmd.CombinedOutput(); err != nil {
        log.Printf("Package verification failed: %s", output)
        // Continue anyway, the build might still work
    }
    
    cmd := exec.Command("go", "build",
        "-o", outputPath,
        "-tags", "dev",
        "-ldflags", "-s -w",
        mainPackage,
    )
    
    cmd.Dir = s.workDir
    
    // Set environment variables
    env := os.Environ()
    var cleanEnv []string
    
    // Filter and set environment
    for _, e := range env {
        // Keep most env vars but override some
        if !strings.HasPrefix(e, "GOOS=") &&
           !strings.HasPrefix(e, "GOARCH=") &&
           !strings.HasPrefix(e, "GO111MODULE=") &&
           !strings.HasPrefix(e, "GOPROXY=") {
            cleanEnv = append(cleanEnv, e)
        }
    }
    
    // Add required env vars for wasm build
    cleanEnv = append(cleanEnv,
        "GOOS=js",
        "GOARCH=wasm",
        "GO111MODULE=on",
        "GOPROXY=direct",
    )
    
    cmd.Env = cleanEnv
    
    var stdout, stderr bytes.Buffer
    cmd.Stdout = &stdout
    cmd.Stderr = &stderr
    
    log.Printf("Building WebAssembly from %s (package: %s)", s.workDir, mainPackage)
    log.Printf("Command: go build -o %s -tags dev -ldflags \"-s -w\" %s", outputPath, mainPackage)
    
    if err := cmd.Run(); err != nil {
        // Try an alternative approach if the first fails
        log.Printf("Build failed: %v", err)
        log.Printf("Stderr: %s", stderr.String())
        log.Printf("Stdout: %s", stdout.String())
        
        return s.buildWasmAlternative(outputPath)
    }
    
    log.Printf("✅ Built: %s", outputPath)
    return outputPath, nil
}

func (s *Server) buildWasmAlternative(outputPath string) (string, error) {
    // Try building with the full path to main.go
    mainFile := filepath.Join(s.workDir, "cmd", "wasm", "main.go")
    
    log.Printf("Trying alternative build with main file: %s", mainFile)
    
    cmd := exec.Command("go", "build",
        "-o", outputPath,
        "-tags", "dev",
        "-ldflags", "-s -w",
        mainFile,
    )
    
    cmd.Dir = s.workDir
    cmd.Env = append(os.Environ(),
        "GOOS=js",
        "GOARCH=wasm",
        "GO111MODULE=on",
    )
    
    var stderr bytes.Buffer
    cmd.Stderr = &stderr
    
    if err := cmd.Run(); err != nil {
        return "", fmt.Errorf("alternative build failed: %v\n%s", err, stderr.String())
    }
    
    return outputPath, nil
}

// monitorClients periodically updates the connected client count
func (s *Server) monitorClients() {
    ticker := time.NewTicker(2 * time.Second)
    defer ticker.Stop()
    
    for range ticker.C {
        count := s.liveReload.GetClientCount()
        s.dashboardData.SetConnectedClients(count)
    }
}

// onFileChange handles file change events from the watcher
func (s *Server) onFileChange(changedFiles []string) {
    log.Printf("Files changed: %v", changedFiles)
    
    s.dashboardData.AddFileChanges(changedFiles)
    s.dashboardData.SetBuildStatus("building")
    
    // Filter only Go files
    var goFiles []string
    for _, file := range changedFiles {
        if strings.HasSuffix(file, ".go") {
            goFiles = append(goFiles, file)
        }
    }
    
    if len(goFiles) == 0 {
        log.Println("No Go files changed, skipping rebuild")
        return
    }
    
    // Rebuild the wasm
    wasmPath, err := s.buildWasm()
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
        
        // Clear errors on successful build
        s.dashboardData.Clear()
        
        // Force browser cache invalidation with timestamp
        timestamp := time.Now().Unix()
        
        // Notify clients to reload with cache busting
        s.liveReload.BroadcastMessage("reload", map[string]interface{}{
            "reason":    "file change",
            "timestamp": timestamp,
            "files":     goFiles,
        })
        
        log.Printf("Rebuilt successfully, notifying clients")
    }
}

// forceRebuild forces a complete rebuild
func (s *Server) forceRebuild() {
    s.dashboardData.SetBuildStatus("building")
    
    wasmPath, err := s.buildWasm()
    if err != nil {
        log.Printf("Force rebuild failed: %v", err)
        s.dashboardData.AddError(err.Error())
        s.dashboardData.SetBuildStatus("error")
    } else {
        s.wasmMu.Lock()
        s.currentWasm = wasmPath
        s.wasmMu.Unlock()
        
        s.dashboardData.SetBuildStatus("success")
        s.dashboardData.SetLastBuildTime(time.Now())
        s.dashboardData.Clear()
        
        s.liveReload.BroadcastReload("manual rebuild")
    }
}

// clearCache clears the build cache
func (s *Server) clearCache() {
    // Implementation depends on your cache structure
    log.Println("Cache cleared (placeholder)")
    s.dashboardData.AddFileChanges([]string{"Cache cleared manually"})
}

// Start starts the development server
func (s *Server) Start() error {
    mux := http.NewServeMux()
    
    // Serve static files from web folder
    webDir := filepath.Join(s.workDir, "web")
    mux.Handle("/", http.FileServer(http.Dir(webDir)))
    
    // Override specific routes
    mux.HandleFunc("/app.wasm", s.serveWasm)
    mux.Handle("/ws", s.liveReload)
    
    // Development dashboard API (if enabled)
    if s.enableDashboard {
        mux.HandleFunc("/api/dashboard", s.serveDashboardData)
        mux.HandleFunc("/api/rebuild", s.handleRebuild)
        mux.HandleFunc("/api/clear-cache", s.handleClearCache)
    }
    
    addr := fmt.Sprintf(":%d", s.port)
    log.Printf("Development server starting on http://localhost%s", addr)
    log.Printf("Serving from: %s", webDir)
    
    if s.enableDashboard {
        log.Printf("Dashboard API: http://localhost%s/api/dashboard", addr)
    }
    
    return http.ListenAndServe(addr, mux)
}

// serveWasm serves the WebAssembly binary
func (s *Server) serveWasm(w http.ResponseWriter, r *http.Request) {
    s.wasmMu.RLock()
    wasmPath := s.currentWasm
    s.wasmMu.RUnlock()
    
    if wasmPath == "" {
        // Try to build on-demand
        log.Println("No WebAssembly binary available, attempting to build...")
        wasmPath, err := s.buildWasm()
        if err != nil {
            http.Error(w, "No WebAssembly binary available and build failed: "+err.Error(), http.StatusNotFound)
            return
        }
        s.wasmMu.Lock()
        s.currentWasm = wasmPath
        s.wasmMu.Unlock()
    }
    
    // Add cache busting headers
    w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
    w.Header().Set("Pragma", "no-cache")
    w.Header().Set("Expires", "0")
    
    http.ServeFile(w, r, wasmPath)
}

// serveDashboardData serves dashboard data as JSON
func (s *Server) serveDashboardData(w http.ResponseWriter, r *http.Request) {
    data := s.dashboardData.GetData()
    
    w.Header().Set("Content-Type", "application/json")
    w.Header().Set("Access-Control-Allow-Origin", "*")
    
    if err := json.NewEncoder(w).Encode(data); err != nil {
        http.Error(w, "Failed to encode dashboard data", http.StatusInternalServerError)
    }
}

// handleRebuild handles manual rebuild requests
func (s *Server) handleRebuild(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }
    
    go s.forceRebuild()
    
    w.Header().Set("Content-Type", "application/json")
    w.Header().Set("Access-Control-Allow-Origin", "*")
    json.NewEncoder(w).Encode(map[string]string{"status": "rebuild_started"})
}

// handleClearCache handles cache clearing requests
func (s *Server) handleClearCache(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }
    
    go s.clearCache()
    
    w.Header().Set("Content-Type", "application/json")
    w.Header().Set("Access-Control-Allow-Origin", "*")
    json.NewEncoder(w).Encode(map[string]string{"status": "cache_cleared"})
}

// main is the entry point for the dev-server
func main() {
    var (
        port          = flag.Int("port", 8080, "Port to listen on")
        workDir       = flag.String("dir", ".", "Working directory")
        watch         = flag.Bool("watch", true, "Enable file watching")
        dashboard     = flag.Bool("dashboard", false, "Enable dashboard")
        profile       = flag.Bool("profile", false, "Enable profiling")
    )
    
    flag.Parse()
    
    server, err := NewServer(*port, *workDir, *dashboard, *profile)
    if err != nil {
        log.Fatalf("Failed to create server: %v", err)
    }
    
    log.Printf("Starting development server...")
    log.Printf("  Port: %d", *port)
    log.Printf("  Directory: %s", *workDir)
    log.Printf("  Watch: %v", *watch)
    log.Printf("  Dashboard: %v", *dashboard)
    
    if err := server.Start(); err != nil {
        log.Fatalf("Server failed: %v", err)
    }
}
