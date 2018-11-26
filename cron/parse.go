package cron

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
)

const (
	multiBlankReg = `[\s|\t]+`
)

const (
	secPattern   = `^((\*)|([0-5]?\d(\-[0-5]?\d)*))(\/(\d+))*$`
	minPattern   = `^((\*)|([0-5]?\d(\-[0-5]?\d)*))(\/(\d+))*$`
	hourPattern  = `^((\*)|(([01]?\d|[2][0-3])(\-([01]?\d|[2][0-3]))*))(\/(\d+))*$`
	domPattern   = `(^\?$)|(^((\*)|(([0]?[1-9]|[12]?\d|[3][01])(\-([0]?[1-9]|[12]?\d|[3][01]))*))(\/(\d+))*$)`
	monthPattern = `^((\*)|(([0]?[1-9]|[1][012])(\-([0]?[1-9]|[1][012]))*))(\/(\d+))*$`
	dowPattern   = `(^\?$)|(^((\*)|(([0]?[0-6])(\-([0]?[0-6]))*))(\/\d+)*)$`
	yearPattern  = `^((\*)|(([2]\d{3})(\-([2]\d{3}))*))(\/\d+)*$`
)

const (
	ptypeSec  = "second"
	ptypeMin  = "minute"
	ptypeHour = "hour"
	ptypeDom  = "dom"
	ptypeMon  = "month"
	ptypeDow  = "dow"
	ptypeYear = "year"
)

var patternMap = map[string]string{
	ptypeSec:  secPattern,
	ptypeMin:  minPattern,
	ptypeHour: hourPattern,
	ptypeDom:  domPattern,
	ptypeMon:  monthPattern,
	ptypeDow:  dowPattern,
	ptypeYear: yearPattern,
}

func Parse(express string) (*Builder, error) {
	//多个空格处理
	reg := regexp.MustCompile(multiBlankReg)
	express = reg.ReplaceAllString(express, " ")
	//拆解成数组切片,并填入builder中
	es := strings.Split(express, " ")
	cron := NewCron()
	switch len(es) {
	case 5:
		cron.SetMinute(es[0])
		cron.SetHour(es[1])
		cron.SetDayOfMonth(es[2])
		cron.SetMonth(es[3])
		cron.SetDayOfWeek(es[4])
	case 6:
		cron.SetMinute(es[0])
		cron.SetHour(es[1])
		cron.SetDayOfMonth(es[2])
		cron.SetMonth(es[3])
		cron.SetDayOfWeek(es[4])
		cron.SetYear(es[5])
	case 7:
		cron.SetSecond(es[0])
		cron.SetMinute(es[1])
		cron.SetHour(es[2])
		cron.SetDayOfMonth(es[3])
		cron.SetMonth(es[4])
		cron.SetDayOfWeek(es[5])
		cron.SetYear(es[6])
	default:
		return nil, errors.New("parse error: illegal element count " + strconv.Itoa(len(es)))
	}
	if err := cron.Valid(); err != nil {
		return nil, err
	}
	return cron, nil

}

func validExpress(ptype string, express []string) error {
	var pattern string
	var ok bool
	if pattern, ok = patternMap[ptype]; !ok {
		return errors.New("match pattern type error :" + ptype)
	}
	if len(express) == 0 {
		return errors.New("express can not be empty")
	}

	reg := regexp.MustCompile(pattern)
	for _, v := range express {
		v = strings.TrimSpace(v)
		if !reg.MatchString(v) {
			return errors.New(ptype + " express invalid : `" + v + "`")
		}
	}
	return nil
}
