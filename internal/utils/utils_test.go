package utils

import "testing"

func TestCleanTextRemovesVoteEmoji(t *testing.T) {
	input := "Refer to the exhibit. Which type of route does R1 use to reach host 10.10.13.10/32?\n🗳️"

	got := CleanText(input)
	want := "Refer to the exhibit. Which type of route does R1 use to reach host 10.10.13.10/32?"
	if got != want {
		t.Fatalf("unexpected cleaned text\nwant: %q\ngot:  %q", want, got)
	}
}

func TestCleanTextRemovesVoteEmojiWithoutVariationSelector(t *testing.T) {
	input := "Suggested Answer: D 🗳"

	got := CleanText(input)
	want := "\nSuggested Answer: D"
	if got != want {
		t.Fatalf("unexpected cleaned text\nwant: %q\ngot:  %q", want, got)
	}
}

