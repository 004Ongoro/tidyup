package cmd

import (
	"testing"
)

func TestIsSafe(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		expected bool
	}{
		{"System folder Windows", "C:\\Windows\\System32", false},
		{"AppData Windows", "C:\\Users\\User\\AppData\\Local", false},
		{"Library macOS", "/Users/user/Library/Application Support", false},
		{"Standard Project", "C:\\Users\\User\\Projects\\my-web-app", true},
		{"Nested safe project", "/home/user/code/go/src/tidyup", true},
		{"Blocked parent", "/Users/user/.vscode/extensions/some-plugin", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isSafe(tt.path); got != tt.expected {
				t.Errorf("isSafe() = %v, want %v for path %s", got, tt.expected, tt.path)
			}
		})
	}
}

func TestFormatSize(t *testing.T) {
	tests := []struct {
		bytes    int64
		expected string
	}{
		{500, "500 B"},
		{1024, "1.0 KB"},
		{1048576, "1.0 MB"},
		{1073741824, "1.0 GB"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			if got := formatSize(tt.bytes); got != tt.expected {
				t.Errorf("formatSize() = %v, want %v", got, tt.expected)
			}
		})
	}
}