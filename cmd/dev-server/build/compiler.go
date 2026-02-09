// cmd/dev-server/build/compiler.go
package build

import (
    "bytes"
    "context"
    "fmt"
    "log"
    "os"
    "os/exec"
    "path/filepath"
    "strings"
    "time"
)

type Compiler struct {
    goBinary      string
    workDir       string
    outputDir     string
    cache         *BuildCache
    buildTags     []string
    ldflags       string
}

func NewCompiler(workDir, outputDir string) (*Compiler, error) {
    cache, err := NewBuildCache(filepath.Join(outputDir, ".buildcache"))
    if err != nil {
        return nil, fmt.Errorf("failed to create build cache: %v", err)
    }
    
    return &Compiler{
        goBinary:  "go",
        workDir:   workDir,
        outputDir: outputDir,
        cache:     cache,
        buildTags: []string{"dev"},
        ldflags:   "-s -w",
    }, nil
}

func (c *Compiler) BuildWasm(ctx context.Context, mainFile string, changedFiles []string) (string, error) {
    // Check cache
    if c.cache.IsValid(mainFile, changedFiles) {
        if entry, exists := c.cache.Get(mainFile); exists {
            return entry.OutputPath, nil
        }
    }
    
    outputPath := filepath.Join(c.outputDir, fmt.Sprintf("app-%d.wasm", time.Now().UnixNano()))
    
    // Use the WORKDIR as the module root
    absMainFile, err := filepath.Abs(mainFile)
    if err != nil {
        return "", fmt.Errorf("failed to get absolute path: %v", err)
    }
    
    // Build with explicit module mode
    cmd := exec.CommandContext(ctx, c.goBinary, "build",
        "-o", outputPath,
        "-tags", joinTags(c.buildTags),
        "-ldflags", c.ldflags,
    )
    
    // For local packages, use the directory containing go.mod
    cmd.Args = append(cmd.Args, "./"+filepath.ToSlash(strings.TrimPrefix(absMainFile, c.workDir+string(filepath.Separator))))
    
    cmd.Dir = c.workDir  // This is CRITICAL - must be the module root
    
    // Set environment to help with local module resolution
    env := os.Environ()
    
    // Remove any GOPROXY that might interfere
    var cleanEnv []string
    for _, e := range env {
        if !strings.HasPrefix(e, "GOPROXY=") && !strings.HasPrefix(e, "GOPRIVATE=") {
            cleanEnv = append(cleanEnv, e)
        }
    }
    
    cleanEnv = append(cleanEnv,
        "GOOS=js",
        "GOARCH=wasm",
        "GO111MODULE=on",
        "GOPROXY=direct",           // Don't use proxy for local builds
        "GOPRIVATE=github.com/mmcnicol/go-app-component-library", // Treat our package as private
    )
    
    cmd.Env = cleanEnv
    
    var stdout, stderr bytes.Buffer
    cmd.Stdout = &stdout
    cmd.Stderr = &stderr
    
    start := time.Now()
    err = cmd.Run()
    buildTime := time.Since(start)
    
    if err != nil {
        // Try one more approach - build from current directory
        return c.buildFromCurrentDir(mainFile, outputPath)
    }
    
    // Update cache
    entry := CacheEntry{
        OutputPath:   outputPath,
        Hash:         c.cache.calculateHash(mainFile, changedFiles),
        Dependencies: changedFiles,
        Timestamp:    time.Now(),
    }
    
    if err := c.cache.Set(mainFile, entry); err != nil {
        log.Printf("Warning: Failed to update cache: %v", err)
    }
    
    log.Printf("Built %s in %v", filepath.Base(outputPath), buildTime)
    return outputPath, nil
}

// tryVendorBuild tries to build using vendor directory
func (c *Compiler) tryVendorBuild(ctx context.Context, mainFile, outputPath string) (string, error) {
    cmd := exec.CommandContext(ctx, c.goBinary, "build",
        "-o", outputPath,
        "-tags", joinTags(c.buildTags),
        "-ldflags", c.ldflags,
        "-mod=vendor",
        mainFile,
    )
    
    cmd.Dir = c.workDir
    cmd.Env = append(os.Environ(),
        "GOOS=js",
        "GOARCH=wasm",
        "GO111MODULE=on",
    )
    
    var stdout, stderr bytes.Buffer
    cmd.Stdout = &stdout
    cmd.Stderr = &stderr
    
    log.Printf("Trying vendor build...")
    
    if err := cmd.Run(); err != nil {
        return "", fmt.Errorf("vendor build also failed: %v\n%s", err, stderr.String())
    }
    
    return outputPath, nil
}

