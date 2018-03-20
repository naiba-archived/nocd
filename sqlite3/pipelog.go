/*
 * Copyright (c) 2018, 奶爸<1@5.nu>
 * All rights reserved.
 */

package sqlite3

import (
	"github.com/jinzhu/gorm"
	"git.cm/naiba/gocd"
)

type PipeLogService struct {
	DB *gorm.DB
}

func (ps *PipeLogService) Create(log *gocd.PipeLog) error {
	return ps.DB.Create(log).Error
}

func (ps *PipeLogService) Pipeline(log *gocd.PipeLog) error {
	return ps.DB.Model(log).Related(&log.Pipeline).Error
}

func (ps *PipeLogService) LastServerLog(sid uint) gocd.PipeLog {
	var pipelines []gocd.Pipeline
	var pl gocd.PipeLog
	id := make([]uint, 0)
	ps.DB.Model(&gocd.Server{ID: sid}).Select("id").Association("Pipelines").Find(&pipelines)
	for _, p := range pipelines {
		id = append(id, p.ID)
	}
	ps.DB.Where("pipeline_id IN (?)", id).Order("id desc").First(&pl)
	return pl
}

func (ps *PipeLogService) UserLogs(uid uint) []gocd.PipeLog {
	var pipelines []gocd.Pipeline
	var pl []gocd.PipeLog
	var user gocd.User
	user.ID = uid
	ps.DB.Model(&user).Select("id").Association("Pipelines").Find(&pipelines)
	id := make([]uint, 0)
	for _, p := range pipelines {
		id = append(id, p.ID)
	}
	// 控制显示的历史 log 数
	ps.DB.Limit(20).Select("id,started_at,stopped_at,pipeline_id,pusher,status").Order("id desc").Where("pipeline_id IN (?)", id).Find(&pl)
	return pl
}

func (ps *PipeLogService) GetByUid(uid, lid uint) (gocd.PipeLog, error) {
	var pipelines []gocd.Pipeline
	var user gocd.User
	user.ID = uid
	ps.DB.Model(&user).Select("id").Association("Pipelines").Find(&pipelines)
	id := make([]uint, 0)
	for _, p := range pipelines {
		id = append(id, p.ID)
	}
	var log gocd.PipeLog
	err := ps.DB.Select("id,log").Where("pipeline_id IN (?)", id).First(&log).Error
	return log, err
}

func (ps *PipeLogService) LastPipelineLog(pid uint) gocd.PipeLog {
	var pl gocd.PipeLog
	ps.DB.Where("pipeline_id = ?", pid).Order("id desc").First(&pl)
	return pl
}
