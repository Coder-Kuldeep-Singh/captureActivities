package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/joho/godotenv"
)

// func executeCronJob() {
// 	gocron.Every(1).Seconds().Do(CaptureScreen)
// 	<-gocron.Start()
// }

func createRepo() string {
	name := "storage"
	_, err := os.Stat(name)
	// if err != nil {
	// 	log.Println(err)
	// }
	if os.IsNotExist(err) {
		errDir := os.MkdirAll(name, 0755)
		if errDir != nil {
			log.Println(err)
		}

	}
	return name
}

func generateFullFilePath(filepath, filename string) string {
	return fmt.Sprintf("%s/%s", filepath, filename)
}

func remove(filepath string) {
	err := os.Remove(filepath)
	if err != nil {
		log.Println("error to remove file/directory")
		return
	}
}

func handleForcecontrol(folderName string) {
	ctx := context.Background()

	// trap Ctrl+C and call cancel on the context
	ctx, cancel := context.WithCancel(ctx)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	defer func() {
		signal.Stop(c)
		cancel()

	}()
	go func() {
		select {
		case <-c:
			cancel()
			remove(folderName)
			os.Exit(1)
		case <-ctx.Done():
		}
	}()
}

func runContinuesly(folderName string) {
	ticker := time.NewTicker(15 * time.Minute)

	go func() {
		for {
			select {
			case <-ticker.C:
				//Call the periodic function here.
				fileName := CaptureScreen()
				mailsPrepared(folderName, fileName)
				remove(generateFullFilePath(folderName, fileName))
			}
		}
	}()

	quit := make(chan bool, 1)
	// main will continue to wait until there is an entry in quit
	<-quit

}

func main() {
	// go executeCronJob()
	// CaptureScreen()
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file -> ", err)
		return
	}
	folderName := createRepo()
	log.Println(folderName)
	handleForcecontrol(folderName)
	now := time.Now()
	hour := now.Hour()
	if hour >= 10 && hour < 18 {
		fmt.Println("Running..")
		runContinuesly(folderName)
	}
	// else {
	// 	fmt.Println("sleeping..")
	// 	log.Println("sleeping for ", time.Hour)
	// 	// time.Sleep(time.Hour)
	// }
}
