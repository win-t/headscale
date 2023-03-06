package loghandler

import (
	"context"
	"fmt"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/rs/zerolog/log"
)

type logHandlerCtxKeyType struct{}

var LogHandlerCtxKey logHandlerCtxKeyType

type LogHandlerCtx struct {
	ID     int64
	Mode   string
	Method string
	Path   string
	Start  time.Time
	Node   string
}

var logHandlerIDCounter int64

func LogHandler(mode string, h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := &LogHandlerCtx{
			ID:     atomic.AddInt64(&logHandlerIDCounter, 1),
			Method: r.Method,
			Path:   r.URL.EscapedPath(),
			Mode:   mode,
		}

		ctx.begin()
		defer ctx.end()

		h.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), LogHandlerCtxKey, ctx)))
	})
}

func (l *LogHandlerCtx) msg() string {
	return fmt.Sprintf("%s %s\n", l.Method, l.Path)
}

func (l *LogHandlerCtx) begin() {
	l.Start = time.Now().UTC()
	log.Info().
		Str("handler", "LogHandler").
		Str("start", l.Start.Format(time.RFC3339)).
		Int64("id", l.ID).
		Str("method", l.Method).
		Str("path", l.Path).
		Str("mode", l.Mode).
		Msg(l.msg())
}

func (l *LogHandlerCtx) end() {
	end := time.Now()
	elapsed := end.Sub(l.Start).Truncate(time.Second)
	log.Info().
		Str("handler", "LogHandler").
		Str("end", end.Format(time.RFC3339)).
		Str("elapsed", elapsed.String()).
		Int64("id", l.ID).
		Str("method", l.Method).
		Str("path", l.Path).
		Str("mode", l.Mode).
		Str("node", l.Node).
		Msg(l.msg())
}
