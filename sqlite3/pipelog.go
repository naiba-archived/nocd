/*
 * Copyright (c) 2018, 奶爸<1@5.nu>
 * All rights reserved.
 */

package sqlite3

import (
	"git.cm/naiba/gocd"
	"github.com/jinzhu/gorm"
)

//PipeLogService 日志服务
type PipeLogService struct {
	DB *gorm.DB
}

//Create 创建日志
func (ps *PipeLogService) Create(log *gocd.PipeLog) error {
	return ps.DB.Create(log).Error
}

//Update 创建日志
func (ps *PipeLogService) Update(log *gocd.PipeLog) error {
	return ps.DB.Save(log).Error
}

//Pipeline 获取部署流程信息
func (ps *PipeLogService) Pipeline(log *gocd.PipeLog) error {
	return ps.DB.Model(log).Related(&log.Pipeline).Error
}

//LastServerLog 服务器的最后一次部署
func (ps *PipeLogService) LastServerLog(sid uint) gocd.PipeLog {
	var pipelines []gocd.Pipeline
	var pl gocd.PipeLog
	id := make([]uint, 0)
	ps.DB.Model(&gocd.Server{ID: sid}).Select("id").Association("Pipelines").Find(&pipelines)
	for _, p := range pipelines {
		id = append(id, p.ID)
	}
	if len(id) > 0 {
		ps.DB.Where("pipeline_id IN (?)", id).Order("stopped_at desc").Take(&pl)
	}
	return pl
}

//UserLogs 用户的所有日志
func (ps *PipeLogService) UserLogs(uid uint, page, size int64) ([]gocd.PipeLog, int64) {
	var pipelines []gocd.Pipeline
	var pl []gocd.PipeLog
	var user gocd.User
	var num int64
	user.ID = uid
	ps.DB.Model(&user).Select("id").Association("Pipelines").Find(&pipelines)
	id := make([]uint, 0)
	for _, p := range pipelines {
		id = append(id, p.ID)
	}
	// 控制显示的历史 log 数
	ps.DB.Offset(page * size).Limit(size).Select("id,started_at,stopped_at,pipeline_id,pusher,status").Order("id desc").Where("pipeline_id IN (?)", id).Find(&pl)
	ps.DB.Model(&gocd.PipeLog{}).Where("pipeline_id IN (?)", id).Count(&num)
	if num%size == 0 {
		num = num / size
	} else {
		num = num/size + 1
	}
	return pl, num
}

//GetByUID 通过用户ID和部署流程ID查找日志
func (ps *PipeLogService) GetByUID(uid, lid uint) (gocd.PipeLog, error) {
	var pipelines []gocd.Pipeline
	var user gocd.User
	user.ID = uid
	ps.DB.Model(&user).Select("id").Association("Pipelines").Find(&pipelines)
	id := make([]uint, 0)
	for _, p := range pipelines {
		id = append(id, p.ID)
	}
	var log gocd.PipeLog
	err := ps.DB.Where("pipeline_id IN (?) AND id = ?", id, lid).First(&log).Error
	return log, err
}

//LastPipelineLog 部署流程最后一次部署
func (ps *PipeLogService) LastPipelineLog(pid uint) gocd.PipeLog {
	var pl gocd.PipeLog
	ps.DB.Where("pipeline_id = ?", pid).Order("id desc").Take(&pl)
	return pl
}

//LastLogs 全站最后部署记录
func (ps *PipeLogService) LastLogs(num uint) []gocd.PipeLog {
	var pl []gocd.PipeLog
	ps.DB.Order("id desc").Limit(num).Find(&pl)
	return pl
}
