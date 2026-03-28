package main

import (
	"strings"
	"testing"
)

// --- measureWidth ---

func TestMeasureWidth_ASCII(t *testing.T) {
	if got := measureWidth("hello"); got != 5 {
		t.Errorf("want 5, got %d", got)
	}
}

func TestMeasureWidth_Multibyte(t *testing.T) {
	// Each kanji counts as 2
	if got := measureWidth("東京"); got != 4 {
		t.Errorf("want 4, got %d", got)
	}
}

func TestMeasureWidth_Mixed(t *testing.T) {
	if got := measureWidth("AB東"); got != 4 {
		t.Errorf("want 4, got %d", got)
	}
}

func TestMeasureWidth_Empty(t *testing.T) {
	if got := measureWidth(""); got != 0 {
		t.Errorf("want 0, got %d", got)
	}
}

// --- matchHeaders ---

func TestMatchHeaders_ExcludeEmpty(t *testing.T) {
	headers := []string{"a", "b", "c"}
	got := matchHeaders(headers, "", true)
	if len(got) != 3 {
		t.Errorf("want all 3 headers, got %v", got)
	}
}

func TestMatchHeaders_IncludeEmpty(t *testing.T) {
	headers := []string{"a", "b", "c"}
	got := matchHeaders(headers, "", false)
	if len(got) != 0 {
		t.Errorf("want 0 headers, got %v", got)
	}
}

func TestMatchHeaders_IncludeSpecific(t *testing.T) {
	headers := []string{"a", "b", "c"}
	got := matchHeaders(headers, "c,a", false)
	if len(got) != 2 || got[0] != "c" || got[1] != "a" {
		t.Errorf("want [c a], got %v", got)
	}
}

func TestMatchHeaders_ExcludeSpecific(t *testing.T) {
	headers := []string{"a", "b", "c"}
	got := matchHeaders(headers, "b", true)
	if len(got) != 2 || got[0] != "a" || got[1] != "c" {
		t.Errorf("want [a c], got %v", got)
	}
}

func TestMatchHeaders_IncludeWildcard(t *testing.T) {
	headers := []string{"a", "b", "c"}
	got := matchHeaders(headers, "*", false)
	if len(got) != 3 {
		t.Errorf("want 3 headers, got %v", got)
	}
}

func TestMatchHeaders_ExcludeWildcard(t *testing.T) {
	headers := []string{"a", "b", "c"}
	got := matchHeaders(headers, "*", true)
	if len(got) != 0 {
		t.Errorf("want 0 headers, got %v", got)
	}
}

func TestMatchHeaders_PrefixWildcard(t *testing.T) {
	headers := []string{"col_a", "col_b", "other"}
	got := matchHeaders(headers, "col_*", false)
	if len(got) != 2 {
		t.Errorf("want 2 headers, got %v", got)
	}
}

func TestMatchHeaders_ExplicitThenWildcard(t *testing.T) {
	headers := []string{"a", "b", "c"}
	// "b" first, then remaining via "*"
	got := matchHeaders(headers, "b,*", false)
	if len(got) != 3 || got[0] != "b" {
		t.Errorf("want b first then rest, got %v", got)
	}
}

func TestMatchHeaders_NonexistentColumn(t *testing.T) {
	headers := []string{"a", "b"}
	got := matchHeaders(headers, "z", false)
	if len(got) != 0 {
		t.Errorf("want empty, got %v", got)
	}
}

// --- parseJSON ---

func TestParseJSON_Basic(t *testing.T) {
	input := `[{"name":"Alice","age":30},{"name":"Bob","age":24}]`
	table, err := parseJSON(strings.NewReader(input), "", "")
	if err != nil {
		t.Fatal(err)
	}
	if len(table) != 3 { // header + 2 rows
		t.Fatalf("want 3 rows, got %d", len(table))
	}
	// Headers should be sorted alphabetically
	if table[0][0] != "age" || table[0][1] != "name" {
		t.Errorf("unexpected headers: %v", table[0])
	}
}

func TestParseJSON_Empty(t *testing.T) {
	table, err := parseJSON(strings.NewReader(`[]`), "", "")
	if err != nil {
		t.Fatal(err)
	}
	if len(table) != 0 {
		t.Errorf("want empty table, got %v", table)
	}
}

func TestParseJSON_MissingKey(t *testing.T) {
	input := `[{"a":"1","b":"2"},{"a":"3"}]`
	table, err := parseJSON(strings.NewReader(input), "", "")
	if err != nil {
		t.Fatal(err)
	}
	// Second row should have empty string for missing "b"
	if table[2][1] != "" {
		t.Errorf("want empty string for missing key, got %q", table[2][1])
	}
}

func TestParseJSON_ColumnOrder(t *testing.T) {
	input := `[{"a":"1","b":"2","c":"3"}]`
	table, err := parseJSON(strings.NewReader(input), "c,a", "")
	if err != nil {
		t.Fatal(err)
	}
	if table[0][0] != "c" || table[0][1] != "a" {
		t.Errorf("want [c a], got %v", table[0])
	}
	if len(table[0]) != 2 {
		t.Errorf("want 2 columns, got %d", len(table[0]))
	}
}

func TestParseJSON_ExcludeColumn(t *testing.T) {
	input := `[{"a":"1","b":"2","c":"3"}]`
	table, err := parseJSON(strings.NewReader(input), "", "b")
	if err != nil {
		t.Fatal(err)
	}
	for _, h := range table[0] {
		if h == "b" {
			t.Error("excluded column 'b' should not appear in headers")
		}
	}
}

func TestParseJSON_InvalidJSON(t *testing.T) {
	_, err := parseJSON(strings.NewReader(`not json`), "", "")
	if err == nil {
		t.Error("want error for invalid JSON")
	}
}
