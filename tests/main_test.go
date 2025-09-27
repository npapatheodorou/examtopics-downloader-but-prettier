package tests

import (
	"os"
	"reflect"
	"strings"
	"path/filepath"
	"fmt"
	"testing"

	"examtopics-downloader/internal/fetch"
	"examtopics-downloader/internal/models"
	"examtopics-downloader/internal/utils"
)

var links []models.QuestionData
func TestGetAllPages(t *testing.T) {
	links = fetch.GetAllPages("lpi", "010-160")
	if len(links) == 0 {
		t.Fatalf("Expected non-empty data for provider 'lpi', but got: %v", links)
	}

	expectedType := reflect.TypeOf(models.QuestionData{})
	for _, link := range links {
		if reflect.TypeOf(link) != expectedType {
			t.Fatalf("Incorrect type for link, expected %v, got %v", expectedType, reflect.TypeOf(link))
		}
	}

	t.Logf("Data len of %d for provider 'lpi'", len(links))
}

func TestValidateExamsOutput(t *testing.T) {
	outputPath := "test.txt"

	utils.SaveLinks(outputPath, links)

	data, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("Expected file at %s but got error: %v", outputPath, err)
	}

	content := string(data)
	expectedContent := "https://www.examtopics.com/discussions/lpi/view"
	if !strings.Contains(content, expectedContent) {
		t.Errorf("Expected file content to contain %q but got:\n%s", expectedContent, content)
	}

	err = os.Remove(outputPath)
	if err != nil {
		t.Fatalf("Error when removing file! %v", err)
	}
}

func TestExamProvider(t *testing.T) {
	data := fetch.GetProviderExams("google")
	if len(data) == 0 {
		t.Fatalf("Expected non-empty data for provider 'google', but got: %v", data)
	}

	t.Logf("Got %d exams for provider 'google'", len(data))
}

func TestWriteDataVariants(t *testing.T) {
	baseName := "write_test"

	tests := []struct {
		fileType string
		ext      string
		checkContent bool
	}{
		{"md", ".md", true},
		{"text", ".txt", true},
		{"html", ".html", true},
		{"pdf", ".pdf", false}, // can't easily check PDF content
	}

	for _, tt := range tests {
		outputPath := fmt.Sprintf("%s.%s", baseName, tt.fileType)
		utils.WriteData(links, outputPath, true, tt.fileType)

		info, err := os.Stat(outputPath)
		if err != nil {
			t.Fatalf("Expected file at %s but got error: %v", outputPath, err)
		}
		if info.Size() == 0 {
			t.Errorf("Expected file %s to be non-empty", outputPath)
		}

		if tt.checkContent {
			data, err := os.ReadFile(outputPath)
			if err != nil {
				t.Fatalf("Failed reading %s: %v", outputPath, err)
			}
			content := string(data)
			if !strings.Contains(content, "Comments:") {
				t.Errorf("Expected %s to contain 'Comments:' but got:\n%s", outputPath, content)
			}
		}

		files, err := filepath.Glob("write_test*")
		if err != nil {
			t.Fatalf("Failed to glob files: %v", err)
		}

		for _, f := range files {
			if err := os.Remove(f); err != nil {
				t.Logf("Failed to remove %s: %v", f, err)
			}
		}
	}
}
