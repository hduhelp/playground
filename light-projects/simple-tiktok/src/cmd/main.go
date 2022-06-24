package main

import (
	"sync"

	"github.com/gin-gonic/gin"
)

var once sync.Once

func init() {
	once.Do(func() {
	})
}

func main() {
	r := gin.Default()

	handle(r)

	// r.Run() 需要指定运行端口
}
