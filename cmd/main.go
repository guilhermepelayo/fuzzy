package main

import (
	"bufio"
	"fmt"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type Result struct {
	FilePath   string
	LineNumber int
	Line       string
	Distance   int
}

const (
	colorReset = "\033[0m"
	bold       = "\033[1m"
	colorFile  = "\033[31m" // Red
	colorLine  = "\033[32m" // Green
	colorMatch = "\033[33m" // Yellow
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: fuzzy <search-term> <directory>")
		return
	}

	searchTerm := os.Args[1]
	directory := os.Args[2]

	var results []Result

	err := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		if !isTextFile(path) {
			return nil
		}

		fileResults, err := searchFile(path, searchTerm)
		if err != nil {
			fmt.Printf("Error reading file %s: %v\n", path, err)
			return nil
		}

		results = append(results, fileResults...)
		return nil
	})

	if err != nil {
		fmt.Printf("Error walking the path %s: %v\n", directory, err)
		return
	}

	printResults(results, searchTerm)
}

func printResults(results []Result, searchTerm string) {
	if len(results) == 0 {
		fmt.Println("No matches found.")
		return
	}

	resultsByFile := make(map[string][]Result)
	for _, result := range results {
		resultsByFile[result.FilePath] = append(resultsByFile[result.FilePath], result)
	}

	for filePath, fileResults := range resultsByFile {
		fmt.Printf("%s%s%s%s\n", bold, colorFile, filePath, colorReset)
		for _, result := range fileResults {
			highlightedLine := highlightMatch(result.Line, searchTerm)
			fmt.Printf("  %s%d%s: %s\n", colorLine, result.LineNumber, colorReset, highlightedLine)
		}
		fmt.Println()
	}
}

func isTextFile(filePath string) bool {
	// Check the file extension first
	ext := filepath.Ext(filePath)
	mimeType := mime.TypeByExtension(ext)
	if strings.HasPrefix(mimeType, "text/") {
		return true
	}

	// Fallback to detecting content type from the file content
	file, err := os.Open(filePath)
	if err != nil {
		return false
	}
	defer file.Close()

	// Read the first 512 bytes to detect content type
	buffer := make([]byte, 512)
	n, err := file.Read(buffer)
	if err != nil {
		return false
	}

	mimeType = http.DetectContentType(buffer[:n])
	return strings.HasPrefix(mimeType, "text/")
}

func searchFile(filePath, searchTerm string) ([]Result, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var results []Result
	scanner := bufio.NewScanner(file)
	lineNumber := 0
	for scanner.Scan() {
		line := scanner.Text()
		lineNumber++
		if strings.Contains(line, searchTerm) {
			results = append(results, Result{FilePath: filePath, LineNumber: lineNumber, Line: line, Distance: 0})
			continue
		}
		words := strings.Fields(line)
		for _, word := range words {
			distance := LevenshteinDistance(word, searchTerm)
			if distance <= len(searchTerm)/2 {
				results = append(results, Result{FilePath: filePath, LineNumber: lineNumber, Line: line, Distance: distance})
				break
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

func highlightMatch(line, searchTerm string) string {
	highlighted := line
	if strings.Contains(line, searchTerm) {
		highlighted = strings.ReplaceAll(line, searchTerm, fmt.Sprintf("%s%s%s", colorMatch, searchTerm, colorReset))
	} else {
		words := strings.Fields(line)
		for _, word := range words {
			distance := LevenshteinDistance(word, searchTerm)
			if distance <= len(searchTerm)/2 {
				highlighted = strings.ReplaceAll(line, word, fmt.Sprintf("%s%s%s", colorMatch, word, colorReset))
				break
			}
		}
	}
	return highlighted
}

func min(a, b, c int) int {
	if a < b {
		if a < c {
			return a
		}
		return c
	}
	if b < c {
		return b
	}
	return c
}

func LevenshteinDistance(a, b string) int {
	la := len(a)
	lb := len(b)
	if la == 0 {
		return lb
	}
	if lb == 0 {
		return la
	}

	matrix := make([][]int, la+1)
	for i := range matrix {
		matrix[i] = make([]int, lb+1)
	}

	for i := 0; i <= la; i++ {
		matrix[i][0] = i
	}
	for j := 0; j <= lb; j++ {
		matrix[0][j] = j
	}

	for i := 1; i <= la; i++ {
		for j := 1; j <= lb; j++ {
			cost := 1
			if a[i-1] == b[j-1] {
				cost = 0
			}
			matrix[i][j] = min(
				matrix[i-1][j]+1,
				matrix[i][j-1]+1,
				matrix[i-1][j-1]+cost,
			)
		}
	}

	return matrix[la][lb]
}
