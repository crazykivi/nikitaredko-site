package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"math"
	"net/http"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"nikitaredko-backend/cache"
)

type ArticleHandler struct {
	outlineURL         string
	apiKey             string
	allowedCollections []string
	cache              *cache.Cache
}

func NewArticleHandler(cacheManager *cache.Cache) *ArticleHandler {
	allowedCollections := []string{}
	if collections := os.Getenv("OUTLINE_ALLOWED_COLLECTIONS"); collections != "" {
		for _, c := range strings.Split(collections, ",") {
			trimmed := strings.TrimSpace(c)
			if trimmed != "" {
				allowedCollections = append(allowedCollections, trimmed)
			}
		}
	}

	return &ArticleHandler{
		outlineURL:         os.Getenv("OUTLINE_API_URL"),
		apiKey:             os.Getenv("OUTLINE_API_KEY"),
		allowedCollections: allowedCollections,
		cache:              cacheManager,
	}
}

type OutlineResponse struct {
	OK   bool            `json:"ok"`
	Data json.RawMessage `json:"data"`
}

type OutlineCollection struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Color       string `json:"color"`
	Icon        string `json:"icon"`
}

type OutlineDocument struct {
	ID               string   `json:"id"`
	Title            string   `json:"title"`
	Text             string   `json:"text"`
	Content          string   `json:"content"`
	CreatedAt        string   `json:"createdAt"`
	UpdatedAt        string   `json:"updatedAt"`
	PublishedAt      *string  `json:"publishedAt"`
	Tags             []string `json:"tags"`
	CollectionID     string   `json:"collectionId"`
	ParentDocumentID *string  `json:"parentDocumentId"`
	ArchivedAt       *string  `json:"archivedAt"`
}

type Article struct {
	ID             string    `json:"id"`
	Title          string    `json:"title"`
	Excerpt        string    `json:"excerpt"`
	Content        string    `json:"content,omitempty"`
	CreatedAt      string    `json:"createdAt"`
	PublishedAt    string    `json:"publishedAt"`
	ReadTime       int       `json:"readTime"`
	Tags           []string  `json:"tags"`
	CollectionID   string    `json:"collectionId"`
	CollectionName string    `json:"collectionName"`
	IsDraft        bool      `json:"isDraft"`
	Children       []Article `json:"children,omitempty"`
	Level          int       `json:"level"`
}

type CollectionWithArticles struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Color        string    `json:"color"`
	Icon         string    `json:"icon"`
	Articles     []Article `json:"articles"`
	ArticleCount int       `json:"articleCount"`
}

type FeedResponse struct {
	Articles []Article `json:"articles"`
	Total    int       `json:"total"`
	Page     int       `json:"page"`
	Limit    int       `json:"limit"`
}

func (h *ArticleHandler) callOutlineAPI(endpoint string, body map[string]interface{}) (json.RawMessage, error) {
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	url := h.outlineURL + endpoint
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+h.apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var outlineResp OutlineResponse
	if err := json.Unmarshal(respBody, &outlineResp); err != nil {
		return nil, err
	}

	if !outlineResp.OK {
		log.Printf("Outline API error: %s", string(respBody))
		return nil, err
	}

	return outlineResp.Data, nil
}

func (h *ArticleHandler) isDraft(doc OutlineDocument) bool {
	if doc.PublishedAt == nil {
		return true
	}

	for _, tag := range doc.Tags {
		tagLower := strings.ToLower(tag)
		if tagLower == "draft" || tagLower == "wip" {
			return true
		}
	}

	content := doc.Text
	if content == "" {
		content = doc.Content
	}
	contentLower := strings.ToLower(content)
	if strings.Contains(contentLower, "<!-- draft -->") ||
		strings.Contains(contentLower, "<!-- wip -->") ||
		strings.Contains(contentLower, "[draft]") ||
		strings.Contains(contentLower, "[wip]") {
		return true
	}

	return false
}

func (h *ArticleHandler) isCollectionAllowedByName(name string) bool {
	if len(h.allowedCollections) == 0 {
		return true
	}
	for _, allowed := range h.allowedCollections {
		if strings.EqualFold(allowed, name) {
			return true
		}
	}
	return false
}

var (
	reFencedCode     = regexp.MustCompile("(?s)```.*?```")
	reInlineCode     = regexp.MustCompile("`[^`]*`")
	reMarkdownSyntax = regexp.MustCompile(`[#*_\[\]\(\)!>~\-]`)
)

