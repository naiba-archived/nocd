/*
 * Copyright (c) 2018, 奶爸<1@5.nu>
 * All rights reserved.
 */

package sqlite3

import (
	"github.com/jinzhu/gorm"
	"git.cm/naiba/gocd"
)

type PipelineService struct {
	DB *gorm.DB
}

func (ps *PipelineService) Create(p *gocd.Pipeline) error {
	return ps.DB.Create(p).Error
}

func (ps *PipelineService) Update(p *gocd.Pipeline) error {
	return ps.DB.Save(p).Error
}

func (ps *PipelineService) Delete(pid uint) error {
	return ps.DB.Where("id = ?", pid).Delete(gocd.Pipeline{}).Error
}

func (ps *PipelineService) UserPipelines(u *gocd.User) []gocd.Pipeline {
	var p []gocd.Pipeline
	ps.DB.Model(u).Related(&p)
	return p
}

func (ps *PipelineService) UserPipeline(uid, pid uint) (gocd.Pipeline, error) {
	var p gocd.Pipeline
	err := ps.DB.Where("user_id = ? AND id = ?", uid, pid).First(&p).Error
	return p, err
}
func (ps *PipelineService) RepoPipelines(r *gocd.Repository) []gocd.Pipeline {
	var p []gocd.Pipeline
	ps.DB.Model(r).Related(&p)
	return p
}

func (ps *PipelineService) GetPipelinesByRidAndEventAndBranch(rid uint, event string, branch string) (p []gocd.Pipeline, err error) {
	err = ps.DB.Where("repository_id = ? AND events LIKE ? AND branch = ?", rid, "%\""+event+"\"%", branch).Find(&p).Error
	return
}

func (ps *PipelineService) Server(p *gocd.Pipeline) error {
	return ps.DB.Model(p).Related(&p.Server).Error
}
func (ps *PipelineService) User(p *gocd.Pipeline) error {
	return ps.DB.Model(p).Related(&p.User).Error
}
