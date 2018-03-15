/*
 * Copyright (c) 2018, 奶爸<1@5.nu>
 * All rights reserved.
 */

package sqlite3

import (
	"github.com/jinzhu/gorm"
	"git.cm/naiba/gocd"
)

type RepositoryService struct {
	DB *gorm.DB
}

func (rs *RepositoryService) CreateRepo(r *gocd.Repository) error {
	return rs.DB.Create(r).Error
}

func (rs *RepositoryService) GetRepoByUser(user *gocd.User) (r []gocd.Repository) {
	rs.DB.Model(user).Related(&r)
	return
}
