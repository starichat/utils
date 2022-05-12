package gorm

/**
实现gorm2。0实现 gorm批量更新的操作
 */

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"reflect"
)

func init() {
	db, err := gorm.Open(mysql.Open(""), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db.Set("on duplicate key upadte values(updat_time)").CreateInBatches()

}