func calculateReadTime(content string) int {
	codeMatches := reFencedCode.FindAllString(content, -1)
	codeLines := 0
	for _, block := range codeMatches {
		lines := strings.Count(block, "\n")
		if lines > 2 {
			codeLines += lines - 2
		}
	}
	text := reFencedCode.ReplaceAllString(content, "")
	text = reInlineCode.ReplaceAllString(text, "")
	text = reMarkdownSyntax.ReplaceAllString(text, " ")
	words := strings.Fields(text)
	wordCount := len(words)
	textMinutes := float64(wordCount) / 200.0
	codeMinutes := float64(codeLines) / 20.0

	totalMinutes := textMinutes + codeMinutes

	readTime := int(math.Ceil(totalMinutes))
	if readTime < 1 {
		return 1
	}

	return readTime
}

func (h *ArticleHandler) mapToArticle(doc OutlineDocument, collectionName string, level int) Article {
	content := doc.Text
	if content == "" {
		content = doc.Content
	}

	excerpt := getExcerpt(content)

	publishedAt := doc.CreatedAt
	if doc.PublishedAt != nil {
		publishedAt = *doc.PublishedAt
	}

	readTime := calculateReadTime(content)

	return Article{
		ID:             doc.ID,
		Title:          doc.Title,
		Excerpt:        excerpt,
		Content:        content,
		CreatedAt:      doc.CreatedAt,
		PublishedAt:    publishedAt,
		ReadTime:       readTime,
		Tags:           doc.Tags,
		CollectionID:   doc.CollectionID,
		CollectionName: collectionName,
		IsDraft:        h.isDraft(doc),
		Level:          level,
		Children:       []Article{},
	}
}

func getExcerpt(content string) string {
	lines := strings.Split(content, "\n")
	var excerptLines []string

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		isHeading := false

		hashCount := 0
		for _, ch := range trimmed {
			if ch == '#' {
				hashCount++
			} else {
				break
			}
		}

		if hashCount > 0 && hashCount <= 6 {
			if len(trimmed) == hashCount {
				isHeading = true
			} else {
				nextChar := trimmed[hashCount]
				if nextChar == ' ' || nextChar == '\t' {
					isHeading = true
				}
			}
		}

		if isHeading {
			break
		}
		excerptLines = append(excerptLines, line)
	}

	excerpt := strings.TrimSpace(strings.Join(excerptLines, "\n"))

	runes := []rune(excerpt)
	if len(runes) > 0 {
		if len(runes) > 200 {
			return string(runes[:200]) + "..."
		}
		return excerpt
	}

	contentRunes := []rune(content)
	if len(contentRunes) > 200 {
		return string(contentRunes[:200]) + "..."
	}
	return content
}

func (h *ArticleHandler) buildArticleTree(docs []OutlineDocument, collectionsMap map[string]OutlineCollection) []Article {
	articleMap := make(map[string]*Article)
	for _, doc := range docs {
		if doc.ArchivedAt != nil || h.isDraft(doc) {
			continue
		}
		coll, ok := collectionsMap[doc.CollectionID]
		if !ok || !h.isCollectionAllowedByName(coll.Name) {
			continue
		}
		article := h.mapToArticle(doc, coll.Name, 0)
		articleMap[doc.ID] = &article
	}

	for _, doc := range docs {
		if doc.ArchivedAt != nil || h.isDraft(doc) {
			continue
		}
		if _, ok := collectionsMap[doc.CollectionID]; !ok {
			continue
		}
		if !h.isCollectionAllowedByName(collectionsMap[doc.CollectionID].Name) {
			continue
		}

		article, exists := articleMap[doc.ID]
		if !exists {
			continue
		}

		if doc.ParentDocumentID != nil {
			if parent, ok := articleMap[*doc.ParentDocumentID]; ok {
				article.Level = parent.Level + 1
				parent.Children = append(parent.Children, *article)
			}
		}
	}

	rootArticles := make([]Article, 0)
	for _, doc := range docs {
		if doc.ArchivedAt != nil || h.isDraft(doc) {
			continue
		}
		if _, ok := collectionsMap[doc.CollectionID]; !ok {
			continue
		}
		if !h.isCollectionAllowedByName(collectionsMap[doc.CollectionID].Name) {
			continue
		}
		if doc.ParentDocumentID != nil {
			continue
		}
		if article, ok := articleMap[doc.ID]; ok {
			rootArticles = append(rootArticles, *article)
		}
	}

	return rootArticles
}

func countAllArticles(articles []Article) int {
	count := len(articles)
	for _, article := range articles {
		count += countAllArticles(article.Children)
	}
	return count
}

