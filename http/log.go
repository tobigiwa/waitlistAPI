package http

import (
	"context"
	"log/slog"
)

func (a Application) Log(level slog.Level, msg error, opt string) {
	a.logger.LogAttrs(
		context.TODO(),
		level,
		msg.Error(),
		slog.String("operation", opt),
	)
}
