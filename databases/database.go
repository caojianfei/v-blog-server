package databases

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/lexkong/log"
)

var DB *gorm.DB

func init() {
	var err error
	DB, err = gorm.Open("mysql", "root:123456@tcp(localhost:3306)/v-blog?charset=utf8mb4&parseTime=True&loc=Local")
	if err != nil {
		log.Fatalf(err, "mysql connect error. msg: %s")
	}
	DB.LogMode(true)
}
