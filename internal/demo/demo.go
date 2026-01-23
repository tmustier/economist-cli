package demo

import (
	"embed"
	"encoding/json"
	"fmt"
	"path"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/tmustier/economist-tui/internal/article"
	"github.com/tmustier/economist-tui/internal/rss"
)

//go:embed fixtures/*.txt fixtures/index.json
var fixturesFS embed.FS

const fixturesIndexPath = "fixtures/index.json"

const DefaultSection = "leaders"

var demoSectionTitles = map[string]string{
	"leaders":                "Leaders",
	"briefing":               "Briefing",
	"finance-and-economics":  "Finance & Economics",
	"united-states":          "United States",
	"britain":                "Britain",
	"europe":                 "Europe",
	"middle-east-and-africa": "Middle East & Africa",
	"asia":                   "Asia",
	"china":                  "China",
	"the-americas":           "The Americas",
	"business":               "Business",
	"science-and-technology": "Science & Technology",
	"culture":                "Culture",
	"graphic-detail":         "Graphic Detail",
	"the-world-this-week":    "The World This Week",
}

var demoBaseDate = time.Date(2026, time.January, 22, 9, 0, 0, 0, time.UTC)

type Source struct {
	sections map[string]sectionData
	articles map[string]*article.Article
	loadErr  error
}

type sectionData struct {
	title string
	items []rss.Item
}

type fixtureSpec struct {
	Slug     string `json:"slug"`
	Title    string `json:"title"`
	Subtitle string `json:"subtitle"`
	Section  string `json:"section"`
	File     string `json:"file,omitempty"`
	Date     string `json:"date,omitempty"`
	Source   string `json:"source"`
}

func NewSource() *Source {
	source := &Source{
		sections: make(map[string]sectionData),
		articles: make(map[string]*article.Article),
	}
	if err := source.addFixtures(); err != nil {
		source.loadErr = err
	}
	return source
}

func (s *Source) Section(section string) (string, []rss.Item, error) {
	if s.loadErr != nil {
		return "", nil, s.loadErr
	}
	key := resolveDemoSection(section)
	if data, ok := s.sections[key]; ok {
		return data.title, data.items, nil
	}
	if data, ok := s.sections[DefaultSection]; ok {
		return data.title, data.items, nil
	}
	return "", nil, fmt.Errorf("demo section not found")
}

func (s *Source) Article(url string) (*article.Article, error) {
	if s.loadErr != nil {
		return nil, s.loadErr
	}
	if art, ok := s.articles[url]; ok {
		copy := *art
		return &copy, nil
	}
	return nil, fmt.Errorf("demo article not found")
}

func (s *Source) addFixtures() error {
	specs, err := loadFixtureSpecs()
	if err != nil {
		return err
	}

	type datedItem struct {
		item      rss.Item
		published time.Time
	}

	sectionItems := make(map[string][]datedItem)
	for i, spec := range specs {
		if spec.Slug == "" || spec.Title == "" {
			return fmt.Errorf("fixture %d missing slug or title", i)
		}
		if strings.TrimSpace(spec.Source) == "" {
			return fmt.Errorf("fixture %d missing source", i)
		}
		if strings.TrimSpace(spec.Section) == "" {
			return fmt.Errorf("fixture %d missing section", i)
		}
		sectionKey := resolveDemoSection(spec.Section)
		if sectionKey == "" {
			return fmt.Errorf("fixture %d has invalid section %q", i, spec.Section)
		}
		published := demoBaseDate.AddDate(0, 0, -i)
		if spec.Date != "" {
			parsed, err := parseFixtureDate(spec.Date)
			if err != nil {
				return err
			}
			published = parsed
		}
		content, err := loadFixtureContent(spec.File)
		if err != nil {
			return err
		}
		content = strings.TrimSpace(content)
		if content == "" {
			content = buildContent(spec.Title)
		}
		url := strings.TrimSpace(spec.Source)
		if !strings.Contains(url, "#") {
			url = fmt.Sprintf("%s#%s", url, spec.Slug)
		}
		item := rss.Item{
			Title:       spec.Title,
			Description: spec.Subtitle,
			Link:        url,
			GUID:        url,
			PubDate:     published.Format(time.RFC1123Z),
		}
		sectionItems[sectionKey] = append(sectionItems[sectionKey], datedItem{
			item:      item,
			published: published,
		})
		s.articles[url] = &article.Article{
			Overtitle: fmt.Sprintf("%s | Demo", demoSectionLabel(sectionKey)),
			Title:     spec.Title,
			Subtitle:  spec.Subtitle,
			DateLine:  formatDateLine(published),
			Content:   content,
			URL:       url,
		}
	}

	if len(sectionItems) == 0 {
		return fmt.Errorf("fixtures index is empty")
	}

	for sectionKey, datedItems := range sectionItems {
		sort.SliceStable(datedItems, func(i, j int) bool {
			return datedItems[i].published.After(datedItems[j].published)
		})

		items := make([]rss.Item, 0, len(datedItems))
		for _, entry := range datedItems {
			items = append(items, entry.item)
		}

		s.sections[sectionKey] = sectionData{title: demoSectionTitle(sectionKey), items: items}
	}

	if _, ok := s.sections[DefaultSection]; !ok {
		return fmt.Errorf("fixtures missing default section %q", DefaultSection)
	}

	return nil
}

