package handlers

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	json "github.com/goccy/go-json"
	"golang.org/x/sync/singleflight"

	"github.com/gin-gonic/gin"

	"nikitaredko-backend/cache"
)

var ErrOutlineNotFound = errors.New("outline resource not found")

type ArticleHandler struct {
	outlineURL            string
	apiKey                string
	allowedCollectionsMap map[string]struct{}
	allowAllCollections   bool
	cache                 *cache.Cache
	httpClient            *http.Client
	sfGroup               singleflight.Group
	indexHTML             []byte
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
	DeletedAt        *string  `json:"deletedAt"`
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

type AttachmentCache struct {
	Body        []byte
	ContentType string
	Headers     map[string]string
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

	resp, err := h.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == http.StatusNotFound {
		return nil, ErrOutlineNotFound
	}

	var outlineResp OutlineResponse
	if err := json.Unmarshal(respBody, &outlineResp); err != nil {
		return nil, err
	}

	if !outlineResp.OK {
		log.Printf("Outline API error: %s", string(respBody))

		if bytes.Contains(respBody, []byte("not_found")) ||
			bytes.Contains(respBody, []byte("Not Found")) {
			return nil, ErrOutlineNotFound
		}

		return nil, fmt.Errorf("outline API returned not ok")
	}

	return outlineResp.Data, nil
}

func (h *ArticleHandler) isHidden(doc OutlineDocument) bool {
	return doc.ArchivedAt != nil || doc.DeletedAt != nil || h.isDraft(doc)
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
		if h.isHidden(doc) {
			continue
		}
		coll, ok := collectionsMap[doc.CollectionID]
		if !ok || !h.isCollectionAllowed(coll.Name) {
			continue
		}
		article := h.mapToArticle(doc, coll.Name, 0)
		articleMap[doc.ID] = &article
	}

	for _, doc := range docs {
		if h.isHidden(doc) {
			continue
		}
		if _, ok := collectionsMap[doc.CollectionID]; !ok {
			continue
		}
		if !h.isCollectionAllowed(collectionsMap[doc.CollectionID].Name) {
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
		if h.isHidden(doc) {
			continue
		}
		if _, ok := collectionsMap[doc.CollectionID]; !ok {
			continue
		}
		if !h.isCollectionAllowed(collectionsMap[doc.CollectionID].Name) {
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
		m, ok := cached.(map[string]OutlineCollection)
		if !ok {
			log.Printf("[Cache] CORRUPTED: %s (type %T), deleting and re-fetching", cacheKey, cached)
			h.cache.Delete(cacheKey)
			return h.fetchCollectionsMap()
		}
		return m, nil
	}

	log.Printf("[Cache] MISS: %s", cacheKey)
	result, err, _ := h.sfGroup.Do("fetch_collections", func() (interface{}, error) {
		if cached, found := h.cache.Get(cacheKey); found {
			log.Printf("[Cache] HIT (double-check): %s", cacheKey)
			return cached, nil
		}

		body := map[string]interface{}{"limit": 100}
		data, err := h.callOutlineAPI("/api/collections.list", body)
		if err != nil {
			return nil, fmt.Errorf("fetch collections from outline: %w", err)
		}

		var collections []OutlineCollection
		if err := json.Unmarshal(data, &collections); err != nil {
			return nil, fmt.Errorf("unmarshal collections: %w", err)
		}

		m := make(map[string]OutlineCollection, len(collections))
		for _, c := range collections {
			m[c.ID] = c
		}

		h.cache.Set(cacheKey, m)
		log.Printf("[Cache] SET: %s (%d items)", cacheKey, len(m))
		return m, nil
	})

	if err != nil {
		return nil, err
	}

	m, ok := result.(map[string]OutlineCollection)
	if !ok {
		return nil, fmt.Errorf("unexpected type from singleflight: %T", result)
	}
	return m, nil
}

func (h *ArticleHandler) fetchAllDocs() ([]OutlineDocument, error) {
	cacheKey := "all_docs_raw"
	if cached, found := h.cache.Get(cacheKey); found {
		log.Printf("[Cache] HIT: %s", cacheKey)
		docs, ok := cached.([]OutlineDocument)
		if !ok {
			log.Printf("[Cache] CORRUPTED: %s (type %T), deleting and re-fetching", cacheKey, cached)
			h.cache.Delete(cacheKey)
			return h.fetchAllDocs()
		}
		return docs, nil
	}

	log.Printf("[Cache] MISS: %s", cacheKey)
	result, err, _ := h.sfGroup.Do("fetch_docs", func() (interface{}, error) {
		if cached, found := h.cache.Get(cacheKey); found {
			log.Printf("[Cache] HIT (double-check): %s", cacheKey)
			return cached, nil
		}

		body := map[string]interface{}{"limit": 100}
		data, err := h.callOutlineAPI("/api/documents.list", body)
		if err != nil {
			return nil, fmt.Errorf("fetch docs from outline: %w", err)
		}

		var docs []OutlineDocument
		if err := json.Unmarshal(data, &docs); err != nil {
			return nil, fmt.Errorf("unmarshal docs: %w", err)
		}

		h.cache.Set(cacheKey, docs)
		log.Printf("[Cache] SET: %s (%d items)", cacheKey, len(docs))
		return docs, nil
	})

	if err != nil {
		return nil, err
	}
	docs, ok := result.([]OutlineDocument)
	if !ok {
		return nil, fmt.Errorf("unexpected type from singleflight: %T", result)
	}
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
		if h.isCollectionAllowed(coll.Name) {
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
		if h.isHidden(doc) {
			continue
		}
		if _, ok := collectionsMap[doc.CollectionID]; !ok {
			continue
		}
		if !h.isCollectionAllowed(collectionsMap[doc.CollectionID].Name) {
			continue
		}
		collectionDocs[doc.CollectionID] = append(collectionDocs[doc.CollectionID], doc)
	}

	result := []CollectionWithArticles{}
	for _, coll := range collectionsMap {
		if !h.isCollectionAllowed(coll.Name) {
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

	body := map[string]interface{}{
		"id": id,
	}
	data, err := h.callOutlineAPI("/api/documents.info", body)
	if err != nil {
		if errors.Is(err, ErrOutlineNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Article not found",
			})
			return
		}
		c.JSON(http.StatusBadGateway, gin.H{
			"error": "Failed to fetch article",
		})
		return
	}

	var doc OutlineDocument
	if err := json.Unmarshal(data, &doc); err != nil {
		c.JSON(http.StatusBadGateway, gin.H{
			"error": "Failed to parse article",
		})
		return
	}
	if doc.ID == "" {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Article not found",
		})
		return
	}

	allowed, collectionName, err := h.canServeDocument(doc)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to verify article",
		})
		return
	}
	if !allowed {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Article not found",
		})
		return
	}

	article := h.mapToArticle(doc, collectionName, 0)

	h.cache.Set(cacheKey, article)
	c.JSON(http.StatusOK, article)
}

