// cmd/dev-server/build/compiler.go
package build

import (
    "bytes"
    "context"
    "fmt"
    "io"
    "os"
    "os/exec"
    "path/filepath"
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

func NewCompiler(workDir, outputDir string) *Compiler {
    return &Compiler{
        goBinary:  "go",
        workDir:   workDir,
        outputDir: outputDir,
        cache:     NewBuildCache(filepath.Join(outputDir, ".buildcache")),
        buildTags: []string{"dev"},
        ldflags:   "-s -w",
    }
}

// BuildWasm compiles Go to WebAssembly with caching
func (c *Compiler) BuildWasm(ctx context.Context, mainFile string, changedFiles []string) (string, error) {
    // Check cache for unchanged files
    if c.cache.IsValid(mainFile, changedFiles) {
        cachedPath := c.cache.GetCachedPath(mainFile)
        if cachedPath != "" {
            return cachedPath, nil
        }
    }
    
    outputPath := filepath.Join(c.outputDir, 
        fmt.Sprintf("app-%d.wasm", time.Now().UnixNano()))
    
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
    )
    
    var stdout, stderr bytes.Buffer
    cmd.Stdout = &stdout
    cmd.Stderr = &stderr
    
    start := time.Now()
    err := cmd.Run()
    buildTime := time.Since(start)
    
    if err != nil {
        return "", fmt.Errorf("build failed: %v\n%s", err, stderr.String())
    }
    
    // Update cache
    c.cache.Update(mainFile, outputPath, changedFiles)
    
    log.Printf("Built %s in %v", filepath.Base(outputPath), buildTime)
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
    tempMain := c.createSelectiveMain(affectedPackages)
    defer os.Remove(tempMain)
    
    return c.BuildWasm(ctx, tempMain, changedFiles)
}