// BuildOnlyChanged implements incremental compilation
func (c *Compiler) BuildOnlyChanged(ctx context.Context, changedFiles []string) (string, error) {
    // Analyze which packages are affected
    affectedPackages := c.analyzeDependencies(changedFiles)
    
    if len(affectedPackages) == 0 {
        // No relevant changes
        return "", nil
    }
    
    // Create a temporary main.go that imports only affected packages
    tempMain, err := c.createSelectiveMain(affectedPackages)
    if err != nil {
        return "", fmt.Errorf("failed to create selective main: %v", err)
    }
    defer os.Remove(tempMain)
    
    return c.BuildWasm(ctx, tempMain, changedFiles)
}

func (c *Compiler) analyzeDependencies(changedFiles []string) []string {
    // Parse Go files to build import graph
    importGraph := c.buildImportGraph()
    
    affectedPackages := make(map[string]bool)
    for _, changedFile := range changedFiles {
        pkg := c.fileToPackage(changedFile)
        if pkg == "" {
            continue
        }
        
        // Find all packages that depend on this package
        c.findDependents(pkg, importGraph, affectedPackages)
    }
    
    // Convert to slice
    result := make([]string, 0, len(affectedPackages))
    for pkg := range affectedPackages {
        result = append(result, pkg)
    }
    
    return result
}

// joinTags joins build tags with commas
func joinTags(tags []string) string {
    return strings.Join(tags, ",")
}

// buildImportGraph builds a map of package -> packages that import it
func (c *Compiler) buildImportGraph() map[string][]string {
    // This is a simplified implementation
    // In a real implementation, you would parse .go files to find imports
    
    // For now, return an empty graph
    return make(map[string][]string)
}

// fileToPackage converts a file path to its package name
func (c *Compiler) fileToPackage(filePath string) string {
    relPath, err := filepath.Rel(c.workDir, filePath)
    if err != nil {
        return ""
    }
    
    // Get the directory containing the file
    dir := filepath.Dir(relPath)
    if dir == "." {
        dir = ""
    }
    
    // In Go, package name is usually the directory name
    // This is simplified - real implementation would parse the package declaration
    if dir == "" {
        return "main"
    }
    return dir
}

// findDependents recursively finds all packages that depend on the given package
func (c *Compiler) findDependents(pkg string, importGraph map[string][]string, visited map[string]bool) {
    if visited[pkg] {
        return
    }
    visited[pkg] = true
    
    // Add direct dependents
    if dependents, exists := importGraph[pkg]; exists {
        for _, dependent := range dependents {
            c.findDependents(dependent, importGraph, visited)
        }
    }
}

// createSelectiveMain creates a temporary main.go file that imports only affected packages
func (c *Compiler) createSelectiveMain(packages []string) (string, error) {
    // Create temporary directory for the main file
    tempDir, err := os.MkdirTemp("", "selective-build-*")
    if err != nil {
        return "", fmt.Errorf("failed to create temp dir: %v", err)
    }
    
    mainFile := filepath.Join(tempDir, "main.go")
    
    // Create the main.go content
    var content bytes.Buffer
    content.WriteString("// +build dev\n\n")
    content.WriteString("package main\n\n")
    content.WriteString("import (\n")
    
    for _, pkg := range packages {
        // Determine import path
        importPath := pkg
        if !strings.Contains(pkg, "/") {
            // Local package
            importPath = fmt.Sprintf("./%s", pkg)
        }
        fmt.Fprintf(&content, "    _ \"%s\"\n", importPath)
    }
    
    content.WriteString(")\n\n")
    content.WriteString("func main() {\n")
    content.WriteString("    // This main is only for selective compilation\n")
    content.WriteString("    // Real entry point is in the actual application\n")
    content.WriteString("}\n")
    
    // Write the file
    if err := os.WriteFile(mainFile, content.Bytes(), 0644); err != nil {
        os.RemoveAll(tempDir)
        return "", fmt.Errorf("failed to write main.go: %v", err)
    }
    
    return mainFile, nil
}

