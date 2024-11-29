package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"regexp"
	"strings"
)

// Minifier handles JavaScript minification
type Minifier struct {
	input string
}

// NewMinifier creates a new minifier instance
func NewMinifier(input string) *Minifier {
	return &Minifier{input: input}
}

// Minify performs the minification process
func (m *Minifier) Minify() string {
	result := m.input

	// Remove single-line comments
	re := regexp.MustCompile(`//.*`)
	result = re.ReplaceAllString(result, "")

	// Remove multi-line comments
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

	return result
}

func main() {
	// Define command line flags
	inputFile := flag.String("input", "", "Input JavaScript file path")
	outputFile := flag.String("output", "", "Output file path (optional)")
	flag.Parse()

	// Validate input file
	if *inputFile == "" {
		log.Fatal("Please provide an input file using -input flag")
	}

	// Read input file
	content, err := ioutil.ReadFile(*inputFile)
	if err != nil {
		log.Fatalf("Error reading input file: %v", err)
	}

	// Create minifier and process the content
	minifier := NewMinifier(string(content))
	minified := minifier.Minify()

	// Determine output file path
	outPath := *outputFile
	if outPath == "" {
		ext := filepath.Ext(*inputFile)
		outPath = strings.TrimSuffix(*inputFile, ext) + ".min" + ext
	}

	// Write to output file
	err = ioutil.WriteFile(outPath, []byte(minified), 0644)
	if err != nil {
		log.Fatalf("Error writing output file: %v", err)
	}

	// Calculate size reduction
	originalSize := len(content)
	minifiedSize := len(minified)
	reduction := float64(originalSize-minifiedSize) / float64(originalSize) * 100

	fmt.Printf("Minification complete!\n")
	fmt.Printf("Original size: %d bytes\n", originalSize)
	fmt.Printf("Minified size: %d bytes\n", minifiedSize)
	fmt.Printf("Size reduction: %.2f%%\n", reduction)
	fmt.Printf("Output written to: %s\n", outPath)
}
