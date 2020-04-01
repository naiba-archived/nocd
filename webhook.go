package nocd

const (
	_ = iota
	// WebhookRequestTypeJSON json
	WebhookRequestTypeJSON
	// WebhookRequestTypeForm form
	WebhookRequestTypeForm
)

const (
	_ = iota
	// WebhookRequestMethodGET ..
	WebhookRequestMethodGET
	// WebhookRequestMethodPOST ..
	WebhookRequestMethodPOST
)

// Webhook ..
type Webhook struct {
	ID            uint   `form:"id" gorm:"primary_key" json:"id,omitempty"`
	PipelineID    uint   `form:"pipeline_id" binding:"required,min=1" json:"pipeline_id,omitempty"`
	URL           string `form:"url" binding:"url" json:"url,omitempty"`
	RequestMethod int    `form:"request_method" json:"request_method,omitempty"`
	RequestType   int    `form:"request_type" json:"request_type,omitempty"`
	RequestBody   string `form:"request_body" gorm:"type:longtext" json:"request_body,omitempty"`
	VerifySSL     *bool  `form:"verify_ssl" json:"verify_ssl,omitempty"`
	PushSuccess   *bool  `form:"push_success" json:"push_success,omitempty"`
	Enable        *bool  `form:"enable" json:"enable,omitempty"`
}

// WebhookService ..
type WebhookService interface {
	Create(w *Webhook) error
	PipelineWebhooks(p *Pipeline) []Webhook
}
