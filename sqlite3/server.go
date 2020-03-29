/*
 * Copyright (c) 2017 - 2020, 奶爸<hi@nai.ba>
 * All rights reserved.
 */

package sqlite3

import (
	"github.com/jinzhu/gorm"
	"github.com/naiba/nocd"
)

//ServerService 服务器服务
type ServerService struct {
	DB *gorm.DB
}

//CreateServer 创建服务器
func (ss *ServerService) CreateServer(s *nocd.Server) error {
	return ss.DB.Create(s).Error
}

//UpdateServer 更细服务器
func (ss *ServerService) UpdateServer(s *nocd.Server) error {
	return ss.DB.Save(s).Error
}

//DeleteServer 删除服务器
func (ss *ServerService) DeleteServer(sid uint) error {
	return ss.DB.Delete(nocd.Server{}, "id = ?", sid).Error
}

//GetServersByUser 获取用户的所有服务器
func (ss *ServerService) GetServersByUser(user *nocd.User) (us []nocd.Server) {
	ss.DB.Model(user).Related(&us)
	return
}

//GetServersByUserAndSid 根据用户ID和服务器ID获取服务器
func (ss *ServerService) GetServersByUserAndSid(user *nocd.User, sid uint) (s nocd.Server, err error) {
	err = ss.DB.Where("id = ? AND user_id =?", sid, user.ID).First(&s).Error
	return
}
