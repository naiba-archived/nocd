/*
 * Copyright (c) 2018, 奶爸<1@5.nu>
 * All rights reserved.
 */

package sqlite3

import (
	"github.com/jinzhu/gorm"
	"git.cm/naiba/gocd"
)

type PipeLogService struct {
	DB *gorm.DB
}

func (ps *PipeLogService) Create(log *gocd.PipeLog) error {
	return ps.DB.Create(log).Error
}
