package ui

import (
	"fmt"
	"strings"

	"github.com/tmustier/economist-cli/internal/article"
)

func RenderArticleHeader(art *article.Article, styles ArticleStyles) string {
	var sb strings.Builder
	if art.Title != "" {
		sb.WriteString(styles.Title.Render(art.Title))
		sb.WriteString("\n")
	}
	if art.Subtitle != "" {
		sb.WriteString(styles.Subtitle.Render(art.Subtitle))
		sb.WriteString("\n")
	}
	if art.DateLine != "" {
		sb.WriteString(styles.Date.Render(art.DateLine))
		sb.WriteString("\n")
	}
	sb.WriteString("\n")
	sb.WriteString(styles.Rule.Render("--------"))
	sb.WriteString("\n\n")
	return sb.String()
}

func ArticleBodyMarkdown(art *article.Article) string {
	var sb strings.Builder
	if art.Content != "" {
		sb.WriteString(art.Content)
	}
	sb.WriteString("\n\n---\n\n")
	sb.WriteString(fmt.Sprintf("🔗 %s\n", art.URL))
	return sb.String()
}
