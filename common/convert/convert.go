package convert

import (
	"time"

	"github.com/shopspring/decimal"
	"github.com/zeromicro/go-zero/core/logx"
)

const (
	Time_LAYOUT = "2006-01-02 15:04:05"
)

func String2Decimal(s string) decimal.Decimal {
	d, _ := decimal.NewFromString(s)
	return d
}

func StringToTime(s string) time.Time {
	timeObj, err := time.Parse(Time_LAYOUT, s)
	if err != nil {
		logx.Errorw("time.Parse(Time_LAYOUT, s)", logx.LogField{Key: "Error: ", Value: err.Error()})
		return time.Time{}
	}
	return timeObj
}

func TimeToString(t time.Time) string {
	return t.Format(Time_LAYOUT)
}
