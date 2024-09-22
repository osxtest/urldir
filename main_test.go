package main

import (
	"net/url"
	"reflect"
	"sort"
	"testing"
)

func Test_dirs(t *testing.T) {
	tests := []struct {
		in    string
		depth int
		want  []string
	}{
		{"http://example.com/a/b/c/d", -1, []string{
			"http://example.com",
			"http://example.com/a",
			"http://example.com/a/b",
			"http://example.com/a/b/c",
		}},
		{"http://example.com/a/b/c/d", 0, []string{
			"http://example.com",
		}},
		{"http://example.com/a/b/c/d", 1, []string{
			"http://example.com",
			"http://example.com/a",
		}},
		{"http://example.com/a/b/c/d", 2, []string{
			"http://example.com",
			"http://example.com/a",
			"http://example.com/a/b",
		}},
		{"http://example.com/a/b/c/d", 3, []string{
			"http://example.com",
			"http://example.com/a",
			"http://example.com/a/b",
			"http://example.com/a/b/c",
		}},
		{"http://example.com/a/b/c/d", 100, []string{
			"http://example.com",
			"http://example.com/a",
			"http://example.com/a/b",
			"http://example.com/a/b/c",
		}},
		{"http://example.com", -1, []string{"http://example.com"}},
		{"http://example.com/", -1, []string{"http://example.com"}},
		{"http://example.com/.", -1, []string{"http://example.com"}},
		{"http://example.com/a", -1, []string{"http://example.com"}},
		{"http://example.com/////", -1, []string{"http://example.com"}},
	}
	for _, tt := range tests {
		u, _ := url.Parse(tt.in)
		got := dirs(u, tt.depth)
		sort.Strings(got)
		sort.Strings(tt.want)
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("dirs(%q, %d) = %#v; want %#v", tt.in, tt.depth, got, tt.want)
		}
	}
}

func Test_build(t *testing.T) {
	base := "http://example.com"
	parsed, _ := url.Parse(base)
	tests := []struct {
		path, want string
	}{
		{"", "http://example.com"},
		{"a/b", "http://example.com/a/b"},
		{"/a/b", "http://example.com/a/b"},
	}
	for _, tt := range tests {
		if got := build(parsed, tt.path); got != tt.want {
			t.Errorf("build(%q, %q) = %v, want %v", base, tt.path, got, tt.want)
		}
	}
}

func Test_parse_edgeCase(t *testing.T) {
	tests := []string{
		"",
		"/a/b/c",
		"file:///etc/passwd",
		"http:example.com",
		"http:/example.com",
		" http://example.com",
	}
	for _, tt := range tests {
		if len(parse(tt, -1)) != 0 {
			t.Errorf("parse(%q, -1) = %v, want []", tt, parse(tt, -1))
		}
	}
}
func Test_parse(t *testing.T) {
	url := "http://example.com/a/b/c/d"
	got := parse(url, -1)
	want := []string{
		"http://example.com",
		"http://example.com/a",
		"http://example.com/a/b",
		"http://example.com/a/b/c",
	}
	sort.Strings(got)
	sort.Strings(want)
	if !reflect.DeepEqual(got, want) {
		t.Errorf("parse(%q, -1) = %#v; want %#v", url, got, want)
	}
}
