package databases

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"log"
	"v-blog/config"
)

var DB *gorm.DB

func New() *gorm.DB {
	conf, err := config.Get()
	if err != nil {
		log.Fatalf(err.Error(), "config has not loaded")
	}
	connectConf := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True&loc=Local",
		conf.Db.User,
		conf.Db.Password,
		conf.Db.Host,
		conf.Db.Port,
		conf.Db.Database,
		conf.Db.Charset)

	db, err := gorm.Open("mysql", connectConf)
	if err != nil {
		log.Fatalf(err.Error(), "mysql connect error. msg: %s \n")
	}

	db.LogMode(true)

	return db
}

func InitDatabase() {
	DB = New()
}
