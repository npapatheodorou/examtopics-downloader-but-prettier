package templates

import (
	"strings"
	"testing"
)

func TestEmbeddedTemplatePresent(t *testing.T) {
	content := strings.TrimSpace(EmbeddedTemplate)
	if content == "" {
		t.Fatal("embedded template is empty")
	}

	if !strings.Contains(content, "<!DOCTYPE html>") {
		t.Fatal("embedded template does not look like an HTML document")
	}
}
