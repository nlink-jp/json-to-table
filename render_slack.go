package main

import (
	"encoding/json"
	"errors"
	"fmt"
)

type slackTextStyle struct {
	Bold bool `json:"bold,omitempty"`
}

type slackTextElement struct {
	Type  string          `json:"type"`
	Text  string          `json:"text"`
	Style *slackTextStyle `json:"style,omitempty"`
}

type slackRichTextSection struct {
	Type     string             `json:"type"`
	Elements []slackTextElement `json:"elements"`
}

type slackRichText struct {
	Type     string                 `json:"type"`
	Elements []slackRichTextSection `json:"elements"`
}

type slackTableBlock struct {
	Type string           `json:"type"`
	Rows [][]slackRichText `json:"rows"`
}

// newRichTextCell constructs a single Slack rich_text cell.
func newRichTextCell(text string, bold bool) slackRichText {
	var style *slackTextStyle
	if bold {
		style = &slackTextStyle{Bold: true}
	}
	return slackRichText{
		Type: "rich_text",
		Elements: []slackRichTextSection{
			{
				Type: "rich_text_section",
				Elements: []slackTextElement{
					{Type: "text", Text: text, Style: style},
				},
			},
		},
	}
}

// renderAsSlackBlockKit formats the table as Slack Block Kit JSON.
func renderAsSlackBlockKit(table [][]string) ([]byte, error) {
	if len(table) == 0 {
		return nil, errors.New("cannot generate Slack Block Kit from empty data")
	}

	var tableRows [][]slackRichText

	// Header row (bold)
	headerRow := make([]slackRichText, len(table[0]))
	for i, header := range table[0] {
		headerRow[i] = newRichTextCell(header, true)
	}
	tableRows = append(tableRows, headerRow)

	// Data rows
	for _, rowData := range table[1:] {
		dataRow := make([]slackRichText, len(rowData))
		for i, cell := range rowData {
			dataRow[i] = newRichTextCell(cell, false)
		}
		tableRows = append(tableRows, dataRow)
	}

	output := map[string]any{
		"blocks": []slackTableBlock{
			{Type: "table", Rows: tableRows},
		},
	}

	jsonData, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal Slack Block Kit JSON: %w", err)
	}

	return jsonData, nil
}
