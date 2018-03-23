/*
 * Copyright (c) 2018, 奶爸<1@5.nu>
 * All rights reserved.
 */

package sqlite3

import (
	"github.com/jinzhu/gorm"
	"git.cm/naiba/gocd"
)

//RepositoryService 项目服务
type RepositoryService struct {
	DB *gorm.DB
}

//Create 创建项目
func (rs *RepositoryService) Create(r *gocd.Repository) error {
	return rs.DB.Create(r).Error
}

//Delete 删除项目
func (rs *RepositoryService) Delete(rid uint) error {
	ids := make([]uint, 0)
	var pipelines []gocd.Pipeline
	if err := rs.DB.Select("id").Where("repository_id = ?", rid).Find(&pipelines).Error; err != nil {
		return err
	}
	for _, p := range pipelines {
		ids = append(ids, p.ID)
	}
	db := rs.DB.Begin()
	// 删除关联 PipeLog
	err := db.Where("pipeline_id IN (?)", ids).Delete(gocd.PipeLog{}).Error
	if err != nil {
		db.Rollback()
		return err
	}
	// 删除关联 Pipeline
	err = db.Where("id IN (?)", ids).Delete(gocd.Pipeline{}).Error
	if err != nil {
		db.Rollback()
		return err
	}
	// 删除项目
	err = db.Where("id = ?", rid).Delete(gocd.Repository{}).Error
	if err != nil {
		db.Rollback()
		return err
	}
	db.Commit()
	return nil
}

//Update 更新项目
func (rs *RepositoryService) Update(r *gocd.Repository) error {
	return rs.DB.Save(r).Error
}

//GetRepoByUser 获取用户的所有项目
func (rs *RepositoryService) GetRepoByUser(user *gocd.User) (r []gocd.Repository) {
	rs.DB.Model(user).Related(&r)
	return
}

//GetRepoByUserAndID 通过用户ID和项目ID寻找
func (rs *RepositoryService) GetRepoByUserAndID(user *gocd.User, rid uint) (r gocd.Repository, err error) {
	err = rs.DB.Where("id = ? AND user_id = ?", rid, user.ID).First(&r).Error
	return
}

//GetRepoByID 通过ID获取项目
func (rs *RepositoryService) GetRepoByID(id uint) (r gocd.Repository, err error) {
	err = rs.DB.First(&r, id).Error
	return
}
