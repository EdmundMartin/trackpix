package db

import (
	"database/sql"
	"log"

	"github.com/EdmundMartin/trackpix/pkg/pipeline"
)

func createTable(conn *sql.DB) {
	_, err := conn.Exec(`
		CREATE TABLE IF NOT EXISTS tracking (
			ip_addr 	  String,
			visited_page  String,
			visited_time  DateTime,
			user_agent    String,
			language      String
		) engine=Memory
	`)
	if err != nil {
		log.Fatal(err)
	}
}

func execCommand(conn *sql.DB, objs []*pipeline.DataRecord) {
	tx, _ := conn.Begin()
	stmt, _ := tx.Prepare("INSERT INTO tracking (ip_addr, visited_page, visited_time, user_agent, language) VALUES (?, ?, ?, ?, ?)")
	defer stmt.Close()
	for _, record := range objs {
		stmt.Exec(
			record.IPAddr,
			record.VisitedPage,
			record.Datetime,
			record.UserAgent,
			record.Language,
		)
	}
	err := tx.Commit()
	if err != nil {
		log.Printf("failed to insert data, error %s", err)
	}
}

// BulkProcess reads from the pipeline channel and writes to DB at set intervals
func BulkProcess(conn *sql.DB, data chan *pipeline.DataRecord) {
	createTable(conn)
	objs := []*pipeline.DataRecord{}
	for {
		datum := <-data
		objs = append(objs, datum)
		if len(objs) >= 100 {
			execCommand(conn, objs)
			objs = []*pipeline.DataRecord{}
		}
	}
}
