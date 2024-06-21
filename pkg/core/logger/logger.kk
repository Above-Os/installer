/*
 Copyright 2021 The KubeSphere Authors.

 Licensed under the Apache License, Version 2.0 (the "License");
 you may not use this file except in compliance with the License.
 You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.
*/

package logger

import (
	"fmt"
	"path/filepath"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"

	"bytetrade.io/web3os/installer/pkg/core/common"
)

var Log *Logger

type Logger struct {
	logrus.FieldLogger
	OutputPath string
	LogLevel   string
}

func getLevel(level string) (l logrus.Level) {
	switch level {
	case "debug":
		l = logrus.DebugLevel
	case "info":
		l = logrus.InfoLevel
	case "warn":
		l = logrus.WarnLevel
	case "error":
		l = logrus.ErrorLevel
	case "panic":
		l = logrus.PanicLevel
	case "fatal":
		l = logrus.FatalLevel
	default:
		l = logrus.InfoLevel
	}
	return
}

func NewLogger(outputPath string, verbose string) *Logger {
	logger := logrus.New()

	formatter := &Formatter{
		HideKeys:               true,
		TimestampFormat:        "15:04:05 MST",
		NoColors:               true,
		ShowLevel:              logrus.WarnLevel,
		FieldsDisplayWithOrder: []string{common.Pipeline, common.Module, common.Task, common.Node},
	}
	logger.SetFormatter(formatter)

	path := filepath.Join(outputPath, "./terminus_install.log")
	writer, _ := rotatelogs.New(
		path+".%Y%m%d",
		rotatelogs.WithLinkName(path),
		rotatelogs.WithRotationTime(24*time.Hour),
	)

	logWriters := lfshook.WriterMap{
		logrus.InfoLevel:  writer,
		logrus.WarnLevel:  writer,
		logrus.ErrorLevel: writer,
		logrus.FatalLevel: writer,
		logrus.PanicLevel: writer,
	}
	l := getLevel(verbose)
	logger.SetLevel(l)
	logWriters[l] = writer
	logger.Hooks.Add(lfshook.NewHook(logWriters, formatter))
	return &Logger{logger, outputPath, verbose}
}

func (k *Logger) Message(node, str string) {
	Log.Infof("message: [%s]\n%s", node, str)
}

func (k *Logger) Messagef(node, format string, args ...interface{}) {
	Log.Infof("message: [%s]\n%s", node, fmt.Sprintf(format, args...))
}
