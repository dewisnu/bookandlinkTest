package helper

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

func ContainsPathTraversal(filename string) bool {
	if filename == "." || filename == ".." ||
		filepath.IsAbs(filename) ||
		filepath.Clean(filename) != filename {
		return true
	}
	return false
}

func IsImage(filename string) bool {
	ext := filepath.Ext(filename)
	switch ext {
	case ".jpg", ".jpeg", ".png":
		return true
	default:
		return false
	}
}

// generateUniqueFilename generates a clean, unique filename
func GenerateUniqueFilename(original string) string {
	ext := filepath.Ext(original)
	name := strings.TrimSuffix(original, ext)
	slug := slugify(name)
	timestamp := time.Now().Unix() // or use uuid.New().String() for full uniqueness
	return fmt.Sprintf("%s-%d%s", slug, timestamp, ext)
}

// Slugify makes the filename lowercase, alphanumeric, and hyphen-separated
func slugify(name string) string {
	// Convert to lowercase
	name = strings.ToLower(name)
	// Replace spaces with dashes
	name = strings.ReplaceAll(name, " ", "-")
	// Remove any character that's not alphanumeric or dash
	reg := regexp.MustCompile("[^a-z0-9\\-]+")
	name = reg.ReplaceAllString(name, "")
	return name
}
