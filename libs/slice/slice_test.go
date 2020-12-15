package slice

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"testing"
	"v-blog/models"
)


func createTestStruct() (articles []models.Article) {
	connectConf := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True&loc=Local",
		"root",
		"123456",
		"127.0.0.1",
		"3306",
		"v-blog",
		"utf8mb4")

	db, err := gorm.Open("mysql", connectConf)
	if err != nil {
		panic("database connect error")
	}

	db.LogMode(false)

	articles = []models.Article{}
	db.Model(&models.Article{}).Limit(10).Find(&articles)
	return
}

func BenchmarkColumn(b *testing.B) {
	s := createTestStruct()
	for i := 0; i < b.N; i++ {
		_, _ = ToSlice(s).Column("Title").CovertToString()
	}
}

func BenchmarkColumnNormal(b *testing.B) {
	s := createTestStruct()
	for i := 0; i < b.N; i++ {
		r := make([]string, 0)
		for j := 0; j < len(s); j++ {
			r = append(r, s[j].Title)
		}
	}
}