func (h *ArticleHandler) fetchCollectionsMap() (map[string]OutlineCollection, error) {
	cacheKey := "collections_map_raw"
	if cached, found := h.cache.Get(cacheKey); found {
		log.Printf("[Cache] HIT: %s", cacheKey)
		return cached.(map[string]OutlineCollection), nil
	}
	log.Printf("[Cache] MISS: %s", cacheKey)

	body := map[string]interface{}{"limit": 100}
	data, err := h.callOutlineAPI("/api/collections.list", body)
	if err != nil {
		return nil, err
	}
	var collections []OutlineCollection
	if err := json.Unmarshal(data, &collections); err != nil {
		return nil, err
	}
	result := make(map[string]OutlineCollection)
	for _, c := range collections {
		result[c.ID] = c
	}

	h.cache.Set(cacheKey, result)
	return result, nil
}

func (h *ArticleHandler) fetchAllDocs() ([]OutlineDocument, error) {
	cacheKey := "all_docs_raw"
	if cached, found := h.cache.Get(cacheKey); found {
		log.Printf("[Cache] HIT: %s", cacheKey)
		return cached.([]OutlineDocument), nil
	}
	log.Printf("[Cache] MISS: %s", cacheKey)
	body := map[string]interface{}{"limit": 100}
	data, err := h.callOutlineAPI("/api/documents.list", body)
	if err != nil {
		return nil, err
	}

	var docs []OutlineDocument
	if err := json.Unmarshal(data, &docs); err != nil {
		return nil, err
	}

	h.cache.Set(cacheKey, docs)
	return docs, nil
}

func (h *ArticleHandler) ListCollections(c *gin.Context) {
	cacheKey := "collections_list"

	if cached, found := h.cache.Get(cacheKey); found {
		log.Printf("[Cache] HIT: %s", cacheKey)
		c.JSON(http.StatusOK, cached)
		return
	}
	log.Printf("[Cache] MISS: %s", cacheKey)

	collectionsMap, err := h.fetchCollectionsMap()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch collections"})
		return
	}

	filtered := []OutlineCollection{}
	for _, coll := range collectionsMap {
		if h.isCollectionAllowedByName(coll.Name) {
			filtered = append(filtered, coll)
		}
	}

	h.cache.Set(cacheKey, filtered)
	c.JSON(http.StatusOK, filtered)
}

func (h *ArticleHandler) ListArticles(c *gin.Context) {
	cacheKey := "articles_list"

	if cached, found := h.cache.Get(cacheKey); found {
		log.Printf("[Cache] HIT: %s", cacheKey)
		c.JSON(http.StatusOK, cached)
		return
	}
	log.Printf("[Cache] MISS: %s", cacheKey)

	collectionsMap, err := h.fetchCollectionsMap()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch collections"})
		return
	}

	docs, err := h.fetchAllDocs()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch articles"})
		return
	}

	articles := h.buildArticleTree(docs, collectionsMap)
	stripContent(articles)
	h.cache.Set(cacheKey, articles)
	c.JSON(http.StatusOK, articles)
}

func flattenArticles(articles []Article) []Article {
	var result []Article
	for _, article := range articles {
		result = append(result, article)
		if len(article.Children) > 0 {
			result = append(result, flattenArticles(article.Children)...)
		}
	}
	return result
}

func stripContent(articles []Article) {
	for i := range articles {
		articles[i].Content = ""
		if len(articles[i].Children) > 0 {
			stripContent(articles[i].Children)
		}
	}
}

func (h *ArticleHandler) ListArticlesStructured(c *gin.Context) {
	cacheKey := "articles_structured"
	if cached, found := h.cache.Get(cacheKey); found {
		log.Printf("[Cache] HIT: %s", cacheKey)
		c.JSON(http.StatusOK, cached)
		return
	}
	log.Printf("[Cache] MISS: %s", cacheKey)

	collectionsMap, err := h.fetchCollectionsMap()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch collections"})
		return
	}

	docs, err := h.fetchAllDocs()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch articles"})
		return
	}

	collectionDocs := make(map[string][]OutlineDocument)
	for _, doc := range docs {
		if doc.ArchivedAt != nil || h.isDraft(doc) {
			continue
		}
		if _, ok := collectionsMap[doc.CollectionID]; !ok {
			continue
		}
		if !h.isCollectionAllowedByName(collectionsMap[doc.CollectionID].Name) {
			continue
		}
		collectionDocs[doc.CollectionID] = append(collectionDocs[doc.CollectionID], doc)
	}

	result := []CollectionWithArticles{}
	for _, coll := range collectionsMap {
		if !h.isCollectionAllowedByName(coll.Name) {
			continue
		}
		docs := collectionDocs[coll.ID]
		articles := h.buildArticleTree(docs, collectionsMap)
		stripContent(articles)
		allArticles := flattenArticles(articles)
		articleCount := len(allArticles)

		result = append(result, CollectionWithArticles{
			ID:           coll.ID,
			Name:         coll.Name,
			Color:        coll.Color,
			Icon:         coll.Icon,
			Articles:     articles,
			ArticleCount: articleCount,
		})
	}

	h.cache.Set(cacheKey, result)
	c.JSON(http.StatusOK, result)
}

