package handlers

import (
	"fmt"
	"log"
	"net/http"
	"sort"
	"strings"
	"text/template"
	"time"

	"github.com/gin-gonic/gin"
)

type SitemapURL struct {
	Loc        string
	LastMod    string
	ChangeFreq string
	Priority   string
}

type SitemapData struct {
	URLs []SitemapURL
}

const sitemapTemplate = `<?xml version="1.0" encoding="UTF-8"?>
<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">
{{range .URLs}}
  <url>
    <loc>{{.Loc}}</loc>
    <lastmod>{{.LastMod}}</lastmod>
    <changefreq>{{.ChangeFreq}}</changefreq>
    <priority>{{.Priority}}</priority>
  </url>
{{end}}
</urlset>`

func (h *ArticleHandler) GetSitemap(c *gin.Context) {
	cacheKey := "sitemap_feed"

	if cached, found := h.cache.Get(cacheKey); found {
		if s, ok := cached.(string); ok {
			log.Printf("[Cache] HIT: %s", cacheKey)
			c.Header("Content-Type", "application/xml; charset=utf-8")
			c.String(http.StatusOK, s)
			return
		}
		log.Printf("[Cache] CORRUPTED: %s (type %T), deleting", cacheKey, cached)
		h.cache.Delete(cacheKey)
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

	var urls []SitemapURL

	urls = append(urls, SitemapURL{
		Loc:        baseURL + "/",
		LastMod:    time.Now().UTC().Format("2006-01-02"),
		ChangeFreq: "daily",
		Priority:   "1.0",
	})
	urls = append(urls, SitemapURL{
		Loc:        baseURL + "/articles",
		LastMod:    time.Now().UTC().Format("2006-01-02"),
		ChangeFreq: "daily",
		Priority:   "0.9",
	})
	urls = append(urls, SitemapURL{
		Loc:        baseURL + "/about",
		LastMod:    time.Now().UTC().Format("2006-01-02"),
		ChangeFreq: "monthly",
		Priority:   "0.7",
	})
	urls = append(urls, SitemapURL{
		Loc:        baseURL + "/uses",
		LastMod:    time.Now().UTC().Format("2006-01-02"),
		ChangeFreq: "monthly",
		Priority:   "0.6",
	})

	type articleEntry struct {
		id        string
		createdAt string
	}
	var articles []articleEntry

	for _, doc := range docs {
		if h.isHidden(doc) {
			continue
		}
		coll, ok := collectionsMap[doc.CollectionID]
		if !ok || !h.isCollectionAllowed(coll.Name) {
			continue
		}
		articles = append(articles, articleEntry{
			id:        doc.ID,
			createdAt: doc.CreatedAt,
		})
	}

	sort.Slice(articles, func(i, j int) bool {
		return articles[i].createdAt > articles[j].createdAt
	})

	for _, article := range articles {
		lastMod := parseDateForSitemap(article.createdAt)
		urls = append(urls, SitemapURL{
			Loc:        fmt.Sprintf("%s/articles/%s", baseURL, article.id),
			LastMod:    lastMod,
			ChangeFreq: "weekly",
			Priority:   "0.8",
		})
	}

	data := SitemapData{URLs: urls}

	tmpl, err := template.New("sitemap").Parse(sitemapTemplate)
	if err != nil {
		c.XML(http.StatusInternalServerError, gin.H{"error": "Failed to parse template"})
		return
	}

	var buf strings.Builder
	if err := tmpl.Execute(&buf, data); err != nil {
		c.XML(http.StatusInternalServerError, gin.H{"error": "Failed to execute template"})
		return
	}

	sitemapXML := buf.String()
	h.cache.Set(cacheKey, sitemapXML)
	c.Header("Content-Type", "application/xml; charset=utf-8")
	c.String(http.StatusOK, sitemapXML)
}

func parseDateForSitemap(dateStr string) string {
	t, err := time.Parse(time.RFC3339, dateStr)
	if err != nil {
		return time.Now().UTC().Format("2006-01-02")
	}
	return t.UTC().Format("2006-01-02")
}
