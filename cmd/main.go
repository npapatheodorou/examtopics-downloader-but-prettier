package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"math"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"examtopics-downloader/internal/fetch"
	"examtopics-downloader/internal/utils"
)

func main() {
	debug := flag.Bool("debug", false, "Enable debug logs")
	flag.Parse()
	fetch.SetDebug(*debug)

	reader := bufio.NewReader(os.Stdin)

	providers := fetch.GetAllProviders()
	if len(providers) == 0 {
		log.Fatal("no providers found")
	}

	providerIdx, err := promptSelection(reader, "Available Providers", providers, formatProviderName)
	if err != nil {
		log.Fatalf("failed reading provider selection: %v", err)
	}
	selectedProvider := providers[providerIdx]

	examLinks := fetch.GetProviderExams(selectedProvider)
	examSlugs := extractExamSlugs(selectedProvider, examLinks)
	if len(examSlugs) == 0 {
		log.Fatalf("no exams found for provider %q", selectedProvider)
	}

	examIdx, err := promptSelection(reader, fmt.Sprintf("Available Exams for %s", formatProviderName(selectedProvider)), examSlugs, func(s string) string {
		return s
	})
	if err != nil {
		log.Fatalf("failed reading exam selection: %v", err)
	}
	selectedExam := examSlugs[examIdx]

	fmt.Printf("\nStarting extraction for %s / %s...\n", formatProviderName(selectedProvider), selectedExam)
	links := fetch.GetAllPages(selectedProvider, selectedExam)
	if len(links) == 0 {
		log.Fatal("no matching questions were extracted")
	}

	outputPath := defaultOutputPath(selectedProvider, selectedExam)
	savedFiles, err := utils.WriteData(links, outputPath, true)
	if err != nil {
		log.Fatalf("failed writing output: %v", err)
	}

	fmt.Printf("Successfully saved output: %s\n", strings.Join(savedFiles, ", "))
}

func promptSelection(reader *bufio.Reader, title string, options []string, formatter func(string) string) (int, error) {
	all := make([]selectionOption, 0, len(options))
	for i, opt := range options {
		all = append(all, selectionOption{
			RawIndex: i,
			Label:    formatter(opt),
		})
	}

	filter := ""
	for {
		filtered := filterOptions(all, filter)
		if len(filtered) == 0 {
			fmt.Printf("\n%s\n", title)
			fmt.Printf("No results for filter %q. Enter / to clear filter.\n", filter)
		} else {
			fmt.Printf("\n%s (%d shown of %d)\n", title, len(filtered), len(all))
			printOptionsInColumns(filtered)
		}

		fmt.Printf("Type number to select, /text to filter, / to clear: ")
		raw, err := reader.ReadString('\n')
		if err != nil {
			return 0, err
		}

		raw = strings.TrimSpace(raw)
		if raw == "" {
			continue
		}
		if strings.HasPrefix(raw, "/") {
			filter = strings.TrimSpace(strings.TrimPrefix(raw, "/"))
			continue
		}

		choice, err := strconv.Atoi(raw)
		if err != nil || choice < 1 || choice > len(filtered) {
			fmt.Println("Invalid selection. Please enter a valid number.")
			continue
		}

		return filtered[choice-1].RawIndex, nil
	}
}

type selectionOption struct {
	RawIndex int
	Label    string
}

func filterOptions(options []selectionOption, filter string) []selectionOption {
	if strings.TrimSpace(filter) == "" {
		return options
	}

	filter = strings.ToLower(strings.TrimSpace(filter))
	filtered := make([]selectionOption, 0, len(options))
	for _, opt := range options {
		if strings.Contains(strings.ToLower(opt.Label), filter) {
			filtered = append(filtered, opt)
		}
	}
	return filtered
}

func printOptionsInColumns(options []selectionOption) {
	if len(options) == 0 {
		return
	}

	lines := make([]string, 0, len(options))
	maxWidth := 0
	for i, opt := range options {
		line := fmt.Sprintf("%3d) %s", i+1, strings.TrimSpace(opt.Label))
		lines = append(lines, line)
		if len(line) > maxWidth {
			maxWidth = len(line)
		}
	}

	colWidth := maxWidth + 4
	if colWidth < 24 {
		colWidth = 24
	}
	targetWidth := 120
	cols := targetWidth / colWidth
	if cols < 1 {
		cols = 1
	}
	if cols > 4 {
		cols = 4
	}

	rows := int(math.Ceil(float64(len(lines)) / float64(cols)))
	for r := 0; r < rows; r++ {
		var row strings.Builder
		for c := 0; c < cols; c++ {
			idx := c*rows + r
			if idx >= len(lines) {
				continue
			}
			if c > 0 {
				row.WriteString("  ")
			}
			row.WriteString(fmt.Sprintf("%-*s", colWidth, lines[idx]))
		}
		fmt.Println(strings.TrimRight(row.String(), " "))
	}
}

func extractExamSlugs(provider string, examLinks []string) []string {
	pattern := regexp.MustCompile(fmt.Sprintf(`(?i)^/exams/%s/([^/]+)/?$`, regexp.QuoteMeta(strings.ToLower(strings.TrimSpace(provider)))))
	seen := map[string]struct{}{}
	out := make([]string, 0, len(examLinks))

	for _, link := range examLinks {
		matches := pattern.FindStringSubmatch(strings.ToLower(strings.TrimSpace(link)))
		if len(matches) != 2 {
			continue
		}

		examSlug := strings.TrimSpace(matches[1])
		if examSlug == "" {
			continue
		}
		if _, exists := seen[examSlug]; exists {
			continue
		}

		seen[examSlug] = struct{}{}
		out = append(out, examSlug)
	}

	sort.Strings(out)
	return out
}

func defaultOutputPath(provider, examSlug string) string {
	baseProvider := sanitizeFilenameSegment(provider)
	baseExamCode := sanitizeFilenameSegment(examSlug)
	if baseProvider == "" {
		baseProvider = "examtopics"
	}
	if baseExamCode == "" {
		baseExamCode = "output"
	}

	return fmt.Sprintf("%s_%s.html", baseProvider, baseExamCode)
}

func sanitizeFilenameSegment(input string) string {
	segment := strings.TrimSpace(strings.ToLower(input))
	if segment == "" {
		return ""
	}

	segment = strings.ReplaceAll(segment, " ", "-")
	invalidChars := regexp.MustCompile(`[^a-z0-9._-]+`)
	segment = invalidChars.ReplaceAllString(segment, "-")
	segment = strings.Trim(segment, "-._")

	return segment
}

func formatProviderName(provider string) string {
	provider = strings.TrimSpace(strings.ToLower(provider))
	if provider == "" {
		return "Unknown"
	}

	overrides := map[string]string{
		"aws":                "AWS",
		"ec-council":         "EC-Council",
		"eccouncil":          "EC-Council",
		"isc2":               "ISC2",
		"isaca":              "ISACA",
		"paloalto-networks":  "Palo Alto Networks",
		"palo-alto-networks": "Palo Alto Networks",
		"servicenow":         "ServiceNow",
		"vmware":             "VMware",
		"lpi":                "LPI",
	}
	if label, ok := overrides[provider]; ok {
		return label
	}

	parts := strings.Fields(strings.ReplaceAll(provider, "-", " "))
	for i, part := range parts {
		if part == "" {
			continue
		}
		parts[i] = strings.ToUpper(part[:1]) + part[1:]
	}
	return strings.Join(parts, " ")
}
