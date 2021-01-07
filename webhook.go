package nocd

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

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

func (n *Webhook) reqURL(status string, pipeline *Pipeline, pipelog *PipeLog) string {
	return replaceParamsInString(n.URL, status, pipeline, pipelog)
}

func (n *Webhook) reqBody(status string, pipeline *Pipeline, pipelog *PipeLog) (string, error) {
	if n.RequestMethod == WebhookRequestMethodGET {
		return "", nil
	}
	switch n.RequestType {
	case WebhookRequestTypeJSON:
		return replaceParamsInJSON(n.RequestBody, status, pipeline, pipelog), nil
	case WebhookRequestTypeForm:
		var data map[string]string
		if err := json.Unmarshal([]byte(n.RequestBody), &data); err != nil {
			return "", err
		}
		params := url.Values{}
		for k, v := range data {
			params.Add(k, replaceParamsInString(v, status, pipeline, pipelog))
		}
		return params.Encode(), nil
	}
	return "", errors.New("不支持的请求类型")
}

func (n *Webhook) reqContentType() string {
	if n.RequestMethod == WebhookRequestMethodGET {
		return ""
	}
	if n.RequestType == WebhookRequestTypeForm {
		return "application/x-www-form-urlencoded"
	}
	return "application/json"
}

func (n *Webhook) Send(status string, pipeline *Pipeline, pipelog *PipeLog) error {
	var verifySSL bool

	if n.VerifySSL != nil && *n.VerifySSL {
		verifySSL = true
	}

	transCfg := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: verifySSL},
	}
	client := &http.Client{Transport: transCfg, Timeout: time.Minute * 10}

	reqBody, err := n.reqBody(status, pipeline, pipelog)

	var resp *http.Response

	if err == nil {
		if n.RequestMethod == WebhookRequestMethodGET {
			resp, err = client.Get(n.reqURL(status, pipeline, pipelog))
		} else {
			resp, err = client.Post(n.reqURL(status, pipeline, pipelog), n.reqContentType(), strings.NewReader(reqBody))
		}
	}

	if err == nil && (resp.StatusCode < 200 || resp.StatusCode > 299) {
		err = fmt.Errorf("%d %s", resp.StatusCode, resp.Status)
	}

	// defer resp.Body.Close()
	// body, _ := ioutil.ReadAll(resp.Body)
	return err
}

func replaceParamsInString(str string, status string, pipeline *Pipeline, pipelog *PipeLog) string {
	str = strings.ReplaceAll(str, "#Pusher#", pipelog.Pusher)
	str = strings.ReplaceAll(str, "#Log#", pipelog.Log)
	str = strings.ReplaceAll(str, "#Status#", status)
	str = strings.ReplaceAll(str, "#PipelineName#", pipeline.Name)
	str = strings.ReplaceAll(str, "#PipelineID#", fmt.Sprintf("%d", pipeline.ID))
	str = strings.ReplaceAll(str, "#StartedAt#", pipelog.StartedAt.String())
	str = strings.ReplaceAll(str, "#StoppedAt#", pipelog.StoppedAt.String())
	return str
}

func replaceParamsInJSON(str string, status string, pipeline *Pipeline, pipelog *PipeLog) string {
	str = strings.ReplaceAll(str, "#Pusher#", jsonEscape(pipelog.Pusher))
	str = strings.ReplaceAll(str, "#Log#", jsonEscape(pipelog.Log))
	str = strings.ReplaceAll(str, "#Status#", jsonEscape(status))
	str = strings.ReplaceAll(str, "#PipelineName#", jsonEscape(pipeline.Name))
	str = strings.ReplaceAll(str, "#PipelineID#", jsonEscape(pipeline.ID))
	str = strings.ReplaceAll(str, "#StartedAt#", jsonEscape(pipelog.StartedAt.String()))
	str = strings.ReplaceAll(str, "#StoppedAt#", jsonEscape(pipelog.StoppedAt.String()))
	return str
}

func jsonEscape(raw interface{}) string {
	b, _ := json.Marshal(raw)
	strb := string(b)
	if strings.HasPrefix(strb, "\"") {
		return strb[1 : len(strb)-1]
	}
	return strb
}
