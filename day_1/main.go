package main

import (
	"fmt"
	"os"
	"strings"
)

//TIP <p>To run your code, right-click the code and select <b>Run</b>.</p> <p>Alternatively, click
// the <icon src="AllIcons.Actions.Execute"/> icon in the gutter and select the <b>Run</b> menu item from here.</p>

func main() {
	//to check the arguments if its proper
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <filePath>")
		return
	}
	filePath := os.Args[1]
	file, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		return
	}
	//to close the file after all the operations
	defer fmt.Printf("Done with analysing log\n")
	lines := strings.Split(string(file), "\n")
	error := 0
	warning := 0
	info := 0

	for _, line := range lines {

		switch {
		case strings.Contains(line, "ERROR"):
			error++
		case strings.Contains(line, "WARNING"):
			warning++
		case strings.Contains(line, "INFO"):
			info++
		}
	}

	fmt.Println("Log Summary")
	fmt.Printf("Errors: %v\n", error)
	fmt.Printf("Warnings: %v\n", warning)
	fmt.Printf("Infos: %v\n", info)
}
