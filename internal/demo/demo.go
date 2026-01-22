package demo

import (
	_ "embed"
	"fmt"
	"strings"
	"time"

	"github.com/tmustier/economist-tui/internal/article"
	"github.com/tmustier/economist-tui/internal/rss"
)

//go:embed fixtures/fair-exchange.txt
var fairExchangeFixture string

//go:embed fixtures/german-europe.txt
var germanEuropeFixture string

const DefaultSection = "leaders"

const demoSectionTitle = "Leaders - demo"

var demoBaseDate = time.Date(2026, time.January, 22, 9, 0, 0, 0, time.UTC)
var demoArchiveDate = time.Date(1940, time.September, 7, 9, 0, 0, 0, time.UTC)

type Source struct {
	sections map[string]sectionData
	articles map[string]*article.Article
}

type sectionData struct {
	title string
	items []rss.Item
}

type demoArticle struct {
	slug     string
	title    string
	subtitle string
	content  string
	date     time.Time
}

func NewSource() *Source {
	source := &Source{
		sections: make(map[string]sectionData),
		articles: make(map[string]*article.Article),
	}
	source.addLeaders()
	return source
}

func (s *Source) Section(section string) (string, []rss.Item, error) {
	if section == "" {
		section = DefaultSection
	}
	key := strings.ToLower(section)
	if data, ok := s.sections[key]; ok {
		return data.title, data.items, nil
	}
	if data, ok := s.sections[DefaultSection]; ok {
		return data.title, data.items, nil
	}
	return "", nil, fmt.Errorf("demo section not found")
}

func (s *Source) Article(url string) (*article.Article, error) {
	if art, ok := s.articles[url]; ok {
		copy := *art
		return &copy, nil
	}
	return nil, fmt.Errorf("demo article not found")
}

func (s *Source) addLeaders() {
	articles := []demoArticle{
		{
			slug:     "fair-exchange",
			title:    "Fair Exchange",
			subtitle: "Destroyers for bases, and a new alliance",
			content:  strings.TrimSpace(fairExchangeFixture),
			date:     demoArchiveDate,
		},
		{
			slug:     "german-europe",
			title:    "German Europe",
			subtitle: "Conquest and the limits of domination",
			content:  strings.TrimSpace(germanEuropeFixture),
			date:     demoArchiveDate,
		},
		{slug: "imagined-markets", title: "Imaginary markets and measured optimism", subtitle: "A fictional briefing on sentiment and supply"},
		{slug: "soft-landing", title: "Why the demo economy always lands softly", subtitle: "Illustrative data without real-world stakes"},
		{slug: "office-coffee", title: "The quiet revolution of office coffee", subtitle: "Productivity gains in the imaginary workplace"},
		{slug: "bureaucracy", title: "Small experiments in better bureaucracy", subtitle: "A fictional reform agenda for calmer Mondays"},
	}

	items := make([]rss.Item, 0, len(articles))
	for i, entry := range articles {
		published := demoBaseDate.AddDate(0, 0, -i)
		if !entry.date.IsZero() {
			published = entry.date
		}
		url := fmt.Sprintf("https://example.com/demo/%s", entry.slug)
		items = append(items, rss.Item{
			Title:       entry.title,
			Description: entry.subtitle,
			Link:        url,
			GUID:        url,
			PubDate:     published.Format(time.RFC1123Z),
		})
		content := strings.TrimSpace(entry.content)
		if content == "" {
			content = buildContent(entry.title)
		}
		s.articles[url] = &article.Article{
			Overtitle: "Leaders | Demo",
			Title:     entry.title,
			Subtitle:  entry.subtitle,
			DateLine:  formatDateLine(published),
			Content:   content,
			URL:       url,
		}
	}

	s.sections[DefaultSection] = sectionData{title: demoSectionTitle, items: items}
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
