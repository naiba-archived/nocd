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

//Users 获取所有用户
func (us *UserService) Users(page, limit int64) ([]*gocd.User, int64) {
	var ul []*gocd.User
	var num int64
	us.DB.Offset(page * 20).Limit(limit).Order("updated_at DESC").Find(&ul)
	for _, u := range ul {
		us.DB.Model(&u).Select("id").Related(&u.Servers)
		us.DB.Model(&u).Select("id").Related(&u.Repositories)
		us.DB.Model(&u).Select("id").Related(&u.Pipelines)
	}
	us.DB.Model(&gocd.User{}).Count(&num)
	if num%gocd.Pagination == 0 {
		num = num/20
	} else {
		num = num/20 + 1
	}
	return ul, num
}

//Create 创建用户
func (us *UserService) Create(u *gocd.User) error {
	return us.DB.Create(u).Error
}

//Update 更新用户
func (us *UserService) Update(u *gocd.User) error {
	return us.DB.Save(u).Error
}

//Verify 校验用户
func (us *UserService) Verify(uid, token string) (*gocd.User, error) {
	var u gocd.User
	return &u, us.DB.Where("id = ? AND token = ?", uid, token).First(&u).Error
}
