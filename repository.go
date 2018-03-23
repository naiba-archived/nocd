/*
 * Copyright (c) 2018, 奶爸<1@5.nu>
 * All rights reserved.
 */

package gocd

import (
	"github.com/naiba/webhooks/github"
	"github.com/naiba/webhooks/bitbucket"
	"github.com/naiba/webhooks/gitlab"
	"github.com/naiba/webhooks/gogs"
)

const (
	_ = iota
	//RepoPlatGitHub GitHub
	RepoPlatGitHub
	//RepoPlatBitBucket BitBucket
	RepoPlatBitBucket
	//RepoPlatGitlab Gitlab
	RepoPlatGitlab
	//RepoPlatGogs Gogs
	RepoPlatGogs
)

//RepoPlatforms 平台信息索引
var RepoPlatforms map[int]string
//RepoEvents 各平台支持的事件索引
var RepoEvents map[int]map[string]string

//Repository 项目
type Repository struct {
	ID       uint       `form:"id"`
	UserID   uint
	User     User       `form:"-" binding:"-"`
	Secret   string
	Name     string     `form:"name" binding:"required,min=1,max=12"`
	Platform int        `form:"platform" binding:"required,min=1,max=4"`
	Pipeline []Pipeline `form:"-" binding:"-"`
}

//RepositoryService 项目服务
type RepositoryService interface {
	Create(repo *Repository) error
	Update(repo *Repository) error
	Delete(rid uint) error
	GetRepoByUser(user *User) []Repository
	GetRepoByID(id uint) (Repository, error)
	GetRepoByUserAndID(user *User, id uint) (Repository, error)
}

func init() {
	RepoPlatforms = map[int]string{
		RepoPlatGitHub:    "GitHub",
		RepoPlatBitBucket: "BitBucket",
		RepoPlatGitlab:    "Gitlab",
		RepoPlatGogs:      "Gogs",
	}
	RepoEvents = map[int]map[string]string{
		RepoPlatGitHub: {
			string(github.PushEvent): "推送(Push)",
		},
		RepoPlatBitBucket: {
			string(bitbucket.PullRequestMergedEvent): "合并(Merge)",
		},
		RepoPlatGitlab: {
			string(gitlab.PushEvents): "推送(Push)",
		},
		RepoPlatGogs: {
			string(gogs.PushEvent): "推送(Push)",
		},
	}
}
