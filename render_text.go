package main

import (
	"fmt"
	"strings"
)

// renderAsText formats the table as plain text with aligned columns.
func renderAsText(table [][]string) (string, error) {
	if len(table) == 0 {
		return "", nil
	}

	colWidths := make([]int, len(table[0]))
	for _, row := range table {
		for i, cell := range row {
			if w := measureWidth(cell); w > colWidths[i] {
				colWidths[i] = w
			}
		}
	}

	var builder strings.Builder
	// Header
	for i, header := range table[0] {
		padding := colWidths[i] - measureWidth(header)
		builder.WriteString(fmt.Sprintf("| %s%s ", header, strings.Repeat(" ", padding)))
	}
	builder.WriteString("|\n")

	// Separator
	for _, width := range colWidths {
		builder.WriteString(fmt.Sprintf("|-%s-", strings.Repeat("-", width)))
	}
	builder.WriteString("|\n")

	// Body
	for _, row := range table[1:] {
		for i, cell := range row {
			padding := colWidths[i] - measureWidth(cell)
			builder.WriteString(fmt.Sprintf("| %s%s ", cell, strings.Repeat(" ", padding)))
		}
		builder.WriteString("|\n")
	}

	return builder.String(), nil
}

// renderAsMarkdown formats the table as a GitHub-Flavored Markdown table.
func renderAsMarkdown(table [][]string) (string, error) {
	if len(table) == 0 {
		return "", nil
	}

	var builder strings.Builder
	// Header
	builder.WriteString("| " + strings.Join(table[0], " | ") + " |\n")
	// Separator
	separator := make([]string, len(table[0]))
	for i := range separator {
		separator[i] = "---"
	}
	builder.WriteString("| " + strings.Join(separator, " | ") + " |\n")
	// Body
	for _, row := range table[1:] {
		builder.WriteString("| " + strings.Join(row, " | ") + " |\n")
	}

	return builder.String(), nil
}
