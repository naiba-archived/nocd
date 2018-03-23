/*
 * Copyright (c) 2018, 奶爸<1@5.nu>
 * All rights reserved.
 */

package sqlite3

import (
	"git.cm/naiba/gocd"
	"github.com/jinzhu/gorm"
)

//UserService 用户服务
type UserService struct {
	DB *gorm.DB
}

//UserByGID 根据GitHubID获取用户
func (us *UserService) UserByGID(gid int64) (*gocd.User, error) {
	var u gocd.User
	return &u, us.DB.Where("g_id = ?", gid).First(&u).Error
}

//CreateUser 创建用户
func (us *UserService) CreateUser(u *gocd.User) error {
	return us.DB.Create(u).Error
}

//UpdateUser 更新用户
func (us *UserService) UpdateUser(u *gocd.User, cols ... string) error {
	return us.DB.Model(u).Select(cols).Updates(u).Error
}

//VerifyUser 校验用户
func (us *UserService) VerifyUser(uid, token string) (*gocd.User, error) {
	var u gocd.User
	return &u, us.DB.Where("id = ? AND token = ?", uid, token).First(&u).Error
}
