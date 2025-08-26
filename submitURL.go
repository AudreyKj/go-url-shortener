package main 

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"crypto/sha1"
	"encoding/binary"
	//"github.com/mattheath/base62"
	"fmt"
)

var m = make(map[string]string)

type UrlRequest struct {
	url string `json:"url"`
}

func shortHash(s string) uint64 {
	hash := sha1.Sum([]byte(s))
	num := binary.BigEndian.Uint64(hash[:8])
	//fmt.Println(base62.EncodeInt64(int64(num)))
	return num
}

func submitURL(c *gin.Context) {
	var req UrlRequest

	if err := c.BindJSON(&req); err != nil{
		return 
	}

	encoded := shortHash(req.url)

	fmt.Println("encoded", encoded)

	c.JSON(http.StatusOK, gin.H{"url": encoded})
}

