package handlers

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

const fullUsesMarkdownContent = `# My Setup

## Hardware

Моё железо.

### MacBook Pro 16"

Основная рабочая машина.

https://apple.com

### Монитор LG 27"

Внешний монитор для работы.

## Software

### VS Code

Редактор кода.

### GoLand

IDE для Go.
`

func TestParseUsesMarkdown(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		expected []UsesCategory
	}{
		{
			name:     "Empty content",
			content:  "",
			expected: []UsesCategory{},
		},
		{
			name:    "Full document parsing",
			content: fullUsesMarkdownContent,
			expected: []UsesCategory{
				{
					ID:          "hardware",
					Title:       "Hardware",
					Description: "Моё железо.",
					Items: []UsesItem{
						{
							Name:        "MacBook Pro 16\"",
							Description: "Основная рабочая машина.",
							URL:         "https://apple.com",
						},
						{
							Name:        "Монитор LG 27\"",
							Description: "Внешний монитор для работы.",
							URL:         "",
						},
					},
				},
				{
					ID:          "software",
					Title:       "Software",
					Description: "",
					Items: []UsesItem{
						{
							Name:        "VS Code",
							Description: "Редактор кода.",
							URL:         "",
						},
						{
							Name:        "GoLand",
							Description: "IDE для Go.",
							URL:         "",
						},
					},
				},
			},
		},
		{
			name: "Category without items",
			content: `# Uses

## Hardware

Just a description without items.
`,
			expected: []UsesCategory{
				{
					ID:          "hardware",
					Title:       "Hardware",
					Description: "Just a description without items.",
					Items:       []UsesItem{},
				},
			},
		},
		{
			name: "Category ID generation",
			content: `# Uses

## Dev Tools

### Terminal
`,
			expected: []UsesCategory{
				{
					ID:          "dev-tools",
					Title:       "Dev Tools",
					Description: "",
					Items: []UsesItem{
						{
							Name:        "Terminal",
							Description: "",
							URL:         "",
						},
					},
				},
			},
		},
	}

	opts := cmpopts.EquateEmpty()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseUsesMarkdown(tt.content)

			if diff := cmp.Diff(tt.expected, got, opts); diff != "" {
				t.Errorf("parseUsesMarkdown() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
