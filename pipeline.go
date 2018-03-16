/*
 * Copyright (c) 2018, 奶爸<1@5.nu>
 * All rights reserved.
 */

package gocd

type Pipeline struct {
	ID           uint       `form:"id" binding:"min=0"`
	Name         string     `form:"name" binding:"required,min=1,max=12"`
	Branch       string     `form:"branch" binding:"required,alphanum,min=1,max=30"`
	Events       string
	EventsSlice  []string   `gorm:"-" form:"events[]" binding:"required,min=1"`
	Shell        string     `form:"shell" binding:"required,min=3,max=1000"`
	UserID       uint
	User         User       `form:"-"`
	ServerID     uint       `form:"server" binding:"required,min=1"`
	Server       Server     `form:"-" binding:"-"`
	RepositoryID uint       `form:"repo" binding:"required,min=1"`
	Repository   Repository `form:"-" binding:"-"`
}

type PipelineService interface {
	CreatePipeline(p *Pipeline) error
	RepoPipelines(r *Repository) []Pipeline
	UserPipelines(u *User) []Pipeline
}
