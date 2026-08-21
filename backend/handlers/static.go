package handlers

import (
	"encoding/json"
	"html"
	"log"
	"net/http"
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
	seoRegexes = []struct {
		name string
		re   *regexp.Regexp
	}{
		{"title", reTitle},
		{"meta_desc", reMetaDesc},
		{"og_title", reOgTitle},
		{"og_desc", reOgDesc},
		{"og_url", reOgURL},
		{"og_type", reOgType},
		{"og_image", reOgImage},
	}
)

var staticPagesSEO = map[string][2]string{
	"/about": {"Обо мне | Nikita Redko", "Кто я такой: карьерный путь, технологии и факты."},
	"/uses":  {"Uses | Nikita Redko", "Моё рабочее место: железо, софт и инструменты."},
}

func (h *ArticleHandler) ServeFrontend(c *gin.Context) {
	if len(h.indexHTML) == 0 {
		c.Status(http.StatusNotFound)
		return
	}
	body := make([]byte, len(h.indexHTML))
	copy(body, h.indexHTML)

	path := c.Request.URL.Path
	if strings.HasPrefix(path, "/articles/") {
		id := strings.TrimPrefix(path, "/articles/")
		if article, ok := h.GetPublicArticle(id); ok {
			body = h.injectArticleSEO(body, article, c)
		}
	} else if meta, ok := staticPagesSEO[path]; ok {
		body = injectStaticSEO(body, meta, c)
	}

	c.Data(http.StatusOK, "text/html; charset=utf-8", body)
}

func ValidateSEOTemplate(htmlContent []byte) {
	missing := []string{}
	for _, r := range seoRegexes {
		if !r.re.Match(htmlContent) {
			missing = append(missing, r.name)
		}
	}
	if len(missing) > 0 {
		log.Printf("[SEO] WARNING: index.html missing SEO patterns: %v — regex replacement will be silently skipped for these tags", missing)
	} else {
		log.Printf("[SEO] index.html validated: all %d SEO patterns found", len(seoRegexes))
	}
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
	s = reOgImage.ReplaceAllString(s, `<meta property="og:image" content="`+esc(scheme+"://"+host+"/favicon.svg")+`" />`)

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

func injectStaticSEO(htmlBody []byte, meta [2]string, c *gin.Context) []byte {
	scheme := "http"
	if c.Request.TLS != nil || c.GetHeader("X-Forwarded-Proto") == "https" {
		scheme = "https"
	}
	pageURL := scheme + "://" + c.Request.Host + c.Request.URL.Path

	s := string(htmlBody)
	s = reTitle.ReplaceAllString(s, "<title>"+html.EscapeString(meta[0])+"</title>")
	s = reMetaDesc.ReplaceAllString(s, `<meta name="description" content="`+html.EscapeString(meta[1])+`" />`)
	s = reOgTitle.ReplaceAllString(s, `<meta property="og:title" content="`+html.EscapeString(meta[0])+`" />`)
	s = reOgDesc.ReplaceAllString(s, `<meta property="og:description" content="`+html.EscapeString(meta[1])+`" />`)
	s = reOgURL.ReplaceAllString(s, `<meta property="og:url" content="`+html.EscapeString(pageURL)+`" />`)

	return []byte(s)
}
