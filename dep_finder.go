package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const repoPath = "/home/soda/STM32Cube/Repository/STM32Cube_FW_WB0_V1.4.0"
const projectInc = "Core/Inc"
const projectSrc = "Core/Src"

var processedHeaders = make(map[string]bool)
var foundSources = make(map[string]bool)

func findFileInRepo(filename string) string {
	var foundPath string
	err := filepath.Walk(repoPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && info.Name() == filename {
			foundPath = path
			return fmt.Errorf("Found") // Stop searching
		}
		return nil
	})
	if err != nil && err.Error() == "Found" {
		return foundPath
	}
	return ""
}

func scanFile(path string) {
	file, err := os.Open(path)
	if err != nil {
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "#include") {
			parts := strings.Split(line, "\"")
			if len(parts) > 1 {
				header := parts[1]
				processHeader(header)
			}
		}
	}
}

func processHeader(header string) {
	if processedHeaders[header] {
		return
	}
	
	// Check if exists locally
	if _, err := os.Stat(filepath.Join(projectInc, header)); err == nil {
		processedHeaders[header] = true
		// Scan it for recursive includes
		scanFile(filepath.Join(projectInc, header))
		return
	}
	
	// Check if it exists in Middlewares (simplified check)
	// In a real tool we'd scan include paths. Here we assume manual copy targets Core/Inc.
	
	// Find in Repo
	path := findFileInRepo(header)
	if path != "" {
		fmt.Printf("Found missing header: %s -> Copying to %s\n", header, projectInc)
		// Copy command
		// In a real tool we'd do the copy. Here we output shell commands or do it.
		// Let's do it!
		copyFile(path, filepath.Join(projectInc, header))
		processedHeaders[header] = true
		
		// Scan the new header
		scanFile(filepath.Join(projectInc, header))

		// Check for corresponding .c file
		cFile := strings.Replace(header, ".h", ".c", 1)
		cPath := filepath.Join(filepath.Dir(path), cFile) // Check same dir first
		if _, err := os.Stat(cPath); os.IsNotExist(err) {
			// Try finding globally
			cPath = findFileInRepo(cFile)
		}
		
		if cPath != "" && !foundSources[cFile] {
			fmt.Printf("Found associated source: %s -> Copying to %s\n", cFile, projectSrc)
			copyFile(cPath, filepath.Join(projectSrc, cFile))
			foundSources[cFile] = true
			// Scan the new source
			scanFile(filepath.Join(projectSrc, cFile))
		}
	}
}

func copyFile(src, dst string) {
	input, err := os.ReadFile(src)
	if err != nil {
		fmt.Println("Error reading", src)
		return
	}
	err = os.WriteFile(dst, input, 0644)
	if err != nil {
		fmt.Println("Error writing", dst)
		return
	}
}

func main() {
	// Seed with existing source files in Core/Src to find their deps
	files, _ := filepath.Glob(filepath.Join(projectSrc, "*.c"))
	for _, f := range files {
		scanFile(f)
	}
	// Also scan bleplat.c explicitly if glob missed it or to be sure
	scanFile("Core/Src/bleplat.c")
}
