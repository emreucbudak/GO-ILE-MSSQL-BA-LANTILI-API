package main

import (
	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

type Product struct {
	ID    uint    `json:"id" gorm:"primaryKey"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

var db *gorm.DB

func initDB() {
	// MSSQL bağlantı dizesi
	dsn := "sqlserver://sa:emreraftongame63@emre:1433?database=GoDeneme"
	var err error
	db, err = gorm.Open(sqlserver.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Veritabanına bağlanılamadı: " + err.Error())
	}

	// Tabloyu otomatik oluştur (migration)
	db.AutoMigrate(&Product{})
}

func main() {

	initDB()

	r := gin.Default()

	r.GET("/products", getProducts)

	r.Run(":8080")
}
func getProducts(c *gin.Context) {
	var products []Product
	if err := db.Find(&products).Error; err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, products)
}
