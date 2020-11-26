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
	ResolutionCoordinates, ClickedCoordinates   *Coordinates
	ClickedTime, ClickedDay, RunningApplication string
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
		ClickedTime:        fmt.Sprintf("%d-%d-%d:%d:%d:%d", click.Day(), click.Month(), click.Year(), click.Hour(), click.Minute(), click.Second()),
		ClickedDay:         fmt.Sprintf("%s", click.Month().String()),
		RunningApplication: getUsedProduct(),
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
