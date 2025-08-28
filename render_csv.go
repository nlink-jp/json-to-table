package main

import (
	"encoding/csv"
	"strings"
)

// renderAsCSV formats the table as a CSV string.
func renderAsCSV(table [][]string) (string, error) {
	if len(table) == 0 {
		return "", nil
	}

	var builder strings.Builder
	writer := csv.NewWriter(&builder)

	err := writer.WriteAll(table)
	if err != nil {
		return "", err
	}

	return builder.String(), nil
}
