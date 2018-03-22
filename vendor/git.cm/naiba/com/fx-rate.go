package com

import (
	"github.com/parnurzeal/gorequest"
	"log"
	"encoding/json"
)

type FXRate struct {
	Base  string
	Date  string
	Rates map[string]float64
}

var rateInstance *FXRate

func GetRate() *FXRate {
	if rateInstance == nil {
		req := gorequest.New()
		_, body, err := req.Get("https://api.fixer.io/latest").End()
		if err != nil {
			log.Println("NewRate", err)
			return nil
		}
		var rate FXRate
		errs := json.Unmarshal([]byte(body), &rate)
		if errs != nil {
			log.Println("NewRate", errs)
			return nil
		}
		rateInstance = &rate
	}
	return rateInstance
}

func (m *FXRate) Convert(src string, dist string, num float64) float64 {
	rate, has := m.Rates[dist]
	rate1, has1 := m.Rates[src]
	if has && (has1 || src == m.Base) {
		if src == m.Base {
			return num * rate
		} else {
			return num / rate1 * rate
		}
	} else {
		log.Println("rate.Convert", "原始或目标汇率不支持")
	}
	return num
}
