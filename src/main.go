package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/user"
	"strings"

	"github.com/k3y0708/otter/terminal"
)

const (
	Prompt       = "otter> "
	OtterHistory = ".otterhistory"
	OtterRC      = ".otterrc"
	NAME         = "minify_tf"
	VERSION      = "v0.0.1"
)

func executeCommand(cmd string, args []string) {
	// Execute command
	switch cmd {
	case "exit":
		os.Exit(0)
	case "source":
		if len(args) == 0 {
			os.Stderr.WriteString("source: missing file operand\n")
			return
		}

		// Check if file exists
		if _, err := os.Stat(args[0]); os.IsNotExist(err) {
			os.Stderr.WriteString("source: file does not exist\n")
			return
		}

		sourceFile(args[0])
	case "about":
		nv := NAME + " " + VERSION
		nvLen := len(nv)
		spacesLen := 26 - nvLen
		for i := 0; i < spacesLen; i++ {
			nv += " "
		}

		fmt.Println("+----------------------------+")
		fmt.Println("| " + nv + "|")
		fmt.Println("| by @k3y0708                |")
		fmt.Println("| https://github.com/k3y0708 |")
		fmt.Println("+----------------------------+")
	case "version":
		fmt.Println(VERSION)
	default:
		//Check if command is a file
		cmd, err := findCommand(cmd)
		if err != nil {
			os.Stderr.WriteString(err.Error() + "\n")
			return
		}

		// If last argument is &, execute command in background
		if len(args) > 0 && args[len(args)-1] == "&" {
			// Execute command in background
		} else {
			out, err := exec.Command(cmd, args...).Output()
			if err != nil {
				os.Stderr.WriteString(err.Error() + "\n")
			}
			fmt.Printf("%s", out)
		}
	}
}

func findCommand(cmd string) (string, error) {
	if _, err := os.Stat(cmd); err == nil {
		// Execute file
		return cmd, nil
	} else {
		// Search for command in PATH
		path := os.Getenv("PATH")
		for _, dir := range strings.Split(path, ":") {
			// Check if command is in directory
			// If it is, execute it
			if _, err := os.Stat(dir + "/" + cmd); err == nil {
				return dir + "/" + cmd, nil
			}
		}
		return "", fmt.Errorf("command not found")
	}
}

func replaceEnvVars(str string) string {

	// If string is single quote, return string
	if len(str) > 1 && str[0] == '\'' && str[len(str)-1] == '\'' {
		return str
	}

	// Replace with $HOME if last character is ~
	if len(str) > 0 && str[len(str)-1] == '~' {
		userhome, err := os.UserHomeDir()
		if err == nil {
			str = strings.Replace(str, "~", userhome, -1)
		}
	}

	// Replace ~ with $HOME if ~ is not followed by a character
	if len(str) > 1 && str[0] == '~' && (str[1] < 'a' || str[1] > 'z') && (str[1] < 'A' || str[1] > 'Z') {
		userhome, err := os.UserHomeDir()
		if err == nil {
			str = strings.Replace(str, "~", userhome, 1)
		}
	}

	// Replace ~user with home directory of user
	if len(str) > 2 && str[0] == '~' && (str[1] >= 'a' && str[1] <= 'z' || str[1] >= 'A' && str[1] <= 'Z') {
		var username string
		for i := 1; i < len(str); i++ {
			if str[i] == '/' {
				break
			} else {
				username += string(str[i])
			}
		}
		usr, err := user.Lookup(username)
		if err == nil {
			str = strings.Replace(str, "~"+username, usr.HomeDir, 1)
		}
	}

	// Replace $VAR with value of VAR
	// (first char is $, second char is [a-zA-Z_] and rest are [a-zA-Z0-9_])
	for i := 0; i < len(str); i++ {
		if str[i] == '$' && str[i+1] != '{' {
			var varName string
			for j := i + 1; j < len(str); j++ {
				if (str[j] < 'a' || str[j] > 'z') && (str[j] < 'A' || str[j] > 'Z') && (str[j] < '0' || str[j] > '9') && str[j] != '_' {
					break
				} else {
					varName += string(str[j])
				}
			}
			varValue := os.Getenv(varName)
			str = strings.Replace(str, "$"+varName, varValue, 1)
		}
	}

	// Replace ${VAR} with value of VAR
	// (first char is $, second char is {, third char is [a-zA-Z_] and rest are [a-zA-Z0-9_], last char is })
	for i := 0; i < len(str); i++ {
		if str[i] == '$' && str[i+1] == '{' {
			var varName string
			for j := i + 2; j < len(str); j++ {
				if str[j] == '}' {
					break
				} else {
					varName += string(str[j])
				}
			}
			varValue := os.Getenv(varName)
			str = strings.Replace(str, "${"+varName+"}", varValue, 1)
		}
	}

	return str
}

