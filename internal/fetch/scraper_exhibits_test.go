package fetch

import (
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"
)

func TestExtractExhibitImageURLs(t *testing.T) {
	html := `
<div class="card-text">
  <img src="//img.examtopics.com/200-301/image1.png">
  <img data-src="/media/exam/image2.jpg">
  <img srcset="https://img.examtopics.com/200-301/image3.webp 1x, https://img.examtopics.com/200-301/image3@2x.webp 2x">
  <img src="data:image/png;base64,abc123">
  <img src="//img.examtopics.com/200-301/image1.png">
</div>`

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		t.Fatalf("failed parsing test html: %v", err)
	}

	got := extractExhibitImageURLs(doc)
	want := []string{
		"https://img.examtopics.com/200-301/image1.png",
		"https://www.examtopics.com/media/exam/image2.jpg",
		"https://img.examtopics.com/200-301/image3.webp",
	}

	if len(got) != len(want) {
		t.Fatalf("expected %d urls, got %d: %#v", len(want), len(got), got)
	}

	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("unexpected url at index %d: want %q, got %q", i, want[i], got[i])
		}
	}
}
