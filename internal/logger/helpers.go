package logger

import (
	"context"
	"time"

	"distributed_calculator/internal/constants"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type ctxKey string

const (
	TraceIDKey       ctxKey = "trace_id"
	RequestIDKey     ctxKey = "request_id"
	CorrelationIDKey ctxKey = "correlation_id"
)

func newEncoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		TimeKey:       constants.LogFieldTimestamp,
		LevelKey:      constants.LogFieldLevel,
		NameKey:       constants.LogFieldLogger,
		CallerKey:     constants.LogFieldCaller,
		FunctionKey:   zapcore.OmitKey,
		MessageKey:    constants.LogFieldMessage,
		StacktraceKey: constants.LogFieldStacktrace,
		LineEnding:    zapcore.DefaultLineEnding,
		EncodeLevel:   zapcore.LowercaseLevelEncoder,
		EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format("02-01-2006 15:04:05"))
		},
		EncodeDuration: zapcore.MillisDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
}

func extractContextFields(ctx context.Context) []zapcore.Field {
	var fields []zapcore.Field

	if traceID := ctx.Value(TraceIDKey); traceID != nil {
		fields = append(fields, zap.String(constants.FieldTraceID, traceID.(string)))
	} else if traceID := ctx.Value(string(TraceIDKey)); traceID != nil {
		fields = append(fields, zap.String(constants.FieldTraceID, traceID.(string)))
	}

	if requestID := ctx.Value(RequestIDKey); requestID != nil {
		fields = append(fields, zap.String(constants.FieldRequestID, requestID.(string)))
	} else if requestID := ctx.Value(string(RequestIDKey)); requestID != nil {
		fields = append(fields, zap.String(constants.FieldRequestID, requestID.(string)))
	}

	if correlationID := ctx.Value(CorrelationIDKey); correlationID != nil {
		fields = append(fields, zap.String(constants.FieldCorrelationID, correlationID.(string)))
	} else if correlationID := ctx.Value(string(CorrelationIDKey)); correlationID != nil {
		fields = append(fields, zap.String(constants.FieldCorrelationID, correlationID.(string)))
	}

	return fields
}
