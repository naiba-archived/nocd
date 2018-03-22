/*
 * Copyright (c) 2018, 奶爸<1@5.nu>
 * All rights reserved.
 */

package gocd

import "time"

const (
	_                            = iota
	PipeLogStatusSuccess
	PipeLogStatusErrorServerConn
	PipeLogStatusErrorShellExec
	PipeLogStatusRunning
)

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

type PipeLogService interface {
	Create(plog *PipeLog) error
	LastServerLog(sid uint) PipeLog
	LastPipelineLog(pid uint) PipeLog
	UserLogs(uid uint) []PipeLog
	Pipeline(log *PipeLog) error
	GetByUid(uid, lid uint) (PipeLog, error)
}
