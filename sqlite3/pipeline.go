/*
 * Copyright (c) 2018, 奶爸<1@5.nu>
 * All rights reserved.
 */

package sqlite3

import (
	"github.com/jinzhu/gorm"
	"git.cm/naiba/gocd"
)

//PipelineService Sqlite3的用户服务实现
type PipelineService struct {
	DB *gorm.DB
}

//Create 创建用户
func (ps *PipelineService) Create(p *gocd.Pipeline) error {
	return ps.DB.Create(p).Error
}

//Update 更新用户
func (ps *PipelineService) Update(p *gocd.Pipeline) error {
	return ps.DB.Save(p).Error
}

//Delete 删除用户
func (ps *PipelineService) Delete(pid uint) error {
	db := ps.DB.Begin()
	err := db.Where("pipeline_id = ?", pid).Delete(gocd.PipeLog{}).Error
	if err != nil {
		db.Rollback()
		return err
	}
	err = db.Where("id = ?", pid).Delete(gocd.Pipeline{}).Error
	if err != nil {
		db.Rollback()
		return err
	}
	db.Commit()
	return nil
}

//UserPipelines 获取用户的所有部署流
func (ps *PipelineService) UserPipelines(u *gocd.User) []gocd.Pipeline {
	var p []gocd.Pipeline
	ps.DB.Model(u).Related(&p)
	return p
}

//UserPipeline 根据用户ID和部署流ID获取部署流，通常用来检测用户对部署流的操作权限
func (ps *PipelineService) UserPipeline(uid, pid uint) (gocd.Pipeline, error) {
	var p gocd.Pipeline
	err := ps.DB.Where("user_id = ? AND id = ?", uid, pid).First(&p).Error
	return p, err
}

//RepoPipelines 获取项目下面的所有部署流
func (ps *PipelineService) RepoPipelines(r *gocd.Repository) []gocd.Pipeline {
	var p []gocd.Pipeline
	ps.DB.Model(r).Related(&p)
	return p
}

//GetPipelinesByRidAndEventAndBranch 根据流ID和事件及分支获取部署流
func (ps *PipelineService) GetPipelinesByRidAndEventAndBranch(rid uint, event string, branch string) (p []gocd.Pipeline, err error) {
	err = ps.DB.Where("repository_id = ? AND events LIKE ? AND branch = ?", rid, "%\""+event+"\"%", branch).Find(&p).Error
	return
}

//Server 读取所属服务器信息
func (ps *PipelineService) Server(p *gocd.Pipeline) error {
	return ps.DB.Model(p).Related(&p.Server).Error
}

//User 读取所属用户信息
func (ps *PipelineService) User(p *gocd.Pipeline) error {
	return ps.DB.Model(p).Related(&p.User).Error
}
