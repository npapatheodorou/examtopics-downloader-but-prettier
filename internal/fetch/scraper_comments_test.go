package fetch

import (
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"
)

func TestExtractDiscussionCommentsStructured(t *testing.T) {
	html := `
<div class="discussion-container">
  <div id="comment-1" class="media comment-container" data-comment-id="1">
    <div class="media-body">
      <div class="comment-head">
        <h5 class="comment-username">alice</h5>
      </div>
      <div class="comment-body">
        <div class="comment-selected-answers badge badge-warning">Selected Answer: <span><strong>C</strong></span></div>
        <div class="comment-content">first line
second line</div>
      </div>
    </div>
  </div>
  <div id="comment-2" class="media comment-container" data-comment-id="2">
    <div class="media-body">
      <div class="comment-head">
        <h5 class="comment-username">bob</h5>
      </div>
      <div class="comment-body">
        <div class="comment-content">plain comment only</div>
      </div>
    </div>
  </div>
</div>`

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		t.Fatalf("failed parsing test html: %v", err)
	}

	comments := extractDiscussionComments(doc)
	if len(comments) != 2 {
		t.Fatalf("expected 2 comments, got %d", len(comments))
	}

	if comments[0].User != "alice" {
		t.Fatalf("expected first user alice, got %q", comments[0].User)
	}
	if comments[0].Answer != "C" {
		t.Fatalf("expected first answer C, got %q", comments[0].Answer)
	}
	if comments[0].Text != "first line\nsecond line" {
		t.Fatalf("expected multiline comment text preserved, got %q", comments[0].Text)
	}

	if comments[1].User != "bob" {
		t.Fatalf("expected second user bob, got %q", comments[1].User)
	}
	if comments[1].Answer != "" {
		t.Fatalf("expected second answer empty, got %q", comments[1].Answer)
	}
	if comments[1].Text != "plain comment only" {
		t.Fatalf("expected second comment text, got %q", comments[1].Text)
	}
}
