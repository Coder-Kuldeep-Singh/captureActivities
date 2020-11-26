package main

import (
	"fmt"
	"log"
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
