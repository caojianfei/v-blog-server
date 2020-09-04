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
		conf.Db.User,
		conf.Db.Password,
		conf.Db.Host,
		conf.Db.Port,
		conf.Db.Database,
		conf.Db.Charset)

	db, err := gorm.Open("mysql", connectConf)
	if err != nil {
		log.Fatalf("mysql connect error. msg: %s", err)
	}

	if conf.AppEnv != "release" {
		db.LogMode(true)
	}

	return db
}

func InitDatabase() {
	DB = New()
}
