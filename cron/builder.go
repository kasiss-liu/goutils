package cron

import (
	"errors"
	"strings"
)

//Cron 结构体
type Cron struct {
	Second     string
	Minute     string
	Hour       string
	DayOfMonth string
	Month      string
	DayOfWeek  string
	Year       string
	isSec      bool
}

//构建一个新的空 cron  秒级默认关闭
func NewCron() *Cron {
	return &Cron{isSec: false}
}

//根据表达式创建一个cron结构
func NewCronWithExpress(exp string) (cron *Cron, err error) {
	//解析表达式
	cron, err = Parse(exp)
	if err != nil {
		return
	}
	//表达式项值验证
	err = cron.ValidExpress()
	if err != nil {
		return
	}
	return
}

//设置秒表达式
func (b *Cron) SetSecond(s string) *Cron {
	b.isSec = true
	b.Second = s
	return b
}

//设置是否开启秒级锁
func (b *Cron) SetIsSec(i bool) *Cron {
	b.isSec = i
	return b
}

//返回是否是秒级cron
func (b *Cron) IsSec() bool {
	return b.isSec
}

//设置分钟表达式
func (b *Cron) SetMinute(m string) *Cron {
	b.Minute = m
	return b
}

//设置小时表达式
func (b *Cron) SetHour(h string) *Cron {
	b.Hour = h
	return b
}

//设置每月天表达式
func (b *Cron) SetDayOfMonth(d string) *Cron {
	b.DayOfMonth = d
	return b
}

//设置月份表达式
func (b *Cron) SetMonth(m string) *Cron {
	b.Month = m
	return b
}

//设置星期表达式
func (b *Cron) SetDayOfWeek(dow string) *Cron {
	b.DayOfWeek = dow
	return b
}

//设置年表达式
func (b *Cron) SetYear(y string) *Cron {
	if y == "" {
		y = "*"
	}
	b.Year = y
	return b
}

//将cron转化为字符串表达式
func (b *Cron) ToExpress() string {
	b.fillAttr()
	s := b.getExpressSlice()
	return strings.Join(s, " ")
}

//拼凑表达式切片
func (b *Cron) getExpressSlice() []string {
	var s = make([]string, 0, 10)
	if b.isSec {
		s = append(s, b.Second)
	}
	s = append(s, b.Minute)
	s = append(s, b.Hour)
	s = append(s, b.DayOfMonth)
	s = append(s, b.Month)
	s = append(s, b.DayOfWeek)
	s = append(s, b.Year)
	return s
}

//自动补齐cron空缺项 默认设置为*
func (b *Cron) fillAttr() {
	if b.isSec {
		if b.Second == "" {
			b.Second = "*"
		}
	}
	if b.Minute == "" {
		b.Minute = "*"
	}
	if b.Hour == "" {
		b.Hour = "*"
	}
	if b.DayOfMonth == "" {
		b.DayOfMonth = "*"
	}
	if b.Month == "" {
		b.Month = "*"
	}
	if b.DayOfWeek == "" {
		b.DayOfWeek = "*"
	}
	if b.Year == "" {
		b.Year = "*"
	}
}

//验证表达式是否正确
func (b *Cron) ValidExpress() error {
	var err error

	if b.isSec {
		err = b.validSecond()
		if err != nil {
			return err
		}
	}

	err = b.validDDConflict()
	if err != nil {
		return err
	}

	err = b.validMinute()
	if err != nil {
		return err
	}
	err = b.validHour()
	if err != nil {
		return err
	}
	err = b.validDom()
	if err != nil {
		return err
	}
	err = b.validMonth()
	if err != nil {
		return err
	}
	err = b.validDow()
	if err != nil {
		return err
	}
	err = b.validYear()
	if err != nil {
		return err
	}

	return nil
}

//验证秒表达式
func (b *Cron) validSecond() error {
	sl := strings.Split(b.Second, ",")
	return validExpress(ptypeSec, sl)
}

//验证分钟表达式
func (b *Cron) validMinute() error {
	sl := strings.Split(b.Minute, ",")
	return validExpress(ptypeMin, sl)
}

//验证小时表达式
func (b *Cron) validHour() error {
	sl := strings.Split(b.Hour, ",")
	return validExpress(ptypeHour, sl)
}

//验证天表达式
func (b *Cron) validDom() error {
	sl := strings.Split(b.DayOfMonth, ",")
	return validExpress(ptypeDom, sl)
}

//验证月份表达式
func (b *Cron) validMonth() error {
	sl := strings.Split(b.Month, ",")
	return validExpress(ptypeMon, sl)
}

//验证星期表达式
func (b *Cron) validDow() error {
	sl := strings.Split(b.DayOfWeek, ",")
	return validExpress(ptypeDow, sl)
}

//验证年表达式
func (b *Cron) validYear() error {
	if b.Year == "" {
		return nil
	}
	sl := strings.Split(b.Year, ",")
	return validExpress(ptypeYear, sl)
}

//校验天与星期的冲突
func (b *Cron) validDDConflict() error {
	if b.DayOfMonth == "?" && b.DayOfWeek == "?" {
		return errors.New("dow and dom can not be ? in one express")
	}
	if b.DayOfMonth == "*" && (b.DayOfWeek != "*" && b.DayOfWeek != "?") {
		return errors.New("dow and dom can not be conflict")
	}
	if (b.DayOfMonth != "*" && b.DayOfMonth != "?") && b.DayOfWeek == "*" {
		return errors.New("dow and dom can not be conflict")
	}
	return nil
}
