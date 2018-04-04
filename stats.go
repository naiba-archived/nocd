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
	RunningCount  int64
	PipeLogCount  int64
	Update        time.Time
}

var db *gorm.DB
var ss *Stats

//GetStats 获取系统统计
func GetStats() Stats {
	update()
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
	db.Model(&PipeLog{}).Where("status = ?", PipeLogStatusRunning).Count(&ss.RunningCount)
	db.Model(&PipeLog{}).Count(&ss.PipeLogCount)
	ss.Update = time.Now()
}
