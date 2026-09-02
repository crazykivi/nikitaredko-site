package handlers

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/gobold"
	"golang.org/x/image/font/gofont/goregular"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"
)

const (
	ogWidth      = 1200
	ogHeight     = 630
	ogPadding    = 80
	ogFontTitle  = 72
	ogFontMeta   = 30
	ogLineHeight = 88
	ogTitleLines = 4
	ogTitleMaxW  = ogWidth - ogPadding*2
	ogCacheKey   = "og_"
)

var ogFonts struct {
	once  sync.Once
	bold  font.Face
	reg   font.Face
	valid bool
}

func loadOGFonts() {
	ogFonts.once.Do(func() {
		parse := func(data []byte, size float64) font.Face {
			f, err := opentype.Parse(data)
			if err != nil {
				log.Printf("[OG] font parse error: %v", err)
				return nil
			}
			face, err := opentype.NewFace(f, &opentype.FaceOptions{
				Size:    size,
				DPI:     72,
				Hinting: font.HintingFull,
			})
			if err != nil {
				log.Printf("[OG] font face error: %v", err)
				return nil
			}
			return face
		}

		ogFonts.bold = parse(gobold.TTF, ogFontTitle)
		ogFonts.reg = parse(goregular.TTF, ogFontMeta)
		ogFonts.valid = ogFonts.bold != nil && ogFonts.reg != nil
		if ogFonts.valid {
			log.Println("[OG] fonts loaded (Go Bold + Go Regular, Cyrillic OK)")
		}
	})
}

func GenerateOGImage(a *Article) []byte {
	loadOGFonts()

	img := image.NewRGBA(image.Rect(0, 0, ogWidth, ogHeight))
	drawOGBackground(img)

	if !ogFonts.valid {
		return encodeOGPNG(img)
	}

	meta := &font.Drawer{Dst: img, Src: image.NewUniform(ogColorMuted), Face: ogFonts.reg}
	meta.Dot = fixed.P(ogPadding+22, 94)
	meta.DrawString("nikitaredko.ru")

	lines := wrapOGTitle(ogFonts.bold, a.Title, ogTitleMaxW, ogTitleLines)
	lineHeight := ogLineHeight
	freeTop := 170
	freeBottom := ogHeight - 130
	blockH := (len(lines) - 1) * lineHeight
	startY := freeTop + (freeBottom-freeTop-blockH)/2
	baseline := startY + ogFontTitle*4/5

	title := &font.Drawer{Dst: img, Src: image.NewUniform(ogColorTitle), Face: ogFonts.bold}
	for i, line := range lines {
		title.Dot = fixed.P(ogPadding, baseline+i*lineHeight)
		title.DrawString(line)
	}

	if len(lines) <= 3 {
		fillRect(img, ogPadding, baseline+blockH+34, 120, 4, ogColorAccent)
	}

	if d := formatOGDate(a.PublishedAt); d != "" {
		meta.Dot = fixed.P(ogPadding, ogHeight-ogPadding+10)
		meta.DrawString(d)
	}
	if a.CollectionName != "" {
		coll := "· " + a.CollectionName
		w := meta.MeasureString(coll)
		meta.Dot = fixed.P(ogWidth-ogPadding-w.Ceil(), ogHeight-ogPadding+10)
		meta.DrawString(coll)
	}

	return encodeOGPNG(img)
}

var ogColorTitle = color.RGBA{0xfa, 0xfa, 0xfa, 0xff}  // zinc-50
var ogColorMuted = color.RGBA{0xa1, 0xa1, 0xaa, 0xff}  // zinc-400
var ogColorAccent = color.RGBA{0x22, 0xd3, 0xee, 0xff} // cyan-400
var ogColorBGFrom = color.RGBA{0x09, 0x09, 0x0b, 0xff} // zinc-950
var ogColorBGTo = color.RGBA{0x1c, 0x2e, 0x44, 0xff}   // zinc-950 → blue tint

func encodeOGPNG(img image.Image) []byte {
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		log.Printf("[OG] png encode error: %v", err)
		return nil
	}
	return buf.Bytes()
}

func drawOGBackground(img *image.RGBA) {
	for y := 0; y < ogHeight; y++ {
		for x := 0; x < ogWidth; x++ {
			t := float64(x+y) / float64(ogWidth+ogHeight)
			img.Set(x, y, lerpOG(ogColorBGFrom, ogColorBGTo, t))
		}
	}

	glow(img, 1050, 80, 260, 0.10) // top-right
	glow(img, 140, 560, 300, 0.08) // bottom-left
	glow(img, 950, 590, 180, 0.05) // bottom-right
}

func glow(img *image.RGBA, cx, cy, r int, alpha float64) {
	c := ogColorAccent
	for y := cy - r; y <= cy+r; y++ {
		if y < 0 || y >= ogHeight {
			continue
		}
		for x := cx - r; x <= cx+r; x++ {
			if x < 0 || x >= ogWidth {
				continue
			}
			dx, dy := float64(x-cx), float64(y-cy)
			d := dx*dx + dy*dy
			if d > float64(r*r) {
				continue
			}
			falloff := 1 - d/float64(r*r) // 1 в центре -> 0 по краям
			a := alpha * falloff
			base := img.RGBAAt(x, y)
			img.SetRGBA(x, y, color.RGBA{
				R: mixCh(base.R, c.R, a),
				G: mixCh(base.G, c.G, a),
				B: mixCh(base.B, c.B, a),
				A: 0xff,
			})
		}
	}
}

