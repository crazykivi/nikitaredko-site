package handlers

import (
	"encoding/json"
	"html"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
)

var (
	reTitle    = regexp.MustCompile(`<title>[^<]*</title>`)
	reMetaDesc = regexp.MustCompile(`<meta name="description" content="[^"]*"[^>]*/?>`)
	reOgTitle  = regexp.MustCompile(`<meta property="og:title" content="[^"]*"[^>]*/?>`)
	reOgDesc   = regexp.MustCompile(`<meta property="og:description" content="[^"]*"[^>]*/?>`)
	reOgURL    = regexp.MustCompile(`<meta property="og:url" content="[^"]*"[^>]*/?>`)
	reOgType   = regexp.MustCompile(`<meta property="og:type" content="[^"]*"[^>]*/?>`)
	reOgImage  = regexp.MustCompile(`<meta property="og:image" content="[^"]*"[^>]*/?>`)
)

func (h *ArticleHandler) ServeFrontend(c *gin.Context) {
	body, err := os.ReadFile("./dist/index.html")
	if err != nil {
		c.Status(http.StatusNotFound)
		return
	}
	if strings.HasPrefix(c.Request.URL.Path, "/articles/") {
		id := strings.TrimPrefix(c.Request.URL.Path, "/articles/")
		if article, ok := h.GetPublicArticle(id); ok {
			body = h.injectArticleSEO(body, article, c)
		}
	}

	c.Data(http.StatusOK, "text/html; charset=utf-8", body)
}

func (h *ArticleHandler) injectArticleSEO(htmlBody []byte, a *Article, c *gin.Context) []byte {
	scheme := "http"
	if c.Request.TLS != nil || c.GetHeader("X-Forwarded-Proto") == "https" {
		scheme = "https"
	}
	host := c.Request.Host
	pageURL := scheme + "://" + host + "/articles/" + a.ID

	title := a.Title + " | Nikita Redko"
	desc := strings.TrimSpace(strings.ReplaceAll(a.Excerpt, "\n", " "))
	if r := []rune(desc); len(r) > 200 {
		desc = string(r[:200]) + "..."
	}
	esc := html.EscapeString

	s := string(htmlBody)
	s = reTitle.ReplaceAllString(s, "<title>"+esc(title)+"</title>")
	s = reMetaDesc.ReplaceAllString(s, `<meta name="description" content="`+esc(desc)+`" />`)
	s = reOgTitle.ReplaceAllString(s, `<meta property="og:title" content="`+esc(title)+`" />`)
	s = reOgDesc.ReplaceAllString(s, `<meta property="og:description" content="`+esc(desc)+`" />`)
	s = reOgURL.ReplaceAllString(s, `<meta property="og:url" content="`+esc(pageURL)+`" />`)
	s = reOgType.ReplaceAllString(s, `<meta property="og:type" content="article" />`)
	s = reOgImage.ReplaceAllString(s, `<meta property="og:image" content="`+esc(scheme+"://"+host+"/cover.png")+`" />`)

	ld, _ := json.Marshal(map[string]interface{}{
		"@context":         "https://schema.org",
		"@type":            "Article",
		"headline":         a.Title,
		"description":      desc,
		"datePublished":    a.PublishedAt,
		"author":           map[string]string{"@type": "Person", "name": "Nikita Redko"},
		"mainEntityOfPage": pageURL,
	})
	s = strings.Replace(s, "</head>",
		"<script type=\"application/ld+json\">"+string(ld)+"</script>\n</head>", 1)

	return []byte(s)
}
