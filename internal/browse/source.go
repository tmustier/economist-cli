package browse

import (
	"strings"

	"github.com/tmustier/economist-tui/internal/article"
	"github.com/tmustier/economist-tui/internal/fetch"
	"github.com/tmustier/economist-tui/internal/rss"
)

type DataSource interface {
	Section(section string) (string, []rss.Item, error)
	Article(url string) (*article.Article, error)
}

type rssSource struct {
	debug bool
}

func (s rssSource) Section(section string) (string, []rss.Item, error) {
	feed, err := rss.FetchSection(section)
	if err != nil {
		return "", nil, err
	}
	return strings.TrimSpace(feed.Channel.Title), feed.Channel.Items, nil
}

func (s rssSource) Article(url string) (*article.Article, error) {
	return fetch.FetchArticle(url, fetch.Options{Debug: s.debug})
}
