// Package spotlight to scan system for system health.
//
// Examples/readme can be found on the GitHub page at https://github.com/ewilan-riviere/spotlight
//
// If you want to use it as CLI, you can install it with:
//
//	go install github.com/ewilan-riviere/spotlight
//
// Then you can use it like this:
//
// Check disk usage.
//
//	spotlight disk
//
// Check RAM usage.
//
//	spotlight ram
//
// Check CPU usage.
//
//	spotlight cpu
//
// Check disk, RAM and CPU usage.
//
//	spotlight health
//
// Check big files.
//
//	spotlight files -e=jpg -e=png -s=50
//
// Check websites
//
//	spotlight ping -d=example.com -d=example2.com
package main

import (
	"fmt"

	"github.com/ewilan-riviere/notifier/notify"
	"github.com/ewilan-riviere/spotlight/pkg/dotenv"
	"github.com/ewilan-riviere/spotlight/pkg/health"
	"github.com/ewilan-riviere/spotlight/pkg/ping"
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
			output := health.DiskUsage(notify)
			fmt.Println(output)
		},
	}

	var cmdRamUsage = &cobra.Command{
		Use:   "ram",
		Short: "Check RAM usage.",
		Long:  `Check RAM usage.`,
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			notify, _ := cmd.Flags().GetBool("notify")
			output := health.RamUsage(notify)

			fmt.Println(output)
		},
	}

	var cmdCpuUsage = &cobra.Command{
		Use:   "cpu",
		Short: "Check CPU usage.",
		Long:  `Check CPU usage.`,
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			notify, _ := cmd.Flags().GetBool("notify")
			output := health.CpuUsage(notify)

			fmt.Println(output)
		},
	}

	var cmdHealth = &cobra.Command{
		Use:   "health",
		Short: "Check disk, RAM and CPU usage.",
		Long:  `Check disk, RAM and CPU usage.`,
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			notify, _ := cmd.Flags().GetBool("notify")
			disk := health.DiskUsage(notify)
			ram := health.RamUsage(notify)
			cpu := health.CpuUsage(notify)

			fmt.Println(disk)
			fmt.Println(ram)
			fmt.Println(cpu)
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

			output := health.BigFiles(size, exts, notify)
			fmt.Println(output)
		},
	}

	var cmdPing = &cobra.Command{
		Use:   "ping",
		Short: "Ping a host.",
		Long:  `Ping a host.`,
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			notify, _ := cmd.Flags().GetBool("notify")
			domains, _ := cmd.Flags().GetStringArray("domains")

			if len(domains) == 0 {
				var domainsDotenv = dotenv.Make().Domains
				domains = domainsDotenv
			}
			websites := ping.Make(domains)

			for _, website := range websites {
				var output = ""
				fmt.Println(website.Domain)
				fmt.Println("Ping:", website.PingSuccess)
				fmt.Println("Curl:", website.CurlSuccess)
				fmt.Println("Http:", website.HttpCode)
				fmt.Println("HttpRedirect:", website.HttpRedirect)
				fmt.Println("UseHttp2:", website.UseHttp2)
				fmt.Println("UseHttps:", website.UseHttps)
				fmt.Println("Server:", website.Server)
				fmt.Println("IpAdress:", website.IpAdress)
				fmt.Println("Time:", website.Time)

				if website.Ok {
					output += "\n"
					output += "✅ " + website.Domain + "\n"
				} else {
					output += "\n"
					output += "❌ " + website.Domain + "\n"
					output += "Ping: " + fmt.Sprintf("%t", website.PingSuccess) + "\n"
					output += "Curl: " + fmt.Sprintf("%t", website.CurlSuccess) + "\n"
					output += "Http: " + fmt.Sprintf("%d", website.HttpCode) + "\n"
					output += "HttpRedirect: " + fmt.Sprintf("%t", website.HttpRedirect) + "\n"
					output += "UseHttp2: " + fmt.Sprintf("%t", website.UseHttp2) + "\n"
					output += "UseHttps: " + fmt.Sprintf("%t", website.UseHttps) + "\n"
					output += "Server: " + website.Server + "\n"
					output += "IpAdress: " + website.IpAdress + "\n"
					output += "Time: " + website.Time + "\n"
				}

				notifier(notify, "```"+output+"```")
				fmt.Println("")
			}
		},
	}

	cmdDiskUsage.Flags().BoolP("notify", "n", false, "Send notification to Discord or Slack.")
	cmdRamUsage.Flags().BoolP("notify", "n", false, "Send notification to Discord or Slack.")
	cmdCpuUsage.Flags().BoolP("notify", "n", false, "Send notification to Discord or Slack.")
	cmdHealth.Flags().BoolP("notify", "n", false, "Send notification to Discord or Slack.")
	cmdBigFiles.Flags().IntP("size", "s", 100, "File size to consider in M, default 100M.")
	cmdBigFiles.Flags().StringArrayP("exts", "e", []string{""}, "File extensions to consider. Can be repeated for multiple extensions.")
	cmdBigFiles.Flags().BoolP("notify", "n", false, "Send notification to Discord or Slack.")
	cmdPing.Flags().BoolP("notify", "n", false, "Send notification to Discord or Slack.")
	cmdPing.Flags().StringArrayP("domains", "d", []string{""}, "Domains to ping. Can be repeated for multiple domains.")

	var rootCmd = &cobra.Command{Use: "spotlight"}
	rootCmd.AddCommand(cmdDiskUsage)
	rootCmd.AddCommand(cmdRamUsage)
	rootCmd.AddCommand(cmdCpuUsage)
	rootCmd.AddCommand(cmdHealth)
	rootCmd.AddCommand(cmdBigFiles)
	rootCmd.AddCommand(cmdPing)
	rootCmd.Execute()
}

func notifier(send bool, output string) {
	if send {
		de := dotenv.Make()
		for _, service := range de.Services {
			if service == "discord" {
				notify.Notifier(output, de.DiscordWebhook)
			}
			if service == "slack" {
				notify.Notifier(output, de.SlackWebhook)
			}
		}
	}
}
