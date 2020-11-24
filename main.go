package main

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

// func executeCronJob() {
// 	gocron.Every(1).Seconds().Do(CaptureScreen)
// 	<-gocron.Start()
// }

func createRepo() string {
	fileInfo, err := os.Stat("storage")

	if os.IsNotExist(err) {
		errDir := os.MkdirAll("storage", 0755)
		if errDir != nil {
			log.Fatal(err)
		}

	}
	return fileInfo.Name()
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
	ticker := time.NewTicker(1 * time.Minute)

	go func() {
		for {
			select {
			case <-ticker.C:
				//Call the periodic function here.
				fileName := CaptureScreen()
				mailsPrepared(folderName, fileName)
			}
		}
	}()

	quit := make(chan bool, 1)
	// main will continue to wait until there is an entry in quit
	<-quit
}
