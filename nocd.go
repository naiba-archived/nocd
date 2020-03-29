/*
 * Copyright (c) 2017 - 2020, 奶爸<hi@nai.ba>
 * All rights reserved.
 */

package nocd

import (
	"runtime"
	"time"

	"github.com/evalphobia/logrus_sentry"
	log "github.com/sirupsen/logrus"
	"gopkg.in/ini.v1"
)

//mLog of sentry logger
var mLog *log.Logger

//Conf of NoCD config
var Conf *ini.File

//Debug debuggable
var Debug bool

//Loc system time location
var Loc *time.Location

//InitSysConfig system: load common config
func InitSysConfig(file string) {
	var err error
	if mLog == nil {
		mLog = log.New()
	}
	Conf, err = ini.Load(file)
	if err != nil {
		mLog.Panicln(err)
	}
	// initial sentry dsn
	hook, err := logrus_sentry.NewSentryHook(Conf.Section("third_party").Key("sentry_dsn").String(), []log.Level{
		log.PanicLevel,
		log.FatalLevel,
		log.ErrorLevel,
	})
	if err == nil {
		mLog.Hooks.Add(hook)
	} else {
		mLog.Panicln(err)
	}
	// set timezone
	Loc, err = time.LoadLocation(Conf.Section("nocd").Key("loc").String())
	if err != nil {
		panic(err)
	}
	// set debuggable
	Debug, err = Conf.Section("nocd").Key("debug").Bool()
	if err != nil {
		panic(err)
	}
	if Debug {
		mLog.SetLevel(log.DebugLevel)
	} else {
		mLog.SetLevel(log.InfoLevel)
	}
}

//Logger 带行号文件名方法名的Logger
func Logger() *log.Entry {
	logger := log.NewEntry(mLog)
	if pc, file, line, ok := runtime.Caller(1); ok {
		fName := runtime.FuncForPC(pc).Name()
		return logger.WithField("file", file).WithField("line", line).WithField("func", fName)
	}
	return logger
}
