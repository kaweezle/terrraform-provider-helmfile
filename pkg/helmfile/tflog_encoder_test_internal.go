// Copyright Antoine Martin 2026
// SPDX-License-Identifier: MIT

// cSpell: words stretchr
package helmfile

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestNewTFLogEncoder(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	baseEncoder := zapcore.NewJSONEncoder(zapcore.EncoderConfig{})

	encoder := NewTFLogEncoder(ctx, baseEncoder)

	require.NotNil(t, encoder)
	assert.NotNil(t, encoder.encoder)
	assert.NotNil(t, encoder.ctx)
	assert.NotNil(t, encoder.fields)
}

func TestNewTFLogEncoderWithSubsystem(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	baseEncoder := zapcore.NewJSONEncoder(zapcore.EncoderConfig{})

	encoder := NewTFLogEncoderWithSubsystem(ctx, baseEncoder, "test-subsystem")

	require.NotNil(t, encoder)
	assert.NotNil(t, encoder.encoder)
	assert.NotNil(t, encoder.ctx)
	assert.NotNil(t, encoder.fields)
}

func TestTFLogEncoder_SetField(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	baseEncoder := zapcore.NewJSONEncoder(zapcore.EncoderConfig{})
	encoder := NewTFLogEncoder(ctx, baseEncoder)

	encoder.SetField("key1", "value1")
	encoder.SetField("key2", 42)

	assert.Equal(t, "value1", encoder.fields["key1"])
	assert.Equal(t, 42, encoder.fields["key2"])
}

func TestTFLogEncoder_Clone(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	baseEncoder := zapcore.NewJSONEncoder(zapcore.EncoderConfig{})
	encoder := NewTFLogEncoder(ctx, baseEncoder)

	encoder.SetField("key1", "value1")

	clone := encoder.Clone()
	require.NotNil(t, clone)

	tfEncoder, ok := clone.(*TFLogEncoder)
	require.True(t, ok, "Clone should return a *TFLogEncoder")

	// Fields should be copied
	assert.Equal(t, encoder.fields["key1"], tfEncoder.fields["key1"])

	// Modifying clone should not affect original
	tfEncoder.SetField("key2", "value2")
	assert.Contains(t, tfEncoder.fields, "key2")
	assert.NotContains(t, encoder.fields, "key2")
}

func TestTFLogEncoder_EncodeEntry(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	// Create a base encoder with proper config
	encoderConfig := zapcore.EncoderConfig{
		MessageKey:     "msg",
		LevelKey:       "level",
		TimeKey:        "ts",
		NameKey:        "logger",
		CallerKey:      "caller",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
	baseEncoder := zapcore.NewJSONEncoder(encoderConfig)

	encoder := NewTFLogEncoder(ctx, baseEncoder)

	// Test encoding entries at different levels
	testCases := []struct {
		name  string
		msg   string
		level zapcore.Level
	}{
		{name: "debug", level: zapcore.DebugLevel, msg: "debug message"},
		{name: "info", level: zapcore.InfoLevel, msg: "info message"},
		{name: "warn", level: zapcore.WarnLevel, msg: "warn message"},
		{name: "error", level: zapcore.ErrorLevel, msg: "error message"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			entry := zapcore.Entry{
				Level:   tc.level,
				Message: tc.msg,
			}

			fields := []zapcore.Field{
				zap.String("field1", "value1"),
				zap.Int("field2", 42),
			}

			buf, err := encoder.EncodeEntry(entry, fields)
			require.NoError(t, err)
			require.NotNil(t, buf)
			assert.NotEmpty(t, buf.String())
		})
	}
}

func TestTFLogEncoder_ObjectEncoderMethods(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	baseEncoder := zapcore.NewJSONEncoder(zapcore.EncoderConfig{
		MessageKey: "msg",
	})
	encoder := NewTFLogEncoder(ctx, baseEncoder)

	// Test various Add methods
	t.Run("AddString", func(_ *testing.T) {
		t.Parallel()
		encoder.AddString("key", "value")
	})

	t.Run("AddInt", func(_ *testing.T) {
		t.Parallel()
		encoder.AddInt("key", 42)
	})

	t.Run("AddBool", func(_ *testing.T) {
		t.Parallel()
		encoder.AddBool("key", true)
	})

	t.Run("AddFloat64", func(_ *testing.T) {
		t.Parallel()
		encoder.AddFloat64("key", 3.14)
	})

	t.Run("OpenNamespace", func(_ *testing.T) {
		t.Parallel()
		encoder.OpenNamespace("namespace")
	})
}

func TestTFLogEncoder_InterfaceCompliance(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	baseEncoder := zapcore.NewJSONEncoder(zapcore.EncoderConfig{})
	encoder := NewTFLogEncoder(ctx, baseEncoder)

	// Verify interface compliance at compile time
	var _ zapcore.Encoder = encoder
	var _ zapcore.ObjectEncoder = encoder
}
