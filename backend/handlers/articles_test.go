package handlers

import (
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestCalculateReadTime(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		expected int
	}{
		{
			name:     "empty content",
			content:  "",
			expected: 1,
		},
		{
			name:     "1 word",
			content:  "word",
			expected: 1,
		},
		{
			name:     "199 words",
			content:  strings.Repeat("word ", 199),
			expected: 1,
		},
		{
			name:     "200 words = 1 min",
			content:  strings.Repeat("word ", 200),
			expected: 1,
		},
		{
			name:     "201 words = 2 min",
			content:  strings.Repeat("word ", 201),
			expected: 2,
		},
		{
			name:     "1000 words = 5 min",
			content:  strings.Repeat("word ", 1000),
			expected: 5,
		},
		{
			name: "with code blocks (10 lines = 0.5 min + 200 words = 1 min)",
			content: strings.Repeat("word ", 200) + "\n\n" +
				"```go\n" + strings.Repeat("fmt.Println(\"hello\")\n", 10) + "```",
			expected: 2, // ceil(1.0 + 0.5) = 2
		},
		{
			name:     "code only (40 lines = 2 min)",
			content:  "```python\n" + strings.Repeat("print('line')\n", 40) + "```",
			expected: 2,
		},
		{
			name:     "inline code stripped",
			content:  "This is `some inline code here` and regular text",
			expected: 1,
		},
		{
			name:     "markdown syntax stripped",
			content:  "# Heading\n\n**bold** and _italic_ and [link](url) and ![img](url)",
			expected: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := calculateReadTime(tt.content)
			if got != tt.expected {
				t.Errorf("calculateReadTime() = %d, want %d", got, tt.expected)
			}
		})
	}
}

func TestGetExcerpt(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		expected string
		contains []string // если не пустой, проверка наличие подстрок
		excludes []string // если не пустой, проверка отсутствие подстрок
	}{
		{
			name:     "empty content",
			content:  "",
			expected: "",
		},
		{
			name:     "basic text",
			content:  "Первый абзац.\nВторой абзац.",
			expected: "Первый абзац.\nВторой абзац.",
		},
		{
			name:     "stops at H2 heading",
			content:  "Вводный текст.\n\n## Заголовок\n\nТекст после заголовка.",
			expected: "Вводный текст.",
			contains: []string{"Вводный текст"},
			excludes: []string{"Заголовок"},
		},
		{
			name:     "stops at H1 heading",
			content:  "Intro\n# Title\nBody",
			expected: "Intro",
			excludes: []string{"Title"},
		},
		{
			name:     "hash in text not heading",
			content:  "C# is a language",
			expected: "C# is a language",
		},
		{
			name:     "truncates at 200 runes",
			content:  strings.Repeat("а", 250),
			expected: strings.Repeat("а", 200) + "...",
		},
		{
			name:     "only heading fallback",
			content:  "## Just a heading",
			expected: "## Just a heading",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getExcerpt(tt.content)

			if tt.expected != "" {
				if got != tt.expected {
					t.Errorf("getExcerpt() = %q, want %q", got, tt.expected)
				}
			}

			for _, substr := range tt.contains {
				if !strings.Contains(got, substr) {
					t.Errorf("getExcerpt() should contain %q, got %q", substr, got)
				}
			}

			for _, substr := range tt.excludes {
				if strings.Contains(got, substr) {
					t.Errorf("getExcerpt() should NOT contain %q, got %q", substr, got)
				}
			}
		})
	}
}

