package handlers

import (
	"bytes"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	json "github.com/goccy/go-json"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/text"

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
		data, err := h.articleHandler.callOutlineAPI(c.Request.Context(), "/api/documents.info", body)
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

	md := goldmark.New()
	reader := text.NewReader([]byte(content))
	doc := md.Parser().Parse(reader)
	source := []byte(content)
	section := ""
	var introBuf bytes.Buffer
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
				flushStage()
				flushGroup()
				title := strings.ToLower(renderNodeText(v, source))
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
				return ast.WalkSkipChildren, nil

			case 3:
				txt := renderNodeText(v, source)
				switch section {
				case "career":
					flushStage()
					currentStage = parseCareerHeader(txt)
				case "stack":
					flushGroup()
					currentGroup = &StackGroup{
						ID:    strings.ToLower(strings.ReplaceAll(txt, " ", "-")),
						Title: txt,
						Items: []StackItem{},
					}
				}
				return ast.WalkSkipChildren, nil
			}

		case *ast.ListItem:
			switch section {
			case "career":
				if currentStage != nil {
					currentStage.Highlights = append(currentStage.Highlights, renderNodeText(v, source))
				}
			case "stack":
				if currentGroup != nil {
					if p := getFirstParagraph(v); p != nil {
						currentGroup.Items = append(currentGroup.Items, extractStackItem(p, source))
					}
				}
			case "facts":
				if p := getFirstParagraph(v); p != nil {
					resp.Facts = append(resp.Facts, extractFact(p, source))
				}
			}
			return ast.WalkSkipChildren, nil

		case *ast.Paragraph, *ast.TextBlock:
			if _, isListItem := n.Parent().(*ast.ListItem); isListItem {
				return ast.WalkSkipChildren, nil
			}

			switch section {
			case "":
				introBuf.WriteString(renderNodeText(v, source))
				introBuf.WriteString("\n")
			case "career":
				if currentStage != nil {
					if currentStage.Description != "" {
						currentStage.Description += " "
					}
					currentStage.Description += renderNodeText(v, source)
				}
			case "stack":
				if currentGroup != nil && currentGroup.Description == "" {
					currentGroup.Description = renderNodeText(v, source)
				}
			}

			return ast.WalkSkipChildren, nil
		}

		return ast.WalkContinue, nil
	})
	flushStage()
	flushGroup()

	resp.Intro = strings.TrimSpace(introBuf.String())
	return resp
}

func extractStackItem(n ast.Node, source []byte) StackItem {
	item := StackItem{}
	var descBuilder strings.Builder
	foundName := false

	separators := []struct {
		str string
		len int
	}{
		{" — ", 3},
		{" - ", 3},
		{" : ", 3},
		{"—", 1},
	}

	for c := n.FirstChild(); c != nil; c = c.NextSibling() {
		if !foundName {
			if link, ok := c.(*ast.Link); ok {
				item.Name = renderNodeText(link, source)
				item.URL = string(link.Destination)
				foundName = true
				continue
			}
			if emp, ok := c.(*ast.Emphasis); ok && emp.Level == 2 {
				item.Name = renderNodeText(emp, source)
				foundName = true
				continue
			}

			txt := renderNodeText(c, source)
			sepFound := false
			for _, sep := range separators {
				if idx := strings.Index(txt, sep.str); idx > 0 {
					item.Name = strings.TrimSpace(txt[:idx])
					descBuilder.WriteString(txt[idx+sep.len:])
					foundName = true
					sepFound = true
					break
				}
			}

			if !sepFound {
				descBuilder.WriteString(txt)
			}
			continue
		}
		descBuilder.WriteString(renderNodeText(c, source))
	}

	if !foundName {
		item.Name = strings.TrimSpace(descBuilder.String())
		descBuilder.Reset()
	}

	desc := descBuilder.String()
	desc = strings.TrimLeft(desc, " —-:")
	item.Description = strings.TrimSpace(desc)

	return item
}

func extractFact(n ast.Node, source []byte) AboutFact {
	fact := AboutFact{}
	var descBuilder strings.Builder
	foundValue := false

	for c := n.FirstChild(); c != nil; c = c.NextSibling() {
		if !foundValue {
			if emp, ok := c.(*ast.Emphasis); ok && emp.Level == 2 {
				fact.Value = renderNodeText(emp, source)
				foundValue = true
				continue
			}
		}
		descBuilder.WriteString(renderNodeText(c, source))
	}

	desc := descBuilder.String()
	desc = strings.TrimLeft(desc, " —-:")
	fact.Label = strings.TrimSpace(desc)

	return fact
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

func normalizeStageType(s string) string {
	clean := strings.ToLower(s)
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
	text := strings.ToLower(stage.Role + " " + stage.Company)
	for typeName, patterns := range typeRules {
		for _, pattern := range patterns {
			if strings.Contains(text, pattern) {
				return typeName
			}
		}
	}
	return "work"
}
