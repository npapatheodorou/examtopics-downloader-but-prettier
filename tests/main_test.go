package tests

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"examtopics-downloader/internal/fetch"
	"examtopics-downloader/internal/models"
	"examtopics-downloader/internal/utils"
)

var links []models.QuestionData

func TestGetAllPages(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test with -short flag (requires network)")
	}
	links = fetch.GetAllPages("lpi", "010-160")
	if len(links) == 0 {
		t.Skip("Skipping integration assertion: provider 'lpi' returned no data (likely remote blocking/rate limiting)")
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
	if testing.Short() {
		links = getMockQuestionData()
	} else if len(links) == 0 {
		t.Skip("Skipping - depends on TestGetAllPages running first")
	}
	outputPath := "test_links.html"

	utils.SaveLinks(outputPath, links)

	data, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("Expected file at %s but got error: %v", outputPath, err)
	}

	content := string(data)
	expectedContent := "https://www.examtopics.com/discussions/"
	if !strings.Contains(content, expectedContent) {
		t.Errorf("Expected file content to contain %q but got:\n%s", expectedContent, content)
	}

	err = os.Remove(outputPath)
	if err != nil {
		t.Fatalf("Error when removing file! %v", err)
	}
}

func TestExamProvider(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test with -short flag (requires network)")
	}
	data := fetch.GetProviderExams("google")
	if len(data) == 0 {
		t.Skip("Skipping integration assertion: provider 'google' returned no exams (likely remote blocking/rate limiting)")
	}

	t.Logf("Got %d exams for provider 'google'", len(data))
}

func TestWriteDataVariants(t *testing.T) {
	testLinks := getMockQuestionData()

	baseName := "write_test"
	tests := []struct {
		outputPath   string
		expectedExts []string
		checkContent bool
	}{
		{fmt.Sprintf("%s.html", baseName), []string{".html"}, true},
		{fmt.Sprintf("%s.custom", baseName), []string{".html"}, false},
	}

	for _, tt := range tests {
		generatedFiles, err := utils.WriteData(testLinks, tt.outputPath, true)
		if err != nil {
			t.Fatalf("WriteData failed for %s: %v", tt.outputPath, err)
		}

		if len(generatedFiles) != len(tt.expectedExts) {
			t.Fatalf("Expected %d generated files for %s but got %d", len(tt.expectedExts), tt.outputPath, len(generatedFiles))
		}

		for idx, ext := range tt.expectedExts {
			expectedFile := fmt.Sprintf("%s%s", baseName, ext)
			if generatedFiles[idx] != expectedFile {
				t.Fatalf("Expected generated file %d to be %s but got %s", idx, expectedFile, generatedFiles[idx])
			}

			info, err := os.Stat(expectedFile)
			if err != nil {
				t.Fatalf("Expected file at %s but got error: %v", expectedFile, err)
			}
			if info.Size() == 0 {
				t.Errorf("Expected file %s to be non-empty", expectedFile)
			}
		}

		if tt.checkContent {
			data, err := os.ReadFile(fmt.Sprintf("%s.html", baseName))
			if err != nil {
				t.Fatalf("Failed reading html output: %v", err)
			}
			content := string(data)
			if !strings.Contains(content, `<div class="q-card open" id="q1"`) {
				t.Errorf("Expected first question card with open class")
			}
			if !strings.Contains(content, `id="q1-preview"`) || !strings.Contains(content, `id="q1-text"`) {
				t.Errorf("Expected q1 preview/text ids in output")
			}
			if !strings.Contains(content, `id="q1-submit"`) || !strings.Contains(content, `id="q2-submit"`) {
				t.Errorf("Expected submit button ids for generated questions")
			}
			if !strings.Contains(content, `data-comments='[`) {
				t.Errorf("Expected data-comments JSON attribute in question cards")
			}
			if !strings.Contains(content, `<div class="q-exhibit">`) {
				t.Errorf("Expected exhibit block for question with image URL")
			}
			if strings.Count(content, `<div class="q-exhibit">`) != 1 {
				t.Errorf("Expected exactly one exhibit block, got %d", strings.Count(content, `<div class="q-exhibit">`))
			}
			if strings.Contains(content, "Suggested Answer:") {
				t.Errorf("Expected suggested answer text to be removed from rendered question text")
			}
			if !strings.Contains(content, "Cisco 200-301 Exam Simulator") {
				t.Errorf("Expected updated title/header metadata from extracted exam info")
			}
			if !strings.Contains(content, "npapatheodorou") {
				t.Errorf("Expected template footer attribution to remain intact")
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

func getMockQuestionData() []models.QuestionData {
	return []models.QuestionData{
		{
			Title:        "Question 123",
			Header:       "What is the purpose of the /etc/fstab file?",
			Content:      "Refer to the exhibit and choose the best answer.\nSuggested Answer: B\nhttps://img.examtopics.com/200-301/image348.png",
			Questions:    []string{"A. To configure network interfaces", "B. To define filesystem mount points", "C. To manage user accounts"},
			Answer:       "B. To define filesystem mount points",
			Timestamp:    "2024-01-15",
			QuestionLink: "https://www.examtopics.com/discussions/cisco/view/138394-exam-200-301-topic-1-question-1-discussion/",
			Comments: []models.CommentData{
				{User: "User1", Answer: "B", Text: "Answer: B (+5 votes)"},
				{User: "User2", Answer: "B", Text: "B is correct"},
			},
		},
		{
			Title:        "Question 456",
			Header:       "Which command displays running processes?",
			Content:      "Which command displays running processes?",
			Questions:    []string{"A. ls", "B. ps", "C. df"},
			Answer:       "B. ps",
			Timestamp:    "2024-01-16",
			QuestionLink: "https://www.examtopics.com/discussions/cisco/view/138395-exam-200-301-topic-1-question-2-discussion/",
			Comments: []models.CommentData{
				{User: "User2", Answer: "", Text: "ps aux is most common (+3 votes)"},
			},
		},
	}
}