func (h *ArticleHandler) GetArticle(c *gin.Context) {
	id := c.Param("id")
	cacheKey := "article_" + id

	if cached, found := h.cache.Get(cacheKey); found {
		log.Printf("[Cache] HIT: %s", cacheKey)
		c.JSON(http.StatusOK, cached)
		return
	}
	log.Printf("[Cache] MISS: %s", cacheKey)

	body := map[string]interface{}{"id": id}
	data, err := h.callOutlineAPI("/api/documents.info", body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch article"})
		return
	}

	var doc OutlineDocument
	if err := json.Unmarshal(data, &doc); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse article"})
		return
	}

	collectionName := ""
	collectionsMap, err := h.fetchCollectionsMap()
	if err == nil {
		if coll, ok := collectionsMap[doc.CollectionID]; ok {
			collectionName = coll.Name
		}
	}

	article := h.mapToArticle(doc, collectionName, 0)

	h.cache.Set(cacheKey, article)
	c.JSON(http.StatusOK, article)
}

func (h *ArticleHandler) SearchArticles(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusOK, []Article{})
		return
	}

	collectionsMap, err := h.fetchCollectionsMap()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch collections"})
		return
	}
	docs, err := h.fetchAllDocs()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch articles"})
		return
	}

	queryLower := strings.ToLower(query)
	var results []Article

	for _, doc := range docs {
		if doc.ArchivedAt != nil || h.isDraft(doc) {
			continue
		}
		coll, ok := collectionsMap[doc.CollectionID]
		if !ok || !h.isCollectionAllowedByName(coll.Name) {
			continue
		}

		titleMatch := strings.Contains(strings.ToLower(doc.Title), queryLower)
		contentMatch := strings.Contains(strings.ToLower(doc.Text), queryLower) ||
			strings.Contains(strings.ToLower(doc.Content), queryLower)

		if titleMatch || contentMatch {
			article := h.mapToArticle(doc, coll.Name, 0)
			article.Content = ""
			if titleMatch {
				article.Tags = append(article.Tags, "title-match")
			}
			results = append(results, article)
		}
	}

	sort.Slice(results, func(i, j int) bool {
		iTitleMatch := false
		jTitleMatch := false
		for _, tag := range results[i].Tags {
			if tag == "title-match" {
				iTitleMatch = true
				break
			}
		}
		for _, tag := range results[j].Tags {
			if tag == "title-match" {
				jTitleMatch = true
				break
			}
		}
		if iTitleMatch && !jTitleMatch {
			return true
		}
		if !iTitleMatch && jTitleMatch {
			return false
		}
		return results[i].CreatedAt > results[j].CreatedAt
	})

	c.JSON(http.StatusOK, results)
}

func (h *ArticleHandler) GetArticlesFeed(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")
	collectionID := c.Query("collection")

	page, _ := strconv.Atoi(pageStr)
	limit, _ := strconv.Atoi(limitStr)
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	collectionsMap, err := h.fetchCollectionsMap()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch collections"})
		return
	}
	docs, err := h.fetchAllDocs()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch articles"})
		return
	}

	var flatArticles []Article
	for _, doc := range docs {
		if doc.ArchivedAt != nil || h.isDraft(doc) {
			continue
		}
		coll, ok := collectionsMap[doc.CollectionID]
		if !ok || !h.isCollectionAllowedByName(coll.Name) {
			continue
		}
		if collectionID != "" && doc.CollectionID != collectionID {
			continue
		}

		article := h.mapToArticle(doc, coll.Name, 0)
		article.Content = ""
		flatArticles = append(flatArticles, article)
	}

	sort.Slice(flatArticles, func(i, j int) bool {
		return flatArticles[i].CreatedAt > flatArticles[j].CreatedAt
	})

	total := len(flatArticles)

	start := (page - 1) * limit
	end := start + limit
	if start > total {
		start = total
	}
	if end > total {
		end = total
	}

	pagedArticles := flatArticles[start:end]

	c.JSON(http.StatusOK, FeedResponse{
		Articles: pagedArticles,
		Total:    total,
		Page:     page,
		Limit:    limit,
	})
}
