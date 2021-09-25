package main

import (
	"embed"
	"fmt"
	"html/template"
	"net/http"
	"os"

	"github.com/dominikbraun/timetrace/config"
	"github.com/dominikbraun/timetrace/core"
	"github.com/dominikbraun/timetrace/fs"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/memstore"
	"github.com/gin-gonic/gin"
)

var timetrace *core.Timetrace

//go:embed html/* images/*
var f embed.FS

func main() {
	config := config.Get()
	file := fs.New(config)
	timetrace = core.New(config, file)
	router := SetupRouter()
	router.Run(":8090")
}

func SetupRouter() *gin.Engine {
	store := memstore.NewStore([]byte("secret"))
	session := sessions.Sessions("timetrace", store)
	options := sessions.Options{MaxAge: 7200}
	fmt.Println("options\n", options)
	router := gin.Default()
	router.Use(session)
	//router.LoadHTMLGlob("html/*")
	templates := template.Must(template.New("").ParseFS(f, "html/*"))
	router.SetHTMLTemplate(templates)
	//router.StaticFile("favicon.ico", "./images/favicon.ico")
	//router.Static("images", "./images")
	router.StaticFS("/images", http.FS(f))
	router.POST("/newuser", NewUser)
	router.POST("/login", ProcessLogin)
	router.GET("/logout", Logout)
	restricted := router.Group("/", AuthRequired)
	{
		restricted.GET("/", DisplayLanding)
		restricted.POST("/", StartStop)
		restricted.POST("/create_project", CreateProject)
		restricted.POST("/delete_project", DeleteProject)
	}

	return router
}

func AuthRequired(c *gin.Context) {
	session := sessions.Default(c)
	options := session.Options

	fmt.Println("checking authorization\n", options)
	fmt.Printf("type %T value %s\n", options, options)
	loggedIn := session.Get("loggedIn")
	fmt.Println("loggedIn status: ", loggedIn)
	if loggedIn != true {
		home := os.Getenv("HOME")
		_, err := os.Open(home + "/.timetrace/users.json")
		if err != nil {
			fmt.Println("error opening user file: ", err)
			c.HTML(http.StatusOK, "New", gin.H{"message": err})
			c.Abort()
			return
		}
		message := session.Get("error")
		fmt.Println("user exists --- message\n", message)
		c.HTML(http.StatusOK, "Login", gin.H{"messge": message})
		c.Abort()
		return
	}
	fmt.Println("authorized - good to go")
}
