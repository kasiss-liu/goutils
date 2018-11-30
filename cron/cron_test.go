package cron

import (
	"testing"
	"time"
)

func TestBuilder(t *testing.T) {
	// `*/2 * ? * * *`  crontab表达式
	express := NewCron().SetMinute("*/2").SetHour("*").SetDayOfMonth("?").SetMonth("*").SetDayOfWeek("*").SetYear("*").ValidExpress()
	t.Log(express)
	//<nil>
	// `*/2 * ? * ?/3 *` crontab表达式 且?/2项存在错误
	cron2 := NewCron().SetMinute("*/2").SetHour("*").SetDayOfMonth("?").SetMonth("*").SetYear("*").SetDayOfWeek("?/3")
	express2valid := cron2.ValidExpress()
	if express2valid == nil {
		express2 := cron2.ToExpress()
		t.Log(express2)
	}
	t.Log(express2valid)
	//dow express invalid : `?/3`

	cronExpress, err := NewCronWithExpress("1,10,20 */5 10-32 * * ? *")
	if err != nil {
		t.Log(err)
		//hour express invalid : `10-32`
	} else {
		t.Log(cronExpress.ToExpress())
	}

	cronExpress, err = NewCronWithExpress("10,20 */5 10-20 * * ? *")
	if err != nil {
		t.Log(err)
	} else {
		t.Log(cronExpress.ToExpress())
		// 10,20 */5 10-20 * * ? *
	}
	cronExpress.SetIsSec(false)

	t.Log(cronExpress.ToExpress())
	// */5 10-20 * * ? *
}

func TestValid(t *testing.T) {
	validDebug = true
	t.Log(time.Now().Format("2006-01-02 15:04:05"))
	res := ValidExpressNow("* */5 10-20 * * ? *")
	t.Log(res)
	now := time.Date(int(2018), time.November, int(26), int(17), int(20), int(0), int(0), time.Local)
	t.Log(now)
	cron, _ := NewCronWithExpress("* 10-20 * * ?")
	t.Log(cron.ToExpress())
	res = Valid(cron, now)
	t.Log(res)

	res = ValidExpress("* 10-20 * * *", now)
	t.Log(res)
}
