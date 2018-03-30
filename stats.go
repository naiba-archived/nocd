/*
 * Copyright (c) 2018, 奶爸<1@5.nu>
 * All rights reserved.
 */

package gocd

import (
	"github.com/jinzhu/gorm"
	"time"
)

//Stats 系统统计信息
type Stats struct {
	UserCount     int64
	ServerCount   int64
	PipelineCount int64
	RepoCount     int64
	PipeLogCount  int64
	LastLog       time.Time
	Update        time.Time
}

var db *gorm.DB
var ss *Stats

//GetStats 获取系统统计
func GetStats() Stats {
	if ss.Update.Add(time.Minute * 10).Before(time.Now()) {
		update()
		return *ss
	}
	return *ss
}

//InitStats 初始化系统统计
func InitStats(d *gorm.DB) {
	db = d
	ss = new(Stats)
}

func update() {
	db.Model(&User{}).Count(&ss.UserCount)
	db.Model(&Server{}).Count(&ss.ServerCount)
	db.Model(&Pipeline{}).Count(&ss.PipelineCount)
	db.Model(&Repository{}).Count(&ss.RepoCount)
	db.Model(&PipeLog{}).Count(&ss.PipeLogCount)
	var l PipeLog
	db.Select("stopped_at").Order("id DESC").Take(&l)
	ss.LastLog = l.StoppedAt
	ss.Update = time.Now()
}
