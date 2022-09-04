package main

import (
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Product struct {
	gorm.Model

	Name string

	Views uint32

	Amount uint32

	Folder string

	ImageNameList string
}

type ProductWrapper struct {
	Id uint

	Name string

	Views uint32

	Amount uint32

	FirstImagePath string
}

type ProductDetails struct {
	Name string

	Views uint32

	Amount uint32

	ImageInfoList []map[string]string
}

type SlideshowUri struct {
	Id uint64 `uri:"id" binding:"required,uuid"`
}

func mock() {
	db, err := gorm.Open(sqlite.Open("slideshow.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&Product{})

	db.Create(&Product{
		Name:          "动漫图片",
		Views:         100,
		Amount:        3,
		Folder:        "/assets/samples",
		ImageNameList: "01.jpeg,02.jpeg,03.jpeg",
	})
}

func listImages() []ProductWrapper {
	db, err := gorm.Open(sqlite.Open("slideshow.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
		return []ProductWrapper{}
	}

	db.AutoMigrate(&Product{})

	var products []Product
	result := db.Find(&products)
	fmt.Println("RowsAffected: ", result.RowsAffected, " Error: ", result.Error)

	var productWrapperList []ProductWrapper
	for _, product := range products {
		imageNameList := strings.Split(product.ImageNameList, ",")
		if len(imageNameList) != 0 {
			var productWrapper ProductWrapper
			productWrapper.Id = product.ID
			productWrapper.Name = product.Name
			productWrapper.Views = product.Views
			productWrapper.Amount = product.Amount
			productWrapper.FirstImagePath = product.Folder + "/" + imageNameList[0]
			productWrapperList = append(productWrapperList, productWrapper)
		}
	}
	return productWrapperList
}

func getImage(id uint64) ProductDetails {
	db, err := gorm.Open(sqlite.Open("slideshow.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
		return ProductDetails{}
	}

	db.AutoMigrate(&Product{})

	var product Product
	db.First(&product, id)

	imageNameList := strings.Split(product.ImageNameList, ",")

	var imageInfoList []map[string]string
	for _, imageName := range imageNameList {
		imagePath := product.Folder + "/" + imageName
		file, _ := os.Open("." + imagePath)
		c, _, err := image.DecodeConfig(file)
		if err != nil {
			fmt.Println("get decode config fail!", err)
			continue
		}

		imageInfo := make(map[string]string)
		imageInfo["imagePath"] = imagePath
		imageInfo["width"] = strconv.Itoa(c.Width)
		imageInfo["height"] = strconv.Itoa(c.Height)
		imageInfoList = append(imageInfoList, imageInfo)
	}

	var productDetails ProductDetails
	productDetails.Name = product.Name
	productDetails.Views = product.Views
	productDetails.Amount = product.Amount
	productDetails.ImageInfoList = imageInfoList

	return productDetails
}

func main() {

	// mock()

	r := gin.Default()
	r.Static("/assets", "assets")
	r.LoadHTMLGlob("templates/*")

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"title": "Slidsshow",
		})
	})

	r.GET("/thread/:id", func(c *gin.Context) {
		var slideshowUri SlideshowUri
		if c.ShouldBindUri(&slideshowUri) != nil {
			fmt.Println(slideshowUri.Id)

			productDetails := getImage(slideshowUri.Id)

			c.HTML(http.StatusOK, "thread.html", gin.H{
				"title":  "Slidsshow",
				"name":   productDetails.Name,
				"views":  productDetails.Views,
				"amount": productDetails.Amount,
			})
		}
	})

	r.GET("/api/thread/:id", func(c *gin.Context) {
		var slideshowUri SlideshowUri
		if c.ShouldBindUri(&slideshowUri) != nil {
			fmt.Println(slideshowUri.Id)

			productDetails := getImage(slideshowUri.Id)

			if len(productDetails.ImageInfoList) != 0 {
				c.AsciiJSON(http.StatusOK, productDetails.ImageInfoList)
			}
		}
	})

	r.GET("/api/images", func(c *gin.Context) {
		data := listImages()
		if len(data) != 0 {
			c.AsciiJSON(http.StatusOK, data)
		}
	})

	r.RunTLS(":8085", "./server.crt", "./server.key")
}
