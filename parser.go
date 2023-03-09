package main

import (
	"fmt"

	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/net/html"
)

func main() {
	r := gin.Default()
	gin.SetMode(gin.ReleaseMode)
	r.GET("/", parseURL)
	r.Run()
}

type Preview struct {
	Img         string `json:"img"`
	Description string `json:"description"`
}

func parseURL(c *gin.Context) {
	url := c.Query("url")

	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
	c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
	c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET")

	resp, err := http.Get(url)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	doc, errr := html.Parse(resp.Body)
	if errr != nil {
		fmt.Println(errr)
	}

	var previewImage, description, title string
	var extractPreviewImage func(*html.Node)
	extractPreviewImage = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "meta" {
			for _, attr := range n.Attr {
				if attr.Key == "property" && attr.Val == "og:image" {
					for _, subAttr := range n.Attr {
						if subAttr.Key == "content" {
							previewImage = subAttr.Val
							return
						}
					}
				}
				if attr.Key == "name" && attr.Val == "description" {
					for _, subAttr := range n.Attr {
						if subAttr.Key == "content" {
							description = subAttr.Val
							return
						}
					}
				}
				if attr.Key == "name" && attr.Val == "twitter:title" {
					for _, subAttr := range n.Attr {
						if subAttr.Key == "content" {
							title = subAttr.Val
							return
						}
					}
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			extractPreviewImage(c)
		}
	}
	extractPreviewImage(doc)
	c.IndentedJSON(http.StatusOK, gin.H{"Values": previewImage + "\n", "Descriptions": description, "Title": title})
}