func (h *ArticleHandler) canServeDocument(doc OutlineDocument) (bool, string, error) {
	if h.isHidden(doc) {
		return false, "", nil
	}

	collectionsMap, err := h.fetchCollectionsMap()
	if err != nil {
		if !h.allowAllCollections {
			return false, "", err
		}
		return true, "", nil
	}

	coll, ok := collectionsMap[doc.CollectionID]
	if !ok {
		if !h.allowAllCollections {
			return false, "", nil
		}
		return true, "", nil
	}

	if !h.isCollectionAllowed(coll.Name) {
		return false, "", nil
	}

	return true, coll.Name, nil
}

func hasTag(tags []string, target string) bool {
	for _, t := range tags {
		if t == target {
			return true
		}
	}
	return false
}

func (h *ArticleHandler) SearchArticles(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusOK, []Article{})
		return
	}

	allArticles, err := h.getPublicArticles()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch articles"})
		return
	}

	queryLower := strings.ToLower(query)
	var results []Article

	for _, art := range allArticles {

		titleMatch := strings.Contains(strings.ToLower(art.Title), queryLower)
		contentMatch := strings.Contains(strings.ToLower(art.Content), queryLower)

		if titleMatch || contentMatch {
			res := art
			res.Content = ""
			res.Tags = make([]string, len(art.Tags))
			copy(res.Tags, art.Tags)
			if titleMatch {
				res.Tags = append(res.Tags, "title-match")
			}
			results = append(results, res)
		}
	}

	sort.Slice(results, func(i, j int) bool {
		iTitle := hasTag(results[i].Tags, "title-match")
		jTitle := hasTag(results[j].Tags, "title-match")
		if iTitle && !jTitle {
			return true
		}
		if !iTitle && jTitle {
			return false
		}
		return results[i].CreatedAt > results[j].CreatedAt
	})

	c.JSON(http.StatusOK, results)
}

