// Copyright Antoine Martin 2026
// SPDX-License-Identifier: MIT
// cSpell: words containedctx wrapcheck

package helmfile

import (
	"context"
	"maps"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"go.uber.org/zap/buffer"
	"go.uber.org/zap/zapcore"
)

// Check that TFLogEncoder implements zapcore.Encoder.
var _ zapcore.Encoder = (*TFLogEncoder)(nil)

// TFLogEncoder wraps another encoder and forwards log entries to tflog.
// It acts as a proxy, delegating encoding to the wrapped encoder while also
// logging entries using Terraform's tflog package.
type TFLogEncoder struct {
	encoder zapcore.Encoder
	ctx     context.Context //nolint:containedctx // context needed for tflog calls
	fields  map[string]any
}

// NewTFLogEncoder creates a new TFLogEncoder that wraps the given encoder.
// The context is used for tflog operations and can be enriched with SetField.
func NewTFLogEncoder(ctx context.Context, encoder zapcore.Encoder) *TFLogEncoder {
	return &TFLogEncoder{
		encoder: encoder,
		ctx:     ctx,
		fields:  make(map[string]any),
	}
}

// NewTFLogEncoderWithSubsystem creates a new TFLogEncoder with a subsystem.
// This is useful for categorizing logs from different components.
func NewTFLogEncoderWithSubsystem(
	ctx context.Context,
	encoder zapcore.Encoder,
	subsystem string,
) *TFLogEncoder {
	return &TFLogEncoder{
		encoder: encoder,
		ctx:     tflog.NewSubsystem(ctx, subsystem),
		fields:  make(map[string]any),
	}
}

// SetField adds a field to the tflog context that will be included in all subsequent logs.
func (e *TFLogEncoder) SetField(key string, value any) {
	e.fields[key] = value
	e.ctx = tflog.SetField(e.ctx, key, value)
}

// Clone creates a copy of the encoder with the same context and fields.
func (e *TFLogEncoder) Clone() zapcore.Encoder {
	clone := &TFLogEncoder{
		encoder: e.encoder.Clone(),
		ctx:     e.ctx,
		fields:  make(map[string]any, len(e.fields)),
	}
	maps.Copy(clone.fields, e.fields)
	return clone
}

// EncodeEntry encodes an entry and fields, forwarding to both the wrapped encoder
// and tflog for Terraform Plugin logging.
//
//nolint:gocritic,wrapcheck // Implements interface method
func (e *TFLogEncoder) EncodeEntry(
	entry zapcore.Entry,
	fields []zapcore.Field,
) (*buffer.Buffer, error) {
	// Forward to the wrapped encoder
	buf, err := e.encoder.EncodeEntry(entry, fields)
	if err != nil {
		return buf, err
	}

	// Also log to tflog with appropriate level
	msg := entry.Message

	// Build fields map for tflog
	tfFields := make(map[string]any, len(e.fields)+len(fields))
	// Include stored fields
	maps.Copy(tfFields, e.fields)

	// Add entry fields
	for _, field := range fields {
		tfFields[field.Key] = field.Interface
	}

	// Log to tflog based on level
	switch entry.Level {
	case zapcore.DebugLevel:
		tflog.Debug(e.ctx, msg, tfFields)
	case zapcore.InfoLevel:
		tflog.Info(e.ctx, msg, tfFields)
	case zapcore.WarnLevel:
		tflog.Warn(e.ctx, msg, tfFields)
	case zapcore.ErrorLevel:
		tflog.Error(e.ctx, msg, tfFields)
	case zapcore.DPanicLevel, zapcore.PanicLevel, zapcore.FatalLevel:
		tflog.Error(e.ctx, msg, tfFields)
	case zapcore.InvalidLevel:
	default:
		tflog.Trace(e.ctx, msg, tfFields)
	}

	return buf, nil
}

// ObjectEncoder interface methods - all delegate to wrapped encoder

// AddArray adds an array-valued field to the encoder.
func (e *TFLogEncoder) AddArray(key string, marshaler zapcore.ArrayMarshaler) error {
	//nolint:wrapcheck // proxying to wrapped encoder
	return e.encoder.AddArray(key, marshaler)
}

