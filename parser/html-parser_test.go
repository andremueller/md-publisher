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
