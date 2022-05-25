package main

import (
	"github.com/jinzhu/gorm"
	"gorm.io/driver/sqlite"
)

const (
	host     = "192.168.30.131"
	port     = 3306
	username = "admin"
	password = "admin2022.post"
	dbname   = "mystate"
)

type Product struct {
	gorm.Model
	Code  string
	Price uint
}

var (
	DSN = "admin:admin2022.post@tcp(192.168.30.131:3306)/mystate?multiStatements=true&allowNativePasswords=false&checkConnLiveness=true&maxAllowedPacket=0&charset=utf8&parseTime=True&loc=Local&timeout=10s&readTimeout=30s&writeTimeout=60s"
)

func main() {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database.")
	}

	// 迁移 schema
	db.AutoMigrate(&Product{})

	// Create
	db.Create(&Product{Code: "D42", Price: 100})

	// Read
	var product Product
	db.First(&product, 1)                 // 根据整型 主键查找
	db.First(&product, "code = ?", "D42") // 查找code 字段值为 D42 的记录

	// Update - 将product 的price 更新为200
	db.Model(&product).Update("Price", 200)

	// Update 更新多个字段
	db.Model(&product).Update(Product{Price: 200, Code: "F42"}) // 仅更新非零字段
	db.Model(&product).Update(map[string]interface{}{"Price": 200, "Code": "F42"})
	// Delete - 删除 product
	db.Delete(&product, 1)
}