func TestIsDraft(t *testing.T) {
	now := "2025-01-01T00:00:00Z"

	tests := []struct {
		name     string
		doc      OutlineDocument
		expected bool
	}{
		{
			name:     "no PublishedAt",
			doc:      OutlineDocument{PublishedAt: nil},
			expected: true,
		},
		{
			name: "published with draft tag",
			doc: OutlineDocument{
				PublishedAt: &now,
				Tags:        []string{"draft"},
			},
			expected: true,
		},
		{
			name: "published with WIP tag",
			doc: OutlineDocument{
				PublishedAt: &now,
				Tags:        []string{"WIP"},
			},
			expected: true,
		},
		{
			name: "published normal",
			doc: OutlineDocument{
				PublishedAt: &now,
				Tags:        []string{"go", "backend"},
				Text:        "Normal content",
			},
			expected: false,
		},
		{
			name: "content contains <!-- draft -->",
			doc: OutlineDocument{
				PublishedAt: &now,
				Text:        "Some text <!-- draft --> more text",
			},
			expected: true,
		},
		{
			name: "content contains <!-- wip -->",
			doc: OutlineDocument{
				PublishedAt: &now,
				Text:        "Some text <!-- wip --> more text",
			},
			expected: true,
		},
		{
			name: "content contains [draft]",
			doc: OutlineDocument{
				PublishedAt: &now,
				Text:        "Some text [draft] more text",
			},
			expected: true,
		},
		{
			name: "content contains [wip]",
			doc: OutlineDocument{
				PublishedAt: &now,
				Text:        "Some text [wip] more text",
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &ArticleHandler{}
			got := h.isDraft(tt.doc)
			if got != tt.expected {
				t.Errorf("isDraft() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestIsHidden(t *testing.T) {
	now := "2025-01-01T00:00:00Z"

	tests := []struct {
		name     string
		doc      OutlineDocument
		expected bool
	}{
		{
			name: "archived",
			doc: OutlineDocument{
				PublishedAt: &now,
				ArchivedAt:  &now,
			},
			expected: true,
		},
		{
			name: "deleted",
			doc: OutlineDocument{
				PublishedAt: &now,
				DeletedAt:   &now,
			},
			expected: true,
		},
		{
			name: "draft (not published)",
			doc: OutlineDocument{
				PublishedAt: nil,
			},
			expected: true,
		},
		{
			name: "normal published",
			doc: OutlineDocument{
				PublishedAt: &now,
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &ArticleHandler{}
			got := h.isHidden(tt.doc)
			if got != tt.expected {
				t.Errorf("isHidden() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestBuildArticleTree(t *testing.T) {
	now := "2025-01-01T00:00:00Z"
	parentID := "parent-1"

	tests := []struct {
		name           string
		docs           []OutlineDocument
		collectionsMap map[string]OutlineCollection
		handler        *ArticleHandler
		expected       []Article
	}{
		{
			name: "parent-child relationship",
			docs: []OutlineDocument{
				{
					ID:           parentID,
					Title:        "Parent",
					CollectionID: "coll-1",
					PublishedAt:  &now,
					Text:         "Parent content",
				},
				{
					ID:               "child-1",
					Title:            "Child",
					CollectionID:     "coll-1",
					PublishedAt:      &now,
					ParentDocumentID: &parentID,
					Text:             "Child content",
				},
			},
			collectionsMap: map[string]OutlineCollection{
				"coll-1": {ID: "coll-1", Name: "Blog"},
			},
			handler: &ArticleHandler{allowAllCollections: true},
			expected: []Article{
				{
					ID:             parentID,
					Title:          "Parent",
					CollectionID:   "coll-1",
					CollectionName: "Blog",
					Level:          0,
					Children: []Article{
						{
							ID:             "child-1",
							Title:          "Child",
							CollectionID:   "coll-1",
							CollectionName: "Blog",
							Level:          1,
							Children:       []Article{},
						},
					},
				},
			},
		},
		{
			name: "hides drafts",
			docs: []OutlineDocument{
				{
					ID:           "draft-1",
					Title:        "Draft",
					CollectionID: "coll-1",
					PublishedAt:  nil,
					Text:         "Draft content",
				},
			},
			collectionsMap: map[string]OutlineCollection{
				"coll-1": {ID: "coll-1", Name: "Blog"},
			},
			handler:  &ArticleHandler{allowAllCollections: true},
			expected: []Article{},
		},
	}

	opts := []cmp.Option{
		cmpopts.IgnoreFields(Article{}, "Excerpt", "Content", "CreatedAt", "PublishedAt", "ReadTime", "Tags", "IsDraft"),
		cmpopts.EquateEmpty(),
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.handler.buildArticleTree(tt.docs, tt.collectionsMap)

			if diff := cmp.Diff(tt.expected, got, opts...); diff != "" {
				t.Errorf("buildArticleTree() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestFlattenArticles(t *testing.T) {
	tree := []Article{
		{
			ID:    "1",
			Title: "Root",
			Children: []Article{
				{
					ID:    "2",
					Title: "Child 1",
					Children: []Article{
						{ID: "3", Title: "Grandchild"},
					},
				},
				{ID: "4", Title: "Child 2"},
			},
		},
	}

	flat := flattenArticles(tree)
	if len(flat) != 4 {
		t.Errorf("flattenArticles() = %d items, want 4", len(flat))
	}
}

func TestCountAllArticles(t *testing.T) {
	tree := []Article{
		{
			ID: "1",
			Children: []Article{
				{ID: "2"},
				{ID: "3", Children: []Article{{ID: "4"}}},
			},
		},
	}

	count := countAllArticles(tree)
	if count != 4 {
		t.Errorf("countAllArticles() = %d, want 4", count)
	}
}

func TestStripContent(t *testing.T) {
	articles := []Article{
		{
			ID:      "1",
			Content: "should be removed",
			Children: []Article{
				{ID: "2", Content: "also removed"},
			},
		},
	}

	stripContent(articles)

	if articles[0].Content != "" {
		t.Error("root content should be empty after stripContent")
	}
	if articles[0].Children[0].Content != "" {
		t.Error("child content should be empty after stripContent")
	}
}

func TestHasTag(t *testing.T) {
	tests := []struct {
		name     string
		tags     []string
		target   string
		expected bool
	}{
		{
			name:     "finds existing tag",
			tags:     []string{"go", "backend", "testing"},
			target:   "go",
			expected: true,
		},
		{
			name:     "does not find missing tag",
			tags:     []string{"go", "backend", "testing"},
			target:   "python",
			expected: false,
		},
		{
			name:     "nil tags",
			tags:     nil,
			target:   "go",
			expected: false,
		},
		{
			name:     "empty tags",
			tags:     []string{},
			target:   "go",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := hasTag(tt.tags, tt.target)
			if got != tt.expected {
				t.Errorf("hasTag(%v, %q) = %v, want %v", tt.tags, tt.target, got, tt.expected)
			}
		})
	}
}

func TestIsCollectionAllowed(t *testing.T) {
	tests := []struct {
		name     string
		handler  *ArticleHandler
		collName string
		expected bool
	}{
		{
			name:     "allowAll permits any collection",
			handler:  &ArticleHandler{allowAllCollections: true},
			collName: "Anything",
			expected: true,
		},
		{
			name: "whitelist matches case-insensitive",
			handler: &ArticleHandler{
				allowAllCollections: false,
				allowedCollectionsMap: map[string]struct{}{
					"blog":  {},
					"notes": {},
				},
			},
			collName: "Blog",
			expected: true,
		},
		{
			name: "whitelist exact match",
			handler: &ArticleHandler{
				allowAllCollections: false,
				allowedCollectionsMap: map[string]struct{}{
					"blog":  {},
					"notes": {},
				},
			},
			collName: "notes",
			expected: true,
		},
		{
			name: "whitelist rejects unlisted",
			handler: &ArticleHandler{
				allowAllCollections: false,
				allowedCollectionsMap: map[string]struct{}{
					"blog":  {},
					"notes": {},
				},
			},
			collName: "secret",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.handler.isCollectionAllowed(tt.collName)
			if got != tt.expected {
				t.Errorf("isCollectionAllowed(%q) = %v, want %v", tt.collName, got, tt.expected)
			}
		})
	}
}
