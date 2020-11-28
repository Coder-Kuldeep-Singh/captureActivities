package main

import (
	"fmt"
	"log"
	"os/user"
	"strings"
	"time"

	"github.com/go-vgo/robotgo"
)

// ClickedInfos holds the clicked info
type ClickedInfos struct {
	ResolutionCoordinates, ClickedCoordinates                                                             *Coordinates
	ClickedFullTime, ClickedDay, RunningApplication, CapturedTime, CapturedYearMonth, CapturedCurrentDate string
}

// Coordinates holds the cord
type Coordinates struct {
	X, Y int
}

func getLMouseClick() *ClickedInfos {
	mleft := robotgo.AddEvent("mleft")
	if mleft {
	}
	click := time.Now()
	time.Sleep(1 * time.Millisecond)

	// log.Println()
	x, y := robotgo.GetScreenSize()
	cx, cy := robotgo.GetMousePos()
	// log.Println(x, y)
	return &ClickedInfos{
		ResolutionCoordinates: &Coordinates{
			X: x,
			Y: y,
		},
		ClickedCoordinates: &Coordinates{
			X: cx,
			Y: cy,
		},
		ClickedFullTime:     fmt.Sprintf("%d-%d-%d:%d:%d:%d", click.Day(), click.Month(), click.Year(), click.Hour(), click.Minute(), click.Second()),
		ClickedDay:          fmt.Sprintf("%s", click.Month().String()),
		RunningApplication:  getUsedProduct(),
		CapturedTime:        fmt.Sprintf("%d:%d:%d", click.Hour(), click.Minute(), click.Second()),
		CapturedYearMonth:   fmt.Sprintf("%d-%d", click.Month(), click.Year()),
		CapturedCurrentDate: fmt.Sprintf("%d-%d-%d", click.Day(), click.Month(), click.Year()),
	}
}

func captureClicks() *ClickedInfos {
	return getLMouseClick()
}

func getUsedProduct() string {
	pid := robotgo.GetPID()
	name, err := robotgo.FindPath(pid)
	if err != nil {
		log.Println(err)
	}
	names := strings.Split(name, "/")
	// log.Println(pid)
	// log.Println(name)
	return names[len(names)-1]
}

// Users handle the list of all users will gonna use tool
type Users struct {
	UserID, Username, HomeDirectory string
}

func getUserInfo() *Users {
	user := gatherUserInfo()
	return &Users{
		UserID:        getUserID(user),
		Username:      getUserName(user),
		HomeDirectory: getUserHomeDir(user),
	}
}

func gatherUserInfo() *user.User {
	user, err := user.Current()
	if err != nil {
		log.Println("error to gather user info:")
		log.Println(err)
	}
	return user
}

func getUserID(user *user.User) string {
	return user.Uid
}

func getUserName(user *user.User) string {
	return user.Username
}
func getUserHomeDir(user *user.User) string {
	return user.HomeDir
}
