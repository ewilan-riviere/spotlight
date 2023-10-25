package health

import (
	"bytes"
	"fmt"
	"os/exec"
	"runtime"
	"strings"

	"github.com/ewilan-riviere/notifier/notify"
	"github.com/ewilan-riviere/spotlight/pkg/dotenv"
)

func DiskUsage(send bool) string {
	commands := map[string]string{
		"linux":   "df -hT /",
		"darwin":  "df -h /System/Volumes/Data",
		"windows": "wmic logicaldisk get size,freespace,caption",
	}

	command := selectCommand(commands, "df -h /")
	output := execCommand(command)
	notifier(send, "```"+output+"```")

	return output
}

func RamUsage(send bool) string {
	commands := map[string]string{
		"linux":   `echo "$(free -m | awk 'NR==2 { printf "%.2f GB\n", $3/1024 }') / $(free -m | awk 'NR==2 { printf "%.2f GB\n", $2/1024 }')"`,
		"darwin":  "top -l 1 -s 0 | awk '/PhysMem/'",
		"windows": `wmic OS get FreePhysicalMemory,TotalVisibleMemorySize /Value | awk -F"=" 'NR==1 { printf "%.2f GB\n", $2/1024/1024 } NR==2 { printf "%.2f GB\n", $2/1024/1024 }'`,
	}

	command := selectCommand(commands, "free -m")
	output := execCommand(command)
	notifier(send, "```"+output+"```")

	return output
}

func CpuUsage(send bool) string {
	commands := map[string]string{
		"linux":   `cat /proc/loadavg | awk '{print $1*100 "%"}'`,
		"darwin":  `top -l 2 | grep -E "^CPU"`,
		"windows": `wmic cpu get loadpercentage | awk 'NR==2 {print $1 "%"}'`,
	}

	command := selectCommand(commands, "top -l 2 | grep -E '^CPU'")
	output := execCommand(command)
	notifier(send, "```"+output+"```")

	return output
}

func BigFiles(size int, extensions []string, send bool) string {
	os := runtime.GOOS
	command := "find / -xdev -type f -size +SIZEM EXTENSIONS -exec du -sh {} ';' | sort -rh | head -n50 ;"
	commands := map[string]string{
		"linux": "find / -xdev -type f -size +SIZEM EXTENSIONS -exec du -sh {} ';' | sort -rh | head -n50 ;",
		// "darwin": "mdfind 'kMDItemFSSize > 30000000'",
		"darwin":  "find /Users -xdev -type f -size +SIZEM EXTENSIONS -exec du -sh {} ';' | sort -rh | head -n50 ;",
		"windows": "wmic logicaldisk get size,freespace,caption",
	}

	if _, ok := commands[os]; ok {
		command = commands[os]
	}

	if os == "linux" {
		command = strings.Replace(command, "SIZE", fmt.Sprintf("%d", size), 1)
	} else if os == "darwin" {
		command = strings.Replace(command, "SIZE", fmt.Sprintf("%d", size), 1)
	}

	if len(extensions) == 0 {
		command = strings.Replace(command, "EXTENSIONS ", "", 1)
	} else {
		files := ""
		for _, ext := range extensions {
			files += fmt.Sprintf("-not -name '*.%s' ", ext)
		}
		command = strings.Replace(command, "EXTENSIONS ", files, 1)
	}

	output := execCommand(command)
	notifier(send, output)

	return ""
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

func selectCommand(commands map[string]string, defaultCommand string) string {
	os := runtime.GOOS
	command := defaultCommand

	if _, ok := commands[os]; ok {
		command = commands[os]
	}

	return command
}

func execCommand(command string) string {
	var output string
	if strings.Contains(command, "|") {
		output = execPipeCommand(command)
	} else {
		output = execClassicCommand(command)
	}

	fmt.Print(output)

	return output
}

func execClassicCommand(command string) string {
	split := strings.Split(command, " ")
	name := split[0]
	args := split[1:]

	cmd := exec.Command(name, args...)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
		return ""
	}

	return out.String()
}

func execPipeCommand(command string) string {
	args := []string{"-c", command}
	cmd := exec.Command("bash", args...)

	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()

	if err != nil {
		fmt.Println(err.Error())
		return ""
	}

	return out.String()
}
