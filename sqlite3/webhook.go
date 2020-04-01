/*
 * Copyright (c) 2017 - 2020, 奶爸<hi@nai.ba>
 * All rights reserved.
 */

package sqlite3

import (
	"github.com/jinzhu/gorm"
	"github.com/naiba/nocd"
)

//WebhookService ..
type WebhookService struct {
	DB *gorm.DB
}

//PipelineWebhooks ..
func (ws *WebhookService) PipelineWebhooks(p *nocd.Pipeline) []nocd.Webhook {
	var w []nocd.Webhook
	ws.DB.Model(p).Related(&w)
	return w
}

// Create ..
func (ws *WebhookService) Create(w *nocd.Webhook) error {
	return ws.DB.Create(w).Error
}
