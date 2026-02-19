package fetch

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"
)

const discussionExamCacheTTL = 24 * time.Hour

type discussionExamCacheEntry struct {
	ExamSlugs []string `json:"exam_slugs"`
	UpdatedAt int64    `json:"updated_at_unix"`
}

type discussionExamCacheFile struct {
	Providers map[string]discussionExamCacheEntry `json:"providers"`
}

var (
	discussionExamCacheMu     sync.Mutex
	discussionExamCacheLoaded bool
	discussionExamCache       = discussionExamCacheFile{Providers: map[string]discussionExamCacheEntry{}}
)

func getCachedDiscussionExamSlugs(providerName string) ([]string, bool) {
	providerName = strings.TrimSpace(strings.ToLower(providerName))
	if providerName == "" {
		return nil, false
	}

	discussionExamCacheMu.Lock()
	defer discussionExamCacheMu.Unlock()

	ensureDiscussionExamCacheLoadedLocked()
	entry, exists := discussionExamCache.Providers[providerName]
	if !exists {
		return nil, false
	}

	if entry.UpdatedAt <= 0 || time.Since(time.Unix(entry.UpdatedAt, 0)) > discussionExamCacheTTL {
		delete(discussionExamCache.Providers, providerName)
		saveDiscussionExamCacheLocked()
		return nil, false
	}

	copied := append([]string(nil), entry.ExamSlugs...)
	sort.Strings(copied)
	return copied, len(copied) > 0
}

func setCachedDiscussionExamSlugs(providerName string, examSlugs []string) {
	providerName = strings.TrimSpace(strings.ToLower(providerName))
	if providerName == "" || len(examSlugs) == 0 {
		return
	}

	normalized := append([]string(nil), examSlugs...)
	sort.Strings(normalized)

	discussionExamCacheMu.Lock()
	defer discussionExamCacheMu.Unlock()

	ensureDiscussionExamCacheLoadedLocked()
	discussionExamCache.Providers[providerName] = discussionExamCacheEntry{
		ExamSlugs: normalized,
		UpdatedAt: time.Now().Unix(),
	}
	saveDiscussionExamCacheLocked()
}

func ensureDiscussionExamCacheLoadedLocked() {
	if discussionExamCacheLoaded {
		return
	}
	discussionExamCacheLoaded = true

	cachePath := discussionExamCachePath()
	payload, err := os.ReadFile(cachePath)
	if err != nil {
		return
	}

	var loaded discussionExamCacheFile
	if err := json.Unmarshal(payload, &loaded); err != nil {
		debugf("failed to parse exam cache file %q: %v", cachePath, err)
		return
	}
	if loaded.Providers == nil {
		loaded.Providers = map[string]discussionExamCacheEntry{}
	}
	discussionExamCache = loaded
}

func saveDiscussionExamCacheLocked() {
	cachePath := discussionExamCachePath()
	if cachePath == "" {
		return
	}

	dir := filepath.Dir(cachePath)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		debugf("failed to create cache dir %q: %v", dir, err)
		return
	}

	payload, err := json.MarshalIndent(discussionExamCache, "", "  ")
	if err != nil {
		debugf("failed to marshal exam cache: %v", err)
		return
	}
	if err := os.WriteFile(cachePath, payload, 0o644); err != nil {
		debugf("failed to write exam cache file %q: %v", cachePath, err)
	}
}

func discussionExamCachePath() string {
	baseDir, err := os.UserCacheDir()
	if err == nil && strings.TrimSpace(baseDir) != "" {
		return filepath.Join(baseDir, "examtopics-downloader", "discussion_exam_slugs.json")
	}
	return filepath.Join(".", ".examtopics_discussion_exam_cache.json")
}
