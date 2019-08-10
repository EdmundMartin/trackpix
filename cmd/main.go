package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/EdmundMartin/trackpix/pkg/db"
	"github.com/EdmundMartin/trackpix/pkg/pipeline"
	"github.com/EdmundMartin/trackpix/pkg/server"
	"github.com/gorilla/mux"
	"github.com/kshvakov/clickhouse"
)

func main() {
	connect, err := sql.Open("clickhouse", "tcp://127.0.0.1:9000?debug=true")
	if err := connect.Ping(); err != nil {
		if exception, ok := err.(*clickhouse.Exception); ok {
			fmt.Printf("[%d] %s \n%s\n", exception.Code, exception.Message, exception.StackTrace)
		} else {
			fmt.Println(err)
		}
		return
	}
	r := mux.NewRouter()
	tdCh := make(chan *server.TrackingData, 5000)
	dbCh := make(chan *pipeline.DataRecord, 5000)
	pipe := pipeline.TrxPipe{
		Functions: []func(*server.TrackingData, *pipeline.DataRecord) *pipeline.DataRecord{
			pipeline.LangFromHeader,
			pipeline.MapValues,
		},
	}
	pBytes, err := ioutil.ReadFile("pix.png")
	if err != nil {
		log.Fatal(err)
	}
	trx := &server.TrackingServer{Headers: []string{"User-Agent", "Accept-Language"}, Results: tdCh, PixelBytes: pBytes}
	go pipe.Run(tdCh, dbCh)
	go db.BulkProcess(connect, dbCh)
	r.HandleFunc("/sneaky.png", trx.Pixel).Methods("GET")
	http.ListenAndServe(":8000", r)
	fmt.Printf("running on :8000")
}
