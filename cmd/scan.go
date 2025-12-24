package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

type ScanResult struct {
	Type string
	Path string
	Size int64
	Time time.Time
}

var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Scan for stale project dependencies",
	Run: func(cmd *cobra.Command, args []string) {
		path, _ := cmd.Flags().GetString("path")
		days, _ := cmd.Flags().GetInt("days")
		deep, _ := cmd.Flags().GetBool("deep")
		threshold := time.Duration(days) * 24 * time.Hour

		matchers := getMatchers()

		if deep {
			color.Red("DEEP SCAN ENABLED: Ignoring anchor file checks.\n")
		}
		color.Cyan(" TidyUp Scanning: %s (Older than %d days)\n", path, days)

		results := make(chan ScanResult)
		var wg sync.WaitGroup
		var totalSaved int64
		var scannedCount int64
		count := 0

		stopUI := make(chan bool)
		go func() {
			for {
				select {
				case <-stopUI:
					return
				default:
					fmt.Printf("\r Scanned: %d directories...", scannedCount)
					time.Sleep(100 * time.Millisecond)
				}
			}
		}()

		wg.Add(1)
		go func() {
			defer wg.Done()
			filepath.WalkDir(path, func(p string, d os.DirEntry, err error) error {
				scannedCount++
				if err != nil || !isSafe(p) {
					return filepath.SkipDir
				}
				if !d.IsDir() || d.Name() == ".git" {
					return nil
				}

				for _, m := range matchers {
					if d.Name() == m.TargetDir {
						if deep {
							// Deep Scan: Just check the folder's own age
							info, err := d.Info()
							if err == nil && time.Since(info.ModTime()) > threshold {
								size, _ := dirSize(p)
								results <- ScanResult{m.Name, p, size, info.ModTime()}
							}
							return filepath.SkipDir
						} else {
							// Standard Scan: Verify with Anchor File
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
			close(results)
		}()

		fmt.Println()
		for res := range results {
			count++
			totalSaved += res.Size
			fmt.Print("\r\033[K")
			color.Red("[STALE] ")
			fmt.Printf("%-10s ", res.Type)
			color.Green("%-10s ", formatSize(res.Size))
			color.HiBlack("| %s (%s)\n", res.Path, res.Time.Format("2006-01-02"))
		}

		stopUI <- true
		wg.Wait()
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
}
