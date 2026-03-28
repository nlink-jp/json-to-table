package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

// Version is set at build time using ldflags
var version = "dev"

// --- Main Execution ---

func main() {
	format := flag.String("format", "text", "Output format: text, md, csv, png, html, slack-block-kit, blocks")
	output := flag.String("o", "", "Output file path (default: stdout)")
	title := flag.String("title", "", "Title for the image output")
	fontSize := flag.Float64("font-size", 12, "Font size for the image output")
	columns := flag.String("columns", "", "Comma-separated list of columns in desired order. Use '*' as a wildcard for other columns.")
	flag.StringVar(columns, "c", "", "Shorthand for --columns")
	excludeColumns := flag.String("exclude-columns", "", "Comma-separated list of columns to exclude. Use '*' as a wildcard.")
	flag.StringVar(excludeColumns, "e", "", "Shorthand for --exclude-columns")
	versionFlag := flag.Bool("version", false, "Print version information and exit")

	flag.Parse()

	if *versionFlag {
		fmt.Printf("json-to-table version %s\n", version)
		os.Exit(0)
	}

	// Check if data is being piped
	stat, err := os.Stdin.Stat()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to stat stdin: %v\n", err)
		os.Exit(1)
	}
	if (stat.Mode() & os.ModeCharDevice) != 0 {
		fmt.Fprintln(os.Stderr, "Error: This tool requires JSON data to be piped via stdin.")
		fmt.Fprintln(os.Stderr, "Usage: cat data.json | json-to-table")
		fmt.Fprintln(os.Stderr, "   or: splunk-cli run ... | jq .results | json-to-table --format png -o report.png")
		os.Exit(1)
	}

	table, err := parseJSON(os.Stdin, *columns, *excludeColumns)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing JSON: %+v\n", err)
		os.Exit(1)
	}

	var outData []byte
	var outStr string

	switch strings.ToLower(*format) {
	case "text":
		outStr, err = renderAsText(table)
	case "md", "markdown":
		outStr, err = renderAsMarkdown(table)
	case "png":
		outData, err = renderAsPNG(table, *title, *fontSize)
	case "html":
		outStr, err = renderAsHTML(table)
	case "slack-block-kit", "blocks":
		outData, err = renderAsSlackBlockKit(table)
	case "csv":
		outStr, err = renderAsCSV(table)
	default:
		err = fmt.Errorf("unknown format: %s", *format)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error rendering output: %+v\n", err)
		os.Exit(1)
	}

	// Determine output destination
	var writer io.Writer = os.Stdout
	if *output != "" {
		file, err := os.Create(*output)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating output file: %+v\n", err)
			os.Exit(1)
		}
		defer file.Close()
		writer = file
	}

	// Write the output
	if outData != nil {
		_, err = writer.Write(outData)
	} else {
		_, err = io.WriteString(writer, outStr)
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error writing output: %+v\n", err)
		os.Exit(1)
	}
}