/*
 * Copyright (c) 2018, 奶爸<1@5.nu>
 * All rights reserved.
 */

package gocd

import (
	"gopkg.in/go-playground/webhooks.v3/github"
	"gopkg.in/go-playground/webhooks.v3/bitbucket"
	"gopkg.in/go-playground/webhooks.v3/gitlab"
)

const (
	_                 = iota
	RepoPlatGitHub
	RepoPlatBitBucket
	RepoPlatGitlab
	RepoPlatGogs
)

var RepoPlatforms map[int]string
var RepoEvents map[int]map[string]string

type Repository struct {
	ID       uint       `form:"id" binding:"min=0"`
	UserID   uint
	User     User       `binding:"-"`
	Secret   string
	Name     string     `form:"name" binding:"required,min=1,max=12"`
	Platform int        `form:"platform" binding:"required,min=1,max=4"`
	Pipeline []Pipeline `binding:"-"`
}

type RepositoryService interface {
	CreateRepo(repo *Repository) error
	GetRepoByUser(user *User) []Repository
	GetRepoByUserAndRid(user *User, rid uint) (Repository, error)
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
			string(github.PushEvent):    "推送",
			string(github.ReleaseEvent): "发布",
			string(github.CreateEvent):  "创建（Tag、分支）",
		},
		RepoPlatBitBucket: {
			string(bitbucket.RepoPushEvent):          "推送",
			string(bitbucket.PullRequestMergedEvent): "合并",
		},
		RepoPlatGitlab: {
			string(gitlab.PushEvents): "推送",
			string(gitlab.TagEvents):  "Tag",
		},
	}
}
