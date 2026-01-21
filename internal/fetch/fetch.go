package fetch

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/tmustier/economist-cli/internal/article"
	"github.com/tmustier/economist-cli/internal/browser"
	"github.com/tmustier/economist-cli/internal/config"
	"github.com/tmustier/economist-cli/internal/daemon"
	appErrors "github.com/tmustier/economist-cli/internal/errors"
	"github.com/tmustier/economist-cli/internal/logging"
)

type Options struct {
	Debug bool
}

func FetchArticle(url string, opts Options) (*article.Article, error) {
	logging.Debugf(opts.Debug, "read: start url=%s", url)

	art, err := fetchViaDaemon(url, opts)
	if err == nil {
		logging.Debugf(opts.Debug, "read: daemon fetch ok")
		return validateArticle(art)
	}
	if !errors.Is(err, daemon.ErrNotRunning) {
		logging.Debugf(opts.Debug, "read: daemon fetch error: %v", err)
		return nil, err
	}

	logging.Debugf(opts.Debug, "read: daemon unavailable, using local fetch")
	art, err = fetchLocal(url, opts)
	if err != nil {
		return nil, err
	}

	return validateArticle(art)
}

func fetchViaDaemon(url string, opts Options) (*article.Article, error) {
	ctx, cancel := context.WithTimeout(context.Background(), browser.FetchTimeout)
	defer cancel()

	logging.Debugf(opts.Debug, "read: trying daemon fetch")
	start := time.Now()
	art, err := daemon.Fetch(ctx, url, opts.Debug)
	if err == nil {
		logging.Debugf(opts.Debug, "read: daemon response in %s", time.Since(start))
		return art, nil
	}
	if !errors.Is(err, daemon.ErrNotRunning) {
		return nil, normalizeError(err)
	}

	logging.Debugf(opts.Debug, "read: daemon not running, starting background")
	_ = daemon.EnsureBackground()

	readyCtx, readyCancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer readyCancel()
	if daemon.WaitForReady(readyCtx, 200*time.Millisecond) {
		logging.Debugf(opts.Debug, "read: daemon ready, retry fetch")
		art, err = daemon.Fetch(ctx, url, opts.Debug)
		if err == nil {
			logging.Debugf(opts.Debug, "read: daemon response after wait in %s", time.Since(start))
			return art, nil
		}
		return nil, normalizeError(err)
	}

	logging.Debugf(opts.Debug, "read: daemon not ready after wait")
	return nil, daemon.ErrNotRunning
}

func fetchLocal(url string, opts Options) (*article.Article, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	art, err := article.FetchWithCookies(url, article.FetchOptions{Debug: opts.Debug}, cfg.Cookies)
	if err != nil {
		return nil, normalizeError(err)
	}

	return art, nil
}

func validateArticle(art *article.Article) (*article.Article, error) {
	if art.Content == "" {
		return nil, appErrors.NewUserError("no article content found - try 'economist login'")
	}
	return art, nil
}

func normalizeError(err error) error {
	if appErrors.IsPaywallError(err) {
		return appErrors.NewUserError("paywall detected - run 'economist login' to read full articles")
	}
	return err
}
