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

func printBetween(start, end, filename string) {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error opening file:-", err)
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

func runASH() {
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

	// Start a goroutine that prints "..." in a loop until it receives a signal on stopChan
	stopChan := make(chan struct{})
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
	err = cmd.Start() // Start the command and immediately return
	err = cmd.Wait()  // Wait for the command to finish
	close(stopChan)   // Signal the goroutine to stop printing dots
	if err != nil {
		fmt.Println("Error running ASH:", err)
		return
	}
}

func analyzeFile() {
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

	// Search for the results file again
	var filePath string
	err := filepath.Walk("ash_output/", func(path string, info os.FileInfo, err error) error {
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

	// Print specific sections of the output
	fmt.Println("\n" + "\033[1;31m" + "Git Secrets Output:" + "\033[0m")
	time.Sleep(3 * time.Second)
	printBetween(">>>>>> begin git secrets --scan result >>>>>>", "<<<<<< end git secrets --scan result <<<<<<", filePath)

	fmt.Println("\n" + "\033[1;31m" + "Grype Output:" + "\033[0m")
	time.Sleep(3 * time.Second)
	printBetween(">>>>>> Begin Grype output", "<<<<<< End Grype output", filePath)

	fmt.Println("\n" + "\033[1;31m" + "Bandit Output:" + "\033[0m")
	time.Sleep(3 * time.Second)
	printBetween(">>>>>> Begin Bandit output", "<<<<<< End Bandit output", filePath)

	fmt.Println("\n" + "\033[1;31m" + "Semgrep Output:" + "\033[0m")
	time.Sleep(3 * time.Second)
	printBetween(">>>>>> Begin Semgrep output", "<<<<<< End Semgrep output", filePath)

	fmt.Println("\n" + "\033[1;31m" + "Checkov Output:" + "\033[0m")
	time.Sleep(3 * time.Second)
	printBetween(">>>>>> begin checkov", "<<<<<< end checkov", filePath)

	fmt.Println("\n" + "\033[1;31m" + "npm-audit Output:" + "\033[0m")
	time.Sleep(3 * time.Second)
	printBetween(">>>>>> Begin npm audit", "<<<<<< End npm audit", filePath)
}

func main() {
	runASH()
	analyzeFile()
}
