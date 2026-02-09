// cmd/dev-server/main.go
package main

import (
    "context"
    "embed"
    "encoding/json"
    "flag"
    "fmt"
    "log"
    "net/http"
    "os"
    "path/filepath"
    "sync"
    "time"
    
    "github.com/mmcnicol/go-app-component-library/cmd/dev-server/build"
    "github.com/mmcnicol/go-app-component-library/cmd/dev-server/handlers"
    "github.com/mmcnicol/go-app-component-library/cmd/dev-server/watch"
)

//go:embed static/*
var staticFiles embed.FS

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
        return nil, fmt.Errorf("failed to create compiler: %v", err)
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
        return nil, fmt.Errorf("initial build failed: %v", err)
    }
    s.currentWasm = initialWasm
    
    // Start monitoring connected clients
    go s.monitorClients()
    
    return s, nil
}

// buildWasm builds the WebAssembly binary
func (s *Server) buildWasm() (string, error) {
    // Create a simple wasm file for development
    tempDir, err := os.MkdirTemp("", "dev-wasm-*")
    if err != nil {
        return "", fmt.Errorf("failed to create temp dir: %v", err)
    }
    defer os.RemoveAll(tempDir)
    
    mainFile := filepath.Join(tempDir, "main.go")
    
    content := `package main

import (
    "github.com/maxence-charriere/go-app/v10/pkg/app"
)

type DevApp struct {
    app.Compo
}

func (d *DevApp) Render() app.UI {
    return app.Div().Body(
        app.H1().Text("Go App Component Library - Development"),
        app.P().Text("✅ Development server is running!"),
        app.P().Text("✨ Hot reload is active - edit your components and see changes instantly."),
        app.P().Text("📁 Serving from web/ folder with wasm_exec.js"),
        app.Hr(),
        app.Div().Style("margin-top", "20px").Body(
            app.H3().Text("Getting Started:"),
            app.Ul().Body(
                app.Li().Text("Edit components in pkg/components/"),
                app.Li().Text("Save changes"),
                app.Li().Text("Watch the browser reload automatically"),
            ),
        ),
        app.Div().Style("margin-top", "20px").Body(
            app.H3().Text("File Structure:"),
            app.Ul().Body(
                app.Li().Text("web/app.wasm - WebAssembly binary (auto-generated)"),
                app.Li().Text("web/wasm_exec.js - WebAssembly runtime (from Go)"),
                app.Li().Text("web/styles.css - Custom styles (optional)"),
            ),
        ),
    )
}

func main() {
    app.Route("/", func() app.Composer { return &DevApp{} })
    app.RunWhenOnBrowser()
}`
    
    if err := os.WriteFile(mainFile, []byte(content), 0644); err != nil {
        return "", fmt.Errorf("failed to write temp main file: %v", err)
    }
    
    // Build to web/app.wasm (fixed name for simplicity)
    outputPath := filepath.Join(s.outputDir, "app.wasm")
    
    wasmPath, err := s.compiler.BuildWasmToPath(context.Background(), mainFile, nil, outputPath)
    if err != nil {
        return "", fmt.Errorf("build failed: %v", err)
    }
    
    log.Printf("Built WebAssembly to: %s", wasmPath)
    return wasmPath, nil
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
    
    // Try incremental build first
    wasmPath, err := s.compiler.BuildOnlyChanged(context.Background(), changedFiles)
    if err != nil || wasmPath == "" {
        // Fall back to full build
        wasmPath, err = s.buildWasm()
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
        
        // Clear errors on successful build
        s.dashboardData.Clear()
        
        // Notify clients to reload
        s.liveReload.BroadcastReload("file change")
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
        http.Error(w, "No WebAssembly binary available", http.StatusNotFound)
        return
    }
    
    // Add cache busting headers
    w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
    w.Header().Set("Pragma", "no-cache")
    w.Header().Set("Expires", "0")
    
    http.ServeFile(w, r, wasmPath)
}

// serveWasmExec serves the wasm_exec.js file
func (s *Server) serveWasmExec(w http.ResponseWriter, r *http.Request) {
    // Try to serve from embedded static files first
    data, err := staticFiles.ReadFile("static/wasm_exec.js")
    if err == nil {
        w.Header().Set("Content-Type", "application/javascript")
        w.Write(data)
        return
    }
    
    // Fallback: try to read from go-app package
    wasmExecPath := filepath.Join(s.workDir, "vendor", "github.com", "maxence-charriere", "go-app", "v10", "cmd", "wasm_exec.js")
    
    if _, err := os.Stat(wasmExecPath); err == nil {
        http.ServeFile(w, r, wasmExecPath)
        return
    }
    
    // Last resort: try Go installation
    goRoot := runtime.GOROOT()
    if goRoot != "" {
        wasmExecPath = filepath.Join(goRoot, "misc", "wasm", "wasm_exec.js")
        if _, err := os.Stat(wasmExecPath); err == nil {
            http.ServeFile(w, r, wasmExecPath)
            return
        }
    }
    
    // Generate a minimal wasm_exec.js if nothing else works
    w.Header().Set("Content-Type", "application/javascript")
    w.Write([]byte(`// Minimal wasm_exec.js for development
    const go = new Go();
    `))
}

