package handlers

import (
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
	json "github.com/goccy/go-json"

	"nikitaredko-backend/cache"
)

type CareerStage struct {
	Period      string   `json:"period"`
	Role        string   `json:"role"`
	Company     string   `json:"company"`
	Description string   `json:"description"`
	Highlights  []string `json:"highlights"`
	Current     bool     `json:"current"`
	Type        string   `json:"type"`
}

type StackItem struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	URL         string `json:"url,omitempty"`
}

type StackGroup struct {
	ID          string      `json:"id"`
	Title       string      `json:"title"`
	Description string      `json:"description"`
	Items       []StackItem `json:"items"`
}

type AboutFact struct {
	Value string `json:"value"`
	Label string `json:"label"`
}

type AboutResponse struct {
	Intro       string        `json:"intro"`
	Facts       []AboutFact   `json:"facts"`
	Career      []CareerStage `json:"career"`
	Stack       []StackGroup  `json:"stack"`
	LastUpdated string        `json:"lastUpdated"`
}

type AboutHandler struct {
	articleHandler *ArticleHandler
	cache          *cache.Cache
}

func NewAboutHandler(articleHandler *ArticleHandler, cacheManager *cache.Cache) *AboutHandler {
	return &AboutHandler{
		articleHandler: articleHandler,
		cache:          cacheManager,
	}
}

func (h *AboutHandler) GetAbout(c *gin.Context) {
	cacheKey := "about_page"
	if cached, found := h.cache.Get(cacheKey); found {
		log.Printf("[Cache] HIT: %s", cacheKey)
		c.JSON(http.StatusOK, cached)
		return
	}
	log.Printf("[Cache] MISS: %s", cacheKey)

	aboutDocID := os.Getenv("ABOUT_DOCUMENT_ID")
	var content string
	var updatedAt string

	if aboutDocID != "" {
		body := map[string]interface{}{"id": aboutDocID}
		data, err := h.articleHandler.callOutlineAPI("/api/documents.info", body)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "About page not found"})
			return
		}
		var doc OutlineDocument
		if err := json.Unmarshal(data, &doc); err == nil {
			content = doc.Text
			if content == "" {
				content = doc.Content
			}
			updatedAt = doc.UpdatedAt
		}
	} else {
		docs, err := h.articleHandler.fetchAllDocs()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch"})
			return
		}
		for _, doc := range docs {
			if strings.EqualFold(doc.Title, "About") || strings.EqualFold(doc.Title, "Обо мне") {
				content = doc.Text
				if content == "" {
					content = doc.Content
				}
				updatedAt = doc.UpdatedAt
				break
			}
		}
	}

	if content == "" {
		c.JSON(http.StatusNotFound, gin.H{"error": "About page not found in Outline"})
		return
	}

	response := parseAboutMarkdown(content)
	response.LastUpdated = updatedAt

	h.cache.Set(cacheKey, response)
	c.JSON(http.StatusOK, response)
}

var (
	reAboutH2    = regexp.MustCompile(`^##\s+(.+)$`)
	reAboutH3    = regexp.MustCompile(`^###\s+(.+)$`)
	reAboutList  = regexp.MustCompile(`^[-*]\s+(.+)$`)
	reAboutBold  = regexp.MustCompile(`^\*\*(.+?)\*\*`)
	reAboutLink  = regexp.MustCompile(`^\[([^\]]+)\]\((https?://[^)\s]+)[^)]*\)`)
	reInlineLink = regexp.MustCompile(`\[([^\]]+)\]\([^)]*\)`)
	reMarkdown   = regexp.MustCompile(`[\*_` + "`" + `~]`)
)

var typeRules = map[string][]string{
	"work":       {"work", "job", "работ", "служба"},
	"pet":        {"pet", "пет", "проект", "personal", "инди", "свои"},
	"freelance":  {"freelance", "фриланс"},
	"education":  {"study", "education", "уч", "стаж", "колледж", "универ"},
	"opensource": {"opensource", "open source", "github", "контрибуц"},
	"community":  {"community", "сообщество", "митап", "доклад"},
}

