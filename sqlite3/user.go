/*
 * Copyright (c) 2018, 奶爸<1@5.nu>
 * All rights reserved.
 */

package sqlite3

import (
	"git.cm/naiba/gocd"
	"github.com/jinzhu/gorm"
)

type UserService struct {
	DB *gorm.DB
}

func (us *UserService) UserByGID(gid int64) (*gocd.User, error) {
	var u gocd.User
	return &u, us.DB.Where("g_id = ?", gid).First(&u).Error
}

func (us *UserService) CreateUser(u *gocd.User) error {
	return us.DB.Create(u).Error
}

func (us *UserService) UpdateUser(u *gocd.User) error {
	return us.DB.Save(u).Error
}

func (us *UserService) VerifyUser(uid, token string) (*gocd.User, error) {
	var u gocd.User
	return &u, us.DB.Where("id = ? AND token = ?", uid, token).First(&u).Error
}
