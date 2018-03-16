/*
 * Copyright (c) 2018, 奶爸<1@5.nu>
 * All rights reserved.
 */

package sqlite3

import (
	"github.com/jinzhu/gorm"
	"git.cm/naiba/gocd"
)

type ServerService struct {
	DB *gorm.DB
}

func (ss *ServerService) CreateServer(s *gocd.Server) error {
	return ss.DB.Create(s).Error
}

func (ss *ServerService) GetServersByUser(user *gocd.User) (us []gocd.Server) {
	ss.DB.Model(user).Related(&us)
	return
}
func (ss *ServerService) GetServersByUserAndSid(user *gocd.User, sid uint) (s gocd.Server, err error) {
	err = ss.DB.Where("id = ? AND user_id =?", sid, user.ID).First(&s).Error
	return
}
