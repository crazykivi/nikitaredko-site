package handlers

import (
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

// Общий HTML-шаблон для тестов
const testHTMLTemplate = `<html><head>
<title>nikitaredko</title>
<meta name="description" content="Default" />
<meta property="og:title" content="Default" />
<meta property="og:description" content="Default" />
<meta property="og:url" content="Default" />
<meta property="og:type" content="website" />
<meta property="og:image" content="Default" />
</head><body></body></html>`

func TestInjectStaticSEO(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		path           string
		host           string
		meta           [2]string
		mustContain    []string
		mustNotContain []string
	}{
		{
			name: "about page with russian meta",
			path: "/about",
			host: "nikitaredko.ru",
			meta: [2]string{"Обо мне | Nikita Redko", "Кто я такой."},
			mustContain: []string{
				"Обо мне | Nikita Redko",
				"Кто я такой.",
				"http://nikitaredko.ru/about",
			},
			mustNotContain: []string{
				">nikitaredko<",
			},
		},
		{
			name: "uses page",
			path: "/uses",
			host: "example.com",
			meta: [2]string{"Uses | Nikita Redko", "Моё рабочее место."},
			mustContain: []string{
				"Uses | Nikita Redko",
				"Моё рабочее место.",
				"http://example.com/uses",
			},
			mustNotContain: []string{
				">nikitaredko<",
			},
		},
		{
			name: "replaces all old meta values",
			path: "/about",
			host: "test.ru",
			meta: [2]string{"New Title", "New Description"},
			mustContain: []string{
				"<title>New Title</title>",
				`<meta name="description" content="New Description" />`,
				`<meta property="og:title" content="New Title" />`,
				`<meta property="og:description" content="New Description" />`,
				`<meta property="og:url" content="http://test.ru/about" />`,
			},
			mustNotContain: []string{
				">nikitaredko<",
				`name="description" content="Default"`,
				`og:title" content="Default"`,
				`og:description" content="Default"`,
				`og:url" content="Default"`,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest("GET", tt.path, nil)
			c.Request.Host = tt.host

			result := string(injectStaticSEO([]byte(testHTMLTemplate), tt.meta, c))

			for _, substr := range tt.mustContain {
				if !strings.Contains(result, substr) {
					t.Errorf("result should contain %q\n\ngot: %s", substr, result)
				}
			}
			for _, substr := range tt.mustNotContain {
				if strings.Contains(result, substr) {
					t.Errorf("result should NOT contain %q\n\ngot: %s", substr, result)
				}
			}
		})
	}
}

func TestInjectArticleSEO(t *testing.T) {
	gin.SetMode(gin.TestMode)

	longExcerpt := strings.Repeat("а", 300)

	tests := []struct {
		name           string
		article        *Article
		host           string
		forwardedProto string
		mustContain    []string
		mustNotContain []string
	}{
		{
			name: "basic article with http",
			article: &Article{
				ID:          "abc123",
				Title:       "Тестовая статья",
				Excerpt:     "Краткое описание статьи.",
				PublishedAt: "2025-01-01T00:00:00Z",
			},
			host: "nikitaredko.ru",
			mustContain: []string{
				"Тестовая статья | Nikita Redko",
				"Краткое описание статьи.",
				`og:type" content="article"`,
				"application/ld+json",
				"http://nikitaredko.ru/articles/abc123",
			},
			mustNotContain: []string{
				">nikitaredko<",
			},
		},
		{
			name: "article with https via X-Forwarded-Proto",
			article: &Article{
				ID:          "secure-article",
				Title:       "Secure Article",
				Excerpt:     "Description.",
				PublishedAt: "2025-01-01T00:00:00Z",
			},
			host:           "nikitaredko.ru",
			forwardedProto: "https",
			mustContain: []string{
				"https://nikitaredko.ru/articles/secure-article",
				`<meta property="og:url" content="https://nikitaredko.ru/articles/secure-article" />`,
			},
			mustNotContain: []string{
				"http://nikitaredko.ru/articles/secure-article",
			},
		},
		{
			name: "truncates long excerpt to 200 chars",
			article: &Article{
				ID:          "long",
				Title:       "Long Article",
				Excerpt:     longExcerpt,
				PublishedAt: "2025-01-01T00:00:00Z",
			},
			host: "test.ru",
			mustNotContain: []string{
				longExcerpt,
			},
			mustContain: []string{
				strings.Repeat("а", 200) + "...",
			},
		},
		{
			name: "JSON-LD contains article schema",
			article: &Article{
				ID:          "jsonld-test",
				Title:       "JSON-LD Test",
				Excerpt:     "Test excerpt",
				PublishedAt: "2025-06-15T10:00:00Z",
			},
			host: "example.com",
			mustContain: []string{
				`"@type":"Article"`,
				`"headline":"JSON-LD Test"`,
				`"description":"Test excerpt"`,
				`"datePublished":"2025-06-15T10:00:00Z"`,
				`"@type":"Person"`,
				`"name":"Nikita Redko"`,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &ArticleHandler{}

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest("GET", "/articles/"+tt.article.ID, nil)
			c.Request.Host = tt.host
			if tt.forwardedProto != "" {
				c.Request.Header.Set("X-Forwarded-Proto", tt.forwardedProto)
			}

			result := string(h.injectArticleSEO([]byte(testHTMLTemplate), tt.article, c))

			for _, substr := range tt.mustContain {
				if !strings.Contains(result, substr) {
					t.Errorf("result should contain %q\n\ngot: %s", substr, result)
				}
			}
			for _, substr := range tt.mustNotContain {
				if strings.Contains(result, substr) {
					t.Errorf("result should NOT contain %q\n\ngot: %s", substr, result)
				}
			}
		})
	}
}

func TestSEORegexMatchRealIndexHTML(t *testing.T) {
	content, err := os.ReadFile("../../index.html")
	if err != nil {
		t.Skipf("index.html not found at ../../index.html, skipping: %v", err)
	}

	for _, r := range seoRegexes {
		if !r.re.Match(content) {
			t.Errorf("SEO regex %q does not match real index.html — update the regex or the template", r.name)
		}
	}
}
