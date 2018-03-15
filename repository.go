/*
 * Copyright (c) 2018, 奶爸<1@5.nu>
 * All rights reserved.
 */

package gocd

const (
	_                 = iota
	RepoPlatGitHub
	RepoPlatBitbucket
	RepoPlatGitlab
	RepoPlatGogs
)

type Repository struct {
	ID       uint   `form:"id" binding:"min=0"`
	UserID   uint
	User     User
	Secret   string
	Name     string `form:"name" binding:"required,min=1,max=12"`
	Platform int    `form:"platform" binding:"required,min=1,max=4"`
}

type RepositoryService interface {
	CreateRepo(repo *Repository) error
	GetRepoByUser(user *User) []Repository
}
