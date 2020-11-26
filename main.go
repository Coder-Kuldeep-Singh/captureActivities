package main

import (
	"context"
	"database/sql"
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

func runContinueslyScreenShots(folderName string) {
	ticker := time.NewTicker(1 * time.Minute)

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

func runContinueslyclicksCapture(db *sql.DB) {
	ticker := time.NewTicker(1 * time.Second)
	go func() {
		for {
			select {
			case <-ticker.C:
				//Call the periodic function here.
				duration := captureClicks()
				getUsedProduct()
				uploadClick(db, duration)
				// go getUsedProduct()
			}
		}
	}()

}

// func capturerunningProcess() {
// 	ticker := time.NewTicker(1 * time.Second)
// 	go func() {
// 		for {
// 			select {
// 			case <-ticker.C:
// 				// Call the periodic function here.
// 				getUsedProduct()

// 			}
// 		}
// 	}()
// }

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
	dbType := "mysql"
	db := loadEnv().connect(dbType)
	if hour >= 10 && hour < 18 {
		fmt.Println("Running..")
		// capturerunningProcess()
		// getUsedProduct()
		runContinueslyclicksCapture(db)
		runContinueslyScreenShots(folderName)

	} else {
		// fmt.Println("sleeping..")
		// log.Println("sleeping for ", time.Hour)
		// time.Sleep(16 + time.Hour)
		log.Println("Capture screen Method are disabled until next office time.")
		runContinueslyclicksCapture(db)
	}
}
