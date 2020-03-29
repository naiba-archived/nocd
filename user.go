/*
 * Copyright (c) 2017 - 2020, 奶爸<hi@nai.ba>
 * All rights reserved.
 */

package nocd

import "github.com/jinzhu/gorm"

const (
	_ = iota
	// RequestMethodGet GET 请求
	RequestMethodGet
	// RequestMethodPost POST 请求
	RequestMethodPost
)

const (
	_ = iota
	// RequestTypeJSON json
	RequestTypeJSON
	// RequestTypeForm form
	RequestTypeForm
)

//User 用户
type User struct {
	gorm.Model

	// 用户 GitHub 信息
	GID        uint `gorm:"unique_index"`
	GName      string
	GLogin     string
	Pubkey     string
	PrivateKey string

	// Webhook 配置
	WebhookURL    string `form:"webhook_url" binding:"url"`
	RequestMethod int    `form:"request_method"`
	RequestType   int    `form:"request_type"`
	RequestBody   string `gorm:"type:longtext" form:"request_body"`
	VerifySSL     bool   `form:"verify_ssl"`
	PushSuccess   bool   `form:"push_success"`

	Servers      []Server     `form:"-"`
	Repositories []Repository `form:"-"`
	Pipelines    []Pipeline   `form:"-"`
	IsBlocked    bool
	IsAdmin      bool
	// 用户Token
	Token string
}

//UserService 用户服务
type UserService interface {
	UserByGID(gid int64) (*User, error)
	Create(u *User) error
	Update(u *User) error
	Verify(uid, token string) (*User, error)
	Users(page, size int64) ([]*User, int64)
}
