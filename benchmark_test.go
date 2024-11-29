package main

import (
	"io/ioutil"
	"path/filepath"
	"testing"
)

func BenchmarkMinification(b *testing.B) {
	testCases := []struct {
		name    string
		options struct {
			preserveLicense bool
			shortenVars    bool
		}
	}{
		{"BasicMinification", struct {
			preserveLicense bool
			shortenVars    bool
		}{false, false}},
		{"WithLicensePreservation", struct {
			preserveLicense bool
			shortenVars    bool
		}{true, false}},
		{"WithVariableShortening", struct {
			preserveLicense bool
			shortenVars    bool
		}{false, true}},
		{"AllOptions", struct {
			preserveLicense bool
			shortenVars    bool
		}{true, true}},
	}

	testFiles := []string{
		"closure.js",   // Tests nested functions and closures
		"complex.js",   // Tests complex class-based code
		"modern.js",    // Tests ES6+ features
		"regex.js",     // Tests regular expressions
		"comments.js",  // Tests comment handling
		"simple.js",    // Tests basic functionality
	}

	for _, tc := range testCases {
		for _, file := range testFiles {
			name := tc.name + "/" + file
			b.Run(name, func(b *testing.B) {
				content, err := ioutil.ReadFile(filepath.Join("test", "testdata", file))
				if err != nil {
					b.Fatalf("Failed to read test file %s: %v", file, err)
				}

				input := string(content)
				b.ResetTimer()

				for i := 0; i < b.N; i++ {
					minifier := NewMinifier(input, tc.options.preserveLicense, tc.options.shortenVars)
					_ = minifier.Minify()
				}
			})
		}
	}
}

func BenchmarkLargeFile(b *testing.B) {
	// Create a large file by repeating the complex test file
	content, err := ioutil.ReadFile(filepath.Join("test", "testdata", "complex.js"))
	if err != nil {
		b.Fatal(err)
	}

	// Repeat the content 100 times to create a ~100KB file
	largeContent := ""
	for i := 0; i < 100; i++ {
		largeContent += string(content)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		minifier := NewMinifier(largeContent, true, true)
		_ = minifier.Minify()
	}
}
