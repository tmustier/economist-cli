package cmd

import (
	"fmt"
	"os"

	"github.com/charmbracelet/glamour"
	"github.com/spf13/cobra"
	"github.com/tmustier/economist-cli/internal/article"
	"github.com/tmustier/economist-cli/internal/config"
	appErrors "github.com/tmustier/economist-cli/internal/errors"
)

var rawOutput bool

var readCmd = &cobra.Command{
	Use:   "read <url>",
	Short: "Read an article",
	Long: `Fetch and display a full article in the terminal.

Requires login first: economist login

Examples:
  economist read https://www.economist.com/leaders/2026/01/15/some-article
  economist read <url> --raw`,
	Args: cobra.ExactArgs(1),
	RunE: runRead,
}

func init() {
	readCmd.Flags().BoolVar(&rawOutput, "raw", false, "Output raw markdown")
}

func runRead(cmd *cobra.Command, args []string) error {
	url := args[0]

	if !config.IsLoggedIn() {
		fmt.Fprintln(os.Stderr, "⚠️  Not logged in. Run 'economist login' first.")
		fmt.Fprintln(os.Stderr, "   (Articles behind the paywall require authentication)")
		fmt.Fprintln(os.Stderr, "")
		fmt.Fprintln(os.Stderr, "   Attempting to fetch anyway (may hit paywall)...")
		fmt.Fprintln(os.Stderr, "")
	}

	art, err := article.Fetch(url, article.FetchOptions{Debug: debugMode})
	if err != nil {
		if appErrors.IsPaywallError(err) {
			return appErrors.NewUserError("paywall detected - run 'economist login' to read full articles")
		}
		return err
	}

	if art.Content == "" {
		return appErrors.NewUserError("no article content found - try 'economist login'")
	}

	if debugMode && art.DebugHTMLPath != "" {
		fmt.Fprintf(os.Stderr, "Debug HTML saved to: %s\n", art.DebugHTMLPath)
	}

	return outputArticle(art)
}

func outputArticle(art *article.Article) error {
	md := art.ToMarkdown()

	if rawOutput {
		fmt.Println(md)
		return nil
	}

	renderer, err := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(100),
	)
	if err != nil {
		fmt.Println(md)
		return nil
	}

	out, err := renderer.Render(md)
	if err != nil {
		fmt.Println(md)
		return nil
	}

	fmt.Println(out)
	return nil
}
