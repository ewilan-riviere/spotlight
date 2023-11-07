package terminal

import (
	"bytes"
	"fmt"
	"os/exec"
	"runtime"
	"strings"
)

func SelectCommand(commands map[string]string, defaultCommand string) string {
	os := runtime.GOOS
	command := defaultCommand

	if _, ok := commands[os]; ok {
		command = commands[os]
	}

	return command
}

func ExecCommand(command string) string {
	var output string
	if strings.Contains(command, "|") {
		output = ExecPipeCommand(command)
	} else {
		output = ExecClassicCommand(command)
	}

	return output
}

func ExecClassicCommand(command string) string {
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

func ExecPipeCommand(command string) string {
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
