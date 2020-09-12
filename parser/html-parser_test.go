package parser

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestGetKeywords(t *testing.T) {
	doc, err := ParseHTML("../data/demo_article.html")
	if err != nil {
		t.Errorf("Failed parsing of html")
	}
	keywords := GetKeywords(doc)
	expected := []string{"pandoc", "demo", "article", "automation"}
	if !cmp.Equal(keywords, expected) {
		t.Error("Keywords not correctly detected")
	}
}

func TestGetTitle(t *testing.T) {
	doc, err := ParseHTML("../data/demo_article.html")
	if err != nil {
		t.Errorf("Failed parsing of html")
	}
	title := GetTitle(doc, "default")
	expected := "A Demo Article Title"
	if !cmp.Equal(title, expected) {
		t.Errorf("Title not correct '%s' <> '%s'", title, expected)
	}
}
