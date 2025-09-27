package utils

import (
	"bufio"
	"bytes"
	"regexp"
	"examtopics-downloader/internal/models"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/mandolyte/mdtopdf"
	"github.com/yuin/goldmark"
)

func writeFile(filename string, content any) {
	file := CreateFile(filename)
	defer file.Close()

	switch v := content.(type) {
	case string:
		fmt.Fprintln(file, v)
	case []string:
		for _, line := range v {
			fmt.Fprintln(file, line)
		}
	default:
		log.Printf("writeFile: unsupported content type %T", v)
		return
	}
}

func WriteData(dataList []models.QuestionData, outputPath string, commentBool bool, fileType string) {
	file := CreateFile(outputPath)
	defer file.Close()

	fmt.Fprintf(file, "# Exam Topics Questions\n\n")
	fmt.Fprintf(file, "@thatonecodes\n\n")

	for _, data := range dataList {
		if data.Title == "" {
			continue
		}

		fmt.Fprintf(file, "## %s\n\n", data.Title)
		fmt.Fprintf(file, "%s\n\n", data.Header)

		if data.Content != "" {
			fmt.Fprintf(file, "%s\n\n", data.Content)
		}

		for _, question := range data.Questions {
			fmt.Fprintf(file, "%s\n\n", question)
		}

		fmt.Fprintf(file, "**Answer: %s**\n\n", data.Answer)
		fmt.Fprintf(file, "**Timestamp: %s**\n\n", data.Timestamp)
		fmt.Fprintf(file, "[View on ExamTopics](%s)\n\n", data.QuestionLink)

		if commentBool {
			fmt.Fprintf(file, "Comments: %s\n", data.Comments)
		}

		fmt.Fprintf(file, "----------------------------------------\n\n")
	}

	switch fileType {
	case "pdf":
		mdContent, err := os.ReadFile(outputPath)
		if err != nil {
			log.Printf("failed to read markdown file: %v", err)
			return
		}

		opts := []mdtopdf.RenderOption{
			mdtopdf.IsHorizontalRuleNewPage(true), // treat --- as new page
		}

		pdfName := strings.TrimSuffix(outputPath, ".md") + ".pdf"
		renderer := mdtopdf.NewPdfRenderer("portrait", "A4", pdfName, "", opts, mdtopdf.LIGHT)
		if err := renderer.Process(mdContent); err != nil {
			log.Printf("mdtopdf conversion failed: %v", err)
			return
		}
		deleteMarkdownFile(outputPath)
	case "html":
		mdBytes, err := os.ReadFile(outputPath)
		if err != nil {
			log.Printf("failed to read file for html conversion: %v", err)
			return
		}

		html, err := mdToHTML(mdBytes)
		if err != nil {
			log.Printf("mdtohtml conversion failed: %v", err)
			return
		}

		fileName := strings.TrimSuffix(outputPath, ".md") + ".html"
		err = os.WriteFile(fileName, html, 0644)
		if err != nil {
			log.Printf("failed to write html file: %v", err)
			return
		}
		deleteMarkdownFile(outputPath)
	case "text":
		mdBytes, err := os.ReadFile(outputPath)
		if err != nil {
			log.Printf("failed to read file for text conversion: %v", err)
			return
		}

		txt := mdToText(string(mdBytes))

		fileName := strings.TrimSuffix(outputPath, ".md") + ".txt"
		err = os.WriteFile(fileName, []byte(txt), 0644)
		if err != nil {
			log.Printf("failed to write text file: %v", err)
			return
		}
		deleteMarkdownFile(outputPath)
	}
}

func mdToHTML(md []byte) ([]byte, error) {
	var buf bytes.Buffer
	if err := goldmark.Convert(md, &buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func deleteMarkdownFile(filePath string) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Delete Markdown file after conversion? (y/n): ")
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(strings.ToLower(input))

	if input == "y" || input == "yes" {
		fmt.Println("Deleting file...")
		os.Remove(filePath)
	} else {
		fmt.Println("Keeping file.")
	}
}

func mdToText(md string) string {
	text := md
	// Markdown headers (#, ##, ###)
	header := regexp.MustCompile(`(?m)^#{1,6}\s*`)
	text = header.ReplaceAllString(text, "")
	// bold/italic symbols (*, **, _)
	formatting := regexp.MustCompile(`(\*\*|\*|__|_)`)
	text = formatting.ReplaceAllString(text, "")
	// links but keep link text [text](url) → text
	link := regexp.MustCompile(`\[(.*?)\]\(.*?\)`)
	text = link.ReplaceAllString(text, "$1")
	// images ![alt](url)
	image := regexp.MustCompile(`!\[.*?\]\(.*?\)`)
	text = image.ReplaceAllString(text, "")

	return text
}

func SaveLinks(filename string, links []models.QuestionData) {
	var fullLinks []string
	for _, link := range links {
		fullLinks = append(fullLinks, link.QuestionLink)
	}
	writeFile(filename, fullLinks)
}
