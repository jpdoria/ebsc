package main

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// dirInterface represents the directory interface.
type dirInterface interface {
	getDateTime() string
	compressDirectory() error
	generateDirectoryName(prefix string) string
	createDirectory(appName string) (string, error)
	directoryExists(appName string) (bool, string)
}

// dirManager struct is used to manage directories.
type dirManager struct {
	time string
}

// getDateTime returns the current date and time in the format "YYYYMMDDHHMMSS".
func (dm *dirManager) getDateTime() string {
	// Get the current date and time.
	tn := time.Now().Format("20060102150405")
	dm.time = tn

	// Return the current date and time.
	return tn
}

// generateDirectoryName returns a directory name in the format "backup-YYYYMMDDHHMMSS".
func (dm *dirManager) generateDirectoryName(prefix string) string {
	// Generate the directory name using the prefix and the current date and time.
	directoryName := fmt.Sprintf("%v-backup-%v", prefix, dm.time)

	// Return the directory name.
	return directoryName
}

// compressDirectory compresses a directory.
func (dm *dirManager) compressDirectory() error {
	zFileName := fmt.Sprintf("./backup-%v.zip", dm.time)
	destFile, err := os.Create(zFileName)

	if err != nil {
		return err
	}

	z := zip.NewWriter(destFile)
	err = filepath.Walk("./backup", func(filePath string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		if err != nil {
			return err
		}

		relPath := strings.TrimPrefix(filePath, filepath.Dir(zFileName))
		zipFile, err := z.Create(relPath)

		if err != nil {
			return err
		}

		fsFile, err := os.Open(filePath)

		if err != nil {
			return err
		}

		_, err = io.Copy(zipFile, fsFile)

		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	err = z.Close()

	if err != nil {
		return err
	}

	return nil
}

// createDirectory creates a directory with the name "backup-YYYYMMDDHHMMSS".
func (dm *dirManager) createDirectory(appName string) (string, error) {
	// Generate the path using the generated directory name.
	path := "./backup/" + dm.generateDirectoryName(appName)

	// Create the directory using the generated directory name.
	if err := os.MkdirAll(path, 0755); err != nil {
		return "", err
	}

	// Return the path.
	return path, nil
}

// directoryExists checks if a directory exists.
func (dm *dirManager) directoryExists(appName string) (bool, string) {
	var exists bool
	var path string

	// List all the directories inside the "backup" directory.
	dirs, err := os.ReadDir("./backup")

	if err != nil {
		os.Mkdir("./backup", 0755)
	}

	// Loop through the directories.
	for _, d := range dirs {
		// Check if the directory name matches the application name.
		if strings.Contains(d.Name(), appName) {
			exists = true
			path = fmt.Sprintf("./backup/%v", d.Name())
		}
	}

	// Return true if the directory exists.
	return exists, path
}
