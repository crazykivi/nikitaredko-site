package handlers

import (
	"log"
	"net/http"
	"os"
	"strings"

	json "github.com/goccy/go-json"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/text"

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

	md := goldmark.New()
	reader := text.NewReader([]byte(content))
	doc := md.Parser().Parse(reader)
	source := []byte(content)

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

	ast.Walk(doc, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}
		switch v := n.(type) {
		case *ast.Heading:
			switch v.Level {
			case 1:
				return ast.WalkSkipChildren, nil
			case 2:
				flushCategory()
				title := renderNodeText(v, source)
				currentCategory = &UsesCategory{
					ID:    strings.ToLower(strings.ReplaceAll(title, " ", "-")),
					Title: title,
					Items: []UsesItem{},
				}
				return ast.WalkSkipChildren, nil
			case 3:
				flushItem()
				if currentCategory != nil {
					currentItem = &UsesItem{
						Name: renderNodeText(v, source),
					}
				}
				return ast.WalkSkipChildren, nil
			}
		case *ast.Paragraph, *ast.TextBlock:
			if _, isListItem := n.Parent().(*ast.ListItem); isListItem {
				return ast.WalkSkipChildren, nil
			}

			url := findFirstURL(v, source)
			cleanText := renderNodeTextSkipAutoLinks(v, source)

			switch {
			case currentCategory != nil && currentItem == nil:
				if currentCategory.Description != "" {
					currentCategory.Description += " "
				}
				currentCategory.Description += cleanText

			case currentItem != nil:
				if url != "" && currentItem.URL == "" {
					currentItem.URL = url
				}
				if cleanText != "" {
					if currentItem.Description != "" {
						currentItem.Description += " "
					}
					currentItem.Description += cleanText
				}
			}
			return ast.WalkSkipChildren, nil
		}

		return ast.WalkContinue, nil
	})

	flushCategory()
	return categories
}
