package main

import (
	"github.com/dominikbraun/timetrace/config"
	"github.com/dominikbraun/timetrace/core"
	"github.com/dominikbraun/timetrace/fs"
	"github.com/gin-gonic/gin"
)

var timetrace *core.Timetrace

func main() {
	config := config.Get()
	file := fs.New(config)
	timetrace = core.New(config, file)
	router := SetupRouter()
	router.Run("127.0.0.1:8090")
}

func SetupRouter() *gin.Engine {
	router := gin.Default()
	router.LoadHTMLGlob("html/*")
	router.StaticFile("favicon.ico", "./images/favicon.ico")
	router.Static("images", "./images")
	router.GET("/", DisplayLanding)
	router.POST("/", StartStop)

	return router
}
