/*
 * Copyright (c) 2018, 奶爸<1@5.nu>
 * All rights reserved.
 */

package gocd

import (
	log "github.com/sirupsen/logrus"
	"gopkg.in/ini.v1"
	"time"
)

var Log *log.Logger
var Conf *ini.File
var Debug bool
var Loc *time.Location

func Initial(file string) {
	var err error
	if Log == nil {
		Log = log.New()
	}
	Conf, err = ini.Load(file)
	if err != nil {
		Log.Panicln(err)
	}
	Loc, err = time.LoadLocation(Conf.Section("gocd").Key("loc").String())
	if err != nil {
		panic(err)
	}
	Debug, err = Conf.Section("gocd").Key("debug").Bool()
	if err != nil {
		panic(err)
	}
	if Debug {
		Log.SetLevel(log.DebugLevel)
	} else {
		Log.SetLevel(log.InfoLevel)
	}
}
