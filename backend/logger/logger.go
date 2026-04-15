package logger

import (
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime/debug"
	"strings"
	"time"
	"web/setting"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

func Init(cfg *setting.LogConfig, mode string) (err error) {
	writeSyncer := getLogWrite(
		cfg.Filename,
		cfg.MaxSize,
		cfg.MaxBackups,
		cfg.MaxAge,
	)
	encoder := getEncoder()

	var level = new(zapcore.Level)
	err = level.UnmarshalText([]byte(viper.GetString("log.level")))
	if err != nil {
		return
	}

	var core zapcore.Core
	if mode == "dev" {
		// 开发模式，日志输出到终端
		consoleEnbcoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
		core = zapcore.NewTee(
			zapcore.NewCore(consoleEnbcoder, zapcore.Lock(os.Stdout), zapcore.DebugLevel),
			zapcore.NewCore(encoder, writeSyncer, level),
		)
	} else {
		// 日志输出到文件
		core = zapcore.NewCore(encoder, writeSyncer, level)
	}

	logger := zap.New(core, zap.AddCaller()) // AddCaller 为日志增加文件位置信息

	zap.ReplaceGlobals(logger) // 替换zap库中的全局的logger
	return
}

func getLogWrite(filename string, maxSize int, maxBackups int, maxAge int) zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   filename,
		MaxSize:    maxSize,    // 最大大小（M）
		MaxBackups: maxBackups, // 最大备份数量
		MaxAge:     maxAge,     // 最大备份天数
		Compress:   false,      // 是否压缩
	}
	return zapcore.AddSync(lumberJackLogger)
}

func getEncoder() zapcore.Encoder {
	// 自定义配置格式
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
	return zapcore.NewConsoleEncoder(encoderConfig)
}

// logger中间件
func GinLogger() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		start := time.Now() // 开始时间
		path := ctx.Request.URL.Path
		query := ctx.Request.URL.RawQuery
		ctx.Next() // 下一个中间件

		cost := time.Since(start)
		zap.L().Info(path, // 配置打印信息
			zap.Int("status", ctx.Writer.Status()), // 通过zap.xx这样结构化的日志减少内存分配，提高运行效率
			zap.String("method", ctx.Request.Method),
			zap.String("path", path),
			zap.String("query", query),
			zap.String("ip", ctx.ClientIP()),
			zap.String("user-agent", ctx.Request.UserAgent()),
			zap.String("errors", ctx.Errors.ByType(gin.ErrorTypePrivate).String()),
			zap.Duration("cost", cost),
		)
	}
}

// recovery(恢复)项目可能出现的panic
func GinRecovery(stack bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}

				httpRequest, _ := httputil.DumpRequest(c.Request, false)
				if brokenPipe {
					zap.L().Error(c.Request.URL.Path,
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
					c.Error(err.(error))
					c.Abort()
					return
				}

				if stack {
					zap.L().Error("[Recovery from panic]",
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
						zap.String("stack", string(debug.Stack())),
					)
				} else {
					zap.L().Error("[Recovery from panic]",
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
				}
				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()
		c.Next()
	}
}
