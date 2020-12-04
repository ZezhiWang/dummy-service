package main

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"sync"
)

var childServices  = []string {
	"dummy2",
	"dummy3",
	"dummy4",
	"dummy5",
}

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

		resp, err := http.Post(os.Getenv("sagas"), "application/json", bytes.NewBuffer(body))
		if err == nil && resp.StatusCode >= 200 && resp.StatusCode <= 299 {
			c.Status(http.StatusOK)
		} else if err != nil {
			c.Status(http.StatusInternalServerError)
		} else {
			c.Status(resp.StatusCode)
		}
	})

	if err := r.Run(":8080"); err != nil {
		panic(err)
	}
}
