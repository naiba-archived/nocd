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

func (ps *PipelineService) CreatePipeline(p *gocd.Pipeline) error {
	return ps.DB.Create(p).Error
}

func (ps *PipelineService) UserPipelines(u *gocd.User) []gocd.Pipeline {
	var p []gocd.Pipeline
	ps.DB.Model(u).Related(&p)
	return p
}
func (ps *PipelineService) RepoPipelines(r *gocd.Repository) []gocd.Pipeline {
	var p []gocd.Pipeline
	ps.DB.Model(r).Related(&p)
	return p
}
