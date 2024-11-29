package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"
)

// MinificationStats holds statistics about the minification process
type MinificationStats struct {
	InputFile     string  `json:"input_file"`
	OutputFile    string  `json:"output_file"`
	OriginalSize  int     `json:"original_size"`
	MinifiedSize  int     `json:"minified_size"`
	Reduction     float64 `json:"reduction_percentage"`
	ProcessTime   float64 `json:"process_time_ms"`
}

// Minifier handles JavaScript minification
type Minifier struct {
	input           string
	preserveLicense bool
	shortenVars     bool
	varMap          map[string]string
	varCounter      int
}

// NewMinifier creates a new minifier instance
func NewMinifier(input string, preserveLicense, shortenVars bool) *Minifier {
	return &Minifier{
		input:           input,
		preserveLicense: preserveLicense,
		shortenVars:    shortenVars,
		varMap:         make(map[string]string),
		varCounter:     0,
	}
}

// generateVarName generates short variable names (a, b, c, ... z, a1, b1, ...)
func (m *Minifier) generateVarName() string {
	alphabet := "abcdefghijklmnopqrstuvwxyz"
	suffix := m.varCounter / 26
	char := alphabet[m.varCounter%26]
	m.varCounter++
	if suffix == 0 {
		return string(char)
	}
	return fmt.Sprintf("%c%d", char, suffix)
}

// shortenVariableNames replaces variable names with shorter versions
func (m *Minifier) shortenVariableNames(code string) string {
	// Preserve strings
	stringLiterals := make(map[string]string)
	re := regexp.MustCompile(`"[^"]*"|'[^']*'`)
	code = re.ReplaceAllStringFunc(code, func(s string) string {
		placeholder := fmt.Sprintf("__STR_%d__", len(stringLiterals))
		stringLiterals[placeholder] = s
		return placeholder
	})

	// Find and replace variable declarations
	re = regexp.MustCompile(`\b(var|let|const)\s+([a-zA-Z_$][a-zA-Z0-9_$]*)\b`)
	code = re.ReplaceAllStringFunc(code, func(s string) string {
		parts := re.FindStringSubmatch(s)
		if len(parts) == 3 {
			original := parts[2]
			if _, exists := m.varMap[original]; !exists {
				m.varMap[original] = m.generateVarName()
			}
			return parts[1] + " " + m.varMap[original]
		}
		return s
	})

	// Replace variable usages
	for original, short := range m.varMap {
		re = regexp.MustCompile(`\b` + original + `\b`)
		code = re.ReplaceAllString(code, short)
	}

	// Restore strings
	for placeholder, str := range stringLiterals {
		code = strings.Replace(code, placeholder, str, -1)
	}

	return code
}

// Minify performs the minification process
func (m *Minifier) Minify() string {
	result := m.input

	// Preserve license comments if requested
	var licenseComment string
	if m.preserveLicense {
		re := regexp.MustCompile(`^/\*![\s\S]*?\*/`)
		license := re.FindString(result)
		if license != "" {
			licenseComment = license + "\n"
			result = re.ReplaceAllString(result, "")
		}
	}

	// Remove single-line comments
	re := regexp.MustCompile(`//.*`)
	result = re.ReplaceAllString(result, "")

	// Remove multi-line comments (except license)
	re = regexp.MustCompile(`/\*[\s\S]*?\*/`)
	result = re.ReplaceAllString(result, "")

	// Remove whitespace at the beginning and end of lines
	re = regexp.MustCompile(`^\s+|\s+$`)
	result = re.ReplaceAllString(result, "")

	// Replace multiple spaces with a single space
	re = regexp.MustCompile(`\s{2,}`)
	result = re.ReplaceAllString(result, " ")

	// Remove spaces around operators
	operators := []string{`\+`, `-`, `\*`, `/`, `=`, `<`, `>`, `!`, `\?`, `:`, `&`, `\|`, `;`, `,`}
	for _, op := range operators {
		re = regexp.MustCompile(`\s*` + op + `\s*`)
		result = re.ReplaceAllString(result, op)
	}

	// Remove unnecessary semicolons
	re = regexp.MustCompile(`;;+`)
	result = re.ReplaceAllString(result, ";")

	// Remove spaces after function keywords and parentheses
	re = regexp.MustCompile(`function\s+`)
	result = re.ReplaceAllString(result, "function")

	// Remove newlines
	re = regexp.MustCompile(`\n+`)
	result = re.ReplaceAllString(result, "")

	// Remove spaces after commas
	re = regexp.MustCompile(`,\s+`)
	result = re.ReplaceAllString(result, ",")

	// Remove spaces around brackets
	re = regexp.MustCompile(`\s*{\s*`)
	result = re.ReplaceAllString(result, "{")
	re = regexp.MustCompile(`\s*}\s*`)
	result = re.ReplaceAllString(result, "}")
	re = regexp.MustCompile(`\s*\[\s*`)
	result = re.ReplaceAllString(result, "[")
	re = regexp.MustCompile(`\s*\]\s*`)
	result = re.ReplaceAllString(result, "]")
	re = regexp.MustCompile(`\s*\(\s*`)
	result = re.ReplaceAllString(result, "(")
	re = regexp.MustCompile(`\s*\)\s*`)
	result = re.ReplaceAllString(result, ")")

	if m.shortenVars {
		result = m.shortenVariableNames(result)
	}

	if m.preserveLicense && licenseComment != "" {
		result = licenseComment + result
	}

	return result
}

