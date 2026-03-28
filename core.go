package main

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"strings"
)

// measureWidth returns the display width of s, counting non-ASCII runes as 2.
func measureWidth(s string) int {
	w := 0
	for _, r := range s {
		if r > 255 {
			w += 2
		} else {
			w++
		}
	}
	return w
}

// matchHeaders applies patterns to a list of headers.
// If isExclusion is true, it returns headers NOT matching the patterns.
// If isExclusion is false, it returns headers matching the patterns, in the order specified by patterns.
func matchHeaders(availableHeaders []string, patterns string, isExclusion bool) []string {
	if patterns == "" {
		if isExclusion {
			return availableHeaders // No patterns to exclude, return all
		}
		return []string{} // No patterns to include, return empty
	}

	userPatterns := strings.Split(patterns, ",")
	matched := make(map[string]bool)
	resultOrder := []string{} // To maintain order for inclusion

	// Create a map for quick lookup of available headers
	availableHeadersMap := make(map[string]bool)
	for _, h := range availableHeaders {
		availableHeadersMap[h] = true
	}

	for _, pattern := range userPatterns {
		trimmedPattern := strings.TrimSpace(pattern)

		if trimmedPattern == "*" {
			// Handle wildcard for remaining headers
			// For inclusion: add all remaining available headers
			// For exclusion: mark all remaining available headers as matched (to be excluded)
			for _, header := range availableHeaders {
				if availableHeadersMap[header] && !matched[header] { // Only consider headers not yet processed by explicit patterns
					if isExclusion {
						matched[header] = true
					} else {
						resultOrder = append(resultOrder, header)
						matched[header] = true
					}
				}
			}
		} else if strings.HasSuffix(trimmedPattern, "*") {
			// Prefix wildcard (e.g., "col*")
			prefix := strings.TrimSuffix(trimmedPattern, "*")
			var currentMatches []string
			for _, header := range availableHeaders {
				if strings.HasPrefix(header, prefix) && availableHeadersMap[header] {
					currentMatches = append(currentMatches, header)
				}
			}
			sort.Strings(currentMatches) // Sort prefix matches for deterministic behavior

			for _, header := range currentMatches {
				if !matched[header] {
					if isExclusion {
						matched[header] = true
					} else {
						resultOrder = append(resultOrder, header)
						matched[header] = true
					}
				}
			}
		} else {
			// Specific column name
			if availableHeadersMap[trimmedPattern] && !matched[trimmedPattern] {
				if isExclusion {
					matched[trimmedPattern] = true
				} else {
					resultOrder = append(resultOrder, trimmedPattern)
					matched[trimmedPattern] = true
				}
			}
		}
	}

	if isExclusion {
		// For exclusion, return headers that were NOT matched
		finalHeaders := []string{}
		for _, header := range availableHeaders {
			if !matched[header] {
				finalHeaders = append(finalHeaders, header)
			}
		}
		return finalHeaders
	} else {
		// For inclusion, return headers that were matched, in order
		return resultOrder
	}
}


// parseJSON reads JSON from an io.Reader and converts it into a table structure,
// respecting the user-defined column order with advanced wildcards.
func parseJSON(r io.Reader, columnOrder string, excludeColumnOrder string) ([][]string, error) {
	var data []map[string]any
	decoder := json.NewDecoder(r)
	if err := decoder.Decode(&data); err != nil {
		return nil, fmt.Errorf("failed to decode json: %w", err)
	}

	if len(data) == 0 {
		return [][]string{}, nil
	}

	// 1. Collect all unique keys from the data and sort them for deterministic order.
	allHeadersSet := make(map[string]bool)
	for _, row := range data {
		for key := range row {
			allHeadersSet[key] = true
		}
	}
	allHeadersList := make([]string, 0, len(allHeadersSet))
	for h := range allHeadersSet {
		allHeadersList = append(allHeadersList, h)
	}

	sort.Strings(allHeadersList) // Ensure initial list is sorted

	// 2. Apply exclusion patterns first
	headersAfterExclusion := matchHeaders(allHeadersList, excludeColumnOrder, true)

	var finalHeaders []string
	if columnOrder == "" {
		// Default behavior: use all headers remaining after exclusion, sorted alphabetically.
		sort.Strings(headersAfterExclusion) // Ensure sorted after exclusion
		finalHeaders = headersAfterExclusion
	} else {
		// Apply inclusion patterns to the headers remaining after exclusion
		finalHeaders = matchHeaders(headersAfterExclusion, columnOrder, false)
	}

	// 3. Create the table data structure (headers + rows) using the final header order.
	table := make([][]string, len(data)+1)
	table[0] = finalHeaders
	for i, rowMap := range data {
		row := make([]string, len(finalHeaders))
		for j, header := range finalHeaders {
			if val, ok := rowMap[header]; ok {
				row[j] = fmt.Sprintf("%v", val)
			} else {
				row[j] = "" // Handle missing keys for a given row
			}
		}
		table[i+1] = row
	}

	return table, nil
}
