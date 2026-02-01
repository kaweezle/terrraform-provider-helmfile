// Copyright Antoine Martin 2026
// SPDX-License-Identifier: MIT
// cSpell: words syncer wrapcheck containedctx

package helmfile

import (
	"bytes"
	"context"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Check that OutputCapture implements io.Writer.
var _ zapcore.WriteSyncer = (*OutputCapture)(nil)

// OutputCapture captures log output from helmfile operations.
type OutputCapture struct {
	buffer *bytes.Buffer
	mutex  sync.Mutex
}

// NewOutputCapture creates a new output capture.
func NewOutputCapture() *OutputCapture {
	return &OutputCapture{
		buffer: &bytes.Buffer{},
	}
}

// Write implements io.Writer.
func (o *OutputCapture) Write(p []byte) (int, error) {
	o.mutex.Lock()
	defer o.mutex.Unlock()
	//nolint:wrapcheck // we want to return the original error
	return o.buffer.Write(p)
}

// String returns the captured output.
func (o *OutputCapture) String() string {
	o.mutex.Lock()
	defer o.mutex.Unlock()
	return o.buffer.String()
}

// Reset clears the captured output.
func (o *OutputCapture) Reset() {
	o.mutex.Lock()
	defer o.mutex.Unlock()
	o.buffer.Reset()
}

func (o *OutputCapture) Sync() error {
	return nil
}

// CreateCaptureLogger creates a zap logger that captures output.
func CreateCaptureLogger(ctx context.Context, capture *OutputCapture) *zap.SugaredLogger {
	// Create encoder config for plain text output
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	encoder := NewTFLogEncoder(ctx, zapcore.NewJSONEncoder(encoderConfig))

	// Create core that writes to our capture buffer
	core := zapcore.NewCore(
		encoder,
		capture,
		zapcore.DebugLevel, // Capture all levels
	)

	// Create logger
	logger := zap.New(core)
	return logger.Sugar()
}
