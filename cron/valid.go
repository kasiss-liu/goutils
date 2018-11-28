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
	//如果
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
func ValidExpress(exp string, t time.Time) bool {
	cron, err := NewCronWithExpress(exp)
	if err != nil {
		return false
	}
	return Valid(cron, t)
}

func ValidExpressNow(exp string) bool {
	return ValidExpress(exp, time.Now())
}

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

func printDebug() {
	if validDebug && validFaildField != "" {
		log.Println(validFaildField + " valid failed")
	}
}
