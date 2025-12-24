package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var blocklist = []string{
	"AppData", "Library", ".vscode", ".antigravity", ".rustup", ".cargo",
	"Program Files", "Windows", "System32", "node_modules",
}

func isSafe(path string) bool {
	parts := strings.Split(path, string(os.PathSeparator))
	for _, part := range parts {
		for _, blocked := range blocklist {
			// Blocking if a PARENT folder is in the blocklist.
			// allowing the folder itself to be 'node_modules' etc.
			if strings.EqualFold(part, blocked) && !isTargetDir(part) {
				return false
			}
		}
	}
	return true
}

func isTargetDir(name string) bool {
	for _, m := range matchers {
		if name == m.TargetDir {
			return true
		}
	}
	return false
}

func dirSize(path string) (int64, error) {
	var size int64
	err := filepath.WalkDir(path, func(_ string, d os.DirEntry, err error) error {
		if err != nil { return err }
		if !d.IsDir() {
			info, err := d.Info()
			if err != nil { return err }
			size += info.Size()
		}
		return nil
	})
	return size, err
}

func formatSize(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}