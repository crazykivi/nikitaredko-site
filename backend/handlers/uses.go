package handlers

import (
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"

	json "github.com/goccy/go-json"

	"nikitaredko-backend/cache"

	"github.com/gin-gonic/gin"
)

type UsesItem struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	URL         string `json:"url,omitempty"`
}

type UsesCategory struct {
	ID          string     `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Items       []UsesItem `json:"items"`
}

type UsesResponse struct {
	Categories  []UsesCategory `json:"categories"`
	LastUpdated string         `json:"lastUpdated"`
}

type UsesHandler struct {
	articleHandler *ArticleHandler
	cache          *cache.Cache
}

func NewUsesHandler(articleHandler *ArticleHandler, cacheManager *cache.Cache) *UsesHandler {
	return &UsesHandler{
		articleHandler: articleHandler,
		cache:          cacheManager,
	}
}

func (h *UsesHandler) GetUses(c *gin.Context) {
	cacheKey := "uses_page"
	if cached, found := h.cache.Get(cacheKey); found {
		log.Printf("[Cache] HIT: %s", cacheKey)
		c.JSON(http.StatusOK, cached)
		return
	}
	log.Printf("[Cache] MISS: %s", cacheKey)

	usesDocID := os.Getenv("USES_DOCUMENT_ID")

	var content string
	var updatedAt string

	if usesDocID != "" {
		body := map[string]interface{}{"id": usesDocID}
		data, err := h.articleHandler.callOutlineAPI("/api/documents.info", body)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Uses page not found"})
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
			if strings.EqualFold(doc.Title, "Uses") || strings.EqualFold(doc.Title, "Setup") {
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
		c.JSON(http.StatusNotFound, gin.H{"error": "Uses page not found in Outline"})
		return
	}

	categories := parseUsesMarkdown(content)
	response := UsesResponse{
		Categories:  categories,
		LastUpdated: updatedAt,
	}

	h.cache.Set(cacheKey, response)
	c.JSON(http.StatusOK, response)
}

func parseUsesMarkdown(content string) []UsesCategory {
	var categories []UsesCategory
	var currentCategory *UsesCategory
	var currentItem *UsesItem

	lines := strings.Split(content, "\n")

	reH2 := regexp.MustCompile(`^##\s+(.+)$`)
	reH3 := regexp.MustCompile(`^###\s+(.+)$`)
	reURL := regexp.MustCompile(`^https?://\S+$`)
	reAngleURLLine := regexp.MustCompile(`^<(https?://[^>]+)>\.?$`)
	reAngleURL := regexp.MustCompile(`<(https?://[^>]+)>`)
	reMdLink := regexp.MustCompile(`\[([^\]]*)\]\((https?://[^)\s]+)[^)]*\)`)

	flushItem := func() {
		if currentItem != nil && currentCategory != nil {
			currentCategory.Items = append(currentCategory.Items, *currentItem)
			currentItem = nil
		}
	}

	flushCategory := func() {
		flushItem()
		if currentCategory != nil {
			categories = append(categories, *currentCategory)
			currentCategory = nil
		}
	}

	extractInlineURL := func(s string) (string, string) {
		url := ""
		if m := reMdLink.FindStringSubmatch(s); m != nil {
			url = m[2]
			s = strings.Replace(s, m[0], m[1], 1)
		}
		for {
			m := reAngleURL.FindStringSubmatch(s)
			if m == nil {
				break
			}
			if url == "" {
				url = m[1]
			}
			s = strings.Replace(s, m[0], "", 1)
		}
		return url, strings.Join(strings.Fields(s), " ")
	}

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}

		h2Match := reH2.FindStringSubmatch(trimmed)
		h3Match := reH3.FindStringSubmatch(trimmed)
		angleURLMatch := reAngleURLLine.FindStringSubmatch(trimmed)
		isPlainURL := reURL.MatchString(trimmed)

		switch {
		case h2Match != nil:
			flushCategory()
			id := strings.ToLower(strings.ReplaceAll(h2Match[1], " ", "-"))
			currentCategory = &UsesCategory{
				ID:    id,
				Title: h2Match[1],
				Items: []UsesItem{},
			}

		case h3Match != nil:
			flushItem()
			if currentCategory != nil {
				currentItem = &UsesItem{Name: h3Match[1]}
			}

		case currentItem != nil && isPlainURL:
			if currentItem.URL == "" {
				currentItem.URL = trimmed
			}

		case currentItem != nil && angleURLMatch != nil:
			if currentItem.URL == "" {
				currentItem.URL = angleURLMatch[1]
			}

		case currentCategory != nil && currentItem == nil && currentCategory.Description == "":
			_, clean := extractInlineURL(trimmed)
			currentCategory.Description = clean

		case currentItem != nil:
			url, clean := extractInlineURL(trimmed)
			if url != "" && currentItem.URL == "" {
				currentItem.URL = url
			}
			if clean != "" {
				if currentItem.Description != "" {
					currentItem.Description += " "
				}
				currentItem.Description += clean
			}
		}
	}

	flushCategory()
	return categories
}
