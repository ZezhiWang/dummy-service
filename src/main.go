package main

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"strings"
	"sync"
)

var childServices = strings.Split(os.Getenv("CHILD_SERVICE"), ",")

var data = make(map[string]string)
var mutex = sync.Mutex{}

func main() {
	r := gin.Default()

	r.POST("/base/:id", func(c *gin.Context) {
		id := c.Param("id")
		mutex.Lock()
		data[id] = id
		mutex.Unlock()
		c.Status(http.StatusOK)
	})

	r.DELETE("/base/:id", func(c *gin.Context) {
		id := c.Param("id")
		mutex.Lock()
		if _, isIn := data[id]; isIn {
			delete(data, id)
			mutex.Unlock()
			c.Status(http.StatusOK)
		} else {
			mutex.Unlock()
			c.Status(http.StatusBadRequest)
		}
	})

	r.POST("/saga/:id", func(c *gin.Context) {
		id := c.Param("id")
		body := generatePostBody(childServices, id)

		resp, err := http.Post("http://coordinator.yac.svc.cluster.local:8080/saga", "application/json", bytes.NewBuffer(body))
		if err == nil && resp.StatusCode >= 200 && resp.StatusCode <= 299 {
			c.Status(http.StatusOK)
		} else if err != nil {
			c.Status(http.StatusInternalServerError)
		} else {
			c.Status(resp.StatusCode)
		}
	})

	r.DELETE("/saga/:id", func(c *gin.Context) {
		id := c.Param("id")
		body := generateDeleteBody(childServices, id)

		resp, err := http.Post("http://coordinator.yac.svc.cluster.local:8080/saga", "application/json", bytes.NewBuffer(body))
		if err == nil && resp.StatusCode >= 200 && resp.StatusCode <= 299 {
			c.Status(http.StatusOK)
		} else if err != nil {
			c.Status(http.StatusInternalServerError)
		} else {
			c.Status(resp.StatusCode)
		}
	})

	if err := r.Run(":8888"); err != nil {
		panic(err)
	}
}
