// Copyright 2024 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cloudlog

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"runtime"

	"cloud.google.com/go/logging"
	"cloud.google.com/go/logging/apiv2/loggingpb"
	"go.opentelemetry.io/otel/trace"
)

// SetGlobal sets the global slog to use output a GCP cloud logging and stderr
func SetGlobalLogger(ctx context.Context, project, logName string, l slog.Level) func(context.Context) error {
	cl, err := logging.NewClient(ctx, fmt.Sprintf("projects/%s", project))
	if err != nil {
		return func(ctx context.Context) error { return nil }
	}
	if err := cl.Ping(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "failed to ping cloud logging, not exporting logs: %v\n", err)
		return func(ctx context.Context) error { return nil }
	}
	h := &cloudLogHandle{
		Handler: slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{}),
		l:       cl.Logger(logName),
		project: project,
		level:   l,
	}

	slog.SetDefault(slog.New(h))
	return func(context.Context) error { return h.l.Flush() }
}

type cloudLogHandle struct {
	slog.Handler
	l       *logging.Logger
	project string
	level   slog.Level
}

func (t *cloudLogHandle) Handle(ctx context.Context, record slog.Record) error {
	payload := map[string]string{
		"message": record.Message,
	}

	entry := logging.Entry{
		Timestamp: record.Time,
	}
	switch record.Level {
	case slog.LevelError:
		entry.Severity = logging.Error
	case slog.LevelWarn:
		entry.Severity = logging.Warning
	case slog.LevelInfo:
		entry.Severity = logging.Info
	case slog.LevelDebug:
		entry.Severity = logging.Debug
	}
	record.Attrs(func(a slog.Attr) bool {
		payload[a.Key] = a.Value.String()
		return true
	})

	entry.Payload = payload

	if s := trace.SpanContextFromContext(ctx); s.IsValid() {
		entry.Trace = fmt.Sprintf("projects/%s/traces/%s", t.project, s.TraceID().String())
		entry.SpanID = s.SpanID().String()
		entry.TraceSampled = s.IsSampled()
	}

	fs := runtime.CallersFrames([]uintptr{record.PC})
	f, _ := fs.Next()
	entry.SourceLocation = &loggingpb.LogEntrySourceLocation{
		File:     f.File,
		Line:     int64(f.Line),
		Function: f.Function,
	}
	if record.Level >= t.level {
		t.l.Log(entry)
	}

	return t.Handler.Handle(ctx, record)
}
