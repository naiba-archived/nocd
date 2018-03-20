/*
 * Copyright (c) 2018, 奶爸<1@5.nu>
 * All rights reserved.
 */

package gocd

import "time"

const (
	PipeLogStatusUnknown         = iota
	PipeLogStatusSuccess
	PipeLogStatusErrorServerConn
	PipeLogStatusErrorShellExec
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
}
