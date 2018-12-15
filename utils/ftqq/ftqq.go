/*
 * Copyright (c) 2018, 奶爸<1@5.nu>
 * All rights reserved.
 */

package ftqq

import (
	"net/url"
	"time"

	"github.com/naiba/nocd"
	"github.com/parnurzeal/gorequest"
)

//SCResp Server酱返回信息
type SCResp struct {
	Dataset string
	Errmsg  string
	Errno   int
	Error   []error
}

//SendMessage 推送消息
func SendMessage(key string, title string, msg string) SCResp {
	var resp SCResp
	msg += "\r\n\r\n(推送时间：" + time.Now().In(nocd.Loc).Format("2006-01-02 15:04:05") + ")"
	// UrlEncode 消息推送不到
	_, _, err := gorequest.New().Post("https://sc.ftqq.com/"+key+".send").
		SendString("text="+url.QueryEscape(title)+"&desp="+url.QueryEscape(msg)).Retry(3, time.Second*3).EndStruct(&resp)
	resp.Error = err
	return resp
}
