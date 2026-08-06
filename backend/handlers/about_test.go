package handlers

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

// TODO: перенести текст в test_about.md
const fullMarkdownContent = `# About Me

Я разработчик из Бурмалды. Люблю Go и Vue.

## Факты

- **5+ лет** — опыт в разработке
- **12** — проектов в продакшене
- **3** — конференции с докладами

## Карьера

### 2023 — н.в. | Senior Backend Developer | TechCorp | type: work

Описание текущей роли.

- Мигрировал монолит на микросервисы
- Настроил CI/CD пайплайн

### 2020 — 2023 | Junior Developer | StartupInc | type: work

Первый опыт работы.

## Стек

### Backend

Основной стек для серверной части.

- **Go** — основной язык
- [PostgreSQL](https://postgresql.org) — реляционная БД

### Frontend

- **Vue 3** — SPA-фреймворк
`

func TestParseAboutMarkdown(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		expected AboutResponse
	}{
		{
			name:    "Empty Content",
			content: "",
			expected: AboutResponse{
				Intro:  "",
				Facts:  []AboutFact{},
				Career: []CareerStage{},
				Stack:  []StackGroup{},
			},
		},
		{
			name:    "Full Document Parsing",
			content: fullMarkdownContent,
			expected: AboutResponse{
				Intro: "Я разработчик из Бурмалды. Люблю Go и Vue.",
				Facts: []AboutFact{
					{Value: "5+ лет", Label: "опыт в разработке"},
					{Value: "12", Label: "проектов в продакшене"},
					{Value: "3", Label: "конференции с докладами"},
				},
				Career: []CareerStage{
					{
						Period:      "2023 — н.в.",
						Role:        "Senior Backend Developer",
						Company:     "TechCorp",
						Type:        "work",
						Current:     true,
						Highlights:  []string{"Мигрировал монолит на микросервисы", "Настроил CI/CD пайплайн"},
						Description: "Описание текущей роли.",
					},
					{
						Period:      "2020 — 2023",
						Role:        "Junior Developer",
						Company:     "StartupInc",
						Type:        "work",
						Current:     false,
						Highlights:  []string{},
						Description: "Первый опыт работы.",
					},
				},
				Stack: []StackGroup{
					{
						ID:          "backend",
						Title:       "Backend",
						Description: "Основной стек для серверной части.",
						Items: []StackItem{
							{Name: "Go", Description: "основной язык"},
							{Name: "PostgreSQL", URL: "https://postgresql.org", Description: "реляционная БД"},
						},
					},
					{
						ID:          "frontend",
						Title:       "Frontend",
						Description: "",
						Items: []StackItem{
							{Name: "Vue 3", Description: "SPA-фреймворк"},
						},
					},
				},
			},
		},
		{
			name:    "Stack Item Separator Dash",
			content: "## Стек\n### Backend\n- **Go** — основной язык разработки",
			expected: AboutResponse{
				Facts:  []AboutFact{},
				Career: []CareerStage{},
				Stack: []StackGroup{
					{
						ID:          "backend",
						Title:       "Backend",
						Description: "",
						Items: []StackItem{
							{Name: "Go", Description: "основной язык разработки"},
						},
					},
				},
			},
		},
	}

	opts := cmpopts.EquateEmpty()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseAboutMarkdown(tt.content)

			if diff := cmp.Diff(tt.expected, got, opts); diff != "" {
				t.Errorf("parseAboutMarkdown() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestParseCareerHeader(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  *CareerStage
	}{
		{
			name:  "Full header",
			input: "2020 — 2023 | Middle Developer | CompanyX | type: freelance",
			want: &CareerStage{
				Period: "2020 — 2023", Role: "Middle Developer",
				Company: "CompanyX", Type: "freelance", Current: false,
				Highlights: []string{},
			},
		},
		{
			name:  "Minimal header (guesses type)",
			input: "2020 | Developer",
			want: &CareerStage{
				Period: "2020", Role: "Developer",
				Company: "", Type: "work", Current: false,
				Highlights: []string{},
			},
		},
		{
			name:  "Current period (russian 'н.в.')",
			input: "2023 — н.в. | Dev | Co",
			want: &CareerStage{
				Period: "2023 — н.в.", Role: "Dev",
				Company: "Co", Type: "work", Current: true,
				Highlights: []string{},
			},
		},
		{
			name:  "Current period (english 'now')",
			input: "2023 — now | Dev | Co",
			want: &CareerStage{
				Period: "2023 — now", Role: "Dev",
				Company: "Co", Type: "work", Current: true,
				Highlights: []string{},
			},
		},
		{
			name:  "Current period (english 'present')",
			input: "2023 — present | Dev | Co",
			want: &CareerStage{
				Period: "2023 — present", Role: "Dev",
				Company: "Co", Type: "work", Current: true,
				Highlights: []string{},
			},
		},
		{
			name:  "Not current period",
			input: "2020 — 2023 | Dev | Co",
			want: &CareerStage{
				Period: "2020 — 2023", Role: "Dev",
				Company: "Co", Type: "work", Current: false,
				Highlights: []string{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseCareerHeader(tt.input)

			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("parseCareerHeader() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestNormalizeStageType(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"work", "work"},
		{"работа в компании", "work"},
		{"pet project", "pet"},
		{"пет-проект", "pet"},
		{"freelance", "freelance"},
		{"фриланс", "freelance"},
		{"education", "education"},
		{"колледж", "education"},
		{"opensource", "opensource"},
		{"github контрибуции", "opensource"},
		{"community", "community"},
		{"", ""},
		{"unknown_type", "other"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			if got := normalizeStageType(tt.input); got != tt.want {
				t.Errorf("normalizeStageType(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestGuessStageType(t *testing.T) {
	tests := []struct {
		name    string
		role    string
		company string
		want    string
	}{
		{"Corporate work", "Backend Developer", "TechCorp", "work"},
		{"Pet project", "Author", "My Pet Project", "pet"},
		{"Education", "Student", "Университет", "education"},
		{"Open Source", "Contributor", "GitHub OSS", "opensource"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stage := &CareerStage{Role: tt.role, Company: tt.company}
			if got := guessStageType(stage); got != tt.want {
				t.Errorf("guessStageType(%q, %q) = %q, want %q", tt.role, tt.company, got, tt.want)
			}
		})
	}
}
