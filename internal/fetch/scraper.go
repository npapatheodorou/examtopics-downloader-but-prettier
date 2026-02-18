package fetch

import (
	"fmt"
	"html"
	"net/url"
	"regexp"
	"strings"
	"sync"

	"examtopics-downloader/internal/constants"
	"examtopics-downloader/internal/models"
	"examtopics-downloader/internal/utils"

	"github.com/PuerkitoBio/goquery"
	"github.com/cheggaaa/pb/v3"
)

func getDataFromLink(link string) *models.QuestionData {
	doc, err := ParseHTML(link, *client)
	if err != nil {
		debugf("failed parsing HTML data from link: %v", err)
		return nil
	}

	var allQuestions []string
	doc.Find("li.multi-choice-item").Each(func(i int, s *goquery.Selection) {
		allQuestions = append(allQuestions, utils.CleanText(s.Text()))
	})

	answer := strings.TrimSpace(doc.Find(".correct-answer").Text())

	return &models.QuestionData{
		Title:        utils.CleanText(doc.Find("h1").Text()),
		Header:       strings.ReplaceAll(strings.TrimSpace(doc.Find(".question-discussion-header").Text()), "\t", ""),
		Content:      utils.CleanText(doc.Find(".card-text").Text()),
		ExhibitURLs:  extractExhibitImageURLs(doc),
		Questions:    allQuestions,
		Answer:       answer,
		Timestamp:    utils.CleanText(doc.Find(".discussion-meta-data > i").Text()),
		QuestionLink: link,
		Comments:     extractDiscussionComments(doc),
	}
}

func extractExhibitImageURLs(doc *goquery.Document) []string {
	var urls []string
	seen := map[string]struct{}{}

	add := func(raw string) {
		normalized := normalizeExhibitURL(raw)
		if normalized == "" {
			return
		}
		if _, exists := seen[normalized]; exists {
			return
		}
		seen[normalized] = struct{}{}
		urls = append(urls, normalized)
	}

	doc.Find(".card-text img").Each(func(i int, s *goquery.Selection) {
		if src, ok := s.Attr("src"); ok {
			add(src)
		}
		if src, ok := s.Attr("data-src"); ok {
			add(src)
		}
		if src, ok := s.Attr("data-original"); ok {
			add(src)
		}
		if src, ok := s.Attr("data-lazy-src"); ok {
			add(src)
		}
		if srcSet, ok := s.Attr("srcset"); ok {
			add(firstURLFromSrcset(srcSet))
		}
	})

	return urls
}

func normalizeExhibitURL(raw string) string {
	raw = strings.TrimSpace(html.UnescapeString(raw))
	if raw == "" || strings.HasPrefix(raw, "data:") {
		return ""
	}

	if strings.HasPrefix(raw, "//") {
		raw = "https:" + raw
	} else if strings.HasPrefix(raw, "/") {
		raw = "https://www.examtopics.com" + raw
	}

	u, err := url.Parse(raw)
	if err != nil {
		return ""
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return ""
	}

	return u.String()
}

func firstURLFromSrcset(srcset string) string {
	items := strings.Split(srcset, ",")
	if len(items) == 0 {
		return ""
	}

	first := strings.TrimSpace(items[0])
	if first == "" {
		return ""
	}
	parts := strings.Fields(first)
	if len(parts) == 0 {
		return ""
	}

	return parts[0]
}

func extractDiscussionComments(doc *goquery.Document) []models.CommentData {
	var comments []models.CommentData
	answerLetterPattern := regexp.MustCompile(`\b([A-F])\b`)

	doc.Find(".discussion-container .comment-container").Each(func(i int, s *goquery.Selection) {
		user := strings.TrimSpace(s.Find(".comment-username").First().Text())
		if user == "" {
			user = "Anonymous"
		}

		answer := ""
		answerText := strings.TrimSpace(s.Find(".comment-selected-answers strong").First().Text())
		if answerText == "" {
			answerText = strings.TrimSpace(s.Find(".comment-selected-answers").First().Text())
		}
		if m := answerLetterPattern.FindStringSubmatch(strings.ToUpper(answerText)); len(m) == 2 {
			answer = m[1]
		}

		content := normalizeCommentText(s.Find(".comment-content").First().Text())
		if content == "" {
			return
		}

		comments = append(comments, models.CommentData{
			User:   user,
			Answer: answer,
			Text:   content,
		})
	})

	return comments
}

func normalizeCommentText(raw string) string {
	raw = strings.ReplaceAll(raw, "\r\n", "\n")
	raw = strings.ReplaceAll(raw, "\r", "\n")
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return ""
	}

	lines := strings.Split(raw, "\n")
	cleaned := make([]string, 0, len(lines))
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}
		cleaned = append(cleaned, trimmed)
	}

	return strings.Join(cleaned, "\n")
}

func fetchAllPageLinksConcurrently(providerName, grepStr string, numPages, concurrency int, onPageProcessed func()) []string {
	var wg sync.WaitGroup
	sem := make(chan struct{}, concurrency)
	results := make(chan []string, numPages)

	rateLimiter := utils.CreateRateLimiter(constants.RequestsPerSecond)
	defer rateLimiter.Stop()

	for i := 1; i <= numPages; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()

			<-rateLimiter.C

			url := fmt.Sprintf("https://www.examtopics.com/discussions/%s/%d", providerName, i)
			results <- getLinksFromPage(url, grepStr)
			if onPageProcessed != nil {
				onPageProcessed()
			}
		}(i)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	// about 10 questions per examtopics page, we can preallocate
	all := make([]string, 0, numPages*10)
	for res := range results {
		all = append(all, res...)
	}

	return all
}

// Main concurrent page scraping logic
func GetAllPages(providerName string, grepStr string) []models.QuestionData {
	baseURL := fmt.Sprintf("https://www.examtopics.com/discussions/%s/", providerName)
	numPages := getMaxNumPages(baseURL)
	startTime := utils.StartTime()
	bar := pb.StartNew(numPages)

	allLinks := fetchAllPageLinksConcurrently(providerName, grepStr, numPages, constants.MaxConcurrentRequests, func() {
		bar.Increment()
	})

	unique := utils.DeduplicateLinks(allLinks)
	sortedLinks := utils.SortLinksByQuestionNumber(unique)
	bar.SetTotal(int64(numPages + len(sortedLinks)))

	if len(sortedLinks) == 0 {
		bar.Finish()
		fmt.Println("No matching questions were found.")
		return nil
	}

	var wg sync.WaitGroup
	sem := make(chan struct{}, constants.MaxConcurrentRequests)
	results := make([]*models.QuestionData, len(sortedLinks))

	rateLimiter := utils.CreateRateLimiter(constants.RequestsPerSecond)
	defer rateLimiter.Stop()

	for i, link := range sortedLinks {
		wg.Add(1)
		url := utils.AddToBaseUrl(link)

		go func(i int, url string) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()

			<-rateLimiter.C

			data := getDataFromLink(url)
			if data != nil {
				results[i] = data
			}
			bar.Increment()
		}(i, url)
	}

	wg.Wait()
	bar.Finish()
	// Filter out nil entries
	var finalData []models.QuestionData
	for _, entry := range results {
		if entry != nil {
			finalData = append(finalData, *entry)
		}
	}

	fmt.Printf("Extraction complete in %s.\n", utils.TimeSince(startTime))

	return finalData
}
