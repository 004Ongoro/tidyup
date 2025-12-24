package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

type ScanResult struct {
	Type string    `json:"type"`
	Path string    `json:"path"`
	Size int64     `json:"size"`
	Time time.Time `json:"last_modified"`
}

var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Scan for stale project dependencies",
	Run: func(cmd *cobra.Command, args []string) {
		path, _ := cmd.Flags().GetString("path")
		days, _ := cmd.Flags().GetInt("days")
		deep, _ := cmd.Flags().GetBool("deep")
		jsonMode, _ := cmd.Flags().GetBool("json")
		threshold := time.Duration(days) * 24 * time.Hour

		matchers := getMatchers()

		results := make(chan ScanResult, 100)
		var resultsSlice []ScanResult
		var wg sync.WaitGroup
		var totalSaved int64
		var scannedCount int64
		count := 0

		// UI Goroutine
		stopUI := make(chan bool)
		if !jsonMode {
			if deep {
				color.Red("DEEP SCAN ENABLED: Ignoring anchor file checks.\n")
			}
			color.Cyan("TidyUp Scanning: %s (Older than %d days)\n", path, days)

			go func() {
				for {
					select {
					case <-stopUI:
						return
					default:
						fmt.Printf("\rScanned: %d directories...", atomic.LoadInt64(&scannedCount))
						time.Sleep(100 * time.Millisecond)
					}
				}
			}()
		}

		// Find immediate subdirectories to distribute to workers
		entries, err := os.ReadDir(path)
		if err != nil {
			if !jsonMode {
				color.Red("Error reading path: %v", err)
			}
			return
		}

		// Parallel Worker Pool
		workerLimit := runtime.NumCPU() * 2
		semaphore := make(chan struct{}, workerLimit)

		for _, entry := range entries {
			if !entry.IsDir() || !isSafe(entry.Name()) {
				continue
			}

			wg.Add(1)
			go func(subDir string) {
				defer wg.Done()
				semaphore <- struct{}{}
				defer func() { <-semaphore }()

				fullSubPath := filepath.Join(path, subDir)
				filepath.WalkDir(fullSubPath, func(p string, d os.DirEntry, err error) error {
					atomic.AddInt64(&scannedCount, 1)
					if err != nil || !isSafe(p) {
						return filepath.SkipDir
					}
					if !d.IsDir() || d.Name() == ".git" {
						return nil
					}

					for _, m := range matchers {
						if d.Name() == m.TargetDir {
							if deep {
								info, err := d.Info()
								if err == nil && time.Since(info.ModTime()) > threshold {
									size, _ := dirSize(p)
									results <- ScanResult{m.Name, p, size, info.ModTime()}
								}
								return filepath.SkipDir
							} else {
								anchor := filepath.Join(filepath.Dir(p), m.AnchorFile)
								if info, err := os.Stat(anchor); err == nil {
									if time.Since(info.ModTime()) > threshold {
										size, _ := dirSize(p)
										results <- ScanResult{m.Name, p, size, info.ModTime()}
									}
									return filepath.SkipDir
								}
							}
						}
					}
					return nil
				})
			}(entry.Name())
		}

		// Result processing
		go func() {
			wg.Wait()
			close(results)
			if !jsonMode {
				stopUI <- true
			}
		}()

		if !jsonMode {
			fmt.Println()
		}

		for res := range results {
			count++
			atomic.AddInt64(&totalSaved, res.Size)

			if jsonMode {
				resultsSlice = append(resultsSlice, res)
			} else {
				fmt.Print("\r\033[K")
				color.Red("[STALE] ")
				fmt.Printf("%-10s ", res.Type)
				color.Green("%-10s ", formatSize(res.Size))
				color.HiBlack("| %s (%s)\n", res.Path, res.Time.Format("2006-01-02"))
			}
		}

		if jsonMode {
			output, _ := json.MarshalIndent(resultsSlice, "", "  ")
			fmt.Println(string(output))
			return
		}

		fmt.Printf("\r\033[K")
		fmt.Println(strings.Repeat("-", 60))
		color.HiYellow("Summary: Found %d folders | Total Space: %s", count, formatSize(totalSaved))
	},
}

func init() {
	rootCmd.AddCommand(scanCmd)
	scanCmd.Flags().StringP("path", "p", ".", "Path to scan")
	scanCmd.Flags().IntP("days", "d", 30, "Age threshold")
	scanCmd.Flags().Bool("deep", false, "Perform deep scan (ignore anchor files)")
	scanCmd.Flags().Bool("json", false, "Output results in JSON format")
}
