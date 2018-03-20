/*
 * Copyright (c) 2018, 奶爸<1@5.nu>
 * All rights reserved.
 */

package gocd

type Server struct {
	ID      uint   `form:"id" binding:"min=0"`
	UserID  uint
	User    User   `form:"-" binding:"-"`
	Pipelines []Pipeline `form:"-" binding:"-"`
	Name    string `form:"name" binding:"required,min=1,max=12"`
	Address string `form:"address" binding:"required,address,min=1,max=30"`
	Port    int    `form:"port" binding:"required,min=1"`
	Login   string `form:"login" binding:"required,alphanum,min=1,max=30"`
	Status  int
}

type ServerService interface {
	CreateServer(s *Server) error
	GetServersByUser(user *User) []Server
	GetServersByUserAndSid(user *User, sid uint) (Server, error)
}
