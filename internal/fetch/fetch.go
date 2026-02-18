package fetch

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"examtopics-downloader/internal/constants"
	"examtopics-downloader/internal/utils"

	"github.com/PuerkitoBio/goquery"
)

var client = utils.NewHTTPClient()

var (
	providerHrefPattern = regexp.MustCompile(`(?i)^/exams/([a-z0-9-]+)/?$`)
)

func FetchURL(url string, client http.Client) []byte {
	backoff := constants.InitalBackoff

	for attempt := 0; attempt <= constants.MaxRetries; attempt++ {
		if attempt > 0 {
			delay := utils.DelayTime(backoff)
			debugf("Retry attempt %d for URL: %s after waiting %v", attempt, url, delay)
			utils.Sleep(delay)
			backoff = utils.BackoffTime(backoff, constants.BackoffFactor)
		}

		resp, err := client.Get(url)
		if err != nil {
			debugf("failed to fetch URL (attempt %d): %v", attempt, err)
			continue
		}

		if resp.StatusCode == http.StatusOK {
			body, err := io.ReadAll(resp.Body)
			resp.Body.Close()
			if err != nil {
				debugf("failed to read response body: %v", err)
				return nil
			}
			return body
		}
		resp.Body.Close()

		if resp.StatusCode != http.StatusServiceUnavailable {
			debugf("request failed with status code: %d", resp.StatusCode)
			return nil
		}
	}

	debugf("exhausted retries for URL: %s", url)
	return nil
}

func ParseHTML(url string, client http.Client) (*goquery.Document, error) {
	body := FetchURL(url, client)
	if body == nil {
		return nil, fmt.Errorf("empty response body from URL %q", url)
	}

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML from URL %q: %w", url, err)
	}

	return doc, nil
}

// Fetches total number of pages
func getMaxNumPages(url string) int {
	doc, err := ParseHTML(url, *client)
	if err != nil {
		debugf("failed parsing HTML for number of pages: %v", err)
		return 1
	}

	var pageCount int
	doc.Find(".discussion-list-page-indicator strong").Each(func(i int, s *goquery.Selection) {
		if i == 1 {
			pageCount, _ = strconv.Atoi(strings.TrimSpace(s.Text()))
		}
	})

	// Handle the null case
	if pageCount == 0 {
		pageCount = 1
	}

	return pageCount
}

func GetAllProviders() []string {
	doc, err := ParseHTML("https://www.examtopics.com/exams/", *client)
	if err != nil {
		debugf("failed to parse HTML for providers: %v", err)
		return nil
	}

	seen := map[string]struct{}{}
	providers := make([]string, 0, 32)

	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		href, exists := s.Attr("href")
		if !exists {
			return
		}

		href = strings.TrimSpace(strings.ToLower(href))
		matches := providerHrefPattern.FindStringSubmatch(href)
		if len(matches) != 2 {
			return
		}

		provider := strings.TrimSpace(matches[1])
		if provider == "" {
			return
		}
		if _, exists := seen[provider]; exists {
			return
		}
		seen[provider] = struct{}{}
		providers = append(providers, provider)
	})

	sort.Strings(providers)
	return providers
}

func GetProviderExams(providerName string) []string {
	providerName = strings.TrimSpace(strings.ToLower(providerName))
	baseURL := fmt.Sprintf("https://www.examtopics.com/exams/%s/", providerName)
	doc, err := ParseHTML(baseURL, *client)
	if err != nil {
		debugf("failed to parse HTML for provider exams: %v", err)
		return nil
	}

	examHrefPattern := regexp.MustCompile(fmt.Sprintf(`(?i)^/exams/%s/([a-z0-9-]+)/?$`, regexp.QuoteMeta(providerName)))
	seen := map[string]struct{}{}
	allExams := make([]string, 0, 32)

	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		href, exists := s.Attr("href")
		if !exists {
			return
		}

		cleanHref := strings.TrimSpace(strings.ToLower(href))
		matches := examHrefPattern.FindStringSubmatch(cleanHref)
		if len(matches) != 2 {
			return
		}

		examSlug := strings.TrimSpace(matches[1])
		if examSlug == "" {
			return
		}

		normalized := fmt.Sprintf("/exams/%s/%s/", providerName, examSlug)
		if _, exists := seen[normalized]; exists {
			return
		}
		seen[normalized] = struct{}{}
		allExams = append(allExams, normalized)
	})

	sort.Strings(allExams)
	return allExams
}

// Extracts matching links from a single page
func getLinksFromPage(url string, grepStr string) []string {
	doc, err := ParseHTML(url, *client)
	if err != nil {
		debugf("failed to parse HTML for %s: %v", url, err)
		return nil
	}

	var matchingLinks []string
	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		href, exists := s.Attr("href")
		if exists && utils.GrepString(href, "/discussions") && utils.GrepString(href, grepStr) {
			matchingLinks = append(matchingLinks, href)
		}
	})

	return matchingLinks
}
