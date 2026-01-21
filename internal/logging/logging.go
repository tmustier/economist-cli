package logging

import (
	"fmt"
	"os"
	"time"
)

func Debugf(enabled bool, format string, args ...any) {
	if !enabled {
		return
	}
	ts := time.Now().Format(time.RFC3339Nano)
	fmt.Fprintf(os.Stderr, "debug %s "+format+"\n", append([]any{ts}, args...)...)
}
