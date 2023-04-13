package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/kr/pretty"
	"github.com/mattkasun/timetrace-gui/database"
	"github.com/mattkasun/timetrace-gui/models"
	"github.com/mattkasun/timetrace-gui/tracking"
	"golang.org/x/crypto/bcrypt"
)

func DisplayLanding(c *gin.Context) {
	var page PageData
	session := sessions.Default(c)
	//options := session.Options
	//fmt.Println("---------options", &options.MaxAge)
	page.Init("main", c)
	pretty.Println(page)
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
	if action == "start" {
		if err := tracking.Start(project); err != nil {
			fmt.Println("err starting", project, err)
		}

	} else if action == "stop" {
		err = tracking.Stop()
	} else {
		err = errors.New("invalid request")
	}
	if err != nil {
		fmt.Println("stop command errors ", err)
		session.Set("message", err.Error())
	} else {
		session.Set("message", "")
	}
	session.Save()
	location := url.URL{Path: "/"}
	c.Redirect(http.StatusFound, location.RequestURI())
}

func CreateProject(c *gin.Context) {
	project := models.Project{}
	project.Name = c.PostForm("name")
	project.ID = uuid.New()
	session := sessions.Default(c)
	if err := database.SaveProject(&project); err != nil {
		session.Set("message", err.Error())
	} else {
		session.Set("message", "New Project "+project.Name+" Created")
	}
	session.Save()
	location := url.URL{Path: "/"}
	c.Redirect(http.StatusFound, location.RequestURI())
}

func DeleteProject(c *gin.Context) {
	success := true
	name := c.PostForm("project")
	deleteRecords := c.PostForm("records")
	session := sessions.Default(c)
	if err := tracking.BackupProject(name); err != nil {
		session.Set("message", err.Error())
		success = false
	}
	if deleteRecords == "on" && success {
		if err := tracking.DeleteRecordsByProject(name); err != nil {
			session.Set("message", err.Error())
			success = false
		}
	}
	if success {
		if err := database.DeleteProject(name); err != nil {
			session.Set("message", err.Error())
		} else {
			session.Set("message", "Project "+name+" has been Deleted")
		}
	}
	session.Save()
	location := url.URL{Path: "/"}
	c.Redirect(http.StatusFound, location.RequestURI())
}

func NewUser(c *gin.Context) {
	var user models.User
	var err error
	user.Username = c.PostForm("user")
	password := c.PostForm("pass")
	user.Admin = true
	user.Password, err = hashPassword(password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err})
		return
	}
	if err := SaveUser(&user); err != nil {
		log.Fatal("err saving user: ", err)
	}
	c.HTML(http.StatusOK, "Login", nil)
}

func ProcessLogin(c *gin.Context) {
	var user models.User
	user.Username = c.PostForm("user")
	user.Password = c.PostForm("pass")
	valid, isadmim, err := ValidateUser(&user)
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

func SaveUser(user *models.User) error {
	return database.SaveUser(user)
}

func ValidateUser(visitor *models.User) (bool, bool, error) {
	user, err := database.GetUser(visitor.Username)

	if err != nil {
		return false, false, err
	}
	if visitor.Username == user.Username && CheckPassword(visitor, &user) {
		if user.Admin {
			return true, true, nil
		}
		return true, false, nil
	}
	return false, false, errors.New("invalid username or password")
}

func CheckPassword(plain, hash *models.User) bool {
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
	//billable := c.PostForm("billable")
	project := c.PostForm("project")
	startDate, err = time.Parse("2006-01-02", start)
	if err != nil {
		ProcessError(c, err, http.StatusBadRequest)
	}
	//set up filters
	endDate, err = time.Parse("2006-01-02", end)
	if err != nil {
		ProcessError(c, err, http.StatusBadRequest)
	}
	pretty.Println(start, startDate, end, endDate)
	rows, err := database.GetAllrecords()
	if err != nil {
		ProcessError(c, err, http.StatusInternalServerError)
	}
	records := []models.Record{}
	for _, row := range rows {
		if row.Start.After(startDate) && row.Start.Before(endDate) {
			if project != "" && row.Project != project {
				continue
			}
			records = append(records, row)
		}
	}
	//get raw report
	//convert raw report to Indented JSON (some fields of raw not exported
	reportData := models.ConvertToReport(records)
	session.Set("message", "")
	session.Save()
	pretty.Println(records)
	c.HTML(http.StatusOK, "ReportData", reportData)
}

func EditRecord(c *gin.Context) {
	action := c.PostForm("action")
	if action == "update" {
		UpdateRecord(c)
		return
	}
	session := sessions.Default(c)
	record := c.PostForm("record")
	edit, err := database.GetRecord(record)
	if err != nil {
		ProcessError(c, err, http.StatusBadRequest)
	}
	session.Set("message", "")
	session.Save()
	c.HTML(http.StatusOK, "EditRecord", edit)
}

func UpdateRecord(c *gin.Context) {
	session := sessions.Default(c)
	id := c.PostForm("record")
	start := c.PostForm("start")
	end := c.PostForm("end")
	record, err := database.GetRecord(id)
	if err != nil {
		ProcessError(c, err, http.StatusBadRequest)
	}
	record.Start, err = time.Parse("2006-01-02-15-04-05", start)
	if err != nil {
		ProcessError(c, err, http.StatusBadRequest)
	}
	record.End, err = time.Parse("2006-01-02-15-04-05", end)
	if err != nil {
		ProcessError(c, err, http.StatusBadRequest)
	}
	if err := database.Saverecord(&record); err != nil {
		ProcessError(c, err, http.StatusBadRequest)
	}
	session.Set("message", "")
	session.Save()
	location := url.URL{Path: "/"}
	c.Redirect(http.StatusFound, location.RequestURI())
}

func GetProjects(c *gin.Context) {
	projects, err := database.GetAllProjects()
	if err != nil {
		ProcessError(c, err, http.StatusBadRequest)
		return
	}
	c.JSON(http.StatusOK, projects)
}
