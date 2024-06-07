package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

// This program runs the ASH tool and processes its output to extract and print
func printBetween(start, end, filename string) {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	inSection := false

	for scanner.Scan() {
		line := scanner.Text()

		if inSection && strings.Contains(line, end) {
			break
		}

		if inSection {
			fmt.Println(line)
		}

		if strings.Contains(line, start) {
			inSection = true
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
	}
}

func main() {
	// Check if ASH is installed
	_, err := exec.LookPath("ash")
	if err != nil {
		fmt.Println("Warning: ASH not installed, please install.")
		return
	}

	// Run 'ash -version'
	cmd := exec.Command("ash", "-v")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		fmt.Println("Error running 'ash -version':", err)
		return
	}

	// Define the patterns to look for
	patterns := []string{
		"ERROR: .*",
		"WARNING: .*",
		"Code Findings",
		// "High",
		// "Medium",
		// "Low",
		// Add more patterns here...
	}

	// Compile the patterns into regular expressions
	regexps := make([]*regexp.Regexp, len(patterns))
	for i, pattern := range patterns {
		regexps[i] = regexp.MustCompile(pattern)
	}

	// Create a channel to signal the loading goroutine to stop
	stopChan := make(chan bool)

	// Start a goroutine that prints "..." in a loop until it receives a signal on stopChan
	go func() {
		fmt.Print("Running ASH ")
		for {
			for i := 0; i < 3; i++ {
				select {
				case <-stopChan:
					fmt.Println()
					return
				default:
					fmt.Print(".")
					time.Sleep(500 * time.Millisecond)
				}
			}
			// Print spaces to overwrite the dots
			fmt.Print("\b\b\b   \b\b\b")
		}
	}()

	// Run the ASH tool
	cmd = exec.Command("ash")
	err = cmd.Run()

	// When the ASH tool is done, send a signal on stopChan to stop the loading goroutine
	stopChan <- true

	// Search for the results file again
	var filePath string
	err = filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if info.Name() == "aggregated_results.txt" {
			filePath = path
			return filepath.SkipDir
		}
		return nil
	})
	if err != nil {
		fmt.Println("Error searching for file:", err)
		return
	}

	// Open the results file
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	// Read the file line by line
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		// Check each line against each pattern
		for _, re := range regexps {
			if re.MatchString(line) {
				// If the line matches the pattern, capture and process it
				fmt.Println(line)
				break
			}
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
	}
	fmt.Println("Grype Output:")
	printBetween(">>>>>> Begin Grype output", "<<<<<< End Grype output", "./ash_output/aggregated_results.txt") // Need to change so file path doesn't have to be specified.
}
