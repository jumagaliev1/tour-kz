package logger

import (
	"context"
	"go.uber.org/zap"
)

type correlationIdType int

const (
	requestIdKey correlationIdType = iota
	sessionIdKey
)

type RequestLogger struct {
	logger zap.SugaredLogger
}

var logger zap.SugaredLogger

func init() {
	zapLogger, _ := zap.NewProduction()
	defer logger.Sync()
	logger = *zapLogger.Sugar()
}

func WithRqId(ctx context.Context, reqID string) context.Context {
	return context.WithValue(ctx, requestIdKey, reqID)
}

func WithSessionId(ctx context.Context, sessionId string) context.Context {
	return context.WithValue(ctx, sessionIdKey, sessionId)
}

func (l *RequestLogger) Logger(ctx context.Context) *zap.SugaredLogger {
	newLogger := &logger
	if ctx != nil {
		if ctxRqId, ok := ctx.Value(requestIdKey).(string); ok {
			newLogger = newLogger.With(zap.String("requestID", ctxRqId))
		}
		if ctxSessionId, ok := ctx.Value(sessionIdKey).(string); ok {
			newLogger = newLogger.With(zap.String("sessionID", ctxSessionId))
		}
	}

	return newLogger
}
