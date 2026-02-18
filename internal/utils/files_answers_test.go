package utils

import (
	"reflect"
	"strings"
	"testing"
)

func TestRemoveSuggestedAnswerText(t *testing.T) {
	input := "Refer to the exhibit.\nSuggested Answer: BD"
	got := removeSuggestedAnswerText(input)

	if strings.Contains(strings.ToLower(got), "suggested answer") {
		t.Fatalf("expected suggested answer text to be removed, got %q", got)
	}
	if got != "Refer to the exhibit." {
		t.Fatalf("unexpected cleaned text, got %q", got)
	}
}

func TestExtractCorrectAnswersMultiple(t *testing.T) {
	options := []answerOption{
		{Letter: "A", Text: "Option A"},
		{Letter: "B", Text: "Option B"},
		{Letter: "C", Text: "Option C"},
		{Letter: "D", Text: "Option D"},
	}

	got := extractCorrectAnswers("Correct Answer: BD", options)
	want := []string{"B", "D"}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("unexpected answers: want %v, got %v", want, got)
	}
}
