package main

import (
	"fmt"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/mattkasun/timetrace-gui/database"
	"github.com/mattkasun/timetrace-gui/models"
	"github.com/mattkasun/timetrace-gui/tracking"
)

var version = "v0.3"

type Report struct {
	Project string
	Records []models.Record
	Total   time.Duration
	Sum     string
}

type Users struct {
	Username string
	Password string
	IsAdmin  bool
}

type ProjectTime struct {
	Project string
	Time    string
}

type PageData struct {
	Page               string
	Version            string
	Message            string
	Tracking           bool
	CurrentProject     string
	CurrentSession     string
	CurrentProjectTime string
	Today              string
	Breaks             string
	Projects           []models.Project
	Summary            map[string]string
	DefaultDate        string
}

func (data *PageData) Init(page string, c *gin.Context) {
	data.Summary = make(map[string]string)
	var err error
	session := sessions.Default(c)
	data.Message = session.Get("message").(string)
	data.Page = page
	data.Version = version
	data.CurrentProject = "---"
	data.CurrentSession = "---"
	data.CurrentProjectTime = "---"
	data.Tracking = false
	data.DefaultDate = time.Now().Local().Format("2006-01-02")
	report := tracking.Status()
	data.CurrentProject = report.Project
	data.Tracking = true
	data.CurrentSession = report.Session.String()
	data.CurrentProjectTime = report.Today.String()
	data.Today = report.Total.String()
	data.Breaks = report.Breaks.String()
	for k, v := range report.Summary {
		data.Summary[k] = v.String()
	}
	//get all projects
	data.Projects, err = database.GetAllProjects()
	if err != nil {
		fmt.Println("error retrieving projects", err)
		data.Projects = []models.Project{}
	}
}
