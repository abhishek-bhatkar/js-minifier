# JavaScript Minifier

A powerful command-line utility written in Go to minify JavaScript files. This tool uses no external dependencies and relies solely on Go's standard library for minification.

## Features

- Removes comments (both single-line and multi-line)
- Removes unnecessary whitespace and newlines
- Removes spaces around operators and punctuation
- Variable name shortening (optional)
- Directory processing with parallel execution
- Watch mode for automatic minification
- Preservation of license comments (optional)
- JSON output format for statistics
- No external dependencies

## Installation

1. Make sure you have Go installed on your system
2. Clone this repository
3. Build the binary:
```bash
go build -o js-minifier
```

## Usage

### Basic Usage

Minify a single file (creates .min.js):
```bash
./js-minifier -input script.js
```

Specify custom output file:
```bash
./js-minifier -input script.js -output custom.min.js
```

### Advanced Features

Process all JavaScript files in a directory:
```bash
./js-minifier -input ./src
```

Watch directory for changes:
```bash
./js-minifier -input ./src -watch
```

Enable variable name shortening:
```bash
./js-minifier -input script.js -shorten-vars
```

Preserve license comments:
```bash
./js-minifier -input script.js -preserve-license
```

Output statistics in JSON format:
```bash
./js-minifier -input script.js -json
```

### Command Line Options

- `-input`: Input JavaScript file or directory (required)
- `-output`: Output file path (optional, default: [input].min.js)
- `-watch`: Watch mode - monitor directory for changes
- `-preserve-license`: Preserve license comments
- `-shorten-vars`: Enable variable name shortening
- `-json`: Output statistics in JSON format

## Minification Rules

The tool applies the following minification rules:
1. Removes all single-line comments (`// ...`)
2. Removes all multi-line comments (`/* ... */`)
3. Preserves license comments (`/*! ... */`) when `-preserve-license` is enabled
4. Removes extra whitespace and newlines
5. Removes spaces around operators (+, -, *, /, =, etc.)
6. Removes unnecessary semicolons
7. Removes spaces after function keywords
8. Removes spaces around brackets and parentheses
9. Shortens variable names (when `-shorten-vars` is enabled)

## Examples

### Basic Minification

Input file (script.js):
```javascript
function calculateSum(a, b) {
    // This is a comment
    const result = a + b;
    return result;
}
```

Output file (script.min.js):
```javascript
function calculateSum(a,b){const result=a+b;return result;}
```

### With Variable Shortening

Using `-shorten-vars` flag:
```javascript
function calculateSum(a,b){const c=a+b;return c;}
```

### Processing a Directory

Process all .js files in src directory:
```bash
./js-minifier -input ./src -shorten-vars
```

This will create minified versions of all JavaScript files:
- src/file1.js → src/file1.min.js
- src/file2.js → src/file2.min.js

### Watch Mode

Monitor directory for changes:
```bash
./js-minifier -input ./src -watch -preserve-license
```

This will:
- Watch the src directory for changes
- Automatically minify modified files
- Preserve license comments
- Show real-time statistics

## Performance

The minifier typically achieves:
- 40-50% size reduction on standard JavaScript files
- 60-70% reduction with variable name shortening
- Parallel processing for directory operations
- Millisecond-level processing time for most files
