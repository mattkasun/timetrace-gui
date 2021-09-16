package main

import (
	"fmt"
	"time"

	"github.com/dominikbraun/timetrace/core"
)

type Users struct {
	Username string
	Password string
	IsAdmin  bool
}

type PageData struct {
	Page               string
	Tracking           bool
	CurrentProject     string
	CurrentSession     string
	CurrentProjectTime string
	Today              string
	Breaks             string
	Projects           []*core.Project
}

func (data *PageData) Init(page string) {
	//config := config.Get()
	//file := fs.New(config)
	//timetrace := core.New(config, file)
	//current Project
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
		if err == nil {
			for _, record := range records {
				if record.Project.Key == data.CurrentProject {
					elapsed = elapsed + record.Duration()
					fmt.Println(elapsed, record.Duration())
				}
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
