package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func main() {
	//firstArg := flag.String("first", "", "the first value in the world")
	flag.Parse()

	// if *firstArg == "" {
	// 	fmt.Println("hey, give me a value")
	// 	flag.PrintDefaults()
	// 	os.Exit(1)
	// }

	cmd := exec.Command("Get Git Branch")
	cmd.Stdin = strings.NewReader("git branch | grep \\* | cut -d ' ' -f2")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()

	if err != nil {
		fmt.Println("Error running command")
		fmt.Println(err.Error())
		os.Exit(1)
	}
	fmt.Println(out.String())
}
