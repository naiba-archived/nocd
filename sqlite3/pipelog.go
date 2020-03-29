/*
 * Copyright (c) 2017 - 2020, 奶爸<hi@nai.ba>
 * All rights reserved.
 */

package sqlite3

import (
	"github.com/jinzhu/gorm"
	"github.com/naiba/nocd"
)

//PipeLogService 日志服务
type PipeLogService struct {
	DB *gorm.DB
}

//Create 创建日志
func (ps *PipeLogService) Create(log *nocd.PipeLog) error {
	return ps.DB.Create(log).Error
}

//Update 创建日志
func (ps *PipeLogService) Update(log *nocd.PipeLog) error {
	return ps.DB.Save(log).Error
}

//Pipeline 获取部署流程信息
func (ps *PipeLogService) Pipeline(log *nocd.PipeLog) error {
	return ps.DB.Model(log).Related(&log.Pipeline).Error
}

//LastServerLog 服务器的最后一次部署
func (ps *PipeLogService) LastServerLog(sid uint) nocd.PipeLog {
	var pipelines []nocd.Pipeline
	var pl nocd.PipeLog
	id := make([]uint, 0)
	ps.DB.Model(&nocd.Server{ID: sid}).Select("id").Association("Pipelines").Find(&pipelines)
	for _, p := range pipelines {
		id = append(id, p.ID)
	}
	if len(id) > 0 {
		ps.DB.Where("pipeline_id IN (?)", id).Order("stopped_at desc").Take(&pl)
	}
	return pl
}

//UserLogs 用户的所有日志
func (ps *PipeLogService) UserLogs(uid uint, page, size int64) ([]nocd.PipeLog, int64) {
	var pipelines []nocd.Pipeline
	var pl []nocd.PipeLog
	var user nocd.User
	var num int64
	user.ID = uid
	ps.DB.Model(&user).Select("id").Association("Pipelines").Find(&pipelines)
	id := make([]uint, 0)
	for _, p := range pipelines {
		id = append(id, p.ID)
	}
	// 控制显示的历史 log 数
	ps.DB.Offset(page*size).Limit(size).Select("id,started_at,stopped_at,pipeline_id,pusher,status").Order("id desc").Where("pipeline_id IN (?)", id).Find(&pl)
	ps.DB.Model(&nocd.PipeLog{}).Where("pipeline_id IN (?)", id).Count(&num)
	if num%size == 0 {
		num = num / size
	} else {
		num = num/size + 1
	}
	return pl, num
}

//Logs 获取所有日志
func (ps *PipeLogService) Logs(status int, page, size int64) ([]nocd.PipeLog, int64) {
	var pl []nocd.PipeLog
	var num int64

	// 控制显示的 log 数
	ps.DB.Offset(page*size).Limit(size).Select("id,started_at,stopped_at,pipeline_id,pusher,status").Order("id desc").Where("status = ?", status).Find(&pl)
	ps.DB.Model(&nocd.PipeLog{}).Where("status = ?", status).Count(&num)
	if num%size == 0 {
		num = num / size
	} else {
		num = num/size + 1
	}
	return pl, num
}

//GetByUID 通过用户ID和部署流程ID查找日志
func (ps *PipeLogService) GetByUID(uid, lid uint) (nocd.PipeLog, error) {
	var pipelines []nocd.Pipeline
	var user nocd.User
	user.ID = uid
	ps.DB.Model(&user).Select("id").Association("Pipelines").Find(&pipelines)
	id := make([]uint, 0)
	for _, p := range pipelines {
		id = append(id, p.ID)
	}
	var log nocd.PipeLog
	err := ps.DB.Where("pipeline_id IN (?) AND id = ?", id, lid).First(&log).Error
	return log, err
}

//GetByID 通过日志ID获取日志
func (ps *PipeLogService) GetByID(lid uint) (nocd.PipeLog, error) {
	var log nocd.PipeLog
	err := ps.DB.First(&log, lid).Error
	return log, err
}

//LastPipelineLog 部署流程最后一次部署
func (ps *PipeLogService) LastPipelineLog(pid uint) nocd.PipeLog {
	var pl nocd.PipeLog
	ps.DB.Where("pipeline_id = ?", pid).Order("id desc").Take(&pl)
	return pl
}

//LastLogs 全站最后部署记录
func (ps *PipeLogService) LastLogs(num uint) []nocd.PipeLog {
	var pl []nocd.PipeLog
	ps.DB.Order("id desc").Limit(num).Find(&pl)
	return pl
}
