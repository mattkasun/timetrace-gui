package main

import (
	"fmt"
	"time"

	"github.com/dominikbraun/timetrace/core"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

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
	Message            string
	Tracking           bool
	CurrentProject     string
	CurrentSession     string
	CurrentProjectTime string
	Today              string
	Breaks             string
	Projects           []*core.Project
	Summary            map[string]string
}

func (data *PageData) Init(page string, c *gin.Context) {
	session := sessions.Default(c)
	data.Message = session.Get("message").(string)
	data.Page = page
	data.CurrentProject = "---"
	data.CurrentSession = "---"
	data.CurrentProjectTime = "---"
	data.Tracking = false
	status, err := timetrace.Status()
	fmt.Println("error: ", err)
	if err == nil {
		fmt.Println("populating data", status.Current, status.TrackedTimeCurrent)
		if status.Current != nil {
			data.CurrentProject = status.Current.Project.Key
			data.Tracking = true
		}
		if status.TrackedTimeCurrent != nil {
			data.CurrentSession = timetrace.Formatter().FormatDuration(*status.TrackedTimeCurrent)
		}
		//Get Time worked today on Current Project
		records, err := timetrace.ListRecords(time.Now())
		var elapsed time.Duration
		m := make(map[string]time.Duration)
		s := make(map[string]string)
		data.Summary = s
		if err == nil {
			for _, record := range records {
				//Get Time worked today on Current Project
				if record.Project.Key == data.CurrentProject {
					elapsed = elapsed + record.Duration()
					fmt.Println(elapsed, record.Duration())
				}
				m[record.Project.Key] = m[record.Project.Key] + record.Duration()
			}
			fmt.Println(m)
			for key, value := range m {
				fmt.Println("key:value ", key, value)
				data.Summary[key] = timetrace.Formatter().FormatDuration(value)
			}
			data.CurrentProjectTime = timetrace.Formatter().FormatDuration(elapsed)
		}
		data.Today = timetrace.Formatter().FormatDuration(status.TrackedTimeToday)
		data.Breaks = timetrace.Formatter().FormatDuration(status.BreakTimeToday)

	}

	//get all projects
	data.Projects, err = timetrace.ListProjects()
	if err != nil {
		data.Projects = []*core.Project{}
	}

	fmt.Println(data.Projects)

}
