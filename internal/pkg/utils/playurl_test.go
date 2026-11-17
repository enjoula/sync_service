package utils

import "testing"

func TestMaskPlayURL(t *testing.T) {
	cases := []struct {
		in   string
		want string
	}{
		{"https://cdn.example.com/path/abcd/index.m3u8", "https://cdn.example.com/path/****/index.m3u8"},
		{"https://cdn.example.com/ab/index.m3u8", "https://cdn.example.co****/index.m3u8"},
		{"ab/index.m3u8", "ab/index.m3u8"},
		{"https://cdn.example.com/no-marker.m3u8", "https://cdn.example.com/no-marker.m3u8"},
		{"", ""},
	}
	for _, c := range cases {
		got := MaskPlayURL(c.in)
		if got != c.want {
			t.Fatalf("MaskPlayURL(%q)=%q, want %q", c.in, got, c.want)
		}
	}
}

func TestIsMaskedPlayURL(t *testing.T) {
	if !IsMaskedPlayURL("https://cdn.example.com/path/****/index.m3u8") {
		t.Fatal("expected masked url")
	}
	if IsMaskedPlayURL("https://cdn.example.com/path/abcd/index.m3u8") {
		t.Fatal("expected unmasked url")
	}
}