// processFile minifies a single JavaScript file
func processFile(inputPath, outputPath string, preserveLicense, shortenVars bool, stats chan<- MinificationStats) {
	start := time.Now()

	content, err := ioutil.ReadFile(inputPath)
	if err != nil {
		log.Printf("Error reading %s: %v\n", inputPath, err)
		return
	}

	minifier := NewMinifier(string(content), preserveLicense, shortenVars)
	minified := minifier.Minify()

	if outputPath == "" {
		ext := filepath.Ext(inputPath)
		outputPath = strings.TrimSuffix(inputPath, ext) + ".min" + ext
	}

	err = ioutil.WriteFile(outputPath, []byte(minified), 0644)
	if err != nil {
		log.Printf("Error writing %s: %v\n", outputPath, err)
		return
	}

	stats <- MinificationStats{
		InputFile:     inputPath,
		OutputFile:    outputPath,
		OriginalSize:  len(content),
		MinifiedSize:  len(minified),
		Reduction:     float64(len(content)-len(minified)) / float64(len(content)) * 100,
		ProcessTime:   float64(time.Since(start).Microseconds()) / 1000.0,
	}
}

// watchDirectory monitors a directory for changes and minifies modified files
func watchDirectory(dir string, preserveLicense, shortenVars bool) {
	fileModTimes := make(map[string]time.Time)
	
	for {
		files, err := filepath.Glob(filepath.Join(dir, "*.js"))
		if err != nil {
			log.Printf("Error scanning directory: %v\n", err)
			continue
		}

		for _, file := range files {
			if strings.HasSuffix(file, ".min.js") {
				continue
			}

			info, err := os.Stat(file)
			if err != nil {
				continue
			}

			lastMod := fileModTimes[file]
			if info.ModTime().After(lastMod) {
				log.Printf("Processing modified file: %s\n", file)
				stats := make(chan MinificationStats, 1)
				processFile(file, "", preserveLicense, shortenVars, stats)
				stat := <-stats
				log.Printf("Reduced by %.2f%% (%d → %d bytes)\n", 
					stat.Reduction, stat.OriginalSize, stat.MinifiedSize)
				fileModTimes[file] = info.ModTime()
			}
		}

		time.Sleep(1 * time.Second)
	}
}

func main() {
	// Define command line flags
	inputFile := flag.String("input", "", "Input JavaScript file or directory")
	outputFile := flag.String("output", "", "Output file path (optional)")
	watchMode := flag.Bool("watch", false, "Watch mode - monitor directory for changes")
	preserveLicense := flag.Bool("preserve-license", false, "Preserve license comments")
	shortenVars := flag.Bool("shorten-vars", false, "Shorten variable names")
	jsonOutput := flag.Bool("json", false, "Output statistics in JSON format")
	flag.Parse()

	if *inputFile == "" {
		log.Fatal("Please provide an input file or directory using -input flag")
	}

	fileInfo, err := os.Stat(*inputFile)
	if err != nil {
		log.Fatalf("Error accessing input path: %v", err)
	}

	if fileInfo.IsDir() {
		if *watchMode {
			log.Printf("Watching directory: %s\n", *inputFile)
			watchDirectory(*inputFile, *preserveLicense, *shortenVars)
		} else {
			files, err := filepath.Glob(filepath.Join(*inputFile, "*.js"))
			if err != nil {
				log.Fatalf("Error scanning directory: %v", err)
			}

			var wg sync.WaitGroup
			stats := make(chan MinificationStats, len(files))

			for _, file := range files {
				if strings.HasSuffix(file, ".min.js") {
					continue
				}

				wg.Add(1)
				go func(file string) {
					defer wg.Done()
					processFile(file, "", *preserveLicense, *shortenVars, stats)
				}(file)
			}

			go func() {
				wg.Wait()
				close(stats)
			}()

			var allStats []MinificationStats
			for stat := range stats {
				allStats = append(allStats, stat)
				if !*jsonOutput {
					fmt.Printf("Processed %s:\n", stat.InputFile)
					fmt.Printf("  Output: %s\n", stat.OutputFile)
					fmt.Printf("  Reduction: %.2f%% (%d → %d bytes)\n", 
						stat.Reduction, stat.OriginalSize, stat.MinifiedSize)
					fmt.Printf("  Process time: %.2f ms\n\n", stat.ProcessTime)
				}
			}

			if *jsonOutput {
				jsonStats, _ := json.MarshalIndent(allStats, "", "  ")
				fmt.Println(string(jsonStats))
			}
		}
	} else {
		stats := make(chan MinificationStats, 1)
		processFile(*inputFile, *outputFile, *preserveLicense, *shortenVars, stats)
		stat := <-stats

		if *jsonOutput {
			jsonStats, _ := json.MarshalIndent(stat, "", "  ")
			fmt.Println(string(jsonStats))
		} else {
			fmt.Printf("Processed %s:\n", stat.InputFile)
			fmt.Printf("  Output: %s\n", stat.OutputFile)
			fmt.Printf("  Reduction: %.2f%% (%d → %d bytes)\n", 
				stat.Reduction, stat.OriginalSize, stat.MinifiedSize)
			fmt.Printf("  Process time: %.2f ms\n", stat.ProcessTime)
		}
	}
}