// serveApp serves the main application HTML page
func (s *Server) serveApp(w http.ResponseWriter, r *http.Request) {
    html := `<!DOCTYPE html>
<html>
<head>
    <title>Go App Component Library - Dev Server</title>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <script>
        // Live reload WebSocket
        const ws = new WebSocket('ws://' + window.location.host + '/ws');
        
        ws.onmessage = function(event) {
            const data = JSON.parse(event.data);
            if (data.type === 'reload') {
                console.log('Reloading application due to:', data.reason);
                window.location.reload();
            }
        };
        
        ws.onopen = function() {
            console.log('Live reload WebSocket connected');
        };
        
        ws.onerror = function(error) {
            console.error('WebSocket error:', error);
        };
        
        // Auto-refresh when wasm changes
        let currentWasmTime = Date.now();
        
        function checkForUpdates() {
            fetch('/app.wasm?check=' + currentWasmTime)
                .then(response => {
                    const lastModified = response.headers.get('last-modified');
                    if (lastModified) {
                        const serverTime = new Date(lastModified).getTime();
                        if (serverTime > currentWasmTime) {
                            console.log('WebAssembly updated, reloading...');
                            window.location.reload();
                        }
                    }
                })
                .catch(err => console.log('Update check failed:', err));
        }
        
        // Check for updates every 2 seconds
        setInterval(checkForUpdates, 2000);
    </script>
</head>
<body>
    <div id="app">
        <div style="padding: 20px; font-family: sans-serif;">
            <h1>Loading Go WebAssembly...</h1>
            <p>If this message persists, check the browser console for errors.</p>
            <div id="loading-status" style="margin-top: 20px;">
                <div style="display: flex; align-items: center; gap: 10px;">
                    <div id="spinner" style="width: 20px; height: 20px; border: 3px solid #f3f3f3; border-top: 3px solid #3498db; border-radius: 50%; animation: spin 1s linear infinite;"></div>
                    <span>Compiling and loading WebAssembly...</span>
                </div>
                <style>
                    @keyframes spin {
                        0% { transform: rotate(0deg); }
                        100% { transform: rotate(360deg); }
                    }
                </style>
            </div>
        </div>
    </div>
    
    <script src="/wasm_exec.js"></script>
    <script>
        // Load WebAssembly with error handling
        async function loadWasm() {
            const go = new Go();
            const importObject = go.importObject;
            
            try {
                const response = await fetch('/app.wasm?t=' + Date.now());
                if (!response.ok) {
                    throw new Error('Failed to fetch WebAssembly');
                }
                
                const bytes = await response.arrayBuffer();
                const result = await WebAssembly.instantiate(bytes, importObject);
                
                document.getElementById('loading-status').innerHTML = 
                    '<div style="color: green;">✓ WebAssembly loaded, starting app...</div>';
                
                go.run(result.instance);
            } catch (err) {
                console.error('Failed to load WebAssembly:', err);
                document.getElementById('app').innerHTML = 
                    '<div style="padding: 20px; color: red;">' +
                    '<h2>Failed to load application</h2>' +
                    '<p>Error: ' + err.message + '</p>' +
                    '<p>Check the browser console for details.</p>' +
                    '<button onclick="window.location.reload()" style="padding: 8px 16px; background: #3498db; color: white; border: none; border-radius: 4px; cursor: pointer;">Retry</button>' +
                    '</div>';
            }
        }
        
        // Start loading
        loadWasm();
    </script>
</body>
</html>`
    
    w.Header().Set("Content-Type", "text/html; charset=utf-8")
    w.Write([]byte(html))
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
        outputDir     = flag.String("output", "", "Output directory for wasm files")
        watch         = flag.Bool("watch", true, "Enable file watching")
        dashboard     = flag.Bool("dashboard", false, "Enable dashboard")
        profile       = flag.Bool("profile", false, "Enable profiling")
    )
    
    flag.Parse()
    
    server, err := NewServer(*port, *workDir, *outputDir, *dashboard, *profile)
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
