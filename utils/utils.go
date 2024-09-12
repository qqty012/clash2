package utils

import (
	"fmt"
	"github.com/qqty012/clash2/constant"
	"os"
	"time"
)

func LogWriteToFile(format string, v ...any) {
	home := constant.Path.HomeDir()
	file, err := os.OpenFile(fmt.Sprintf("%s/clash.log", home), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {

		}
	}(file)

	_, err = file.WriteString(fmt.Sprintf(format, v...) + "\n")
	if err != nil {
		return
	}
}

func E8Time(t time.Time) string {
	location, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		return t.Format("2006-01-02 15:04:05")
	}
	return t.In(location).Format("2006-01-02 15:04:05")
}

func NowE8Time() string {

	return E8Time(time.Now())
}
