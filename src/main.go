
package main

import (
	"net/http"
	"os"
	"io/ioutil"
	"log"
	"fmt"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

type Slideshow struct {
	Name     string    `form:"name"`
}

func getFolder(path string) []string {
	pwd, _ := os.Getwd()

	fileInfoList,err := ioutil.ReadDir(filepath.Join(pwd, path))
	if err != nil {
		log.Fatal(err)
	}

	var folder []string

	for i := range fileInfoList {
		folder = append(folder, fileInfoList[i].Name())
	}

	return folder
}

func getImages(path string) []string {
	pwd, _ := os.Getwd()

	fileInfoList,err := ioutil.ReadDir(filepath.Join(pwd, path))
	if err != nil {
		log.Fatal(err)
	}

	var images []string

	for i := range fileInfoList {
		uri := path + "/" + fileInfoList[i].Name()
		images = append(images, uri)
	}

	return images
}

func main() {
	r := gin.Default()
	r.Static("/assets", "assets")
	r.LoadHTMLGlob("templates/*")

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"title": "Slidsshow",
		})
	})

	r.GET("/api/folder", func(c *gin.Context) {
		data := getFolder("/assets/images")
		c.AsciiJSON(http.StatusOK, data)
	})

	r.GET("/api/images", func(c *gin.Context) {
		var slideshow Slideshow
		if c.ShouldBind(&slideshow) == nil {
			fmt.Println(slideshow.Name)

			path := fmt.Sprintf("/assets/images/%s", slideshow.Name)
			data := getImages(path)
			c.AsciiJSON(http.StatusOK, data)
		} else {
			c.AsciiJSON(http.StatusBadRequest, gin.H{
				"error": "file does not exist!",
			})
		}
	})

	r.RunTLS(":8085", "./server.crt", "./server.key")
}
