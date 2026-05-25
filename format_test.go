package main

import (
	"strings"
	"testing"
	"unicode/utf8"
)

func TestFormatContent(t *testing.T) {
	url := "https://jbrio.net/posts/202522/"

	t.Run("includes tags that fit", func(t *testing.T) {
		fm := &Frontmatter{Description: "a pc in your browser", Tags: []string{"notes", "browser"}}
		got := formatContent(fm, url, limitBluesky)
		want := "a pc in your browser\n\n#notes #browser\n\n" + url
		if got != want {
			t.Fatalf("got %q, want %q", got, want)
		}
	})

	t.Run("falls back to title when no description", func(t *testing.T) {
		fm := &Frontmatter{Title: "My Title"}
		got := formatContent(fm, url, limitBluesky)
		if !strings.HasPrefix(got, "My Title\n\n") {
			t.Fatalf("expected title prefix, got %q", got)
		}
	})

	t.Run("drops tags that do not fit", func(t *testing.T) {
		fm := &Frontmatter{
			Description: strings.Repeat("x", 250),
			Tags:        []string{"toolongtagwontfit"},
		}
		got := formatContent(fm, url, limitTwitter)
		if strings.Contains(got, "#toolongtagwontfit") {
			t.Fatalf("tag should have been dropped: %q", got)
		}
	})

	t.Run("truncates over-limit description", func(t *testing.T) {
		fm := &Frontmatter{Description: strings.Repeat("x", 400)}
		got := formatContent(fm, url, limitTwitter)
		if utf8.RuneCountInString(got) > limitTwitter {
			t.Fatalf("result %d runes exceeds limit %d", utf8.RuneCountInString(got), limitTwitter)
		}
		if !strings.Contains(got, "...") || !strings.HasSuffix(got, url) {
			t.Fatalf("expected truncation marker and url suffix, got %q", got)
		}
	})
}

func TestSlugFromPath(t *testing.T) {
	if got := slugFromPath("data/posts/202545/index.md"); got != "202545" {
		t.Fatalf("got %q, want 202545", got)
	}
}

func TestAllowedNetworks(t *testing.T) {
	tests := []struct {
		name string
		env  string
		want map[string]bool // nil == no filter (all)
	}{
		{"unset is all", "", nil},
		{"whitespace is all", "   ", nil},
		{"single", "mastodon", map[string]bool{"mastodon": true}},
		{"comma list", "bluesky,mastodon", map[string]bool{"bluesky": true, "mastodon": true}},
		{"space, case and x alias", "Bluesky X", map[string]bool{"bluesky": true, "twitter": true}},
		{"unknown dropped", "mastodon,bogus", map[string]bool{"mastodon": true}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("SOCIAL_NETWORKS", tt.env)
			if got := allowedNetworks(); !sameSet(got, tt.want) {
				t.Fatalf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEnabled(t *testing.T) {
	if !enabled(nil, "bluesky") {
		t.Fatal("nil set should enable every network")
	}
	set := map[string]bool{"mastodon": true}
	if enabled(set, "bluesky") {
		t.Fatal("bluesky absent from set should be disabled")
	}
	if !enabled(set, "mastodon") {
		t.Fatal("mastodon present in set should be enabled")
	}
}

// sameSet treats nil and empty as equal (both mean "no networks selected").
func sameSet(a, b map[string]bool) bool {
	if len(a) != len(b) {
		return false
	}
	for k := range a {
		if !b[k] {
			return false
		}
	}
	return true
}
