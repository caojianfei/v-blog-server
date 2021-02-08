package databases

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	log "github.com/sirupsen/logrus"
	"v-blog/config"
)

var DB *gorm.DB

func New() *gorm.DB {
	conf, err := config.Get()
	if err != nil {
		log.Fatalf("config has not loaded. err: %s", err)
	}
	connectConf := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True&loc=Local",
		conf.Mysql.User,
		conf.Mysql.Password,
		conf.Mysql.Host,
		conf.Mysql.Port,
		conf.Mysql.Database,
		conf.Mysql.Charset)

	db, err := gorm.Open("mysql", connectConf)
	if err != nil {
		log.Fatalf("mysql connect error. msg: %s", err)
	}

	if conf.App.Env != "release" {
		db.LogMode(true)
	}

	return db
}

func InitDatabase() {
	DB = New()
}
