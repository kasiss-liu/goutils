package cron

import (
	"errors"
	"strings"
)

type Builder struct {
	Second     string
	Minute     string
	Hour       string
	DayOfMonth string
	Month      string
	DayOfWeek  string
	Year       string
	isSec      bool
}

func NewCron() *Builder {
	return &Builder{isSec: false}
}
func NewCronWithExpress(exp string) (cron *Builder, err error) {
	cron, err = Parse(exp)
	if err != nil {
		return
	}
	err = cron.Valid()
	if err != nil {
		return
	}
	return
}

func (b *Builder) SetSecond(s string) *Builder {
	b.isSec = true
	b.Second = s
	return b
}
func (b *Builder) SetIsSec(i bool) *Builder {
	b.isSec = i
	return b
}

func (b *Builder) SetMinute(m string) *Builder {
	b.Minute = m
	return b
}

func (b *Builder) SetHour(h string) *Builder {
	b.Hour = h
	return b
}
func (b *Builder) SetDayOfMonth(d string) *Builder {
	b.DayOfMonth = d
	return b
}
func (b *Builder) SetMonth(m string) *Builder {
	b.Month = m
	return b
}
func (b *Builder) SetDayOfWeek(dow string) *Builder {
	b.DayOfWeek = dow
	return b
}
func (b *Builder) SetYear(y string) *Builder {
	b.Year = y
	return b
}
func (b *Builder) ToExpress() string {
	b.fillAttr()
	s := b.getExpressSlice()
	return strings.Join(s, " ")
}

func (b *Builder) getExpressSlice() []string {
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

func (b *Builder) fillAttr() {
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

func (b *Builder) Valid() error {
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
func (b *Builder) validSecond() error {
	sl := strings.Split(b.Second, ",")
	return validExpress(ptypeSec, sl)
}
func (b *Builder) validMinute() error {
	sl := strings.Split(b.Minute, ",")
	return validExpress(ptypeMin, sl)
}
func (b *Builder) validHour() error {
	sl := strings.Split(b.Hour, ",")
	return validExpress(ptypeHour, sl)
}
func (b *Builder) validDom() error {
	sl := strings.Split(b.DayOfMonth, ",")
	return validExpress(ptypeDom, sl)
}
func (b *Builder) validMonth() error {
	sl := strings.Split(b.Month, ",")
	return validExpress(ptypeMon, sl)
}
func (b *Builder) validDow() error {
	sl := strings.Split(b.DayOfWeek, ",")
	return validExpress(ptypeDow, sl)
}
func (b *Builder) validYear() error {
	if b.Year == "" {
		return nil
	}
	sl := strings.Split(b.Year, ",")
	return validExpress(ptypeYear, sl)
}
func (b *Builder) validDDConflict() error {
	if b.DayOfMonth == "?" && b.DayOfWeek == "?" {
		return errors.New("dow and dom can not be ? in one express")
	}
	return nil
}
