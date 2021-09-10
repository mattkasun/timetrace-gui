package main

import (
	"fmt"

	"github.com/dominikbraun/timetrace/core"
)

type PageData struct {
	Page           string
	Tracking       bool
	CurrentProject string
	CurrentTime    string
	Today          string
	Breaks         string
	Projects       []*core.Project
}

func (data *PageData) Init(page string) {
	//config := config.Get()
	//file := fs.New(config)
	//timetrace := core.New(config, file)
	//current Project
	data.Page = page
	data.CurrentProject = "---"
	data.CurrentTime = "---"
	data.Tracking = false
	report, err := timetrace.Status()
	fmt.Println(report)
	fmt.Println("error: ", err)
	if err == nil {
		if report.Current != nil {
			data.CurrentProject = report.Current.Project.Key
			data.Tracking = true
		}
		if report.TrackedTimeCurrent != nil {
			data.CurrentTime = timetrace.Formatter().FormatCurrentTime(report)
		}
		data.Today = timetrace.Formatter().FormatDuration(report.TrackedTimeToday)
		data.Breaks = timetrace.Formatter().FormatDuration(report.BreakTimeToday)
	}
	//get all projects
	data.Projects, err = timetrace.ListProjects()
	if err != nil {
		data.Projects = []*core.Project{}
	}
	fmt.Println(data.Projects)

}
