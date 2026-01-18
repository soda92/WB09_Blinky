package internal

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func FileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func CopyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	// Create destination directory if it doesn't exist
	if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
		return err
	}

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	return err
}

func CopyDir(src, dst string) error {
	// Ensure destination exists
	if err := os.MkdirAll(dst, 0755); err != nil {
		return err
	}

	// Use 'cp -r src/.' to copy contents
	// We assume Linux/Unix environment where cp is available
	srcArg := strings.TrimRight(src, "/") + "/."

	cmd := exec.Command("cp", "-r", srcArg, dst)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func AppendToFile(path, content string) error {
	f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.WriteString(content)
	return err
}

func ReplaceInFile(path, oldString, newString string) error {
	input, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	output := strings.Replace(string(input), oldString, newString, -1)

	err = os.WriteFile(path, []byte(output), 0644)
	return err
}

func FindInPaths(filename string, paths []string) string {
	for _, p := range paths {
		fullPath := filepath.Join(p, filename)
		if FileExists(fullPath) {
			return fullPath
		}
	}
	return ""
}

func FindFileInRepo(root, filename string) string {
	var foundPath string
	filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && info.Name() == filename {
			foundPath = path
			return fmt.Errorf("Found") // Stop searching
		}
		return nil
	})
	return foundPath
}
