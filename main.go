package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"
)

var debugFile *os.File

func init() {
	var err error
	debugFile, err = os.OpenFile("debug.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to open debug file: %v\n", err)
		return
	}
}

func debugLog(format string, args ...interface{}) {
	if debugFile != nil {
		fmt.Fprintf(debugFile, format+"\n", args...)
	}
}

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
	debugLog("DEBUG: Minify function called")
	result := m.input
	debugLog("Initial input: %s", result)

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
	debugLog("After license preservation: %s", result)

	// Remove single-line comments
	re := regexp.MustCompile(`//.*`)
	result = re.ReplaceAllString(result, "")
	debugLog("After removing single-line comments: %s", result)

	// Remove multi-line comments (except license)
	re = regexp.MustCompile(`/\*[\s\S]*?\*/`)
	result = re.ReplaceAllString(result, "")
	debugLog("After removing multi-line comments: %s", result)

	// Remove whitespace at the beginning and end of lines
	re = regexp.MustCompile(`^\s+|\s+$`)
	result = re.ReplaceAllString(result, "")
	debugLog("After trimming whitespace: %s", result)

	// Replace multiple spaces with a single space
	re = regexp.MustCompile(`\s{2,}`)
	result = re.ReplaceAllString(result, " ")
	debugLog("After replacing multiple spaces: %s", result)

	// Remove spaces around operators
	operators := []string{`+`, `-`, `*`, `/`, `=`, `<`, `>`, `!`, `?`, `:`, `&`, `|`, `;`, `,`}
	for _, op := range operators {
		re = regexp.MustCompile(`\s*` + regexp.QuoteMeta(op) + `\s*`)
		result = re.ReplaceAllString(result, op)
	}
	debugLog("After fixing operators: %s", result)

	// Remove unnecessary semicolons
	re = regexp.MustCompile(`;;+`)
	result = re.ReplaceAllString(result, ";")
	debugLog("After removing semicolons: %s", result)

	// Remove spaces after function keywords and parentheses
	re = regexp.MustCompile(`function\s+`)
	result = re.ReplaceAllString(result, "function ")

	// Fix spaces between function name and parentheses
	re = regexp.MustCompile(`([a-zA-Z0-9_$])\s*\(`)
	result = re.ReplaceAllString(result, "$1(")
	debugLog("After fixing function spacing: %s", result)

	// Remove newlines
	re = regexp.MustCompile(`\n+`)
	result = re.ReplaceAllString(result, "")
	debugLog("After removing newlines: %s", result)

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
	debugLog("After removing bracket spaces: %s", result)

	if m.shortenVars {
		result = m.shortenVariableNames(result)
		debugLog("After shortening variables: %s", result)
	}

	if m.preserveLicense && licenseComment != "" {
		result = licenseComment + result
	}

	debugLog("Final result: %s", result)
	return result
}

// processFile minifies a single JavaScript file
func processFile(inputPath, outputPath string, preserveLicense, shortenVars bool, stats chan<- MinificationStats) {
	debugLog("DEBUG: Processing file: %s", inputPath)
	
	start := time.Now()

	// Read input file
	content, err := ioutil.ReadFile(inputPath)
	if err != nil {
		debugLog("Error reading input file: %v", err)
		return
	}
	debugLog("File content: %s", string(content))

	minifier := NewMinifier(string(content), preserveLicense, shortenVars)
	minified := minifier.Minify()

	if outputPath == "" {
		ext := filepath.Ext(inputPath)
		outputPath = strings.TrimSuffix(inputPath, ext) + ".min" + ext
	}

	err = ioutil.WriteFile(outputPath, []byte(minified), 0644)
	if err != nil {
		debugLog("Error writing output file: %v", err)
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
			debugLog("Error scanning directory: %v", err)
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
				debugLog("Processing modified file: %s", file)
				stats := make(chan MinificationStats, 1)
				processFile(file, "", preserveLicense, shortenVars, stats)
				stat := <-stats
				debugLog("Reduced by %.2f%% (%d → %d bytes)", 
					stat.Reduction, stat.OriginalSize, stat.MinifiedSize)
				fileModTimes[file] = info.ModTime()
			}
		}

		time.Sleep(1 * time.Second)
	}
}

func main() {
	// Explicitly write to stderr
	debugLog("DEBUG: Minification process started")
	
	input := flag.String("input", "", "Input JavaScript file or directory")
	output := flag.String("output", "", "Output file or directory")
	preserveLicense := flag.Bool("preserve-license", false, "Preserve license comments")
	shortenVars := flag.Bool("shorten-vars", false, "Shorten variable names")
	jsonOutput := flag.Bool("json", false, "Output statistics in JSON format")
	watchMode := flag.Bool("watch", false, "Watch directory for changes")
	flag.Parse()

	// Debug: Print all flags and their values directly to stderr
	debugLog("DEBUG: Input: %s", *input)
	debugLog("DEBUG: Output: %s", *output)
	debugLog("DEBUG: Preserve License: %v", *preserveLicense)
	debugLog("DEBUG: Shorten Vars: %v", *shortenVars)
	debugLog("DEBUG: JSON Output: %v", *jsonOutput)
	debugLog("DEBUG: Watch Mode: %v", *watchMode)

	if *input == "" {
		debugLog("Please provide an input file or directory using -input flag")
		return
	}

	fileInfo, err := os.Stat(*input)
	if err != nil {
		debugLog("Error accessing input path: %v", err)
		return
	}

	if fileInfo.IsDir() {
		if *watchMode {
			debugLog("Watching directory: %s", *input)
			watchDirectory(*input, *preserveLicense, *shortenVars)
		} else {
			files, err := filepath.Glob(filepath.Join(*input, "*.js"))
			if err != nil {
				debugLog("Error scanning directory: %v", err)
				return
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
					debugLog("Processed %s:", stat.InputFile)
					debugLog("  Output: %s", stat.OutputFile)
					debugLog("  Reduction: %.2f%% (%d → %d bytes)", 
						stat.Reduction, stat.OriginalSize, stat.MinifiedSize)
					debugLog("  Process time: %.2f ms", stat.ProcessTime)
				}
			}

			if *jsonOutput {
				jsonStats, _ := json.MarshalIndent(allStats, "", "  ")
				debugLog("%s", string(jsonStats))
			}
		}
	} else {
		stats := make(chan MinificationStats, 1)
		processFile(*input, *output, *preserveLicense, *shortenVars, stats)
		stat := <-stats

		if *jsonOutput {
			jsonStats, _ := json.MarshalIndent(stat, "", "  ")
			debugLog("%s", string(jsonStats))
		} else {
			debugLog("Processed %s:", stat.InputFile)
			debugLog("  Output: %s", stat.OutputFile)
			debugLog("  Reduction: %.2f%% (%d → %d bytes)", 
				stat.Reduction, stat.OriginalSize, stat.MinifiedSize)
			debugLog("  Process time: %.2f ms", stat.ProcessTime)
		}
	}
}
