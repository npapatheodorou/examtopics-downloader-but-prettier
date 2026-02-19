package utils

import (
	"testing"

	"examtopics-downloader/internal/models"
)

func TestDeriveExamMetaPrefersSelectedExam(t *testing.T) {
	data := []models.QuestionData{
		{
			QuestionLink: "https://www.examtopics.com/discussions/oracle/view/305691-exam-1z0-1042-20-topic-1-question-3-discussion/",
			Title:        "Exam 1z0-1042-20 topic 1 question 3 discussion",
		},
	}

	meta := deriveExamMeta(data, "oracle", "1z0-1042")
	if meta.Company != "Oracle" {
		t.Fatalf("unexpected company: %q", meta.Company)
	}
	if meta.ExamCode != "1Z0-1042" {
		t.Fatalf("unexpected exam code: %q", meta.ExamCode)
	}
	if meta.Badge != "1Z0-1042" {
		t.Fatalf("unexpected badge: %q", meta.Badge)
	}
}

func TestDeriveExamCodeKeepsOraclePrefix(t *testing.T) {
	got := deriveExamCode("1z0-1042-20")
	if got != "1Z0-1042-20" {
		t.Fatalf("unexpected exam code: %q", got)
	}
}
