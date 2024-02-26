package tool

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc/metadata"
	"gopkg.in/natefinch/lumberjack.v2"
	_ "gopkg.in/natefinch/lumberjack.v2"
	"os"
	"runtime"
)

var lg *zap.Logger

var gDebugLevel = DebugLevel{
	Env:   "test",
	Debug: 1,
	Info:  1,
	Error: 1,
}

type DebugLevel struct {
	Env   string `json:"env"`
	Debug int    `json:"debug"`
	Info  int    `json:"info"`
	Error int    `json:"error"`
}

type LogConfig struct {
	DebugFileName string `json:"debugFileName"`
	InfoFileName  string `json:"infoFileName"`
	ErrorFileName string `json:"errorFileName"`
	MaxSize       int    `json:"maxSize"`
	MaxAge        int    `json:"maxAge"`     //存活时间，单位天
	MaxBackups    int    `json:"maxBackups"` //备份的数量
}

func logHeader(ctx context.Context) map[string]string {
	var strMap = make(map[string]string)
	if ctx == nil {
		return strMap
	}
	pc, codePath, codeLine, ok := runtime.Caller(2)
	var code = ""
	var funcName = ""
	if ok {
		code = fmt.Sprintf("%s:%d", codePath, codeLine)
		funcName = runtime.FuncForPC(pc).Name()
	}
	md, _ := metadata.FromIncomingContext(ctx)
	if md == nil {
		if _, ok := ctx.(*gin.Context); ok {
			gctx := ctx.(*gin.Context)
			dd, _ := gctx.Get("metadata")
			if dd != nil {
				md, _ = metadata.FromOutgoingContext(dd.(context.Context))
			} else {
				md = metadata.New(nil)
			}
		} else {
			if ctx.Value("metadata") != nil {
				md = ctx.Value("metadata").(metadata.MD)
			} else {
				md = metadata.New(nil)
			}
		}
	}
	strMap["headData"] = fmt.Sprintf("%+v", md)
	strMap["filePath"] = code
	strMap["funcName"] = funcName
	return strMap
}

func Debug(ctx context.Context, msg string) {
	if gDebugLevel.Debug == 0 {
		return
	}
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()
	logMap := logHeader(ctx)
	logMap["logMsg"] = msg
	lg.Debug(MapToJson(logMap))
}

func Info(ctx context.Context, msg string) {
	if gDebugLevel.Info == 0 {
		return
	}
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()
	logMap := logHeader(ctx)
	logMap["logMsg"] = msg
	lg.Info(MapToJson(logMap))
}

func Error(ctx context.Context, msg string) {
	if gDebugLevel.Error == 0 {
		return
	}
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()
	logMap := logHeader(ctx)
	logMap["logMsg"] = msg
	lg.Error(MapToJson(logMap))
}

func InitLogger(cfg *LogConfig, debugLevel *DebugLevel) (err error) {
	copier.Copy(&gDebugLevel, debugLevel)

	writeSyncerDebug := getLogWriter(cfg.DebugFileName, cfg.MaxSize, cfg.MaxBackups, cfg.MaxAge)
	writeSyncerInfo := getLogWriter(cfg.InfoFileName, cfg.MaxSize, cfg.MaxBackups, cfg.MaxAge)
	writeSyncerError := getLogWriter(cfg.ErrorFileName, cfg.MaxSize, cfg.MaxBackups, cfg.MaxAge)
	encoder := getEncoder()

	if gDebugLevel.Env != "prod" || gDebugLevel.Env != "release" {
		//file out
		debugCore := zapcore.NewCore(encoder, writeSyncerDebug, zapcore.DebugLevel)
		infoCore := zapcore.NewCore(encoder, writeSyncerInfo, zapcore.InfoLevel)
		errorCore := zapcore.NewCore(encoder, writeSyncerError, zapcore.WarnLevel)
		//std out
		consoleEncoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
		std := zapcore.NewCore(consoleEncoder, zapcore.Lock(os.Stdout), zapcore.DebugLevel)
		core := zapcore.NewTee(debugCore, infoCore, errorCore, std)
		lg = zap.New(core, zap.AddCaller())
		zap.ReplaceGlobals(lg)
	} else {
		//std out
		consoleEncoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
		std := zapcore.NewCore(consoleEncoder, zapcore.Lock(os.Stdout), zapcore.DebugLevel)
		core := zapcore.NewTee(std)
		lg = zap.New(core, zap.AddCaller())
		zap.ReplaceGlobals(lg)
	}
	return
}

func getLogWriter(filename string, maxSize int, maxBackups int, maxAge int) zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   filename,
		MaxSize:    maxSize,
		MaxBackups: maxBackups,
		MaxAge:     maxAge,
	}
	return zapcore.AddSync(lumberJackLogger)
}

func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.TimeKey = "time"
	encoderConfig.EncodeDuration = zapcore.SecondsDurationEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	return zapcore.NewConsoleEncoder(encoderConfig)
}
