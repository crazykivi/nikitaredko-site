package handlers

import (
	"bytes"
	"image/png"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"golang.org/x/image/font"
	"golang.org/x/sync/singleflight"

	"nikitaredko-backend/cache"
)

func sfGroupForTest() singleflight.Group { return singleflight.Group{} }

type fontDrawerHelper struct {
	face font.Face
}

func (h *fontDrawerHelper) width(s string) int {
	d := &font.Drawer{Face: h.face}
	return d.MeasureString(s).Ceil()
}

func newOGTestHandler(t *testing.T) *ArticleHandler {
	t.Helper()
	t.Setenv("CACHE_TTL_MINUTES", "30")
	h := &ArticleHandler{
		cache:   cache.New(),
		sfGroup: sfGroupForTest(),
	}
	h.cache.Set("public_articles_processed", []Article{
		{
			ID:             "og-test-article",
			Title:          "Тестовая статья про генерацию OG-картинок",
			Excerpt:        "Описание",
			PublishedAt:    "2026-03-15T10:00:00Z",
			CollectionName: "Блог",
		},
	})
	return h
}

func TestGenerateOGImageValidPNG(t *testing.T) {
	imgBytes := GenerateOGImage(&Article{
		ID:             "x",
		Title:          "Тестовая статья про генерацию OG-картинок",
		PublishedAt:    "2026-03-15T10:00:00Z",
		CollectionName: "Блог",
	})
	if len(imgBytes) == 0 {
		t.Fatal("generated image is empty")
	}

	img, err := png.Decode(bytes.NewReader(imgBytes))
	if err != nil {
		t.Fatalf("not a valid PNG: %v", err)
	}
	bounds := img.Bounds()
	if bounds.Dx() != ogWidth || bounds.Dy() != ogHeight {
		t.Errorf("expected %dx%d, got %dx%d", ogWidth, ogHeight, bounds.Dx(), bounds.Dy())
	}
}

func TestGenerateOGImageHandlesLongTitle(t *testing.T) {
	imgBytes := GenerateOGImage(&Article{
		Title: "Очень длинный заголовок статьи с множеством слов " +
			"который гарантированно не поместится в четыре строки карточки " +
			"и должен быть аккуратно обрезан с многоточием в конце " +
			"без выхода текста за границы изображения",
	})
	if len(imgBytes) == 0 {
		t.Fatal("generated image is empty")
	}
	if _, err := png.Decode(bytes.NewReader(imgBytes)); err != nil {
		t.Fatalf("not a valid PNG: %v", err)
	}
}

func TestWrapOGTitle(t *testing.T) {
	loadOGFonts()
	if !ogFonts.valid {
		t.Skip("fonts not available")
	}

	tests := []struct {
		name     string
		title    string
		maxLines int
	}{
		{"empty title", "", ogTitleLines},
		{"short title", "Короткий", ogTitleLines},
		{"typical title", "Как я настроил CI/CD для пет-проекта на GitHub Actions", ogTitleLines},
		{"very long title", strings.Repeat("длинноеслово ", 40), ogTitleLines},
		{"single huge word", strings.Repeat("а", 500), ogTitleLines},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lines := wrapOGTitle(ogFonts.bold, tt.title, ogTitleMaxW, tt.maxLines)
			if len(lines) > tt.maxLines {
				t.Errorf("got %d lines, want <= %d", len(lines), tt.maxLines)
			}
			measure := &fontDrawerHelper{face: ogFonts.bold}
			for _, line := range lines {
				if measure.width(line) > ogTitleMaxW {
					t.Errorf("line overflows: %q width=%d max=%d", line, measure.width(line), ogTitleMaxW)
				}
			}
			if tt.title == "" && len(lines) != 0 {
				t.Errorf("empty title should produce no lines, got %v", lines)
			}
		})
	}
}

func TestWrapOGTitleTruncatesWithEllipsis(t *testing.T) {
	loadOGFonts()
	if !ogFonts.valid {
		t.Skip("fonts not available")
	}

	title := strings.Repeat("слово ", 60)
	lines := wrapOGTitle(ogFonts.bold, title, ogTitleMaxW, 3)
	if len(lines) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(lines))
	}
	last := lines[len(lines)-1]
	if !strings.HasSuffix(last, "…") {
		t.Errorf("last line should end with ellipsis, got %q", last)
	}
}

func TestFormatOGDate(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"2026-03-15T10:00:00Z", "15 марта 2026"},
		{"2026-01-01", "1 января 2026"},
		{"", ""},
		{"not-a-date", ""},
	}
	for _, tt := range tests {
		if got := formatOGDate(tt.input); got != tt.want {
			t.Errorf("formatOGDate(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestGetOGImageHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := newOGTestHandler(t)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/api/og/og-test-article", nil)
	c.Params = gin.Params{{Key: "id", Value: "og-test-article"}}

	h.GetOGImage(c)

	if w.Code != 200 {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
	if ct := w.Header().Get("Content-Type"); ct != "image/png" {
		t.Errorf("expected image/png, got %q", ct)
	}
	img, err := png.Decode(bytes.NewReader(w.Body.Bytes()))
	if err != nil {
		t.Fatalf("response body is not a valid PNG: %v", err)
	}
	if img.Bounds().Dx() != ogWidth {
		t.Errorf("unexpected width %d", img.Bounds().Dx())
	}

	w2 := httptest.NewRecorder()
	c2, _ := gin.CreateTestContext(w2)
	c2.Request = httptest.NewRequest("GET", "/api/og/og-test-article", nil)
	c2.Params = gin.Params{{Key: "id", Value: "og-test-article"}}
	h.GetOGImage(c2)
	if w2.Code != 200 {
		t.Fatalf("cached request failed: %d", w2.Code)
	}
	if w2.Body.Len() != w.Body.Len() {
		t.Errorf("cached response differs in size: %d vs %d", w2.Body.Len(), w.Body.Len())
	}
}

func TestGetOGImageHandlerNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := newOGTestHandler(t)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/api/og/unknown-id-123", nil)
	c.Params = gin.Params{{Key: "id", Value: "unknown-id-123"}}

	h.GetOGImage(c)

	if w.Code != 404 {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}

func TestGetOGImageHandlerInvalidID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := newOGTestHandler(t)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/api/og/../../etc", nil)
	c.Params = gin.Params{{Key: "id", Value: "../../etc"}}

	h.GetOGImage(c)

	if w.Code != 400 {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestSaveSampleOGImage(t *testing.T) {
	out := os.Getenv("OG_SAMPLE_OUT")
	if out == "" {
		t.Skip("set OG_SAMPLE_OUT=<path.png> to render a sample card")
	}
	img := GenerateOGImage(&Article{
		ID:             "sample",
		Title:          "Как я настроил CI/CD для пет-проекта на GitHub Actions",
		PublishedAt:    "2026-09-02T12:00:00Z",
		CollectionName: "Блог",
	})
	if err := os.WriteFile(out, img, 0o644); err != nil {
		t.Fatal(err)
	}
	t.Logf("sample written to %s", out)
}
