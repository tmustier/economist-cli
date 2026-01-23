package rss

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
	"time"

	"github.com/tmustier/economist-tui/internal/cache"
)

const rssCacheTTL = 120 * time.Second
const rssCachePrefix = "rss-"

type rssCacheEntry struct {
	CachedAt time.Time `json:"cached_at"`
	Body     []byte    `json:"body"`
}

func loadCachedSection(sectionPath string) ([]byte, time.Time, bool, error) {
	path := rssCachePath(sectionPath)
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, time.Time{}, false, nil
		}
		return nil, time.Time{}, false, err
	}

	var entry rssCacheEntry
	if err := json.Unmarshal(data, &entry); err != nil {
		_ = os.Remove(path)
		return nil, time.Time{}, false, err
	}

	return entry.Body, entry.CachedAt, true, nil
}

func saveCachedSection(sectionPath string, body []byte) error {
	if err := os.MkdirAll(rssCacheDir(), 0755); err != nil {
		return err
	}

	entry := rssCacheEntry{
		CachedAt: time.Now().UTC(),
		Body:     body,
	}

	data, err := json.Marshal(entry)
	if err != nil {
		return err
	}

	path := rssCachePath(sectionPath)
	return os.WriteFile(path, data, 0600)
}

func rssCacheDir() string {
	return cache.CacheDir()
}

func rssCachePath(sectionPath string) string {
	hash := sha1.Sum([]byte(sectionPath))
	name := rssCachePrefix + hex.EncodeToString(hash[:]) + ".json"
	return filepath.Join(rssCacheDir(), name)
}