func mixCh(base, target uint8, t float64) uint8 {
	v := float64(base) + (float64(target)-float64(base))*t
	if v < 0 {
		return 0
	}
	if v > 255 {
		return 255
	}
	return uint8(v + 0.5)
}

func lerpOG(a, b color.RGBA, t float64) color.RGBA {
	return color.RGBA{
		R: mixCh(a.R, b.R, t),
		G: mixCh(a.G, b.G, t),
		B: mixCh(a.B, b.B, t),
		A: 0xff,
	}
}

func fillRect(img *image.RGBA, x, y, w, h int, c color.RGBA) {
	for j := y; j < y+h; j++ {
		for i := x; i < x+w; i++ {
			if i >= 0 && i < ogWidth && j >= 0 && j < ogHeight {
				img.SetRGBA(i, j, c)
			}
		}
	}
}

var ruMonths = [12]string{
	"января", "февраля", "марта", "апреля", "мая", "июня",
	"июля", "августа", "сентября", "октября", "ноября", "декабря",
}

func formatOGDate(publishedAt string) string {
	for _, layout := range []string{time.RFC3339, "2006-01-02T15:04:05", "2006-01-02"} {
		if t, err := time.Parse(layout, publishedAt); err == nil {
			return strconv.Itoa(t.Day()) + " " + ruMonths[t.Month()-1] + " " + strconv.Itoa(t.Year())
		}
	}
	return ""
}

func wrapOGTitle(face font.Face, title string, maxWidth int, maxLines int) []string {
	measure := &font.Drawer{Face: face}
	width := func(s string) int { return measure.MeasureString(s).Ceil() }
	ellipsis := "…"

	var words []string
	for _, w := range strings.Fields(title) {
		for width(w) > maxWidth && len([]rune(w)) > 1 {
			runes := []rune(w)
			lo, hi := 1, len(runes)
			for lo < hi {
				mid := (lo + hi + 1) / 2
				if width(string(runes[:mid])) <= maxWidth {
					lo = mid
				} else {
					hi = mid - 1
				}
			}
			words = append(words, string(runes[:lo]))
			w = string(runes[lo:])
		}
		words = append(words, w)
	}

	lines := make([]string, 0, maxLines)
	cur := ""
	for i := 0; i < len(words); i++ {
		cand := words[i]
		if cur != "" {
			cand = cur + " " + words[i]
		}
		if width(cand) <= maxWidth {
			cur = cand
			continue
		}
		if cur != "" {
			lines = append(lines, cur)
			if len(lines) == maxLines {
				lines[maxLines-1] = ellipsizeOG(lines[maxLines-1], width, maxWidth, ellipsis)
				return lines
			}
		}
		cur = words[i]
	}
	if cur != "" {
		lines = append(lines, cur)
	}
	return lines
}

func ellipsizeOG(line string, width func(string) int, maxWidth int, ellipsis string) string {
	if width(line+ellipsis) <= maxWidth {
		return line + ellipsis
	}
	runes := []rune(line)
	for len(runes) > 0 {
		cand := string(runes) + ellipsis
		if width(cand) <= maxWidth {
			return cand
		}
		runes = runes[:len(runes)-1]
	}
	return ellipsis
}

func (h *ArticleHandler) GetOGImage(c *gin.Context) {
	id := c.Param("id")
	if !reAttachmentID.MatchString(id) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id format"})
		return
	}

	cacheKey := ogCacheKey + id
	if h.cache != nil {
		if cached, found := h.cache.Get(cacheKey); found {
			if b, ok := cached.([]byte); ok && len(b) > 0 {
				serveOG(c, b)
				return
			}
			h.cache.Delete(cacheKey)
		}
	}

	result, err, _ := h.sfGroup.Do("generate_"+cacheKey, func() (interface{}, error) {
		if h.cache != nil {
			if cached, found := h.cache.Get(cacheKey); found {
				if b, ok := cached.([]byte); ok && len(b) > 0 {
					return b, nil
				}
			}
		}

		article, ok := h.GetPublicArticle(id)
		if !ok {
			return nil, nil
		}

		b := GenerateOGImage(article)
		if len(b) == 0 {
			return nil, nil
		}
		if h.cache != nil {
			h.cache.Set(cacheKey, b)
		}
		return b, nil
	})

	if err != nil || result == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Article not found"})
		return
	}

	b, ok := result.([]byte)
	if !ok || len(b) == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate image"})
		return
	}

	serveOG(c, b)
}

func serveOG(c *gin.Context, b []byte) {
	c.Header("Cache-Control", "public, max-age=3600")
	c.Data(http.StatusOK, "image/png", b)
}
