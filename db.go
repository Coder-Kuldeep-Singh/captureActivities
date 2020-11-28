package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq" //import postgres driver
)

// DBConfig hold the database connection values
type DBConfig struct {
	Host, Password, Name, User, Port string
}

func loadEnv() *DBConfig {
	return &DBConfig{
		Host:     configString("DBHOST"),
		Password: configString("DBPASSWORD"),
		Name:     configString("DBNAME"),
		User:     configString("DBUSER"),
		Port:     configString("DBPORT"),
	}
}

func configString(name string) string {
	return os.Getenv(name)
}

func (db *DBConfig) dcs(connectionType string) string {
	return fmt.Sprintf("%s:%s@tcp(%s)/%s", db.User, db.Password, db.Host, db.Name)

}

func (db *DBConfig) connect(connectionType string) *sql.DB {
	dbs, err := sql.Open(toLower(connectionType), db.dcs(toLower(connectionType)))
	if err != nil {
		fmt.Printf("Error %s when opening DB\n", err)
		log.Fatalln(err)
	}
	setLimits(dbs)
	return dbs
}
func toLower(str string) string {
	return strings.ToLower(str)
}

func setLimits(db *sql.DB) {
	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(20)
	db.SetConnMaxLifetime(time.Minute * 5)
}

// func getList(rows *sql.Rows) []string {
// 	List := []string{}
// 	for rows.Next() {
// 		var values string
// 		err := rows.Scan(&values)
// 		if err != nil {
// 			fmt.Println("error to scan rows.")
// 			fmt.Println(err)
// 			return nil
// 		}
// 		List = append(List, values)
// 	}
// 	return List
// }

// func postInfo()

func genQuery(db *sql.DB, query string) *sql.Rows {
	res, err := db.Query(query)
	if err != nil {
		fmt.Printf("error to run query %s\n", query)
		fmt.Println(err)
		return nil
	}
	return res
}

// func getDatabaseList(db *sql.DB, query string) []string {
// 	rows := genQuery(db, query)
// 	return getList(rows)

// }

func execute(db *sql.DB, query string) {
	_, err := db.Exec(query)
	if err != nil {
		fmt.Printf("error to execute query %s\n", query)
		fmt.Println(err)
		return
	}
	// log.Println(exec.RowsAffected())
}

func uploadClick(db *sql.DB, click *ClickedInfos) {
	// log.Println(click)
	query := fmt.Sprintf("insert into clicks(screenX,screenY,captureCoordinateX,captureCoordinateY,capturedTime,capturedDay,running_application,captured_time,captured_year_month,currentdate) values(%d,%d,%d,%d,'%s','%s','%s','%s','%s','%s')", click.ResolutionCoordinates.X, click.ResolutionCoordinates.Y, click.ClickedCoordinates.X, click.ClickedCoordinates.Y, click.ClickedFullTime, click.ClickedDay, click.RunningApplication, click.CapturedTime, click.CapturedYearMonth, click.CapturedCurrentDate)
	execute(db, query)
}

// DailyGraph holds the daily clicked and used products
type DailyGraph struct {
	Product string
	Count   int
	Days    string
}

func RunningTime() (string, string) {
	return "10:0:0", "18:0:0"
}

func getUsedProductPerDay(db *sql.DB) *[]DailyGraph {
	t := time.Now()
	day := fmt.Sprintf("%d-%d-%d", t.Day(), t.Month(), t.Year())
	startTime, endTime := RunningTime()
	query := fmt.Sprintf("SELECT distinct(running_application) from clicks where currentdate='%s' and captured_time >= '%s' and captured_time <= '%s'", day, startTime, endTime)
	rows := genQuery(db, query)
	daily := []DailyGraph{}
	for rows.Next() {
		var values string
		err := rows.Scan(&values)
		if err != nil {
			fmt.Println("error to scan rows.")
			fmt.Println(err)
			return nil
		}
		query := fmt.Sprintf("SELECT count(running_application) from clicks where currentdate='%s' and running_application='%s' and captured_time >= '%s' and captured_time <= '%s'", day, values, startTime, endTime)
		daily = append(daily, DailyGraph{
			Product: values,
			Count:   getProductUsedCount(db, query),
		})
	}
	// log.Println(daily)
	return &daily
}

func getProductUsedCount(db *sql.DB, query string) int {
	row := db.QueryRow(query)
	var count int
	err := row.Scan(&count)
	if err != nil {
		log.Println("error to scan count row ", err)
		return 0
	}
	return count
}

func getUsedProductPerDays(db *sql.DB) *[]DailyGraph {
	// query := fmt.Sprintf("SELECT distinct(running_application) from clicks where currentdate='%s'", day)
	query := fmt.Sprintf("SELECT distinct(currentdate) from clicks")
	rows := genQuery(db, query)
	daily := []DailyGraph{}
	for rows.Next() {
		var values string
		err := rows.Scan(&values)
		if err != nil {
			fmt.Println("error to scan rows.")
			fmt.Println(err)
			return nil
		}
		if len(values) == 0 {
			continue
		}

		startTime, endTime := RunningTime()
		query := fmt.Sprintf("SELECT count(running_application) from clicks where currentdate='%s' and captured_time >='%s' and captured_time <= '%s'", values, startTime, endTime)
		daily = append(daily, DailyGraph{
			Count: getProductUsedCount(db, query),
			Days:  values,
		})
	}
	return &daily
}

func getUsedProductPerDaysFull(db *sql.DB) *[]DailyGraph {
	// query := fmt.Sprintf("SELECT distinct(running_application) from clicks where currentdate='%s'", day)
	query := fmt.Sprintf("SELECT distinct(currentdate) from clicks")
	rows := genQuery(db, query)
	daily := []DailyGraph{}
	for rows.Next() {
		var values string
		err := rows.Scan(&values)
		if err != nil {
			fmt.Println("error to scan rows.")
			fmt.Println(err)
			return nil
		}
		if len(values) == 0 {
			continue
		}

		query := fmt.Sprintf("SELECT count(running_application) from clicks where currentdate='%s'", values)
		daily = append(daily, DailyGraph{
			Count: getProductUsedCount(db, query),
			Days:  values,
		})
	}
	return &daily
}

func getUsedProductPerDayFull(db *sql.DB) *[]DailyGraph {
	t := time.Now()
	day := fmt.Sprintf("%d-%d-%d", t.Day(), t.Month(), t.Year())
	query := fmt.Sprintf("SELECT distinct(running_application) from clicks where currentdate='%s'", day)
	rows := genQuery(db, query)
	daily := []DailyGraph{}
	for rows.Next() {
		var values string
		err := rows.Scan(&values)
		if err != nil {
			fmt.Println("error to scan rows.")
			fmt.Println(err)
			return nil
		}
		query := fmt.Sprintf("SELECT count(running_application) from clicks where currentdate='%s' and running_application='%s'", day, values)
		daily = append(daily, DailyGraph{
			Product: values,
			Count:   getProductUsedCount(db, query),
		})
	}
	// log.Println(daily)
	return &daily
}

// CreateUsers function generate user
func CreateUsers(db *sql.DB) {
	user := getUserInfo()
	query := fmt.Sprintf("insert into users(userid,username,user_home_directory) values('%s','%s','%s')", user.UserID, user.Username, user.HomeDirectory)
	execute(db, query)
}