func (h *ArticleHandler) getPublicArticles() ([]Article, error) {
	cacheKey := "public_articles_processed"

	if cached, found := h.cache.Get(cacheKey); found {
		return cached.([]Article), nil
	}

	result, err, _ := h.sfGroup.Do("public_articles", func() (interface{}, error) {
		if cached, found := h.cache.Get(cacheKey); found {
			return cached.([]Article), nil
		}

		collectionsMap, err := h.fetchCollectionsMap()
		if err != nil {
			return nil, err
		}
		docs, err := h.fetchAllDocs()
		if err != nil {
			return nil, err
		}

		articles := make([]Article, 0, len(docs))
		for _, doc := range docs {
			if h.isHidden(doc) {
				continue
			}
			coll, ok := collectionsMap[doc.CollectionID]
			if !ok || !h.isCollectionAllowed(coll.Name) {
				continue
			}
			art := h.mapToArticle(doc, coll.Name, 0)
			articles = append(articles, art)
		}
		sort.Slice(articles, func(i, j int) bool {
			return articles[i].CreatedAt > articles[j].CreatedAt
		})

		h.cache.Set(cacheKey, articles)
		return articles, nil
	})

	if err != nil {
		return nil, err
	}

	return result.([]Article), nil
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
	allArticles, err := h.getPublicArticles()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch articles"})
		return
	}

	var feedArticles []Article
	for _, art := range allArticles {
		if collectionID != "" && art.CollectionID != collectionID {
			continue
		}
		art.Content = ""
		feedArticles = append(feedArticles, art)
	}

	total := len(feedArticles)
	start := (page - 1) * limit
	end := start + limit
	if start > total {
		start = total
	}
	if end > total {
		end = total
	}

	c.JSON(http.StatusOK, FeedResponse{
		Articles: feedArticles[start:end],
		Total:    total,
		Page:     page,
		Limit:    limit,
	})
}

func NewArticleHandler(cacheManager *cache.Cache) *ArticleHandler {
	allowedMap := make(map[string]struct{})
	if collections := os.Getenv("OUTLINE_ALLOWED_COLLECTIONS"); collections != "" {
		for _, c := range strings.Split(collections, ",") {
			trimmed := strings.TrimSpace(c)
			if trimmed != "" {
				allowedMap[strings.ToLower(trimmed)] = struct{}{}
			}
		}
	}

	httpClient := &http.Client{
		Timeout: 15 * time.Second,
		Transport: &http.Transport{
			MaxIdleConns:        100,
			MaxIdleConnsPerHost: 20,
			MaxConnsPerHost:     50,
			IdleConnTimeout:     90 * time.Second,
			DisableCompression:  false,
			TLSHandshakeTimeout: 5 * time.Second,
		},
	}

	indexHTML, err := os.ReadFile("./dist/index.html")
	if err != nil {
		log.Printf("[Static] Warning: index.html not found: %v", err)
		indexHTML = []byte("<html><body>Not built</body></html>")
	} else {
		log.Printf("[Static] index.html cached in memory (%d bytes)", len(indexHTML))
	}

	return &ArticleHandler{
		outlineURL:            os.Getenv("OUTLINE_API_URL"),
		apiKey:                os.Getenv("OUTLINE_API_KEY"),
		allowedCollectionsMap: allowedMap,
		allowAllCollections:   len(allowedMap) == 0,
		cache:                 cacheManager,
		httpClient:            httpClient,
		indexHTML:             indexHTML,
	}
}

func (h *ArticleHandler) isCollectionAllowed(name string) bool {
	if h.allowAllCollections {
		return true
	}
	_, ok := h.allowedCollectionsMap[strings.ToLower(name)]
	return ok
}

func (h *ArticleHandler) ProxyOutlineAttachment(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing id parameter"})
		return
	}

	cacheKey := "attachment_" + id
	if cached, found := h.cache.Get(cacheKey); found {
		data := cached.(*AttachmentCache)
		for k, v := range data.Headers {
			c.Header(k, v)
		}
		c.Data(http.StatusOK, data.ContentType, data.Body)
		return
	}

	url := h.outlineURL + "/api/attachments.redirect?id=" + id
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request"})
		return
	}
	req.Header.Set("Authorization", "Bearer "+h.apiKey)

	resp, err := h.httpClient.Do(req)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "Failed to fetch attachment"})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		c.Status(resp.StatusCode)
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read response"})
		return
	}

	contentType := resp.Header.Get("Content-Type")
	h.cache.Set(cacheKey, &AttachmentCache{
		Body:        body,
		ContentType: contentType,
		Headers: map[string]string{
			"Cache-Control": "public, max-age=3600",
		},
	})

	c.Header("Cache-Control", "public, max-age=3600")
	c.Data(http.StatusOK, contentType, body)
}

func (h *ArticleHandler) GetPublicArticle(id string) (*Article, bool) {
	articles, err := h.getPublicArticles()
	if err != nil {
		return nil, false
	}
	for i := range articles {
		if articles[i].ID == id {
			return &articles[i], true
		}
	}
	return nil, false
}
