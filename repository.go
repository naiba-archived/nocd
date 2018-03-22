/*
 * Copyright (c) 2018, 奶爸<1@5.nu>
 * All rights reserved.
 */

package gocd

import (
	"gopkg.in/go-playground/webhooks.v3/github"
	"gopkg.in/go-playground/webhooks.v3/bitbucket"
	"gopkg.in/go-playground/webhooks.v3/gitlab"
	"github.com/naiba/webhooks/gogs"
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
	ID       uint       `form:"id"`
	UserID   uint
	User     User       `form:"-" binding:"-"`
	Secret   string
	Name     string     `form:"name" binding:"required,min=1,max=12"`
	Platform int        `form:"platform" binding:"required,min=1,max=4"`
	Pipeline []Pipeline `form:"-" binding:"-"`
}

type RepositoryService interface {
	CreateRepo(repo *Repository) error
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
