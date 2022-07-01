package mysql

import "github.com/jinzhu/gorm"

func mysql() {
	db, err := gorm.Open("mysql", "sizu:uuC4W9YyZ6@tcp(am-bp1y98h053h3267bh167320o.ads.aliyuncs.com:3306)/xpx_tms?charset=utf8&parseTime=True&loc=Local")
	defer db.Close()
	if err != nil {
		print(err)
	}
}
