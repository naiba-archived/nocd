/*
 * Copyright (c) 2017 - 2020, 奶爸<hi@nai.ba>
 * All rights reserved.
 */

package nocd

import "time"

const (
	_ = iota
	//PipeLogStatusSuccess 成功
	PipeLogStatusSuccess
	//PipeLogStatusErrorServerConn 服务器链接失败
	PipeLogStatusErrorServerConn
	//PipeLogStatusErrorExec 执行错误
	PipeLogStatusErrorExec
	//PipeLogStatusRunning 正在执行部署
	PipeLogStatusRunning
	//PipeLogStatusHumanStopped 人工停止
	PipeLogStatusHumanStopped
	//PipeLogStatusErrorTimeout 执行超时
	PipeLogStatusErrorTimeout
)

//PipeLog 部署日志
type PipeLog struct {
	ID         uint
	StartedAt  time.Time
	StoppedAt  time.Time
	Pipeline   Pipeline
	PipelineID uint
	Pusher     string
	Log        string
	Status     int
}

//PipeLogService 部署日志服务
type PipeLogService interface {
	Create(plog *PipeLog) error
	Update(plog *PipeLog) error
	LastServerLog(sid uint) PipeLog
	LastPipelineLog(pid uint) PipeLog
	UserLogs(uid uint, page, size int64) ([]PipeLog, int64)
	Pipeline(log *PipeLog) error
	GetByUID(uid, lid uint) (PipeLog, error)
	GetByID(lid uint) (PipeLog, error)
	Logs(status int, page, size int64) ([]PipeLog, int64)
	LastLogs(num uint) []PipeLog
}
