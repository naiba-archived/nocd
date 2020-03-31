package nocd

import "github.com/jinzhu/gorm"

const (
	_ = iota
	// RequestMethodGet GET 请求
	RequestMethodGet
	// RequestMethodPost POST 请求
	RequestMethodPost
)

const (
	_ = iota
	// WebhookRequestTypeJSON json
	WebhookRequestTypeJSON
	// WebhookRequestTypeForm form
	WebhookRequestTypeForm
)

// Webhook ..
type Webhook struct {
	gorm.Model
	PipelineID    uint   `form:"pipeline_id" binding:"required"`
	URL           string `form:"url" binding:"url"`
	RequestMethod int    `form:"request_method"`
	RequestType   int    `form:"request_type"`
	RequestBody   string `gorm:"type:longtext" form:"request_body"`
	VerifySSL     bool   `form:"verify_ssl"`
	PushSuccess   bool   `form:"push_success"`
	Enable        bool   `form:"enable"`
}