// Cleanup removes old build artifacts
func (c *Compiler) Cleanup(maxAge time.Duration) error {
    entries, err := os.ReadDir(c.outputDir)
    if err != nil {
        return fmt.Errorf("failed to read output directory: %v", err)
    }
    
    cutoff := time.Now().Add(-maxAge)
    for _, entry := range entries {
        if entry.IsDir() {
            continue
        }
        
        info, err := entry.Info()
        if err != nil {
            continue
        }
        
        if info.ModTime().Before(cutoff) && strings.HasSuffix(info.Name(), ".wasm") {
            filePath := filepath.Join(c.outputDir, info.Name())
            if err := os.Remove(filePath); err != nil {
                log.Printf("Failed to remove old build %s: %v", info.Name(), err)
            } else {
                log.Printf("Removed old build: %s", info.Name())
            }
        }
    }
    
    return nil
}

// GetGoVersion returns the Go version being used
func (c *Compiler) GetGoVersion() (string, error) {
    cmd := exec.Command(c.goBinary, "version")
    output, err := cmd.Output()
    if err != nil {
        return "", fmt.Errorf("failed to get Go version: %v", err)
    }
    return strings.TrimSpace(string(output)), nil
}

// BuildWasmToPath builds WebAssembly to a specific output path
func (c *Compiler) BuildWasmToPath(ctx context.Context, mainFile string, changedFiles []string, outputPath string) (string, error) {
    // Check cache for unchanged files
    if c.cache.IsValid(mainFile, changedFiles) {
        if entry, exists := c.cache.Get(mainFile); exists {
            return entry.OutputPath, nil
        }
    }
    
    // Build command
    cmd := exec.CommandContext(ctx, c.goBinary, "build",
        "-o", outputPath,
        "-tags", joinTags(c.buildTags),
        "-ldflags", c.ldflags,
        mainFile,
    )
    
    cmd.Dir = c.workDir
    cmd.Env = append(os.Environ(),
        "GOOS=js",
        "GOARCH=wasm",
        "GO111MODULE=on",
    )
    
    var stdout, stderr bytes.Buffer
    cmd.Stdout = &stdout
    cmd.Stderr = &stderr
    
    start := time.Now()
    err := cmd.Run()
    buildTime := time.Since(start)
    
    if err != nil {
        log.Printf("Build failed with error: %v", err)
        log.Printf("Stderr output: %s", stderr.String())
        return "", fmt.Errorf("build failed: %v\n%s", err, stderr.String())
    }
    
    // Update cache
    entry := CacheEntry{
        OutputPath:   outputPath,
        Hash:         c.cache.calculateHash(mainFile, changedFiles),
        Dependencies: changedFiles,
        Timestamp:    time.Now(),
    }
    
    if err := c.cache.Set(mainFile, entry); err != nil {
        log.Printf("Warning: Failed to update cache: %v", err)
    }
    
    log.Printf("Built %s in %v", filepath.Base(outputPath), buildTime)
    return outputPath, nil
}

func (c *Compiler) buildFromCurrentDir(mainFile, outputPath string) (string, error) {
    // Alternative: build using the current directory approach
    relPath, err := filepath.Rel(c.workDir, mainFile)
    if err != nil {
        return "", fmt.Errorf("failed to get relative path: %v", err)
    }
    
    cmd := exec.Command("go", "build",
        "-o", outputPath,
        "-tags", joinTags(c.buildTags),
        "-ldflags", c.ldflags,
        "./"+filepath.ToSlash(relPath),
    )
    
    cmd.Dir = c.workDir
    cmd.Env = append(os.Environ(),
        "GOOS=js",
        "GOARCH=wasm",
        "GO111MODULE=on",
    )
    
    var stderr bytes.Buffer
    cmd.Stderr = &stderr
    
    if err := cmd.Run(); err != nil {
        return "", fmt.Errorf("fallback build failed: %v\n%s", err, stderr.String())
    }
    
    return outputPath, nil
}
