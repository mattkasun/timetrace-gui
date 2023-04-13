package main

import (
	"embed"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/memstore"
	"github.com/gin-gonic/gin"
	"github.com/mattkasun/timetrace-gui/database"
)

//go:embed: images/favicon.ico
var icon embed.FS

//go:embed html/* images/*
var f embed.FS

func PrintDuration(d time.Duration) string {
	return fmt.Sprint(d)
}

func main() {
	//config := config.Get()
	//file := fs.New(config)
	if err := database.InitializeDatabase(); err != nil {
		log.Fatal("could not init database ", err)
	}
	//timetrace = core.New(config, file)
	//if err := timetrace.EnsureDirectories(); err != nil {
	//log.Fatal("unable to create dirs", err)
	//}
	router := SetupRouter()
	router.Run(":8090")
}

func SetupRouter() *gin.Engine {
	funcMap := template.FuncMap{
		"printDuration": PrintDuration,
	}
	store := memstore.NewStore([]byte("secret"))
	session := sessions.Sessions("timetrace", store)
	options := sessions.Options{MaxAge: 7200}
	fmt.Println("options\n", options)
	router := gin.Default()
	router.Use(session)
	//router.LoadHTMLGlob("html/*")
	templates := template.Must(template.New("").Funcs(funcMap).ParseFS(f, "html/*"))
	router.SetHTMLTemplate(templates)
	//router.StaticFile("favicon.ico", "./images/favicon.ico")
	router.StaticFS("/favicon.ico", http.FS(icon))
	//router.Static("images", "./images")
	router.GET("/images/*filepath", func(c *gin.Context) {
		c.FileFromFS(c.Request.URL.Path, http.FS(f))
	})
	//router.StaticFS("/images", http.FS(f))
	router.POST("/newuser", NewUser)
	router.POST("/login", ProcessLogin)
	router.GET("/logout", Logout)
	restricted := router.Group("/", AuthRequired)
	{
		restricted.GET("/", DisplayLanding)
		restricted.POST("/", StartStop)
		restricted.POST("/create_project", CreateProject)
		restricted.POST("/delete_project", DeleteProject)
		restricted.POST("/reports", GenerateReport)
		restricted.POST("/edit", EditRecord)
	}

	api := router.Group("/api/v1")
	{
		api.GET("/projects", GetProjects)
	}

	return router
}

func AuthRequired(c *gin.Context) {
	session := sessions.Default(c)
	loggedIn := session.Get("loggedIn")
	fmt.Println("loggedIn status: ", loggedIn)
	if err := database.InitializeDatabase(); err != nil {
		log.Println(err)
	}
	if loggedIn != true {
		_, err := database.GetAllUsers()
		if err != nil {
			if errors.Is(err, database.ErrNoResults) {
				//if err.Error() == "no results found" {
				log.Println("no users", err)
				c.HTML(http.StatusOK, "New", gin.H{"message": err})
				c.Abort()
				return
			}
			log.Println("getusers error", err)
			c.HTML(http.StatusInternalServerError, "Login", gin.H{"message": "database err"})
			c.Abort()
			return
		}
		c.HTML(http.StatusOK, "Login", nil)
		c.Abort()
		return
	}
	fmt.Println("authorized - good to go")
}
