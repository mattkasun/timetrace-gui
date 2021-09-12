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
	//session := sessions.Default(c)
	//options := session.Options
	//fmt.Println("---------options", &options.MaxAge)
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
		if isadmim {
			session.Set("admin", true)
		}
		sessions.Default(c).Options(sessions.Options{MaxAge: 70})
		session.Save()
		location := url.URL{Path: "/"}
		c.Redirect(http.StatusFound, location.RequestURI())
	}
}

func Logout(c *gin.Context) {
	session := sessions.Default(c)
	//sessions.Default(c).Options(sessions.Options{MaxAge: -1})
	session.Options(sessions.Options{MaxAge: -1})
	session.Save()
	c.JSON(http.StatusOK, gin.H{"message": "Successful logout"})
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
