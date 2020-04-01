/*
 * Copyright (c) 2017 - 2020, 奶爸<hi@nai.ba>
 * All rights reserved.
 */

package nocd

//Pipeline 部署流程
type Pipeline struct {
	ID           uint   `form:"id" binding:"min=0"`
	Name         string `form:"name" binding:"required,min=1,max=12"`
	Branch       string `form:"branch" binding:"required,alphanum,min=1,max=30"`
	Events       string
	Shell        string `form:"shell" binding:"required,min=3,max=1000"`
	UserID       uint
	ServerID     uint `form:"server" binding:"required,min=1"`
	RepositoryID uint `form:"repo" binding:"required,min=1"`

	User       User       `form:"-" binding:"-"`
	Server     Server     `form:"-" binding:"-"`
	Repository Repository `form:"-" binding:"-"`
	PipeLog    []PipeLog  `form:"-" binding:"-"`
	Webhook    []Webhook  `form:"-" binding:"-"`

	EventsSlice []string `gorm:"-" form:"events[]" binding:"required,min=1"`
}

//PipelineService 部署流程服务
type PipelineService interface {
	Create(p *Pipeline) error
	Update(p *Pipeline) error
	Delete(pid uint) error
	RepoPipelines(r *Repository) []Pipeline
	UserPipelines(u *User) []Pipeline
	UserPipeline(uid, pid uint) (Pipeline, error)
	GetPipelinesByRidAndEventAndBranch(rid uint, event string, branch string) ([]Pipeline, error)
	Server(p *Pipeline) error
	User(p *Pipeline) error
	Webhooks(p *Pipeline) error
}
