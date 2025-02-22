package main

import (
	"fmt"
	"net/http"
	"strconv"

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
	r.POST("/products", addProducts)
	r.DELETE("/products/:name", deleteByProductName)
	r.GET("/products/:id", getProductsById)
	r.PUT("products/:id", updateProduct)

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
func addProducts(c *gin.Context) {
	var product Product

	// JSON'u Product modeline çevir
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Geçersiz JSON"})
		return
	}

	// Veritabanına ekle
	result := db.Create(&product)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Veritabanına eklenirken hata oluştu"})
		return
	}

	// Başarıyla eklenen veriyi döndür
	c.JSON(http.StatusCreated, product)
}
func deleteByProductName(c *gin.Context) {
	var prs Product
	prs.Name = c.Param("name")

	res := db.Where("name = ?", prs.Name).Delete(&Product{})
	if res.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "GEÇERSİZ İSTEK"})
	}
	if res.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": fmt.Sprintf("%s adında bir ürün bulunamadı", prs.Name)})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("%s adındaki ürün başarıyla silindi", prs.Name)})

}
func getProductsById(c *gin.Context) {
	var prdct Product
	id64, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(400, gin.H{"error": "Geçersiz ID formatı"})
		return
	}
	// Tek bir kaydı almak için first yoksa find
	result := db.Where("id = ?", uint(id64)).First(&prdct)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {

			c.JSON(http.StatusBadRequest, gin.H{"error": "Kayıt Bulunamadı!"})
		} else {

			c.JSON(http.StatusBadRequest, gin.H{"error": "Veri Tabanına Bağlanılamadı"})
		}
	} else {
		c.JSON(http.StatusOK, &prdct)
	}

}
func updateProduct(c *gin.Context) {
	var prdct Product
	var prds Product
	id64, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(400, gin.H{"error": "Geçersiz ID formatı"})
		return
	}
	result := db.Where("id = ?", uint(id64)).First(&prdct)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {

			c.JSON(http.StatusBadRequest, gin.H{"error": "Kayıt Bulunamadı!"})
		} else {

			c.JSON(http.StatusBadRequest, gin.H{"error": "Veri Tabanına Bağlanılamadı"})
		}
	} else {
		if err := c.ShouldBindJSON(&prds); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Geçersiz JSON"})
			return
		}
		prdct.Name = prds.Name
		prdct.Price = prds.Price
		db.Save(&prdct)
		c.JSON(http.StatusOK, &prds)
	}

}
