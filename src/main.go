package main

import (
	"bytes"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
)

var childServices = strings.Split(os.Getenv("CHILD_SERVICE"), ",")

var data = make(map[string]string)
var mutex = sync.Mutex{}

func main() {
	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	r.POST("/:id", func(c *gin.Context) {
		id := c.Param("id")
		mutex.Lock()
		data[id] = id
		mutex.Unlock()
		c.Status(http.StatusOK)
	})

	r.POST("/:id", func(c *gin.Context) {
		id := c.Param("id")
		if len(childServices) > 0 {
			body := generatePostBody(childServices, id)
			resp, err := http.Post("http://coordinator-baseline.yac.svc.cluster.local:8080/saga", "application/json", bytes.NewBuffer(body))
			if err == nil && resp.StatusCode >= 200 && resp.StatusCode <= 299 {
				c.Status(http.StatusOK)
			} else if err != nil {
				c.Status(http.StatusInternalServerError)
			} else {
				c.Status(resp.StatusCode)
			}
		} else {
			id := c.Param("id")
			mutex.Lock()
			data[id] = id
			mutex.Unlock()
			c.Status(http.StatusOK)
		}
	})

	r.DELETE("/:id", func(c *gin.Context) {
		id := c.Param("id")
		if len(childServices) > 0 {
			body := generateDeleteBody(childServices, id)
			resp, err := http.Post("http://coordinator-baseline.yac.svc.cluster.local:8080/saga", "application/json", bytes.NewBuffer(body))
			if err == nil && resp.StatusCode >= 200 && resp.StatusCode <= 299 {
				c.Status(http.StatusOK)
			} else if err != nil {
				c.Status(http.StatusInternalServerError)
			} else {
				c.Status(resp.StatusCode)
			}
		} else {
			mutex.Lock()
			if _, isIn := data[id]; isIn {
				delete(data, id)
				mutex.Unlock()
				c.Status(http.StatusOK)
			} else {
				mutex.Unlock()
				c.Status(http.StatusBadRequest)
			}
		}
	})

	if err := r.Run(":8888"); err != nil {
		panic(err)
	}
}
