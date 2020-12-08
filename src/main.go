package main

import (
	"bytes"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
)

const responseDelay = 5

var childServices = strings.Split(os.Getenv("CHILD_SERVICE"), ",")

var data = make(map[string]string)
var mutex = sync.Mutex{}

const coordinatorAddr = "http://coordinator.yac-baseline.svc.cluster.local:8080/saga/"

func main() {
	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	r.POST("/:id", func(c *gin.Context) {
		id := c.Param("id")
		time.Sleep(responseDelay * time.Millisecond)
		if len(childServices) == 0 || (len(childServices) == 1 && childServices[0] == "") {
			id := c.Param("id")
			mutex.Lock()
			data[id] = id
			mutex.Unlock()
			c.Status(http.StatusOK)
		} else {
			body := generatePostBody(childServices, id)
			reqId := xid.New().String()
			resp, err := http.Post(coordinatorAddr + reqId, "application/json", bytes.NewBuffer(body))
			if err == nil && resp.StatusCode >= 200 && resp.StatusCode <= 299 {
				c.Status(http.StatusOK)
			} else if err != nil {
				log.Println(err)
				c.Status(http.StatusInternalServerError)
			} else {
				c.Status(resp.StatusCode)
			}
		}
	})

	r.DELETE("/:id", func(c *gin.Context) {
		id := c.Param("id")
		if len(childServices) == 0 || (len(childServices) == 1 && childServices[0] == "") {
			mutex.Lock()
			if _, isIn := data[id]; isIn {
				delete(data, id)
				mutex.Unlock()
				c.Status(http.StatusOK)
			} else {
				mutex.Unlock()
				c.Status(http.StatusBadRequest)
			}
		} else {
			body := generateDeleteBody(childServices, id)
			reqId := xid.New().String()
			resp, err := http.Post(coordinatorAddr + reqId, "application/json", bytes.NewBuffer(body))
			if err == nil && resp.StatusCode >= 200 && resp.StatusCode <= 299 {
				c.Status(http.StatusOK)
			} else if err != nil {
				c.Status(http.StatusInternalServerError)
			} else {
				c.Status(resp.StatusCode)
			}
		}
	})

	if err := r.Run(":8888"); err != nil {
		panic(err)
	}
}
