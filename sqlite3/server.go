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
	D *gorm.DB
}

func (ss *ServerService) CreateServer(s *gocd.Server) error {
	return ss.D.Create(s).Error
}

func (ss *ServerService) GetServersByUser(user *gocd.User) (us []gocd.Server) {
	ss.D.Model(user).Related(&us)
	return
}
