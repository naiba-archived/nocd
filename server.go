/*
 * Copyright (c) 2017 - 2020, 奶爸<hi@nai.ba>
 * All rights reserved.
 */

package nocd

const (
	_ = iota
	// ServerLoginTypePassword 密码登录
	ServerLoginTypePassword
	// ServerLoginTypePrivateKey 私钥登录
	ServerLoginTypePrivateKey
)

//Server 服务器
type Server struct {
	ID        uint   `form:"id" binding:"min=0" json:"id,omitempty"`
	UserID    uint   `json:"user_id,omitempty"`
	Name      string `form:"name" binding:"required,min=1,max=12" json:"name,omitempty"`
	Address   string `form:"address" binding:"required,min=1,max=30" json:"address,omitempty"`
	Port      int    `form:"port" binding:"required,min=1" json:"port,omitempty"`
	Login     string `form:"login" json:"login,omitempty"`
	LoginType uint   `form:"login_type" binding:"required" json:"login_type,omitempty"`
	Password  string `form:"password" gorm:"longtext" binding:"required,min=1" json:"password,omitempty"`

	User      User       `form:"-" binding:"-" json:"-"`
	Pipelines []Pipeline `form:"-" binding:"-" json:"-"`
}

//ServerService 服务器服务
type ServerService interface {
	CreateServer(s *Server) error
	DeleteServer(sid uint) error
	UpdateServer(s *Server) error
	GetServersByUser(user *User) []Server
	GetServersByUserAndSid(user *User, sid uint) (Server, error)
}
