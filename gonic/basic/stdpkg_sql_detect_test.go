package basic_test

import (
	"database/sql"
	"fmt"
	"sync"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mxk/go-sqlite/sqlite3"
)

var (
	mysqlAvailable bool    = true
	mutex                  = sync.Mutex{}
	db             *sql.DB = nil
)

func MysqlDetection(ticker *time.Ticker) {
	for {
		select {
		case <-ticker.C:
			if e := db.Ping(); e != nil {
				fmt.Println("Ping err", e)

				mutex.Lock()
				mysqlAvailable = false
				mutex.Unlock()

			} else {
				fmt.Println("status ok")
			}
		}
	}
}

func MysqlSwitch() {
	for {
		mutex.Lock()
		if !mysqlAvailable {
			fmt.Println("Switch Sqlite3")
			db, _ = sql.Open("sqlite3", "./foo.db")
			mysqlAvailable = true
		}
		mutex.Unlock()

		time.Sleep(time.Second * 4)
	}
}

func Test_SQL_detect(t *testing.T) {
	c := make(chan bool)

	db, _ = sql.Open("mysql", "yeqiang:yeqiang@/test_yeqiang")
	ticker := time.NewTicker(time.Second * 2)

	go MysqlDetection(ticker)
	go MysqlSwitch()

	<-c
}
