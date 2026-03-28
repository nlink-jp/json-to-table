package main

import (
	"strings"
	"testing"
)

var sampleTable = [][]string{
	{"name", "age", "city"},
	{"Alice", "30", "New York"},
	{"Bob", "24", "London"},
}

var multibyteTable = [][]string{
	{"名前", "年齢"},
	{"東京太郎", "30"},
}

// --- renderAsText ---

func TestRenderAsText_Empty(t *testing.T) {
	got, err := renderAsText(nil)
	if err != nil || got != "" {
		t.Errorf("want empty string, got %q, err %v", got, err)
	}
}

func TestRenderAsText_Basic(t *testing.T) {
	got, err := renderAsText(sampleTable)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(got, "Alice") || !strings.Contains(got, "London") {
		t.Errorf("missing expected values in output:\n%s", got)
	}
	// Should have header separator line
	if !strings.Contains(got, "|---") && !strings.Contains(got, "|-") {
		t.Errorf("missing separator in text output:\n%s", got)
	}
}

func TestRenderAsText_Multibyte(t *testing.T) {
	got, err := renderAsText(multibyteTable)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(got, "東京太郎") {
		t.Errorf("missing multibyte value in output:\n%s", got)
	}
}

// --- renderAsMarkdown ---

func TestRenderAsMarkdown_Empty(t *testing.T) {
	got, err := renderAsMarkdown(nil)
	if err != nil || got != "" {
		t.Errorf("want empty string, got %q, err %v", got, err)
	}
}

func TestRenderAsMarkdown_Basic(t *testing.T) {
	got, err := renderAsMarkdown(sampleTable)
	if err != nil {
		t.Fatal(err)
	}
	lines := strings.Split(strings.TrimRight(got, "\n"), "\n")
	if len(lines) != 4 { // header + separator + 2 rows
		t.Errorf("want 4 lines, got %d:\n%s", len(lines), got)
	}
	if !strings.HasPrefix(lines[0], "| name") {
		t.Errorf("unexpected header line: %s", lines[0])
	}
	if !strings.Contains(lines[1], "---") {
		t.Errorf("missing separator: %s", lines[1])
	}
}

// --- renderAsCSV ---

func TestRenderAsCSV_Empty(t *testing.T) {
	got, err := renderAsCSV(nil)
	if err != nil || got != "" {
		t.Errorf("want empty string, got %q, err %v", got, err)
	}
}

func TestRenderAsCSV_Basic(t *testing.T) {
	got, err := renderAsCSV(sampleTable)
	if err != nil {
		t.Fatal(err)
	}
	lines := strings.Split(strings.TrimRight(got, "\n"), "\n")
	if len(lines) != 3 {
		t.Errorf("want 3 lines, got %d", len(lines))
	}
	if lines[0] != "name,age,city" {
		t.Errorf("unexpected header: %s", lines[0])
	}
}

func TestRenderAsCSV_CommaInValue(t *testing.T) {
	table := [][]string{
		{"col"},
		{"a,b"},
	}
	got, err := renderAsCSV(table)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(got, `"a,b"`) {
		t.Errorf("comma in value should be quoted, got: %s", got)
	}
}

// --- renderAsHTML ---

func TestRenderAsHTML_Empty(t *testing.T) {
	got, err := renderAsHTML(nil)
	if err != nil || got != "" {
		t.Errorf("want empty string, got %q, err %v", got, err)
	}
}

func TestRenderAsHTML_Basic(t *testing.T) {
	got, err := renderAsHTML(sampleTable)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(got, "<table") || !strings.Contains(got, "Alice") {
		t.Errorf("missing expected HTML content:\n%s", got)
	}
	if !strings.Contains(got, "<th") {
		t.Errorf("missing <th> header cells:\n%s", got)
	}
}

// --- renderAsSlackBlockKit ---

func TestRenderAsSlackBlockKit_Empty(t *testing.T) {
	_, err := renderAsSlackBlockKit(nil)
	if err == nil {
		t.Error("want error for empty table")
	}
}

func TestRenderAsSlackBlockKit_Basic(t *testing.T) {
	got, err := renderAsSlackBlockKit(sampleTable)
	if err != nil {
		t.Fatal(err)
	}
	s := string(got)
	if !strings.Contains(s, `"type": "table"`) {
		t.Errorf("missing table type in output:\n%s", s)
	}
	if !strings.Contains(s, `"bold": true`) {
		t.Errorf("header should be bold:\n%s", s)
	}
	if !strings.Contains(s, "Alice") {
		t.Errorf("missing data value:\n%s", s)
	}
}