func parseAboutMarkdown(content string) AboutResponse {
	var resp AboutResponse
	resp.Facts = []AboutFact{}
	resp.Career = []CareerStage{}
	resp.Stack = []StackGroup{}

	lines := strings.Split(content, "\n")
	section := ""
	var introLines []string
	var currentStage *CareerStage
	var currentGroup *StackGroup

	flushStage := func() {
		if currentStage != nil {
			resp.Career = append(resp.Career, *currentStage)
			currentStage = nil
		}
	}
	flushGroup := func() {
		if currentGroup != nil {
			resp.Stack = append(resp.Stack, *currentGroup)
			currentGroup = nil
		}
	}

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		if strings.HasPrefix(trimmed, "# ") {
			continue
		}

		if m := reAboutH2.FindStringSubmatch(trimmed); m != nil {
			flushStage()
			flushGroup()
			title := strings.ToLower(m[1])
			switch {
			case strings.Contains(title, "карьер") || strings.Contains(title, "career") || strings.Contains(title, "опыт"):
				section = "career"
			case strings.Contains(title, "стек") || strings.Contains(title, "stack") || strings.Contains(title, "технолог") || strings.Contains(title, "навык") || strings.Contains(title, "skill"):
				section = "stack"
			case strings.Contains(title, "факт") || strings.Contains(title, "fact") || strings.Contains(title, "цифр"):
				section = "facts"
			default:
				section = "other"
			}
			continue
		}

		if m := reAboutH3.FindStringSubmatch(trimmed); m != nil {
			switch section {
			case "career":
				flushStage()
				currentStage = parseCareerHeader(m[1])
			case "stack":
				flushGroup()
				currentGroup = &StackGroup{
					ID:    strings.ToLower(strings.ReplaceAll(m[1], " ", "-")),
					Title: m[1],
					Items: []StackItem{},
				}
			}
			continue
		}

		if trimmed == "" {
			continue
		}

		if m := reAboutList.FindStringSubmatch(trimmed); m != nil {
			itemText := m[1]
			switch section {
			case "career":
				if currentStage != nil {
					currentStage.Highlights = append(currentStage.Highlights, cleanInline(itemText))
				}
			case "stack":
				if currentGroup != nil {
					currentGroup.Items = append(currentGroup.Items, parseStackItem(itemText))
				}
			case "facts":
				resp.Facts = append(resp.Facts, parseFact(itemText))
			}
			continue
		}

		switch section {
		case "":
			introLines = append(introLines, trimmed)
		case "career":
			if currentStage != nil {
				if currentStage.Description != "" {
					currentStage.Description += " "
				}
				currentStage.Description += cleanInline(trimmed)
			}
		case "stack":
			if currentGroup != nil && currentGroup.Description == "" {
				currentGroup.Description = cleanInline(trimmed)
			}
		}
	}
	flushStage()
	flushGroup()

	resp.Intro = strings.Join(introLines, "\n")
	return resp
}

func parseCareerHeader(header string) *CareerStage {
	parts := strings.Split(header, "|")
	stage := &CareerStage{Highlights: []string{}}

	if len(parts) >= 1 {
		stage.Period = strings.TrimSpace(parts[0])
	}
	if len(parts) >= 2 {
		stage.Role = strings.TrimSpace(parts[1])
	}
	if len(parts) >= 3 {
		stage.Company = strings.TrimSpace(parts[2])
	}
	if len(parts) >= 4 {
		raw := strings.TrimSpace(parts[3])
		if idx := strings.Index(strings.ToLower(raw), "type:"); idx != -1 {
			raw = strings.TrimSpace(raw[idx+5:])
		}
		stage.Type = normalizeStageType(raw)
	}

	if stage.Type == "" {
		stage.Type = guessStageType(stage)
	}

	lower := strings.ToLower(stage.Period)
	stage.Current = strings.Contains(lower, "н.в") ||
		strings.Contains(lower, "now") ||
		strings.Contains(lower, "present") ||
		strings.Contains(lower, "по наст")

	return stage
}

func stripMarkdown(s string) string {
	return strings.ToLower(reMarkdown.ReplaceAllString(s, ""))
}

func normalizeStageType(s string) string {
	clean := stripMarkdown(s)

	for typeName, patterns := range typeRules {
		for _, pattern := range patterns {
			if strings.Contains(clean, pattern) {
				return typeName
			}
		}
	}

	if clean == "" {
		return ""
	}
	return "other"
}

func guessStageType(stage *CareerStage) string {
	text := stripMarkdown(stage.Role + " " + stage.Company)

	for typeName, patterns := range typeRules {
		for _, pattern := range patterns {
			if strings.Contains(text, pattern) {
				return typeName
			}
		}
	}

	return "work"
}

func parseStackItem(text string) StackItem {
	item := StackItem{}
	rest := text

	if m := reAboutLink.FindStringSubmatch(rest); m != nil {
		item.Name = m[1]
		item.URL = m[2]
		rest = strings.TrimSpace(rest[len(m[0]):])
	} else if m := reAboutBold.FindStringSubmatch(rest); m != nil {
		item.Name = m[1]
		rest = strings.TrimSpace(rest[len(m[0]):])
	} else if idx := strings.IndexAny(rest, "—"); idx > 0 {
		item.Name = strings.TrimSpace(rest[:idx])
		rest = rest[idx:]
	} else {
		item.Name = strings.TrimSpace(rest)
		rest = ""
	}

	rest = strings.TrimPrefix(rest, "—")
	rest = strings.TrimPrefix(rest, "-")
	rest = strings.TrimPrefix(rest, ":")
	item.Description = cleanInline(rest)
	return item
}

func parseFact(text string) AboutFact {
	if m := reAboutBold.FindStringSubmatch(text); m != nil {
		rest := strings.TrimSpace(text[len(m[0]):])
		rest = strings.TrimPrefix(strings.TrimPrefix(rest, "—"), "-")
		return AboutFact{Value: m[1], Label: cleanInline(rest)}
	}
	return AboutFact{Label: cleanInline(text)}
}

func cleanInline(s string) string {
	s = reInlineLink.ReplaceAllString(s, "$1")
	s = strings.ReplaceAll(s, "**", "")
	s = strings.ReplaceAll(s, "`", "")
	return strings.TrimSpace(s)
}
