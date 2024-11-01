package log

import (
	"context"
	"fmt"
	"log/slog"
	"os"
)

// DiscardHandler Since there is no Discard/Noop functionality in log/slog
// overriding existing enabled function to do it.
// https://github.com/golang/go/issues/62005
type DiscardHandler struct {
	slog.Handler
}

func (h *DiscardHandler) Enabled(_ context.Context, level slog.Level) bool {
	return false
}

var logLevels = map[string]slog.Level{
	"DEBUG": slog.LevelDebug,
	"INFO":  slog.LevelInfo,
	"WARN":  slog.LevelWarn,
	"ERROR": slog.LevelError,
}

// InitializeLogger initializes the logger with the specified log level
func InitializeLogger(logLevel string) *slog.Logger {
	ll, ok := logLevels[logLevel]
	if !ok {
		panic(fmt.Sprintf("unknown log level: %s", logLevel))
	}
	return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: ll}))

}
