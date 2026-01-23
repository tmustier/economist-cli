package browse

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/tmustier/economist-tui/internal/app"
	"github.com/tmustier/economist-tui/internal/rss"
	"github.com/tmustier/economist-tui/internal/ui"
)

type Options struct {
	Debug   bool
	NoColor bool
	Source  DataSource
}

func Run(section string, opts Options) error {
	source := opts.Source
	if source == nil {
		source = rssSource{debug: opts.Debug}
	}

	sectionTitle, items, err := loadSection(source, section)
	if err != nil {
		return err
	}

	if _, ok := source.(rssSource); ok {
		go rss.PrefetchAll()
	}

	ui.InitTheme()
	host, err := app.NewHost(app.ScreenBrowse, map[app.ScreenID]app.ScreenBuilder{
		app.ScreenBrowse: func() tea.Model {
			return NewModel(section, items, sectionTitle, opts, source)
		},
	})
	if err != nil {
		return err
	}
	p := tea.NewProgram(host, tea.WithAltScreen())
	_, err = p.Run()
	return err
}

func loadSection(source DataSource, section string) (string, []rss.Item, error) {
	sectionTitle, items, err := source.Section(section)
	if err != nil {
		return "", nil, err
	}

	if len(items) > 50 {
		items = items[:50]
	}

	sectionTitle = strings.TrimSpace(sectionTitle)
	if sectionTitle == "" {
		sectionTitle = section
	}

	return sectionTitle, items, nil
}
