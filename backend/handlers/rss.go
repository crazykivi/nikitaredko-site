package handlers

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"sort"
	"strings"
	"text/template"
	"time"

	"github.com/gin-gonic/gin"
)

type RSSItem struct {
	Title       string
	Link        string
	Description string
	PubDate     string
	GUID        string
}

type RSSData struct {
	Title       string
	Link        string
	Description string
	LastBuild   string
	Items       []RSSItem
}

const rssTemplate = `<?xml version="1.0" encoding="UTF-8"?>
<rss version="2.0" xmlns:atom="http://www.w3.org/2005/Atom">
  <channel>
    <title>{{.Title}}</title>
    <link>{{.Link}}</link>
    <description>{{.Description}}</description>
    <lastBuildDate>{{.LastBuild}}</lastBuildDate>
    <atom:link href="{{.Link}}/api/rss.xml" rel="self" type="application/rss+xml"/>
    {{range .Items}}
    <item>
      <title><![CDATA[{{.Title}}]]></title>
      <link>{{.Link}}</link>
      <description><![CDATA[{{.Description}}]]></description>
      <pubDate>{{.PubDate}}</pubDate>
      <guid isPermaLink="false">{{.GUID}}</guid>
    </item>
    {{end}}
  </channel>
</rss>`

func (h *ArticleHandler) GetRSS(c *gin.Context) {
	cacheKey := "rss_feed"

	if cached, found := h.cache.Get(cacheKey); found {
		log.Printf("[Cache] HIT: %s", cacheKey)
		c.Header("Content-Type", "application/xml; charset=utf-8")
		c.String(http.StatusOK, cached.(string))
		return
	}

	log.Printf("[Cache] MISS: %s", cacheKey)
	scheme := "http"
	if c.Request.TLS != nil || c.GetHeader("X-Forwarded-Proto") == "https" {
		scheme = "https"
	}
	host := c.Request.Host
	baseURL := fmt.Sprintf("%s://%s", scheme, host)

	collectionsMap, err := h.fetchCollectionsMap()
	if err != nil {
		c.XML(http.StatusInternalServerError, gin.H{"error": "Failed to fetch collections"})
		return
	}

	docs, err := h.fetchAllDocs()
	if err != nil {
		c.XML(http.StatusInternalServerError, gin.H{"error": "Failed to fetch articles"})
		return
	}

	var articles []Article
	for _, doc := range docs {
		if h.isHidden(doc) {
			continue
		}
		coll, ok := collectionsMap[doc.CollectionID]
		if !ok || !h.isCollectionAllowed(coll.Name) {
			continue
		}
		article := h.mapToArticle(doc, coll.Name, 0)
		article.Content = ""
		articles = append(articles, article)
	}

	sort.Slice(articles, func(i, j int) bool {
		return articles[i].CreatedAt > articles[j].CreatedAt
	})

	limit := 20
	if len(articles) > limit {
		articles = articles[:limit]
	}

	var items []RSSItem
	for _, article := range articles {
		link := fmt.Sprintf("%s/articles/%s", baseURL, article.ID)
		items = append(items, RSSItem{
			Title:       article.Title,
			Link:        link,
			Description: article.Excerpt,
			PubDate:     parseDateForRSS(article.CreatedAt),
			GUID:        article.ID,
		})
	}

	siteTitle := os.Getenv("SITE_TITLE")
	if siteTitle == "" {
		siteTitle = "Nikita Redko - Блог"
	}

	siteDesc := os.Getenv("SITE_DESCRIPTION")
	if siteDesc == "" {
		siteDesc = "Блог Никиты Редко: статьи о разработке, проектах, мыслях вслух и всём подряд."
	}

	data := RSSData{
		Title:       siteTitle,
		Link:        baseURL,
		Description: siteDesc,
		LastBuild:   time.Now().UTC().Format(time.RFC1123Z),
		Items:       items,
	}

	tmpl, err := template.New("rss").Parse(rssTemplate)
	if err != nil {
		c.XML(http.StatusInternalServerError, gin.H{"error": "Failed to parse template"})
		return
	}

	var buf strings.Builder
	if err := tmpl.Execute(&buf, data); err != nil {
		c.XML(http.StatusInternalServerError, gin.H{"error": "Failed to execute template"})
		return
	}

	rssXML := buf.String()
	h.cache.Set(cacheKey, rssXML)

	c.Header("Content-Type", "application/xml; charset=utf-8")
	c.String(http.StatusOK, rssXML)
}

func parseDateForRSS(dateStr string) string {
	t, err := time.Parse(time.RFC3339, dateStr)
	if err != nil {
		return time.Now().UTC().Format(time.RFC1123Z)
	}
	return t.UTC().Format(time.RFC1123Z)
}
