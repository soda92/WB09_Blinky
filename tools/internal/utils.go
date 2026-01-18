package internal

import (
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
	// Using 'cp -r' for simplicity and speed on Linux
	cmd := exec.Command("cp", "-r", src, dst)
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
