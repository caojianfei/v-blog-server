package databases

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/spf13/viper"
	"log"
)

var DB *gorm.DB

func New() *gorm.DB {
	dbConf := viper.GetStringMapString("db")
	connectConf := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True&loc=Local",
		dbConf["user"],
		dbConf["password"],
		dbConf["host"],
		dbConf["port"],
		dbConf["database"],
		dbConf["charset"])

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