// AddObject adds an object-valued field to the encoder.
func (e *TFLogEncoder) AddObject(key string, marshaler zapcore.ObjectMarshaler) error {
	//nolint:wrapcheck // proxying to wrapped encoder
	return e.encoder.AddObject(key, marshaler)
}

// AddBinary adds a binary field to the encoder.
func (e *TFLogEncoder) AddBinary(key string, value []byte) {
	e.encoder.AddBinary(key, value)
}

// AddByteString adds a UTF-8 encoded byte string field to the encoder.
func (e *TFLogEncoder) AddByteString(key string, value []byte) {
	e.encoder.AddByteString(key, value)
}

// AddBool adds a boolean field to the encoder.
func (e *TFLogEncoder) AddBool(key string, value bool) {
	e.encoder.AddBool(key, value)
}

// AddComplex128 adds a complex128 field to the encoder.
func (e *TFLogEncoder) AddComplex128(key string, value complex128) {
	e.encoder.AddComplex128(key, value)
}

// AddComplex64 adds a complex64 field to the encoder.
func (e *TFLogEncoder) AddComplex64(key string, value complex64) {
	e.encoder.AddComplex64(key, value)
}

// AddDuration adds a time.Duration field to the encoder.
func (e *TFLogEncoder) AddDuration(key string, value time.Duration) {
	e.encoder.AddDuration(key, value)
}

// AddFloat64 adds a float64 field to the encoder.
func (e *TFLogEncoder) AddFloat64(key string, value float64) {
	e.encoder.AddFloat64(key, value)
}

// AddFloat32 adds a float32 field to the encoder.
func (e *TFLogEncoder) AddFloat32(key string, value float32) {
	e.encoder.AddFloat32(key, value)
}

// AddInt adds an int field to the encoder.
func (e *TFLogEncoder) AddInt(key string, value int) {
	e.encoder.AddInt(key, value)
}

// AddInt64 adds an int64 field to the encoder.
func (e *TFLogEncoder) AddInt64(key string, value int64) {
	e.encoder.AddInt64(key, value)
}

// AddInt32 adds an int32 field to the encoder.
func (e *TFLogEncoder) AddInt32(key string, value int32) {
	e.encoder.AddInt32(key, value)
}

// AddInt16 adds an int16 field to the encoder.
func (e *TFLogEncoder) AddInt16(key string, value int16) {
	e.encoder.AddInt16(key, value)
}

// AddInt8 adds an int8 field to the encoder.
func (e *TFLogEncoder) AddInt8(key string, value int8) {
	e.encoder.AddInt8(key, value)
}

// AddString adds a string field to the encoder.
func (e *TFLogEncoder) AddString(key, value string) {
	e.encoder.AddString(key, value)
}

// AddTime adds a time.Time field to the encoder.
func (e *TFLogEncoder) AddTime(key string, value time.Time) {
	e.encoder.AddTime(key, value)
}

// AddUint adds a uint field to the encoder.
func (e *TFLogEncoder) AddUint(key string, value uint) {
	e.encoder.AddUint(key, value)
}

// AddUint64 adds a uint64 field to the encoder.
func (e *TFLogEncoder) AddUint64(key string, value uint64) {
	e.encoder.AddUint64(key, value)
}

// AddUint32 adds a uint32 field to the encoder.
func (e *TFLogEncoder) AddUint32(key string, value uint32) {
	e.encoder.AddUint32(key, value)
}

// AddUint16 adds a uint16 field to the encoder.
func (e *TFLogEncoder) AddUint16(key string, value uint16) {
	e.encoder.AddUint16(key, value)
}

// AddUint8 adds a uint8 field to the encoder.
func (e *TFLogEncoder) AddUint8(key string, value uint8) {
	e.encoder.AddUint8(key, value)
}

// AddUintptr adds a uintptr field to the encoder.
func (e *TFLogEncoder) AddUintptr(key string, value uintptr) {
	e.encoder.AddUintptr(key, value)
}

// AddReflected uses reflection to serialize arbitrary objects.
func (e *TFLogEncoder) AddReflected(key string, value interface{}) error {
	//nolint:wrapcheck // proxying to wrapped encoder
	return e.encoder.AddReflected(key, value)
}

// OpenNamespace opens an isolated namespace for subsequent fields.
func (e *TFLogEncoder) OpenNamespace(key string) {
	e.encoder.OpenNamespace(key)
}
