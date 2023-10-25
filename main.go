// Package spotlight to scan system like disk usage or big files.
//
// Examples/readme can be found on the GitHub page at https://github.com/ewilan-riviere/spotlight
//
// If you want to use it as CLI, you can install it with:
//
//	go install github.com/ewilan-riviere/spotlight
//
// Then you can use it like this:
//
//	spotlight disk-usage
//	spotlight big-files -e=jpg -e=png -s=50
package main

import (
	"fmt"

	"github.com/ewilan-riviere/spotlight/pkg/health"
	"github.com/spf13/cobra"
)

func main() {
	var cmdDiskUsage = &cobra.Command{
		Use:   "disk",
		Short: "Check disk usage.",
		Long:  `Check disk usage.`,
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			notify, _ := cmd.Flags().GetBool("notify")
			health.DiskUsage(notify)
		},
	}

	var cmdRamUsage = &cobra.Command{
		Use:   "ram",
		Short: "Check RAM usage.",
		Long:  `Check RAM usage.`,
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			notify, _ := cmd.Flags().GetBool("notify")
			health.RamUsage(notify)
		},
	}

	var cmdCpuUsage = &cobra.Command{
		Use:   "cpu",
		Short: "Check CPU usage.",
		Long:  `Check CPU usage.`,
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			notify, _ := cmd.Flags().GetBool("notify")
			health.CpuUsage(notify)
		},
	}

	var cmdHealth = &cobra.Command{
		Use:   "health",
		Short: "Check disk, RAM and CPU usage.",
		Long:  `Check disk, RAM and CPU usage.`,
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			notify, _ := cmd.Flags().GetBool("notify")
			health.DiskUsage(notify)
			fmt.Println("")
			health.RamUsage(notify)
			fmt.Println("")
			health.CpuUsage(notify)
		},
	}

	var cmdBigFiles = &cobra.Command{
		Use:   "files",
		Short: "Check big files.",
		Long:  `Check big files, by default bigger than 100M.`,
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			notify, _ := cmd.Flags().GetBool("notify")
			size, _ := cmd.Flags().GetInt("size")
			exts, _ := cmd.Flags().GetStringArray("exts")

			health.BigFiles(size, exts, notify)
		},
	}

	cmdDiskUsage.Flags().BoolP("notify", "n", false, "Send notification to Discord or Slack.")
	cmdRamUsage.Flags().BoolP("notify", "n", false, "Send notification to Discord or Slack.")
	cmdCpuUsage.Flags().BoolP("notify", "n", false, "Send notification to Discord or Slack.")
	cmdHealth.Flags().BoolP("notify", "n", false, "Send notification to Discord or Slack.")
	cmdBigFiles.Flags().IntP("size", "s", 100, "File size to consider in M, default 100M.")
	cmdBigFiles.Flags().StringArrayP("exts", "e", []string{""}, "File extensions to consider. Can be repeated for multiple extensions.")
	cmdBigFiles.Flags().BoolP("notify", "n", false, "Send notification to Discord or Slack.")

	var rootCmd = &cobra.Command{Use: "spotlight"}
	rootCmd.AddCommand(cmdDiskUsage)
	rootCmd.AddCommand(cmdRamUsage)
	rootCmd.AddCommand(cmdCpuUsage)
	rootCmd.AddCommand(cmdHealth)
	rootCmd.AddCommand(cmdBigFiles)
	rootCmd.Execute()
}
