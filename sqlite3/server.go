/*
 * Copyright (c) 2018, 奶爸<1@5.nu>
 * All rights reserved.
 */

package sqlite3

import (
	"git.cm/naiba/gocd"
	"github.com/jinzhu/gorm"
)

//ServerService 服务器服务
type ServerService struct {
	DB *gorm.DB
}

//CreateServer 创建服务器
func (ss *ServerService) CreateServer(s *gocd.Server) error {
	return ss.DB.Create(s).Error
}

//UpdateServer 更细服务器
func (ss *ServerService) UpdateServer(s *gocd.Server) error {
	return ss.DB.Save(s).Error
}

//DeleteServer 删除服务器
func (ss *ServerService) DeleteServer(sid uint) error {
	return ss.DB.Delete(gocd.Server{}, "id = ?", sid).Error
}

//GetServersByUser 获取用户的所有服务器
func (ss *ServerService) GetServersByUser(user *gocd.User) (us []gocd.Server) {
	ss.DB.Model(user).Related(&us)
	return
}

//GetServersByUserAndSid 根据用户ID和服务器ID获取服务器
func (ss *ServerService) GetServersByUserAndSid(user *gocd.User, sid uint) (s gocd.Server, err error) {
	err = ss.DB.Where("id = ? AND user_id =?", sid, user.ID).First(&s).Error
	return
}