func stringToCommandArgs(str string) (string, []string) {
	var cmd string
	var args []string
	var inSingleQuotes bool
	var inDoubleQuotes bool
	var currentArg string

	for _, char := range str {
		if char == '"' && !inSingleQuotes {
			inDoubleQuotes = !inDoubleQuotes
		} else if char == '\'' && !inDoubleQuotes {
			inSingleQuotes = !inSingleQuotes
		} else if char == ' ' && !inSingleQuotes && !inDoubleQuotes {
			if cmd == "" {
				cmd = currentArg
			} else {
				args = append(args, currentArg)
			}
			currentArg = ""
		} else {
			currentArg += string(char)
		}
	}

	if cmd == "" {
		cmd = currentArg
	} else {
		args = append(args, currentArg)
	}

	for i, arg := range args {
		arg = replaceEnvVars(arg)
		args[i] = arg
	}
	return cmd, args
}

func sourceFile(path string) {
	// Panic if file does not exist
	if _, err := os.Stat(path); os.IsNotExist(err) {
		panic(err)
	}

	file, err := os.Open(OtterRC)
	if err != nil {
		panic(err)
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		cmd, args := stringToCommandArgs(scanner.Text())
		executeCommand(cmd, args)
	}
	file.Close()
}

func main() {
	// Get user home directory
	userhome, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	otterHistory := userhome + OtterHistory
	otterRC := userhome + OtterRC

	term, err := terminal.NewWithStdInOut()
	if err != nil {
		panic(err)
	}

	// If .otterhistory file does not exist, create it else read it
	if _, err := os.Stat(otterHistory); os.IsNotExist(err) {
		file, err := os.Create(otterHistory)
		if err != nil {
			panic(err)
		}
		file.Close()
	} else {
		file, err := os.Open(otterHistory)
		if err != nil {
			panic(err)
		}
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			term.AddToHistory(scanner.Text())
		}
		file.Close()
	}

	// Welcome message
	fmt.Println("Welcome to Otter 🦦 Shell!")
	term.SetPrompt(Prompt)

	// If .otterrc file does not exist, create it else read it
	if _, err := os.Stat(otterRC); os.IsNotExist(err) {
		file, err := os.Create(otterRC)
		if err != nil {
			panic(err)
		}
		file.Close()
	} else {
		sourceFile(otterRC)
	}

	line, err := term.ReadLine()

	for {
		// If command is EOF (but not empty string), exit
		if err == io.EOF {
			term.Write([]byte(line))
			fmt.Println()
			return
		}
		if !((err != nil && strings.Contains(err.Error(), "control-c break")) || len(line) == 0) {
			// Add command to .otterhistory file
			file, err := os.OpenFile(otterHistory, os.O_APPEND|os.O_WRONLY, 0644)
			if err != nil {
				panic(err)
			}
			file.WriteString(line + "\n")
			file.Close()

			// Parse command
			cmd, args := stringToCommandArgs(line)
			executeCommand(cmd, args)
		}
		line, err = term.ReadLine()
	}
}
