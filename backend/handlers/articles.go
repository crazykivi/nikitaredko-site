package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
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
	Content        string    `json:"content"`
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

func (h *ArticleHandler) mapToArticle(doc OutlineDocument, collectionName string, level int) Article {
	content := doc.Text
	if content == "" {
		content = doc.Content
	}

	excerpt := ""
	if len(content) > 150 {
		excerpt = content[:150] + "..."
	} else {
		excerpt = content
	}

	publishedAt := doc.CreatedAt
	if doc.PublishedAt != nil {
		publishedAt = *doc.PublishedAt
	}

	wordCount := len([]rune(content))
	readTime := (wordCount + 199) / 200

	return Article{
		ID:             doc.ID,
		Title:          doc.Title,
		Excerpt:        excerpt,
		Content:        content,
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
	return result, nil
}

func (h *ArticleHandler) fetchAllDocs() ([]OutlineDocument, error) {
	body := map[string]interface{}{"limit": 100}
	data, err := h.callOutlineAPI("/api/documents.list", body)
	if err != nil {
		return nil, err
	}

	var docs []OutlineDocument
	if err := json.Unmarshal(data, &docs); err != nil {
		return nil, err
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
