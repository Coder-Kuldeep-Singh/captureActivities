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

// func genQuery(db *sql.DB, query string) *sql.Rows {
// 	res, err := db.Query(query)
// 	if err != nil {
// 		fmt.Printf("error to run query %s\n", query)
// 		fmt.Println(err)
// 		return nil
// 	}
// 	return res
// }

// func getDatabaseList(db *sql.DB, query string) []string {
// 	rows := genQuery(db, query)
// 	return getList(rows)

// }

func execute(db *sql.DB, query string) {
	exec, err := db.Exec(query)
	if err != nil {
		fmt.Printf("error to execute query %s\n", query)
		fmt.Println(err)
		return
	}
	log.Println(exec.RowsAffected())
}

func uploadClick(db *sql.DB, click *ClickedInfos) {
	log.Println(click)
	query := fmt.Sprintf("insert into clicks(screenX,screenY,capturedCoordinateX,capturedCoordinateY,capturedTime,capturedDay) values(%d,%d,%d,%d,%s,%s)", click.ResolutionCoordinates.X, click.ResolutionCoordinates.Y, click.ClickedCoordinates.X, click.ClickedCoordinates.Y, click.ClickedTime, click.ClickedDay)
	execute(db, query)
}
