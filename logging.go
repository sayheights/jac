package jac

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type txnLogger interface {
	Log(t *transaction)
}

type noopLogger struct{}

func (n *noopLogger) Log(t *transaction) {
	return
}

type logger struct {
	*zap.Logger
}

func newLogger(host, name string) *logger {
	encConf := zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		NameKey:        "name",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "msg",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.EpochTimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	conf := zap.Config{
		Level:             zap.NewAtomicLevelAt(zap.InfoLevel),
		Development:       false,
		Encoding:          "json",
		EncoderConfig:     encConf,
		OutputPaths:       []string{"stderr"},
		ErrorOutputPaths:  []string{"stderr"},
		DisableStacktrace: true,
	}

	l, err := conf.Build()
	if err != nil {
		panic(err)
	}
	l = l.Named("jac")
	l = l.Named(name)
	l = l.With(
		zap.String("host", host),
	)

	return &logger{l}
}

func (l *logger) Log(t *transaction) {
	switch t.state {
	case txnInitial:
		l.Info("Initialized transaction",
			zap.String("id", t.id),
			zap.String("method", t.req.Method),
			zap.String("uri", t.req.URL.RequestURI()),
			zap.Int("max_retry", t.ret.MaxAmount),
		)
	case txnRetryable:
		l.Warn("Retrying transaction",
			zap.String("id", t.id),
			zap.Int("attempt", t.count),
			zap.Duration("backoff", t.wait),
			zap.Error(t.err),
		)
	case txnResponseReady:
		l.Info("Successful transaction",
			zap.String("id", t.id),
			zap.Int("code", t.res.StatusCode),
			zap.Int("attempt", t.response.AttemptCount),
			zap.Duration("duration", t.response.Duration),
			zap.Int("status_code", t.res.StatusCode),
		)
	case txnExhausted, txnUnrecoverable:
		status := "No response"
		if t.res != nil {
			status = t.res.Status
		}
		l.Error("Failed transaction",
			zap.String("id", t.id),
			zap.Int("attempt", t.count),
			zap.String("status", status),
			zap.Error(t.err),
		)
		return
	default:
	}
}
