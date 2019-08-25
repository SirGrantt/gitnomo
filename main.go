package main

import (
	"bytes"
	//	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func main() {
	//firstArg := flag.String("first", "", "the first value in the world")
	// flag.Parse()

	// if *firstArg == "" {
	// 	fmt.Println("hey, give me a value")
	// 	flag.PrintDefaults()
	// 	os.Exit(1)
	// }

	// cmd := exec.Command("git", "branch | grep \\* | cut -d ' ' -f2")
	// args := []string{"branch", "|", "grep", "\\*", "cut", "-d", "' '", "-f2"}
	cmd := exec.Command("git", "branch")
	//cmd.Stdin = strings.NewReader("branch | grep \\* | cut -d ' ' -f2")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()

	result := strings.Split(out.String(), "\n")

	var branchName string

	for i, val := range result {
		found := strings.Index(val, "*")
		if found > -1 {
			selectedBranch := result[i]
			splitFromStar := strings.Split(selectedBranch, "*")
			branchName = strings.TrimSpace(splitFromStar[1])
		}
	}

	if err != nil {
		fmt.Println("Error running command")
		fmt.Println(err.Error())
		os.Exit(1)
	}
	fmt.Println(branchName)
}
