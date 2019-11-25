package log

import (
	"github.com/sirupsen/logrus"
	"os"
)

func init() {
	// 以JSON格式为输出，代替默认的ASCII格式
	logrus.SetFormatter(&logrus.JSONFormatter{})
	// 以Stdout为输出，代替默认的stderr
	logrus.SetOutput(os.Stdout)
	// 设置日志等级
	logrus.SetLevel(logrus.InfoLevel)
}

func Trace(args ...interface{}){
	logrus.Traceln(args)
}

func TraceByFields(fields map[string]interface{},args ...interface{}){
	logrus.WithFields(logrus.Fields(fields)).Traceln(args)
}

func Info(args ...interface{}){
	logrus.Infoln(args)
}

func InfoByFields(fields map[string]interface{},args ...interface{}){
	logrus.WithFields(logrus.Fields(fields)).Infoln(args)
}

func Debug(args ...interface{}){
	logrus.Debugln(args)
}

func DebugByFields(fields map[string]interface{},args ...interface{})  {
	logrus.WithFields(logrus.Fields(fields)).Debugln(args)
}

func Error(args ...interface{}){
	logrus.Errorln(args)
}

func ErrorByFields(fields map[string]interface{},args ...interface{}){
	logrus.WithFields(logrus.Fields(fields)).Errorln(args)
}



func Fatal(args ...interface{}){
	logrus.Fatalln(args)
}

func FatalByFields(fields map[string]interface{},args ...interface{}){
	logrus.WithFields(logrus.Fields(fields)).Fatalln(args)
}