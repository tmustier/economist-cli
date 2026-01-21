package article

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/chromedp"
	"github.com/tmustier/economist-cli/internal/browser"
	"github.com/tmustier/economist-cli/internal/config"
	appErrors "github.com/tmustier/economist-cli/internal/errors"
)

// Content length thresholds
const (
	minParagraphLen = 40  // Skip very short paragraphs
	minContentLen   = 500 // Minimum content to consider article "loaded"
)

type Article struct {
	Title         string
	Subtitle      string
	DateLine      string
	Content       string
	URL           string
	DebugHTMLPath string
}

type FetchOptions struct {
	Debug bool
}

func Fetch(articleURL string, opts FetchOptions) (*Article, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	ctx, cancel := browser.HeadlessContext(context.Background(), opts.Debug)
	defer cancel()

	ctx, cancel = context.WithTimeout(ctx, browser.FetchTimeout)
	defer cancel()

	// Inject saved cookies (ignore errors - will just hit paywall)
	_ = browser.InjectCookies(ctx, cfg.Cookies)

	var html string
	err = chromedp.Run(ctx,
		chromedp.Navigate(articleURL),
		chromedp.WaitVisible("body", chromedp.ByQuery),
		chromedp.Sleep(browser.PageLoadWait),
		chromedp.OuterHTML("html", &html),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to load page: %w", err)
	}

	art, parseErr := parseArticle(html, articleURL)
	if opts.Debug {
		if path, err := writeDebugHTML(html); err == nil {
			if art == nil {
				art = &Article{URL: articleURL}
			}
			art.DebugHTMLPath = path
		}
	}
	if parseErr != nil {
		return art, parseErr
	}

	return art, nil
}

func parseArticle(html, articleURL string) (*Article, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return nil, err
	}

	article := &Article{URL: articleURL}
	article.Title = findFirst(doc, "h1.article__headline", "[data-test-id='headline']", "article h1", "h1")
	article.Subtitle = findFirst(doc, ".article__description", "[data-test-id='subheadline']", ".article__subheadline")
	article.DateLine = strings.TrimSpace(doc.Find("time").First().Text())
	article.Content = extractContent(doc)

	if err := checkPaywall(html, article.Content); err != nil {
		return article, err
	}

	return article, nil
}

func findFirst(doc *goquery.Document, selectors ...string) string {
	for _, sel := range selectors {
		if text := strings.TrimSpace(doc.Find(sel).First().Text()); text != "" {
			return text
		}
	}
	return ""
}

func extractContent(doc *goquery.Document) string {
	var paragraphs []string

	// Primary selectors for article body
	doc.Find(".article__body-text p, [data-component='article-body'] p").Each(func(i int, s *goquery.Selection) {
		if isInsideRelatedSection(s) {
			return
		}
		if text := cleanParagraph(s); text != "" {
			paragraphs = append(paragraphs, text)
		}
	})

	// Fallback to broader selectors
	if len(paragraphs) == 0 {
		doc.Find("article p, main p").Each(func(i int, s *goquery.Selection) {
			if text := cleanParagraph(s); text != "" && !looksLikeTeaser(text) {
				paragraphs = append(paragraphs, text)
			}
		})
	}

	return strings.Join(paragraphs, "\n\n")
}

func isInsideRelatedSection(s *goquery.Selection) bool {
	return s.ParentsFiltered("[class*='related'], [class*='teaser'], [class*='promo']").Length() > 0
}

func cleanParagraph(s *goquery.Selection) string {
	text := strings.TrimSpace(s.Text())
	if len(text) < minParagraphLen || isBoilerplate(text) {
		return ""
	}
	return text
}

// looksLikeTeaser detects short promotional text that isn't article content.
func looksLikeTeaser(text string) bool {
	// Real article paragraphs are typically longer
	return len(text) < 80
}

var boilerplatePatterns = []string{
	"subscribe",
	"sign up",
	"newsletter",
	"keep reading",
	"this article appeared",
	"reuse this content",
	"more from",
	"advertisement",
	"listen to this story",
	"enjoy more audio",
}

func isBoilerplate(text string) bool {
	lower := strings.ToLower(text)
	for _, pattern := range boilerplatePatterns {
		if strings.Contains(lower, pattern) {
			return true
		}
	}
	return false
}

var paywallIndicators = []string{
	"Subscribe to read",
	"Keep reading with a subscription",
	"This article is for subscribers",
	"Sign in to continue",
}

func checkPaywall(html, content string) error {
	for _, indicator := range paywallIndicators {
		if strings.Contains(html, indicator) && len(content) < minContentLen {
			return appErrors.PaywallError{}
		}
	}
	return nil
}

func writeDebugHTML(html string) (string, error) {
	file, err := os.CreateTemp("", "economist-article-*.html")
	if err != nil {
		return "", err
	}
	defer file.Close()

	if _, err := file.WriteString(html); err != nil {
		return "", err
	}

	return file.Name(), nil
}

func (a *Article) ToMarkdown() string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("# %s\n\n", a.Title))

	if a.Subtitle != "" {
		sb.WriteString(fmt.Sprintf("*%s*\n\n", a.Subtitle))
	}

	if a.DateLine != "" {
		sb.WriteString(fmt.Sprintf("📅 %s\n\n", a.DateLine))
	}

	sb.WriteString("---\n\n")
	sb.WriteString(a.Content)
	sb.WriteString("\n\n---\n")
	sb.WriteString(fmt.Sprintf("🔗 %s\n", a.URL))

	return sb.String()
}
