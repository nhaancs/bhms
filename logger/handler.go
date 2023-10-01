/*
 * NOTICE: The code in this file may be updated in the future.
 *
 * Code taken from https://github.com/ardanlabs/service.
 *
 * Copyright Ardan Labs.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package logger

import (
	"context"

	"log/slog"
)

// logHandler provides a wrapper around the slog handler to capture which
// log level is being logged for event handling.
type logHandler struct {
	handler slog.Handler
	events  Events
}

func newLogHandler(handler slog.Handler, events Events) *logHandler {
	return &logHandler{
		handler: handler,
		events:  events,
	}
}

// Enabled reports whether the handler handles records at the given level.
// The handler ignores records whose level is lower.
func (h *logHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.handler.Enabled(ctx, level)
}

// WithAttrs returns a new JSONHandler whose attributes consists
// of h's attributes followed by attrs.
func (h *logHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &logHandler{handler: h.handler.WithAttrs(attrs), events: h.events}
}

// WithGroup returns a new Handler with the given group appended to the receiver's
// existing groups. The keys of all subsequent attributes, whether added by With
// or in a Record, should be qualified by the sequence of group names.
func (h *logHandler) WithGroup(name string) slog.Handler {
	return &logHandler{handler: h.handler.WithGroup(name), events: h.events}
}

// Handle looks to see if an event function needs to be executed for a given
// log level and then formats its argument Record.
func (h *logHandler) Handle(ctx context.Context, r slog.Record) error {
	switch r.Level {
	case slog.LevelDebug:
		if h.events.Debug != nil {
			h.events.Debug(ctx, toRecord(r))
		}

	case slog.LevelError:
		if h.events.Error != nil {
			h.events.Error(ctx, toRecord(r))
		}

	case slog.LevelWarn:
		if h.events.Warn != nil {
			h.events.Warn(ctx, toRecord(r))
		}

	case slog.LevelInfo:
		if h.events.Info != nil {
			h.events.Info(ctx, toRecord(r))
		}
	}

	return h.handler.Handle(ctx, r)
}
