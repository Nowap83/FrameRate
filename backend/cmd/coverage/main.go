package main

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

func main() {
	fmt.Println("Running Go tests with coverage...")

	// 1. Run go test with coverage
	cmd := exec.Command("go", "test", "-coverprofile=coverage.out", "./...")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run() // Ignore errors as some tests might fail but we still want coverage

	// 2. Extract coverage percentage
	fmt.Println("Generating coverage report...")
	coverCmd := exec.Command("go", "tool", "cover", "-func=coverage.out")
	output, err := coverCmd.Output()
	if err != nil {
		fmt.Println("Error running go tool cover:", err)
		return
	}

	lines := strings.Split(string(output), "\n")
	var totalLine string
	for _, line := range lines {
		if strings.HasPrefix(line, "total:") {
			totalLine = line
			break
		}
	}

	if totalLine == "" {
		fmt.Println("Could not find total coverage in output")
		return
	}

	fields := strings.Fields(totalLine)
	if len(fields) < 3 {
		fmt.Println("Unexpected coverage output format")
		return
	}

	coverageStr := fields[len(fields)-1]
	coverage := strings.TrimRight(coverageStr, "%")
	fmt.Printf("Total Backend Coverage: %s%%\n", coverageStr)

	// 3. Update README.md
	readmePath := "../README.md"
	readmeBytes, err := os.ReadFile(readmePath)
	if err != nil {
		fmt.Println("Error reading README.md:", err)
		return
	}
	readmeContent := string(readmeBytes)

	// Determine color based on threshold
	color := "red"
	var covFloat float64
	fmt.Sscanf(coverage, "%f", &covFloat)
	if covFloat >= 80 {
		color = "brightgreen"
	} else if covFloat >= 60 {
		color = "yellow"
	} else if covFloat >= 40 {
		color = "orange"
	}

	badgeURL := fmt.Sprintf("https://img.shields.io/badge/Coverage-%s%%25-%s", coverageStr, color)
	badgeMarkdown := fmt.Sprintf("![Coverage](%s)", badgeURL)

	// Replace existing badge or insert at top
	re := regexp.MustCompile(`\!\[Coverage\]\(https:\/\/img\.shields\.io\/badge\/Coverage-[^\)]+\)`)

	if re.MatchString(readmeContent) {
		readmeContent = re.ReplaceAllString(readmeContent, badgeMarkdown)
	} else {
		readmeContent = badgeMarkdown + "\n\n" + readmeContent
	}

	err = os.WriteFile(readmePath, []byte(readmeContent), 0644)
	if err != nil {
		fmt.Println("Error updating README.md:", err)
		return
	}

	fmt.Println("README.md successfully updated with latest coverage badge!")
}
