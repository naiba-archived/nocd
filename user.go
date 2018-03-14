/*
 * Copyright (c) 2018, 奶爸<1@5.nu>
 * All rights reserved.
 */

package gocd

import "github.com/jinzhu/gorm"

type User struct {
	gorm.Model
	// 用户GitHubID
	GID        int64 `gorm:"unique_index"`
	GName      string
	GLogin     string
	GType      string
	Pubkey     string
	PrivateKey string
	Avatar     int
	// 用户Token
	Token string
}

type UserService interface {
	UserByGID(gid int64) (*User, error)
	CreateUser(u *User) error
	UpdateUser(u *User) error
	VerifyUser(uid, token string) (*User, error)
}
