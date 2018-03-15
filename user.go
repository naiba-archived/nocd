/*
 * Copyright (c) 2018, 奶爸<1@5.nu>
 * All rights reserved.
 */

package gocd

import "github.com/jinzhu/gorm"

type User struct {
	gorm.Model
	// 用户GitHubID
	GID        uint `gorm:"unique_index"`
	GName      string
	GLogin     string
	GType      string
	Pubkey     string
	PrivateKey string
	Avatar     int
	Servers    []Server
	// 用户Token
	Token string
}

type UserService interface {
	UserByGID(gid int64) (*User, error)
	CreateUser(u *User) error
	UpdateUser(u *User, cols ... string) error
	VerifyUser(uid, token string) (*User, error)
}
