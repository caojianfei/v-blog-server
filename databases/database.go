package databases

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/lexkong/log"
)

var DB *gorm.DB

func init() {
	var err error
	DB, err = gorm.Open("mysql", "caojianfei:Caojf@1910@tcp(47.97.196.203)/v-blog?charset=utf8mb4&parseTime=True&loc=Local")
	if err != nil {
		log.Fatalf(err, "mysql connect error. msg: %s")
	}
}
