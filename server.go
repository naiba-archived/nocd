/*
 * Copyright (c) 2017 - 2020, 奶爸<hi@nai.ba>
 * All rights reserved.
 */

package nocd

//Server 服务器
type Server struct {
	ID        uint `form:"id" binding:"min=0"`
	UserID    uint
	User      User       `form:"-" binding:"-"`
	Pipelines []Pipeline `form:"-" binding:"-"`
	Name      string     `form:"name" binding:"required,min=1,max=12"`
	Address   string     `form:"address" binding:"required,address,min=1,max=30"`
	Port      int        `form:"port" binding:"required,min=1"`
	Login     string     `form:"login"`
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
