package cmd

import (
	"fmt"
	"os/exec"
	"runtime"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var scheduleCmd = &cobra.Command{
	Use:   "schedule",
	Short: "Automate TidyUp to run in the background",
	Long: `Sets up a native system task (Task Scheduler on Windows, Cron on Unix) 
to run 'tidyup clean --force' at regular intervals.`,
	Run: func(cmd *cobra.Command, args []string) {
		remove, _ := cmd.Flags().GetBool("remove")
		interval, _ := cmd.Flags().GetString("interval")

		if remove {
			removeSchedule()
			return
		}

		setupSchedule(interval)
	},
}

func setupSchedule(interval string) {
	exe, err := GetExecutablePath()
	if err != nil {
		color.Red("Error finding executable: %v", err)
		return
	}

	// Logic varies by OS
	switch runtime.GOOS {
	case "windows":
		// Windows Task Scheduler command
		// Runs daily at 12:00 PM by default
		args := []string{
			"/Create", "/SC", "DAILY", "/TN", "TidyUpAutoClean",
			"/TR", fmt.Sprintf("%s clean --force", exe), "/ST", "12:00", "/F",
		}
		if interval == "weekly" {
			args[2] = "WEEKLY"
		}

		out, err := exec.Command("schtasks", args...).CombinedOutput()
		if err != nil {
			color.Red("Failed to create Windows Task: %v\nOutput: %s", err, string(out))
		} else {
			color.Green(" Windows Task 'TidyUpAutoClean' created successfully!")
		}

	case "darwin", "linux":
		// Cron implementation
		cronJob := ""
		if interval == "daily" {
			cronJob = fmt.Sprintf("0 12 * * * %s clean --force\n", exe)
		} else {
			cronJob = fmt.Sprintf("0 12 * * 0 %s clean --force\n", exe) // Weekly
		}

		cmd := exec.Command("bash", "-c", fmt.Sprintf("(crontab -l ; echo '%s') | crontab -", cronJob))
		if err := cmd.Run(); err != nil {
			color.Red("Failed to update crontab: %v", err)
		} else {
			color.Green(" Cron job added successfully!")
		}

	default:
		color.Yellow("Scheduling is not yet supported on %s", runtime.GOOS)
	}
}

func removeSchedule() {
	switch runtime.GOOS {
	case "windows":
		out, err := exec.Command("schtasks", "/Delete", "/TN", "TidyUpAutoClean", "/F").CombinedOutput()
		if err != nil {
			color.Red("Failed to remove task: %s", string(out))
		} else {
			color.Green(" Windows Task removed.")
		}
	case "darwin", "linux":
		cmd := exec.Command("bash", "-c", "crontab -l | grep -v 'tidyup clean --force' | crontab -")
		if err := cmd.Run(); err != nil {
			color.Red("Failed to clear crontab: %v", err)
		} else {
			color.Green("Cron job removed.")
		}
	}
}

func init() {
	rootCmd.AddCommand(scheduleCmd)
	scheduleCmd.Flags().BoolP("remove", "r", false, "Remove the automated schedule")
	scheduleCmd.Flags().StringP("interval", "i", "daily", "Interval for cleaning (daily/weekly)")
}