func buildContent(title string) string {
	paragraphs := []string{
		"This demo content is stored locally so screenshots and tests can run without network access.",
		fmt.Sprintf("The headline \"%s\" is a placeholder used to show how headlines wrap and how the reader renders long paragraphs.", title),
		"Use ↑/↓ to scroll, b to go back, and c to toggle columns. Resize the terminal to see the layout adapt.",
		"Demo mode keeps everything local so you can explore the TUI without a subscription.",
		"Paragraph lengths are intentionally varied to show line wrapping, spacing, and the feel of the reading experience.",
		"If you are taking screenshots, this page is designed to be safe for public sharing.",
		"End of the sample article ■",
	}

	return strings.Join(paragraphs, "\n\n")
}

func formatDateLine(t time.Time) string {
	day := t.Day()
	suffix := "th"
	if day%100 < 11 || day%100 > 13 {
		switch day % 10 {
		case 1:
			suffix = "st"
		case 2:
			suffix = "nd"
		case 3:
			suffix = "rd"
		}
	}
	return fmt.Sprintf("%s %d%s %d", t.Format("Jan"), day, suffix, t.Year())
}

func demoSectionLabel(section string) string {
	if title, ok := demoSectionTitles[section]; ok {
		return title
	}
	trimmed := strings.TrimSpace(section)
	if trimmed == "" {
		if fallback, ok := demoSectionTitles[DefaultSection]; ok {
			return fallback
		}
		return DefaultSection
	}
	return strings.ReplaceAll(trimmed, "-", " ")
}

func demoSectionTitle(section string) string {
	return fmt.Sprintf("%s (Demo)", demoSectionLabel(section))
}

func resolveDemoSection(section string) string {
	trimmed := strings.TrimSpace(section)
	if trimmed == "" {
		return DefaultSection
	}
	normalized := strings.ToLower(trimmed)
	if resolved, ok := rss.Sections[normalized]; ok {
		return resolved
	}
	return normalized
}

func loadFixtureSpecs() ([]fixtureSpec, error) {
	data, err := fixturesFS.ReadFile(fixturesIndexPath)
	if err != nil {
		return nil, fmt.Errorf("read fixtures index: %w", err)
	}
	var specs []fixtureSpec
	if err := json.Unmarshal(data, &specs); err != nil {
		return nil, fmt.Errorf("parse fixtures index: %w", err)
	}
	if len(specs) == 0 {
		return nil, fmt.Errorf("fixtures index is empty")
	}
	return specs, nil
}

func loadFixtureContent(name string) (string, error) {
	if name == "" {
		return "", nil
	}
	content, err := fixturesFS.ReadFile(path.Join("fixtures", name))
	if err != nil {
		return "", fmt.Errorf("read fixture %s: %w", name, err)
	}
	return string(content), nil
}

var fixtureDateLayouts = []string{
	"2006-01-02",
	"Jan 2 2006",
	"Jan 2, 2006",
	"January 2 2006",
	"January 2, 2006",
}

var fixtureOrdinalSuffix = regexp.MustCompile(`(\d+)(st|nd|rd|th)`)

func parseFixtureDate(value string) (time.Time, error) {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return time.Time{}, fmt.Errorf("fixture date is empty")
	}
	normalized := fixtureOrdinalSuffix.ReplaceAllString(trimmed, "$1")
	normalized = strings.ReplaceAll(normalized, ",", "")
	for _, layout := range fixtureDateLayouts {
		parsed, err := time.Parse(layout, normalized)
		if err == nil {
			return time.Date(parsed.Year(), parsed.Month(), parsed.Day(), 9, 0, 0, 0, time.UTC), nil
		}
	}
	return time.Time{}, fmt.Errorf("parse fixture date %q: unsupported format", value)
}
