package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/tmustier/economist-cli/internal/rss"
)

var (
	headlinesLimit  int
	headlinesSearch string
)

var headlinesCmd = &cobra.Command{
	Use:   "headlines [section]",
	Short: "Show latest headlines from a section",
	Long: `Show latest headlines from The Economist RSS feeds.

Examples:
  economist headlines leaders
  economist headlines finance -n 5
  economist headlines business -s "AI"`,
	Args: cobra.MaximumNArgs(1),
	RunE: runHeadlines,
}

func init() {
	headlinesCmd.Flags().IntVarP(&headlinesLimit, "number", "n", 10, "Number of headlines to show")
	headlinesCmd.Flags().StringVarP(&headlinesSearch, "search", "s", "", "Search headlines for a term")
}

func runHeadlines(cmd *cobra.Command, args []string) error {
	section := "leaders"
	if len(args) > 0 {
		section = args[0]
	}

	items, title, err := fetchHeadlines(section)
	if err != nil {
		return err
	}

	printHeadlines(items, title)
	return nil
}

func fetchHeadlines(section string) ([]rss.Item, string, error) {
	if headlinesSearch != "" {
		items, err := rss.Search(section, headlinesSearch)
		if err != nil {
			return nil, "", err
		}
		title := fmt.Sprintf("🔍 Search results for \"%s\" in %s:", headlinesSearch, section)
		return items, title, nil
	}

	feed, err := rss.FetchSection(section)
	if err != nil {
		return nil, "", err
	}
	title := fmt.Sprintf("📰 %s", strings.TrimSpace(feed.Channel.Title))
	return feed.Channel.Items, title, nil
}

func printHeadlines(items []rss.Item, title string) {
	fmt.Printf("%s\n\n", title)

	if headlinesLimit > 0 && len(items) > headlinesLimit {
		items = items[:headlinesLimit]
	}

	for i, item := range items {
		fmt.Printf("%d. %s\n", i+1, item.CleanTitle())
		if desc := item.CleanDescription(); desc != "" {
			fmt.Printf("   %s\n", desc)
		}
		fmt.Printf("   📅 %s\n", item.FormattedDate())
		fmt.Printf("   🔗 %s\n\n", item.Link)
	}

	if len(items) == 0 {
		fmt.Println("No articles found.")
	}
}
