package main

import (
	"io/ioutil"
	"path/filepath"
	"strings"
	"testing"
)

// TestCase represents a single minification test case
type TestCase struct {
	Name           string
	Input          string
	ExpectedOutput string
	Options        MinificationOptions
}

// MinificationOptions represents options for minification
type MinificationOptions struct {
	PreserveLicense bool
	ShortenVars     bool
}

func normalizeWhitespace(s string) string {
	// Normalize whitespace for comparison
	s = strings.ReplaceAll(s, "\n", "")
	s = strings.ReplaceAll(s, "\t", "")
	s = strings.ReplaceAll(s, " ", "")
	return s
}

// TestMinifierBasic tests basic minification functionality
func TestMinifierBasic(t *testing.T) {
	input := `function test(a, b) {
		// This is a comment
		return a + b;
	}`
	expected := "function test(a,b){return a+b;}"

	minifier := NewMinifier(input, false, false)
	result := minifier.Minify()

	if normalizeWhitespace(result) != normalizeWhitespace(expected) {
		t.Errorf("Basic minification failed.\nExpected: %s\nGot: %s", expected, result)
	}
}

// TestMinifierPreserveLicense tests license comment preservation
func TestMinifierPreserveLicense(t *testing.T) {
	input := `/*!
 * License
 */
function test() {}`

	minifier := NewMinifier(input, true, false)
	result := minifier.Minify()

	if !strings.Contains(result, "/*!") || !strings.Contains(result, "License") {
		t.Error("License comment was not preserved")
	}

	if !strings.Contains(result, "function") || !strings.Contains(result, "test") {
		t.Error("Function definition was lost during minification")
	}
}

// TestMinifierVariableShortening tests variable name shortening
func TestMinifierVariableShortening(t *testing.T) {
	input := `const longVariableName = 42;
	let anotherLongName = longVariableName + 1;`
	
	minifier := NewMinifier(input, false, true)
	result := minifier.Minify()

	// Check if variables were shortened
	if strings.Contains(result, "longVariableName") || strings.Contains(result, "anotherLongName") {
		t.Error("Variable shortening failed, original names still present")
	}
}

// TestFileProcessing tests processing of actual JavaScript files
func TestFileProcessing(t *testing.T) {
	testFiles := []string{
		"closure.js",
		"comments.js",
		"complex.js",
		"modern.js",
		"regex.js",
		"simple.js",
	}

	for _, file := range testFiles {
		t.Run(file, func(t *testing.T) {
			inputPath := filepath.Join("test", "testdata", file)
			content, err := ioutil.ReadFile(inputPath)
			if err != nil {
				t.Fatalf("Failed to read test file %s: %v", file, err)
			}

			minifier := NewMinifier(string(content), false, false)
			result := minifier.Minify()

			// Basic validation
			if len(result) >= len(string(content)) {
				t.Errorf("Minification did not reduce file size for %s", file)
			}

			// Check for common syntax errors
			if strings.Count(result, "{") != strings.Count(result, "}") {
				t.Errorf("Mismatched braces in output for %s", file)
			}
			if strings.Count(result, "(") != strings.Count(result, ")") {
				t.Errorf("Mismatched parentheses in output for %s", file)
			}
		})
	}
}

// TestEdgeCases tests various edge cases
func TestEdgeCases(t *testing.T) {
	testCases := []TestCase{
		{
			Name:           "Empty Input",
			Input:          "",
			ExpectedOutput: "",
			Options:        MinificationOptions{false, false},
		},
		{
			Name:           "Only Comments",
			Input:          "// Just a comment\n/* Another comment */",
			ExpectedOutput: "",
			Options:        MinificationOptions{false, false},
		},
		{
			Name:           "Complex String Literals",
			Input:          `const str = "This is a \"quoted\" string"`,
			ExpectedOutput: `const str="This is a \"quoted\" string"`,
			Options:        MinificationOptions{false, false},
		},
		{
			Name:           "Regular Expressions",
			Input:          `const regex = /test/g;`,
			ExpectedOutput: `const regex=/test/g;`,
			Options:        MinificationOptions{false, false},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			minifier := NewMinifier(tc.Input, tc.Options.PreserveLicense, tc.Options.ShortenVars)
			result := minifier.Minify()
			if normalizeWhitespace(result) != normalizeWhitespace(tc.ExpectedOutput) {
				t.Errorf("%s failed.\nExpected: %s\nGot: %s", tc.Name, tc.ExpectedOutput, result)
			}
		})
	}
}
