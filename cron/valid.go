package cron

import (
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"
)

//调试开关
var validDebug = false

//验证调试字段名称
var validFaildField = ""

//验证时间点是否符合表达式值
func Valid(b *Cron, now time.Time) bool {
	//debug下打印失败返回
	defer printDebug()
	if err := b.ValidExpress(); err != nil {
		return false
	}
	//如果cron为秒级 判断秒数
	if b.isSec {
		if !validSecond(b.Second, now) {
			return false
		}
	}
	if !validMinute(b.Minute, now) {
		return false
	}
	if !validHour(b.Hour, now) {
		return false
	}
	if !validDom(b.DayOfMonth, now) {
		return false
	}
	if !validMonth(b.Month, now) {
		return false
	}
	if !validDow(b.DayOfWeek, now) {
		return false
	}
	if !validYear(b.Year, now) {
		return false
	}
	return true
}

//验证一个时间点 是否符合cron表达式
func ValidExpress(exp string, t time.Time) bool {
	cron, err := NewCronWithExpress(exp)
	if err != nil {
		return false
	}
	return Valid(cron, t)
}

//判断now是否符合cron表达式
func ValidExpressNow(exp string) bool {
	return ValidExpress(exp, time.Now())
}

//验证秒
func validSecond(exp string, t time.Time) bool {
	reg := regexp.MustCompile(secPattern)
	exps := strings.Split(exp, ",")
	s := t.Second()
	for _, e := range exps {
		matches := reg.FindStringSubmatch(e)
		if len(matches) > 0 {
			exp := matches[1]
			scopes := strings.Split(exp, "-")
			min := 0
			max := 59
			if scopes[0] != "*" {
				min, _ = strconv.Atoi(scopes[0])
			}
			if len(scopes) > 1 {
				max, _ = strconv.Atoi(scopes[1])
			}
			step, _ := strconv.Atoi(matches[6])
			if validPart(min, max, step, s) {
				return true
			}
		}
	}
	validFaildField = "Second"
	return false
}

//验证分钟
func validMinute(exp string, t time.Time) bool {
	reg := regexp.MustCompile(minPattern)
	exps := strings.Split(exp, ",")
	s := t.Minute()
	for _, e := range exps {
		matches := reg.FindStringSubmatch(e)
		if len(matches) > 0 {
			exp := matches[1]
			scopes := strings.Split(exp, "-")
			min := 0
			max := 59
			if scopes[0] != "*" {
				min, _ = strconv.Atoi(scopes[0])
			}
			if len(scopes) > 1 {
				max, _ = strconv.Atoi(scopes[1])
			}
			step, _ := strconv.Atoi(matches[6])
			if validPart(min, max, step, s) {
				return true
			}
		}
	}
	validFaildField = "Minute"
	return false
}

//验证小时
func validHour(exp string, t time.Time) bool {
	reg := regexp.MustCompile(hourPattern)
	exps := strings.Split(exp, ",")
	s := t.Hour()
	for _, e := range exps {
		matches := reg.FindStringSubmatch(e)
		if len(matches) > 0 {
			exp := matches[1]
			scopes := strings.Split(exp, "-")
			min := 0
			max := 23
			if scopes[0] != "*" {
				min, _ = strconv.Atoi(scopes[0])
			}
			if len(scopes) > 1 {
				max, _ = strconv.Atoi(scopes[1])
			}
			step, _ := strconv.Atoi(matches[8])
			if validPart(min, max, step, s) {
				return true
			}
		}
	}
	validFaildField = "Hour"
	return false
}

//验证日期
func validDom(exp string, t time.Time) bool {
	if strings.Contains(exp, "?") {
		return true
	}
	days := calDaysOfMonth(t)
	reg := regexp.MustCompile(domPattern)
	exps := strings.Split(exp, ",")
	s := t.Day()
	for _, e := range exps {
		matches := reg.FindStringSubmatch(e)
		if len(matches) > 0 {
			exp := matches[1]
			scopes := strings.Split(exp, "-")
			min := 1
			max := days
			if scopes[0] != "*" {
				min, _ = strconv.Atoi(scopes[0])
			}
			if len(scopes) > 1 {
				max, _ = strconv.Atoi(scopes[1])
			}
			step, _ := strconv.Atoi(matches[10])
			if validPart(min, max, step, s) {
				return true
			}
		}
	}
	validFaildField = "Dom"
	return false
}

//验证月份
func validMonth(exp string, t time.Time) bool {
	reg := regexp.MustCompile(monthPattern)
	exps := strings.Split(exp, ",")
	s := int(t.Month())
	for _, e := range exps {
		matches := reg.FindStringSubmatch(e)
		if len(matches) > 0 {
			exp := matches[1]
			scopes := strings.Split(exp, "-")
			min := 1
			max := 12
			if scopes[0] != "*" {
				min, _ = strconv.Atoi(scopes[0])
			}
			if len(scopes) > 1 {
				max, _ = strconv.Atoi(scopes[1])
			}
			step, _ := strconv.Atoi(matches[8])
			if validPart(min, max, step, s) {
				return true
			}
		}
	}
	validFaildField = "Month"
	return false
}

//验证星期
func validDow(exp string, t time.Time) bool {

	if strings.Contains(exp, "?") {
		return true
	}
	reg := regexp.MustCompile(dowPattern)
	exps := strings.Split(exp, ",")
	s := int(t.Weekday()) + 1
	if s == 7 {
		s = 0
	}
	for _, e := range exps {
		matches := reg.FindStringSubmatch(e)
		if len(matches) > 0 {
			exp := matches[1]
			scopes := strings.Split(exp, "-")
			min := 0
			max := 6
			if scopes[0] != "*" {
				min, _ = strconv.Atoi(scopes[0])
			}
			if len(scopes) > 1 {
				max, _ = strconv.Atoi(scopes[1])
			}
			step, _ := strconv.Atoi(matches[8])
			if validPart(min, max, step, s) {
				return true
			}
		}
	}
	validFaildField = "Dow"
	return false

}

//验证年
func validYear(exp string, t time.Time) bool {
	if exp == "" {
		return true
	}
	reg := regexp.MustCompile(dowPattern)
	exps := strings.Split(exp, ",")
	s := t.Year()
	for _, e := range exps {
		matches := reg.FindStringSubmatch(e)
		if len(matches) > 0 {
			exp := matches[1]
			scopes := strings.Split(exp, "-")
			min := 1970
			max := 2099
			if scopes[0] != "*" {
				min, _ = strconv.Atoi(scopes[0])
			}
			if len(scopes) > 1 {
				max, _ = strconv.Atoi(scopes[1])
			}
			step, _ := strconv.Atoi(matches[8])
			if validPart(min, max, step, s) {
				return true
			}
		}
	}
	validFaildField = "Year"
	return false
}

//验证通用步进表达式
func validPart(min, max, step, time int) bool {
	if step == 0 {
		step = 1
	}
	for i := min; i <= max; i += step {
		if time == i {
			return true
		}
	}
	return false
}

//计算时间所在的月份天数
func calDaysOfMonth(t time.Time) int {
	maxdef := 31
	switch t.Month() {
	case time.January, time.March, time.May, time.July, time.August, time.October, time.December:
		maxdef = 31
	case time.February:
		if t.YearDay() > 365 {
			maxdef = 29
		} else {
			maxdef = 28
		}
	default:
		maxdef = 30
	}
	return maxdef
}

//debug 打印不符合要求的条件
func printDebug() {
	if validDebug && validFaildField != "" {
		log.Println(validFaildField + " valid failed")
	}
}
