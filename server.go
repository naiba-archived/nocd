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
	ID        uint `form:"id" binding:"min=0"`
	UserID    uint
	User      User       `form:"-" binding:"-"`
	Pipelines []Pipeline `form:"-" binding:"-"`
	Name      string     `form:"name" binding:"required,min=1,max=12"`
	Address   string     `form:"address" binding:"required,min=1,max=30"`
	Port      int        `form:"port" binding:"required,min=1"`
	Login     string     `form:"login"`
	LoginType uint       `form:"login_type" binding:"required"`
	Password  string     `form:"password" gorm:"longtext" binding:"required,min=1"`
	Status    int
}

//ServerService 服务器服务
type ServerService interface {
	CreateServer(s *Server) error
	DeleteServer(sid uint) error
	UpdateServer(s *Server) error
	GetServersByUser(user *User) []Server
	GetServersByUserAndSid(user *User, sid uint) (Server, error)
}
