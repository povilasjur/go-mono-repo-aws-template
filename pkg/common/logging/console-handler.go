// Package text implements a development-friendly textual handler.
package logging

import (
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	"github.com/apex/log"
)

// Default handler outputting to stderr.
var Default = New(os.Stderr)

// Strings mapping.
var Strings = [...]string{
	log.DebugLevel: "DEBUG",
	log.InfoLevel:  "INFO",
	log.WarnLevel:  "WARN",
	log.ErrorLevel: "ERROR",
	log.FatalLevel: "FATAL",
}

// Handler implementation.
type Handler struct {
	mu     sync.Mutex
	Writer io.Writer
}

// New handler.
func New(w io.Writer) *Handler {
	return &Handler{
		Writer: w,
	}
}

// HandleLog implements log.Handler.
func (h *Handler) HandleLog(e *log.Entry) error {
	level := Strings[e.Level]
	names := e.Fields.Names()

	h.mu.Lock()
	defer h.mu.Unlock()

	// Get the current time and format it
	currentTime := time.Now().Format("2006-01-02 15:04:05.000")

	traceId := e.Fields.Get("trace-id")
	if traceId == nil {
		traceId = ""
	}
	spanId := e.Fields.Get("span-id")
	if spanId == nil {
		spanId = ""
	}
	component := e.Fields.Get("component")
	if component == nil {
		component = ""
	}

	// Print the log message without color formatting
	fmt.Fprintf(h.Writer, "%s %10s [%-16s, %-16s] %-30s: %s",
		currentTime, level, spanId, traceId, component, e.Message)

	for _, name := range names {
		if name == "trace-id" || name == "span-id" || name == "component" {
			continue
		}
		fmt.Fprintf(h.Writer, " %s=%v", name, e.Fields.Get(name))
	}

	fmt.Fprintln(h.Writer)

	return nil
}
