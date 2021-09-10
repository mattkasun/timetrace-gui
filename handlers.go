package main

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/dominikbraun/timetrace/core"
	"github.com/gin-gonic/gin"
)

func DisplayLanding(c *gin.Context) {
	var page PageData

	page.Init("main")
	c.HTML(http.StatusOK, "layout", page)
}

func StartStop(c *gin.Context) {
	var page PageData
	page.Init("main")
	action := c.PostForm("action")
	project := c.PostForm("project")
	var err error
	var record *core.Record
	if action == "start" {
		fmt.Println("Staring Project ", project, "tracking status is ", page.Tracking)
		if page.Tracking {
			fmt.Println("Currently Tracking ", page.CurrentProject, " need to stop first ", page.CurrentTime)
			//need to check that current project has been tracked for at least one minute
			//no need to track time of less than minute and allow creation of new record
			if page.CurrentTime == "0h 0min" {
				fmt.Println("need to delete current record")
				record, err = timetrace.LoadRecord(time.Now())
				err = timetrace.DeleteRecord(*record)
			} else {
				err = timetrace.Stop()
			}
			if err != nil {
				fmt.Println("error stopping project ", page.CurrentProject, err)
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				c.Abort()
				return
			}
			fmt.Println(page.CurrentProject, " stop successfully")
		}
		err = timetrace.Start(project, true)
		fmt.Println("Started project ", project, err)
	} else if action == "stop" {
		fmt.Println("stop time tracking")
		err = timetrace.Stop()
		fmt.Println("time tracking stopped", err)
	} else {
		err = errors.New("invalid request")
	}
	fmt.Println("error", err)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.Abort()
		return
	}
	location := url.URL{Path: "/"}
	c.Redirect(http.StatusFound, location.RequestURI())
	//page.Init()
	//c.HTML(http.StatusOK, "layout", page)
}
