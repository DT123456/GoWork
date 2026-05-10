package core

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	ZapLog    *zap.Logger
	ZapMySQL  *zap.Logger
	ZapRedis  *zap.Logger
	ZapGin    *zap.Logger
	ZapSystem *zap.Logger
)

func InitZap() error {
	// 从配置获取日志目录
	logDir := AppConfig.Log.File
	if logDir == "" {
		logDir = "log"
	}

	// 确保日志目录存在
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return err
	}

	// 同时输出到文件和控制台
	file, _ := os.OpenFile(logDir, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	writer := zapcore.AddSync(file)

	// 配置编码器
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		CallerKey:      "line",
		MessageKey:     "msg",
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// 同时输出到文件和控制台
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		zapcore.NewMultiWriteSyncer(writer, zapcore.AddSync(os.Stdout)),
		zapcore.DebugLevel,
	)

	// 创建日志记录器
	ZapLog = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
	ZapMySQL = ZapLog.With(zap.String("module", "mysql"))
	ZapRedis = ZapLog.With(zap.String("module", "redis"))
	ZapGin = ZapLog.With(zap.String("module", "gin"))
	ZapSystem = ZapLog.With(zap.String("module", "system"))

	return nil
}
