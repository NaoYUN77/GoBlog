package logger

import (
	"Blog/settings"
	"fmt"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)



var Logger *zap.Logger
func GetLevel() zapcore.Level {
	switch settings.Cnf.App.Level {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	case "fatal":
		return zapcore.FatalLevel
	default:
		return zapcore.InfoLevel
	}
}

// 定义日志模式
func Init(mode string) error {
	//core 需要满足 encoder  Writer Level
	switch mode {

	//开发模式日志输出到console
	case "development":
		encoderConfig := zap.NewDevelopmentEncoderConfig()

		//encoder console实现, 将日志输出到终端
		consoleEncoder := zapcore.NewConsoleEncoder(
			encoderConfig,
		)

		//console Syncer实现, 这里包装成AddSync 通过os.Stdout输出到终端
		consoleSyncer := zapcore.AddSync(os.Stdout)
		encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		core := zapcore.NewCore(
			consoleEncoder,
			consoleSyncer,
			GetLevel())
		Logger = zap.New(core)
		return nil

		//部署输出到日志文件中
	case "production":
		encoderConfig := zap.NewProductionEncoderConfig()
		fileEncoder := zapcore.NewJSONEncoder(
			encoderConfig,
		)
		file, err := os.OpenFile(

			fmt.Sprintf("%s",settings.Cnf.App.Log),
			os.O_CREATE|os.O_APPEND|os.O_WRONLY,
			0644,
		)
		fileWriteSyncer := zapcore.AddSync(file)
		encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		if err != nil {
			return err
		}
		core := zapcore.NewCore(
			fileEncoder,
			fileWriteSyncer,
			GetLevel(),
		)
		Logger = zap.New(core)
		return nil
	default:
		return fmt.Errorf("invalid log mode: %s", mode)
		}
}