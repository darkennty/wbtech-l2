package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"syscall"
)

func main() {
	fmt.Println("Unix Shell Interpreter")
	fmt.Println("Commands: cd, pwd, echo, kill, ps, exit")

	currentPath, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	homePath := currentPath

	in := bufio.NewReader(os.Stdin)
	var line string
	var commands []string

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT) // Ctrl+C
	var externalCmd *exec.Cmd

	go func() {
		for {
			<-sigChan
			if externalCmd != nil && externalCmd.Process != nil {
				externalCmd.Process.Signal(syscall.SIGINT)
				fmt.Println()
			} else {
				fmt.Fprintln(os.Stdout, "")
				pathToPrint := strings.Replace(currentPath, homePath, "~", -1)
				fmt.Print(pathToPrint + " $ ")
			}
		}
	}()

Outer:
	for line != "exit" {
		externalCmd = nil

		pathToPrint := strings.Replace(currentPath, homePath, "~", -1)
		fmt.Print(pathToPrint + " $ ")

		line, err = in.ReadString('\n')
		if err != nil {
			if err == io.EOF { // Ctrl+D
				break
			}
			log.Println("shell: error scanning command: ", err)
			continue
		}
		line = strings.TrimSpace(line)

		if len(line) == 0 {
			continue
		}

		errorOccurred := false
		orCommands := strings.Split(line, "||")

		for c, orCommand := range orCommands {
			if c != 0 && !errorOccurred {
				break
			}
			errorOccurred = false

			andCommands := strings.Split(orCommand, "&&")

			for _, andCommand := range andCommands {
				if errorOccurred {
					break
				}

				commands = strings.Split(andCommand, "|")

				var prevOutput io.Reader

				for i, command := range commands {
					cmd := strings.Split(strings.TrimSpace(command), " ")

					stdin, stdout, _ := os.Pipe()

					switch cmd[0] {
					case "pwd":
						if i+1 == len(commands) {
							fmt.Fprintln(os.Stdout, currentPath)
						} else {
							fmt.Fprintln(stdout, currentPath)
						}
					case "cd":
						if len(commands) == 1 {
							var newPath string
							if len(cmd) < 2 {
								newPath = homePath
							} else {
								newPath = cmd[1]
							}

							err = os.Chdir(newPath)
							if err != nil {
								fmt.Printf("shell: cd: %s: No such file or directory\n", newPath)
							} else {
								currentPath, _ = os.Getwd()
							}
						}
					case "echo":
						for i := 1; i < len(cmd); i++ {
							cmd[i] = strings.Trim(cmd[i], "\"")
						}

						res := strings.Join(cmd[1:], " ")

						if i+1 == len(commands) {
							fmt.Fprintln(os.Stdout, res)
						} else {
							fmt.Fprintln(stdout, res)
						}
					case "ps":
						var toRun *exec.Cmd
						if runtime.GOOS == "windows" {
							toRun = exec.Command("tasklist")
						} else {
							toRun = exec.Command(cmd[0], cmd[1:]...)
						}

						toRun.Stdin = os.Stdin
						toRun.Stderr = os.Stderr

						if i+1 == len(commands) {
							toRun.Stdout = os.Stdout
						} else {
							toRun.Stdout = stdout
						}

						err = toRun.Run()
						if err != nil {
							log.Println("shell: error getting running processes: ", err)
							errorOccurred = true
							break
						}
						toRun.Wait()
					case "kill":
						if len(cmd) < 2 {
							fmt.Println("usage: kill pid")
							errorOccurred = true
							break
						}
						pid := cmd[1]

						maxPid := 32767
						pidInt, err := strconv.Atoi(pid)

						if err != nil || pidInt > maxPid {
							fmt.Printf("shell: kill: %s: arguments must be process or job IDs\n", pid)
							errorOccurred = true
							break
						}

						process, err := os.FindProcess(pidInt)
						if err != nil {
							fmt.Printf("shell: kill: (%d) - No such process\n", pidInt)
							errorOccurred = true
							break
						}

						err = process.Kill()
						if err != nil {
							log.Println("shell: kill: error killing process:", err)
							errorOccurred = true
							break
						}
					case "exit":
						if len(commands) == 0 {
							break Outer
						}
					default:
						toRun := exec.Command(cmd[0], cmd[1:]...)

						if prevOutput != nil {
							toRun.Stdin = prevOutput
						} else {
							toRun.Stdin = os.Stdin
						}

						toRun.Stderr = os.Stderr

						if i+1 == len(commands) {
							toRun.Stdout = os.Stdout
						} else {
							toRun.Stdout = stdout
						}

						externalCmd = toRun

						cmdError := toRun.Run()
						if cmdError != nil {
							pattern := "executable file not found in"
							matched, _ := regexp.MatchString(pattern, cmdError.Error())
							if matched || runtime.GOOS == "windows" {
								fmt.Printf("shell: %s: not found\n", cmd[0])
							}

							errorOccurred = true
							break
						}
						toRun.Wait()
					}

					stdout.Close()
					prevOutput = stdin
				}
			}
		}
	}
}
