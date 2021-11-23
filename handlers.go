package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/dominikbraun/timetrace/core"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func DisplayLanding(c *gin.Context) {
	var page PageData
	session := sessions.Default(c)
	//options := session.Options
	//fmt.Println("---------options", &options.MaxAge)
	page.Init("main", c)
	session.Set("message", "")
	session.Save()
	c.HTML(http.StatusOK, "layout", page)
}

func StartStop(c *gin.Context) {
	var page PageData
	session := sessions.Default(c)
	page.Init("main", c)
	action := c.PostForm("action")
	project := c.PostForm("project")
	var err error
	var record *core.Record
	if action == "start" {
		if page.Tracking {
			fmt.Println("Currently Tracking ", page.CurrentProject, " need to stop first ", page.CurrentSession)
			//need to check that current project has been tracked for at least one minute
			//no need to track time of less than minute and allow creation of new record
			if page.CurrentSession == "0h 0min" {
				record, err = timetrace.LoadRecord(time.Now())
				err = timetrace.DeleteRecord(*record)
			} else {
				err = timetrace.Stop()
			}
		}
		err = timetrace.Start(project, true, []string{})

	} else if action == "stop" {
		err = timetrace.Stop()
	} else {
		err = errors.New("invalid request")
	}
	if err != nil {
		fmt.Println("------errors ", err)
		session.Set("message", err.Error())
	} else {
		session.Set("message", "")
	}
	session.Save()
	location := url.URL{Path: "/"}
	c.Redirect(http.StatusFound, location.RequestURI())
}

func CreateProject(c *gin.Context) {
	var project core.Project
	project.Key = c.PostForm("name")
	session := sessions.Default(c)
	if err := timetrace.SaveProject(project, false); err != nil {
		session.Set("message", err.Error())
	} else {
		session.Set("message", "New Project "+project.Key+" Created")
	}
	session.Save()
	location := url.URL{Path: "/"}
	c.Redirect(http.StatusFound, location.RequestURI())
}

func DeleteProject(c *gin.Context) {
	var project core.Project
	success := true
	project.Key = c.PostForm("project")
	deleteRecords := c.PostForm("records")
	session := sessions.Default(c)
	if err := timetrace.BackupProject(project.Key); err != nil {
		session.Set("message", err.Error())
		success = false
	}
	if deleteRecords == "on" && success {
		if err := timetrace.DeleteRecordsByProject(project.Key); err != nil {
			session.Set("message", err.Error())
			success = false
		}
	}
	if success {
		if err := timetrace.DeleteProject(project); err != nil {
			session.Set("message", err.Error())
		} else {
			session.Set("message", "Project "+project.Key+" has been Deleted")
		}
	}
	session.Save()
	location := url.URL{Path: "/"}
	c.Redirect(http.StatusFound, location.RequestURI())
}

func NewUser(c *gin.Context) {
	var user Users
	user.Username = c.PostForm("user")
	user.Password = c.PostForm("pass")
	user.IsAdmin = true
	if err := SaveUser(user); err != nil {
		log.Fatal("err saving user: ", err)
	}

	c.HTML(http.StatusOK, "Login", nil)
}

func ProcessLogin(c *gin.Context) {
	var user Users
	user.Username = c.PostForm("user")
	user.Password = c.PostForm("pass")
	valid, isadmim, err := ValidateUser(user)
	if err != nil {
		c.HTML(http.StatusBadRequest, "Login", gin.H{"message": err.Error()})
	}
	if valid {
		session := sessions.Default(c)
		session.Set("loggedIn", true)
		session.Set("message", "")
		if isadmim {
			session.Set("admin", true)
		}
		sessions.Default(c).Options(sessions.Options{MaxAge: 28800, Secure: true, SameSite: http.SameSiteLaxMode})
		session.Save()
		location := url.URL{Path: "/"}
		c.Redirect(http.StatusFound, location.RequestURI())
	}
}

func Logout(c *gin.Context) {
	session := sessions.Default(c)
	session.Set("loggedIn", false)
	session.Save()
	location := url.URL{Path: "/"}
	c.Redirect(http.StatusFound, location.RequestURI())
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 4)
	return string(bytes), err
}

func SaveUser(user Users) error {
	home := os.Getenv("HOME")
	f, err := os.Create(home + "/.timetrace/users.json")
	if err != nil {
		return err
	}
	user.Password, err = hashPassword(user.Password)
	if err != nil {
		return err
	}
	b, err := json.Marshal(user)
	if err != nil {
		return err
	}
	_, err = f.Write(b)
	if err != nil {
		return err
	}
	return nil
}

func ValidateUser(visitor Users) (bool, bool, error) {
	var user Users
	home := os.Getenv("HOME")
	f, err := os.Open(home + "/.timetrace/users.json")
	if err != nil {
		return false, false, err
	}
	decoder := json.NewDecoder(f)
	for decoder.More() {
		err = decoder.Decode(&user)
		if err != nil {
			return false, false, err
		}
		if visitor.Username == user.Username && CheckPassword(visitor, user) {
			if user.IsAdmin {
				return true, true, nil
			}
			return true, false, nil
		}
		return false, false, errors.New("Invalid username or password")
	}
	//shouldn't get here
	return false, false, nil
}

func CheckPassword(plain, hash Users) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash.Password), []byte(plain.Password))
	return err == nil
}

func ProcessError(c *gin.Context, err error, status int) {
	session := sessions.Default(c)
	session.Set("message", err.Error())
	session.Save()
	location := url.URL{Path: "/"}
	c.Redirect(status, location.RequestURI())
}

func GenerateReport(c *gin.Context) {
	var startDate, endDate time.Time
	var err error
	session := sessions.Default(c)
	start := c.PostForm("start")
	end := c.PostForm("end")
	billable := c.PostForm("billable")
	project := c.PostForm("project")

	startDate, err = timetrace.Formatter().ParseDate(start)
	if err != nil {
		ProcessError(c, err, http.StatusBadRequest)
	}
	//set up filters
	endDate, err = timetrace.Formatter().ParseDate(end)
	if err != nil {
		ProcessError(c, err, http.StatusBadRequest)
	}
	var filter = []func(*core.Record) bool{
		core.FilterNoneNilEndTime,
		core.FilterByTimeRange(startDate, endDate),
	}
	if project != "" {
		filter = append(filter, core.FilterByProject(project))
	}
	if billable == "billable" {
		filter = append(filter, core.FilterBillable(true))
	}
	if billable == "nonbillable" {
		filter = append(filter, core.FilterBillable(false))
	}
	//get raw report
	raw, err := timetrace.Report(filter...)
	if err != nil {
		ProcessError(c, err, http.StatusInternalServerError)
	}
	//convert raw report to Indented JSON (some fields of raw not exported
	output, err := raw.Json()
	if err != nil {
		ProcessError(c, err, http.StatusInternalServerError)
	}
	var rawdata = make(map[string]interface{})
	err = json.Unmarshal(output, &rawdata)
	if err != nil {
		ProcessError(c, err, http.StatusInternalServerError)
	}
	var reports []Report
	for key, record := range rawdata {
		jsonbody, err := json.Marshal(record)
		if err != nil {
			ProcessError(c, err, http.StatusInternalServerError)
		}
		var report Report
		err = json.Unmarshal(jsonbody, &report)
		if err != nil {
			ProcessError(c, err, http.StatusInternalServerError)
		}
		report.Project = key
		report.Sum = timetrace.Formatter().FormatDuration(report.Total)
		reports = append(reports, report)
	}
	session.Set("message", "")
	session.Save()
	c.HTML(http.StatusOK, "ReportData", reports)
}
